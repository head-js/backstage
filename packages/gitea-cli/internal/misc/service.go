package misc

import (
	"fmt"

	"code.gitea.io/sdk/gitea"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
)

func CreateBlame(appId, planId, title, context string) (map[string]interface{}, error) {
	const aId = "backstage"
	const pId = "blames"

	blameId, err := plan.GenNextBlameId(aId, pId)
	if err != nil {
		return nil, err
	}

	blameTitle := fmt.Sprintf("%s: %s", blameId, title)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	appLabel, err := adapter.GetLabelByName(aId, pId, "app/"+appId)
	if err != nil {
		appLabel, err = adapter.GetLabelByName(aId, pId, "app/unknown")
		if err != nil {
			return nil, err
		}
	}

	planLabel, err := adapter.GetLabelByName(aId, pId, "plan/"+planId)
	if err != nil {
		planLabel, err = adapter.GetLabelByName(aId, pId, "plan/unknown")
		if err != nil {
			return nil, err
		}
	}

	issue, err := adapter.CreateIssue(aId, pId, blameTitle, context, "")
	if err != nil {
		return nil, err
	}

	issueNo := fmt.Sprintf("%d", issue.Index)
	if _, err := adapter.AddLabelToIssue(aId, pId, issueNo, fmt.Sprintf("%d", appLabel.ID)); err != nil {
		return nil, err
	}
	if _, err := adapter.AddLabelToIssue(aId, pId, issueNo, fmt.Sprintf("%d", planLabel.ID)); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":      blameId,
		"title":   issue.Title,
		"context": context,
	}, nil
}

func ListBlame(appId, planId string) ([]map[string]interface{}, error) {
	const aId = "backstage"
	const pId = "blames"

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	appLabel, err := adapter.GetLabelByName(aId, pId, "app/"+appId)
	if err != nil {
		appLabel, err = adapter.GetLabelByName(aId, pId, "app/unknown")
		if err != nil {
			return nil, err
		}
	}

	planLabel, err := adapter.GetLabelByName(aId, pId, "plan/"+planId)
	if err != nil {
		planLabel, err = adapter.GetLabelByName(aId, pId, "plan/unknown")
		if err != nil {
			return nil, err
		}
	}

	issues, err := adapter.SearchIssueOfRepo(aId, pId, gitea.ListIssueOption{
		Labels: []string{appLabel.Name, planLabel.Name},
	})
	if err != nil {
		return nil, err
	}

	var blames []map[string]interface{}
	for _, issue := range issues {
		blameId, title, _, err := plan.ExtractBlameId(issue.Title)
		if err != nil {
			continue
		}
		blames = append(blames, map[string]interface{}{
			"id":      blameId,
			"title":   title,
			"context": issue.Body,
		})
	}

	return blames, nil
}
