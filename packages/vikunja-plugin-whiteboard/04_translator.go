//go:build ignore

package main

import (
	translatorerrors "errors"
	translatorfmt "fmt"
	translatorsort "sort"
	translatorstrconv "strconv"
	translatorstrings "strings"
)

const taskDescriptionSeparator = "* * *"
const workspacePriorityBase = 200 * 65536

func projectToWorkspace(project VikunjaProject) Workspace {
	return Workspace{
		Id:          project.Id,
		Name:        project.Title,
		Description: project.Description,
		Priority:    workspacePriorityBase - project.Position,
	}
}

// visitNodeId is the assignment formula for new VisitNode Id. It is used only
// when persisting tasks.ext_iid; reads always use the stored value.
func visitNodeId(projectTitle string, taskIndex int64) string {
	return translatorstrings.ToUpper(projectTitle) + "-" + translatorstrconv.FormatInt(taskIndex, 10)
}

func taskToVisitNode(project VikunjaProject, task VikunjaTask, extIid, extSession string) (*VisitNode, error) {
	description, err := translateTaskDescription(task.Description)
	if err != nil {
		return nil, err
	}
	return visitNodeFromDescription(project, task, description, extIid, extSession), nil
}

// Session is sourced from tasks.ext_session only. The description passes
// through untouched; any "Session" section inside it is ordinary data.
func visitNodeFromDescription(project VikunjaProject, task VikunjaTask, description map[string][]string, extIid, extSession string) *VisitNode {
	return &VisitNode{
		Id:          extIid,
		Workspace:   project.Title,
		Session:     extSession,
		Title:       task.Title,
		Progress:    translatorfmt.Sprintf("%.0f%%", task.PercentDone*100),
		Description: description,
	}
}

// translateTaskDescription converts the whiteboard task description format into
// a map of section titles to values. Each section is separated by the
// "* * *" token and starts with an H2 heading.
func translateTaskDescription(description string) (map[string][]string, error) {
	sections := make(map[string][]string)
	blocks := translatorstrings.Split(description, taskDescriptionSeparator)

	for _, rawBlock := range blocks {
		block := translatorstrings.TrimSpace(rawBlock)
		if block == "" {
			continue
		}

		parts := translatorstrings.SplitN(block, "\n", 2)
		heading := translatorstrings.TrimSpace(parts[0])
		if !translatorstrings.HasPrefix(heading, "## ") {
			return nil, translatorerrors.New("task description section must start with an H2 heading")
		}

		key := translatorstrings.TrimSpace(translatorstrings.TrimPrefix(heading, "## "))
		if key == "" {
			return nil, translatorerrors.New("task description section heading must not be empty")
		}

		value := ""
		if len(parts) == 2 {
			value = translatorstrings.TrimSpace(parts[1])
		}
		sections[key] = append(sections[key], value)
	}

	return sections, nil
}

// formatTaskDescription writes one H2 block per value, keeping array order.
func formatTaskDescription(description map[string][]string) (string, error) {
	keys := make([]string, 0, len(description))
	for key := range description {
		keys = append(keys, key)
	}
	translatorsort.Strings(keys)

	blocks := make([]string, 0)
	for _, key := range keys {
		heading := translatorstrings.TrimSpace(key)
		if heading == "" || translatorstrings.ContainsAny(heading, "\r\n") {
			return "", translatorerrors.New("description key must be a non-empty single-line heading")
		}
		for _, value := range description[key] {
			blocks = append(blocks, "## "+heading+"\n\n"+translatorstrings.TrimSpace(value))
		}
	}
	return translatorstrings.Join(blocks, "\n\n"+taskDescriptionSeparator+"\n\n"), nil
}

// formatVisitNodeDescription serializes the description sections as given.
// Session is not part of the description; it lives in tasks.ext_session.
func formatVisitNodeDescription(description map[string][]string) (string, error) {
	return formatTaskDescription(description)
}
