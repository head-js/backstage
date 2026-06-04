package plan

import (
	"fmt"
	"regexp"
	"strings"

	"com.lisitede.backstage.gitea/framework"
)

// ExtractTicketId 解析带有指定前缀的 ticket 标题，返回完整 ticket ID、ticket 名称、数字 ID 部分和错误。
//
// 期望格式："<PREFIX>-<numId>: <name>"
// 例如：ExtractTicketId("TASK", "TASK-101: 修复 bug") -> ("TASK-101", "修复 bug", "101", nil)
func ExtractTicketId(prefix, ticketTitle string) (string, string, string, error) {
	pattern := fmt.Sprintf(`^%s-(\d{3}):\s*(.+)$`, regexp.QuoteMeta(prefix))
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(ticketTitle)

	if len(matches) < 3 {
		return "", "", "", framework.InvalidFormatException(fmt.Sprintf("ticket title must follow format '%s-{id}: {name}': %s", prefix, ticketTitle))
	}

	ticketId := prefix + "-" + matches[1]
	ticketName := matches[2]
	numId := matches[1]
	return ticketId, ticketName, numId, nil
}

func ExtractPlanId(planTitle string) (string, string, string, error) {
	if strings.TrimSpace(planTitle) == "" {
		return "", "", "", framework.InvalidFormatException("title is empty, consider checking repo.description")
	}
	return ExtractTicketId("PLAN", planTitle)
}

func ExtractPhaseId(milestoneTitle string) (string, string, string, error) {
	return ExtractTicketId("PHASE", milestoneTitle)
}

func ExtractTaskId(issueTitle string) (string, string, string, error) {
	return ExtractTicketId("TASK", issueTitle)
}

func ExtractBlameId(issueTitle string) (string, string, string, error) {
	return ExtractTicketId("BLAME", issueTitle)
}
