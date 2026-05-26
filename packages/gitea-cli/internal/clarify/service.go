package clarify

import (
	"fmt"

	"com.lisitede.backstage.gitea/framework"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
)

func GetClarify(appId, planId, clarifyId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}

	issue, err := adapter.ShowIssueByPrefixId(appId, planId, clarifyId)
	if err != nil {
		return "", err
	}

	return issue.Body, nil
}

func AddClarify(appId, planId, clarifyId, context string) (interface{}, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	issue, err := adapter.ShowIssueByPrefixId(appId, planId, clarifyId)
	if err != nil {
		return nil, err
	}

	content, err := framework.LoadContext(context)
	if err != nil {
		return nil, err
	}

	issueNo := fmt.Sprintf("%d", issue.Index)
	_, err = adapter.CreateComment(appId, planId, issueNo, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return framework.RestOK, nil
}

