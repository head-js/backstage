//go:build ignore

package main

import (
	servicecontext "context"
	servicejson "encoding/json"
	servicefmt "fmt"
	servicesort "sort"
	servicestrings "strings"

	serviceecho "github.com/labstack/echo/v5"
)

// Yaegi loads this file after the types, client, store, and translator, before main.go.
// Pass the client explicitly; do not reference main.go's singleton here.

func listWorkspaces(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	authorization string,
) ([]Workspace, error) {
	body, err := client.listProjects(ctx, authorization)
	if err != nil {
		return nil, err
	}

	var projectList VikunjaProjectList
	if err := servicejson.Unmarshal(body, &projectList); err != nil {
		return nil, err
	}

	workspaces := make([]Workspace, 0, len(projectList.Items))
	for _, project := range projectList.Items {
		if project.Id < 0 {
			continue
		}
		workspaces = append(workspaces, projectToWorkspace(project))
	}
	return workspaces, nil
}

func getVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeId string,
	authorization string,
) (*VisitNode, error) {
	project, task, ext, err := resolveVisitNodeTask(client, ctx, visitNodeId, authorization)
	if err != nil {
		return nil, err
	}
	return taskToVisitNode(*project, *task, visitNodeId, ext.Session)
}

func deleteVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeId string,
	authorization string,
) error {
	_, task, _, err := resolveVisitNodeTask(client, ctx, visitNodeId, authorization)
	if err != nil {
		return err
	}
	return client.deleteTask(ctx, task.Id, authorization)
}

// resolveVisitNodeTask resolves a stored VisitNode Id (tasks.ext_iid) to its
// project, task, and extension columns. The Id is looked up in the database,
// never parsed. Callers that need only the task ignore the taskExt return.
func resolveVisitNodeTask(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeId, authorization string,
) (*VikunjaProject, *VikunjaTask, taskExt, error) {
	taskId, ext, found, err := findTaskExtByIid(visitNodeId)
	if err != nil {
		return nil, nil, taskExt{}, err
	}
	if !found {
		return nil, nil, taskExt{}, serviceecho.NewHTTPError(404, servicefmt.Sprintf("visit node %q not found", visitNodeId))
	}

	body, err := client.getTask(ctx, taskId, authorization)
	if err != nil {
		return nil, nil, taskExt{}, err
	}
	var task VikunjaTask
	if err := servicejson.Unmarshal(body, &task); err != nil {
		return nil, nil, taskExt{}, err
	}
	project, err := findProjectById(client, ctx, task.ProjectId, authorization)
	if err != nil {
		return nil, nil, taskExt{}, err
	}
	return project, &task, ext, nil
}

func findProjectById(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	id int64,
	authorization string,
) (*VikunjaProject, error) {
	body, err := client.listProjects(ctx, authorization)
	if err != nil {
		return nil, err
	}
	var projectList VikunjaProjectList
	if err := servicejson.Unmarshal(body, &projectList); err != nil {
		return nil, err
	}
	for _, project := range projectList.Items {
		if project.Id == id {
			return &project, nil
		}
	}
	return nil, serviceecho.NewHTTPError(404, servicefmt.Sprintf("project %d not found", id))
}

func findProject(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	name, authorization string,
) (*VikunjaProject, error) {
	body, err := client.listProjects(ctx, authorization)
	if err != nil {
		return nil, err
	}
	var projectList VikunjaProjectList
	if err := servicejson.Unmarshal(body, &projectList); err != nil {
		return nil, err
	}
	for _, project := range projectList.Items {
		if project.Id >= 0 && servicestrings.EqualFold(project.Title, name) {
			return &project, nil
		}
	}
	return nil, serviceecho.NewHTTPError(404, servicefmt.Sprintf("project %q not found", name))
}

