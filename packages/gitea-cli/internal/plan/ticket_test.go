package plan

import (
	"testing"
)

// go test -v -run TestExtractTaskId ./internal/plan/
func TestExtractTaskId(t *testing.T) {
	tests := []struct {
		name       string
		issueTitle string
		wantId     string
		wantName   string
		wantNumId  string
		wantErr    bool
	}{
		{
			name:       "normal",
			issueTitle: "TASK-101: 修复权限模块 (perm/order)",
			wantId:     "TASK-101",
			wantName:   "修复权限模块 (perm/order)",
			wantNumId:  "101",
			wantErr:    false,
		},
		{
			name:       "mixed case prefix",
			issueTitle: "Task-103: 混合大小写",
			wantId:     "",
			wantName:   "",
			wantNumId:  "",
			wantErr:    true,
		},
		{
			name:       "invalid id format - letters",
			issueTitle: "TASK-ABC: 非法ID",
			wantId:     "",
			wantName:   "",
			wantNumId:  "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotName, gotNumId, err := ExtractTaskId(tt.issueTitle)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTaskId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ExtractTaskId() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotName != tt.wantName {
				t.Errorf("ExtractTaskId() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotNumId != tt.wantNumId {
				t.Errorf("ExtractTaskId() gotNumId = %v, want %v", gotNumId, tt.wantNumId)
			}
		})
	}
}

func TestExtractTicketId(t *testing.T) {
	gotId, gotName, gotNumId, err := ExtractTicketId("BLAME", "BLAME-001: 问题")
	if err != nil {
		t.Fatalf("ExtractTicketId() error = %v", err)
	}
	if gotId != "BLAME-001" || gotName != "问题" || gotNumId != "001" {
		t.Fatalf("ExtractTicketId() = (%q, %q, %q), want (BLAME-001, 问题, 001)", gotId, gotName, gotNumId)
	}
}

func TestExtractPlanId(t *testing.T) {
	gotId, gotName, gotNumId, err := ExtractPlanId("PLAN-102: My Plan")
	if err != nil {
		t.Fatalf("ExtractPlanId() error = %v", err)
	}
	if gotId != "PLAN-102" || gotName != "My Plan" || gotNumId != "102" {
		t.Fatalf("ExtractPlanId() = (%q, %q, %q), want (PLAN-102, My Plan, 102)", gotId, gotName, gotNumId)
	}
}
