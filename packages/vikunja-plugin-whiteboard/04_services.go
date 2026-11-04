//go:build ignore

package main

import (
	servicecontext "context"
	servicejson "encoding/json"
	servicefmt "fmt"
	servicesort "sort"
	servicestrconv "strconv"
	servicestrings "strings"

	serviceecho "github.com/labstack/echo/v5"
)

// Yaegi loads this file after the types, client, and translator, before main.go.
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
		if project.ID < 0 {
			continue
		}
		workspaces = append(workspaces, projectToWorkspace(project))
	}
	return workspaces, nil
}

func getVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeID string,
	authorization string,
) (*VisitNode, error) {
	project, task, err := resolveVisitNodeTask(client, ctx, visitNodeID, authorization)
	if err != nil {
		return nil, err
	}
	return taskToVisitNode(*project, *task)
}

func deleteVisitNode(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeID string,
	authorization string,
) error {
	_, task, err := resolveVisitNodeTask(client, ctx, visitNodeID, authorization)
	if err != nil {
		return err
	}
	return client.deleteTask(ctx, task.ID, authorization)
}

func resolveVisitNodeTask(
	client *vikunjaHTTPClient,
	ctx servicecontext.Context,
	visitNodeID, authorization string,
) (*VikunjaProject, *VikunjaTask, error) {
	separator := servicestrings.LastIndex(visitNodeID, "-")
	if separator <= 0 {
		return nil, nil, servicefmt.Errorf("visit node ID must contain a project title and task index separated by '-'")
	}
	projectName := visitNodeID[:separator]
	index, err := servicestrconv.ParseInt(visitNodeID[separator+1:], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	project, err := findProject(client, ctx, projectName, authorization)
	if err != nil {
		return nil, nil, err
	}
	body, err := client.getTaskByIndex(ctx, project.ID, index, authorization)
	if err != nil {
		return nil, nil, err
	}
	var task VikunjaTask
	if err := servicejson.Unmarshal(body, &task); err != nil {
		return nil, nil, err
	}
	return project, &task, nil
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
		if project.ID >= 0 && servicestrings.EqualFold(project.Title, name) {
			return &project, nil
		}
	}
	return nil, serviceecho.NewHTTPError(404, servicefmt.Sprintf("project %q not found", name))
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
	description, err := formatVisitNodeDescription(input.Description, input.Session)
	if err != nil {
		return nil, serviceecho.NewHTTPError(400, err.Error())
	}
	project, err := findProject(client, ctx, input.Workspace, authorization)
	if err != nil {
		return nil, err
	}
	body, err := client.createTask(ctx, project.ID, input.Title, description, authorization)
	if err != nil {
		return nil, err
	}
	var task VikunjaTask
	if err := servicejson.Unmarshal(body, &task); err != nil {
		return nil, err
	}
	return taskToVisitNode(*project, task)
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

	visitNodes := make([]VisitNode, 0)
	for _, project := range projectList.Items {
		if project.ID < 0 {
			continue
		}

		body, err := client.listProjectTasks(ctx, project.ID, authorization)
		if err != nil {
			return nil, err
		}
		var taskList VikunjaTaskList
		if err := servicejson.Unmarshal(body, &taskList); err != nil {
			return nil, err
		}

		for _, task := range taskList.Items {
			description, err := translateTaskDescription(task.Description)
			if err != nil {
				description = map[string][]string{}
			}
			visitNodes = append(visitNodes, *visitNodeFromDescription(project, task, description))
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
		if project.ID >= 0 {
			projectNames[project.ID] = project.Title
		}
	}

	visitEdges := make([]VisitEdge, 0)
	seen := make(map[VisitEdge]bool)
	for _, project := range projectList.Items {
		if project.ID < 0 {
			continue
		}
		body, err := client.listProjectTasks(ctx, project.ID, authorization)
		if err != nil {
			return nil, err
		}
		var taskList VikunjaTaskList
		if err := servicejson.Unmarshal(body, &taskList); err != nil {
			return nil, err
		}
		for _, task := range taskList.Items {
			for _, target := range task.RelatedTasks.Blocking {
				projectName, ok := projectNames[target.ProjectID]
				if !ok {
					continue
				}
				visitEdge := visitEdgeFromIDs(projectNames[project.ID], task.Index, projectName, target.Index)
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
	sourceProject, sourceTask, err := resolveVisitNodeTask(client, ctx, input.Source, authorization)
	if err != nil {
		return nil, err
	}
	targetProject, targetTask, err := resolveVisitNodeTask(client, ctx, input.Target, authorization)
	if err != nil {
		return nil, err
	}
	if err := client.createTaskRelation(ctx, sourceTask.ID, targetTask.ID, authorization); err != nil {
		return nil, err
	}
	visitEdge := visitEdgeFromIDs(sourceProject.Title, sourceTask.Index, targetProject.Title, targetTask.Index)
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
	_, sourceTask, err := resolveVisitNodeTask(client, ctx, input.Source, authorization)
	if err != nil {
		return err
	}
	_, targetTask, err := resolveVisitNodeTask(client, ctx, input.Target, authorization)
	if err != nil {
		return err
	}
	return client.deleteTaskRelation(ctx, sourceTask.ID, targetTask.ID, authorization)
}