// storedExtIid returns the persisted VisitNode Id of a task. Every task is
// created through this plugin, which always writes ext_iid, so a missing or
// empty value means the data integrity invariant is broken. The other ext
// columns (e.g. ext_session) are optional and read as "" when unset.
func storedExtIid(exts map[int64]taskExt, taskId int64) (string, error) {
	if ext, ok := exts[taskId]; ok && ext.Iid != "" {
		return ext.Iid, nil
	}
	return "", servicefmt.Errorf("task %d has no ext_iid", taskId)
}

func createVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	input VisitNodeCreateRequest,
	authorization string,
) (*VisitNode, error) {
	input.Workspace = servicestrings.TrimSpace(input.Workspace)
	input.Title = servicestrings.TrimSpace(input.Title)
	if input.Workspace == "" || input.Title == "" {
		return nil, serviceecho.NewHTTPError(400, "Workspace and Title are required")
	}
	description, err := formatVisitNodeDescription(input.Description)
	if err != nil {
		return nil, serviceecho.NewHTTPError(400, err.Error())
	}
	project, err := findProject(client, ctx, input.Workspace, authorization)
	if err != nil {
		return nil, err
	}
	body, err := client.createTask(ctx, project.Id, input.Title, description, authorization)
	if err != nil {
		return nil, err
	}
	var task VikunjaTask
	if err := servicejson.Unmarshal(body, &task); err != nil {
		return nil, err
	}
	extIid := visitNodeId(project.Title, task.Index)
	if err := setTaskExtIid(task.Id, extIid); err != nil {
		return nil, err
	}
	input.Session = servicestrings.TrimSpace(input.Session)
	if input.Session != "" {
		if err := setTaskExtSession(task.Id, input.Session); err != nil {
			return nil, err
		}
	}
	return taskToVisitNode(*project, task, extIid, input.Session)
}

// updateVisitNode applies only the fields the caller sent (non-nil pointers)
// and returns the re-read node. Progress is a percentage integer (0-100);
// Session accepts any string, including "" to clear. Fields absent from the
// request are left untouched, so callers can update one field at a time.
func updateVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	input VisitNodeUpdateRequest,
	authorization string,
) (*VisitNode, error) {
	input.Id = servicestrings.TrimSpace(input.Id)
	if input.Id == "" {
		return nil, serviceecho.NewHTTPError(400, "Id is required")
	}
	if input.Progress == nil && input.Session == nil {
		return nil, serviceecho.NewHTTPError(400, "at least one of Progress or Session is required")
	}
	if input.Progress != nil && (*input.Progress < 0 || *input.Progress > 100) {
		return nil, serviceecho.NewHTTPError(400, "Progress must be a percentage integer between 0 and 100")
	}
	taskId, _, found, err := findTaskExtByIid(input.Id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, serviceecho.NewHTTPError(404, servicefmt.Sprintf("visit node %q not found", input.Id))
	}
	if input.Progress != nil {
		if err := client.updateTaskPercentDone(ctx, taskId, float64(*input.Progress)/100, authorization); err != nil {
			return nil, err
		}
	}
	if input.Session != nil {
		if err := setTaskExtSession(taskId, *input.Session); err != nil {
			return nil, err
		}
	}
	return getVisitNode(client, ctx, input.Id, authorization)
}

func listVisitNodes(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	authorization string,
) ([]VisitNode, error) {
	body, err := client.listProjects(ctx, authorization)
	if err != nil {
		return nil, err
	}
	var projectList VikunjaProjectList
	if err := servicejson.Unmarshal(body, &projectList); err != nil {
		return nil, err
	}
	exts, err := listExt()
	if err != nil {
		return nil, err
	}

	visitNodes := make([]VisitNode, 0)
	for _, project := range projectList.Items {
		if project.Id < 0 {
			continue
		}

		body, err := client.listProjectTasks(ctx, project.Id, authorization)
		if err != nil {
			return nil, err
		}
		var taskList VikunjaTaskList
		if err := servicejson.Unmarshal(body, &taskList); err != nil {
			return nil, err
		}

		for _, task := range taskList.Items {
			extIid, err := storedExtIid(exts, task.Id)
			if err != nil {
				return nil, err
			}
			description, err := translateTaskDescription(task.Description)
			if err != nil {
				description = map[string][]string{}
			}
			visitNodes = append(visitNodes, *visitNodeFromDescription(project, task, description, extIid, exts[task.Id].Session))
		}
	}

	return visitNodes, nil
}

