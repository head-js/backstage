package linear

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
)

// Adapter Linear API 适配器
type Adapter struct {
	client graphql.Client
}

// authTransport 在每个请求上注入 Authorization 头
type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	return t.base.RoundTrip(req)
}

// NewAdapter 创建新的 Linear 适配器
func NewAdapter() (*Adapter, error) {
	url := os.Getenv("BACKSTAGE_LINEAR_URL")
	token := os.Getenv("BACKSTAGE_LINEAR_TOKEN")

	if url == "" {
		return nil, fmt.Errorf("BACKSTAGE_LINEAR_URL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("BACKSTAGE_LINEAR_TOKEN is required")
	}

	httpClient := &http.Client{
		Transport: &authTransport{token: token, base: http.DefaultTransport},
	}
	client := graphql.NewClient(url, httpClient)

	return &Adapter{client: client}, nil
}

// GetViewer 获取当前 Token 对应的用户（用于探活）
func (a *Adapter) GetViewer() (*GetViewerViewerUser, error) {
	resp, err := GetViewer(context.Background(), a.client)
	if err != nil {
		return nil, err
	}
	return resp.Viewer, nil
}

// ListTeam 列出当前 Token 可见的所有 Team
func (a *Adapter) ListTeam() ([]*ListTeamTeamsTeamConnectionNodesTeam, error) {
	resp, err := ListTeam(context.Background(), a.client)
	if err != nil {
		return nil, err
	}
	return resp.Teams.Nodes, nil
}

// ListProjectOfTeam 列出指定 Team 下的所有 Project
// teamKey: Team 的 key（如 "GLNT"），Linear API 同时接受 UUID 和 key
func (a *Adapter) ListProjectOfTeam(teamKey string) ([]*ListProjectOfTeamTeamProjectsProjectConnectionNodesProject, error) {
	resp, err := ListProjectOfTeam(context.Background(), a.client, teamKey)
	if err != nil {
		return nil, err
	}
	return resp.Team.Projects.Nodes, nil
}

// ListIssueOfProject 列出指定 Team + Project 下的所有 Issue
// teamKey: Team 的 key（如 "GLNT"）
// projectName: Project 的 name（如 "Alert"）
func (a *Adapter) ListIssueOfProject(teamKey, projectName string) ([]*ListIssueOfProjectIssuesIssueConnectionNodesIssue, error) {
	resp, err := ListIssueOfProject(context.Background(), a.client, teamKey, projectName)
	if err != nil {
		return nil, err
	}
	return resp.Issues.Nodes, nil
}

// GetIssue 获取单个 Issue 详情
// issueId: Issue 的 identifier（如 "GLNT-2"），Linear API 同时接受 UUID 和 identifier
func (a *Adapter) GetIssue(issueId string) (*GetIssueIssue, error) {
	resp, err := GetIssue(context.Background(), a.client, issueId)
	if err != nil {
		return nil, err
	}
	return resp.Issue, nil
}
