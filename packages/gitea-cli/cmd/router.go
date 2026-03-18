package cmd

import (
	"fmt"
	"strings"

	"github.com/ucarion/urlpath"
)

// Route 定义路由结构（通用）
type Route struct {
	Method  string
	Pattern string
	Matcher urlpath.Path
}

// MatchResult 路由匹配结果（通用）
type MatchResult struct {
	Method  string
	Pattern string
	Params  map[string]string
}

// match 通用路由匹配方法
// routes: 路由列表
// method: HTTP 方法
// path: 请求路径
func match(routes []Route, method, path string) (MatchResult, error) {
	normalizedMethod := strings.ToUpper(method)
	if !strings.HasPrefix(path, "/") {
		return MatchResult{}, fmt.Errorf("path must start with /: %s", path)
	}

	for _, route := range routes {
		if route.Method != normalizedMethod {
			continue
		}

		m, ok := route.Matcher.Match(path)
		if !ok {
			continue
		}

		params := m.Params
		if params == nil {
			params = map[string]string{}
		}

		return MatchResult{
			Method:  normalizedMethod,
			Pattern: route.Pattern,
			Params:  params,
		}, nil
	}

	return MatchResult{}, fmt.Errorf("unsupported path: %s %s", normalizedMethod, path)
}