func listVisitEdges(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	authorization string,
) ([]VisitEdge, error) {
	body, err := client.listProjects(ctx, authorization)
	if err != nil {
		return nil, err
	}
	var projectList VikunjaProjectList
	if err := servicejson.Unmarshal(body, &projectList); err != nil {
		return nil, err
	}
	projectNames := make(map[int64]string)
	for _, project := range projectList.Items {
		if project.Id >= 0 {
			projectNames[project.Id] = project.Title
		}
	}
	exts, err := listExt()
	if err != nil {
		return nil, err
	}

	visitEdges := make([]VisitEdge, 0)
	seen := make(map[VisitEdge]bool)
	for _, project := range projectList.Items {
		if project.Id < 0 {
			continue
		}
		body, err := client.listProjectTasks(ctx, project.Id, authorization)
		if err != nil {
			return nil, err
		}
		var taskList VikunjaTaskList
		if err := servicejson.Unmarshal(body, &taskList); err != nil {
			return nil, err
		}
		for _, task := range taskList.Items {
			source, err := storedExtIid(exts, task.Id)
			if err != nil {
				return nil, err
			}
			for _, target := range task.RelatedTasks.Blocking {
				if _, ok := projectNames[target.ProjectId]; !ok {
					continue
				}
				targetExtIid, err := storedExtIid(exts, target.Id)
				if err != nil {
					return nil, err
				}
				visitEdge := VisitEdge{Source: source, Target: targetExtIid}
				if !seen[visitEdge] {
					seen[visitEdge] = true
					visitEdges = append(visitEdges, visitEdge)
				}
			}
		}
	}
	servicesort.Slice(visitEdges, func(i, j int) bool {
		if visitEdges[i].Source == visitEdges[j].Source {
			return visitEdges[i].Target < visitEdges[j].Target
		}
		return visitEdges[i].Source < visitEdges[j].Source
	})
	return visitEdges, nil
}

func createVisitEdge(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	input VisitEdge,
	authorization string,
) (*VisitEdge, error) {
	input.Source = servicestrings.TrimSpace(input.Source)
	input.Target = servicestrings.TrimSpace(input.Target)
	if input.Source == "" || input.Target == "" {
		return nil, serviceecho.NewHTTPError(400, "Source and Target are required")
	}
	_, sourceTask, _, err := resolveVisitNodeTask(client, ctx, input.Source, authorization)
	if err != nil {
		return nil, err
	}
	_, targetTask, _, err := resolveVisitNodeTask(client, ctx, input.Target, authorization)
	if err != nil {
		return nil, err
	}
	if err := client.createTaskRelation(ctx, sourceTask.Id, targetTask.Id, authorization); err != nil {
		return nil, err
	}
	visitEdge := VisitEdge{Source: input.Source, Target: input.Target}
	return &visitEdge, nil
}

func deleteVisitEdge(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	input VisitEdge,
	authorization string,
) error {
	input.Source = servicestrings.TrimSpace(input.Source)
	input.Target = servicestrings.TrimSpace(input.Target)
	if input.Source == "" || input.Target == "" {
		return serviceecho.NewHTTPError(400, "Source and Target are required")
	}
	_, sourceTask, _, err := resolveVisitNodeTask(client, ctx, input.Source, authorization)
	if err != nil {
		return err
	}
	_, targetTask, _, err := resolveVisitNodeTask(client, ctx, input.Target, authorization)
	if err != nil {
		return err
	}
	return client.deleteTaskRelation(ctx, sourceTask.Id, targetTask.Id, authorization)
}
