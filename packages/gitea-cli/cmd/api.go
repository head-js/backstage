package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"com.lisitede.backstage.gitea/internal/gitea"
	"github.com/spf13/cobra"
	"github.com/ucarion/urlpath"
)

type apiRoute struct {
	method  string
	pattern string
	matcher urlpath.Path
}

type apiMatch struct {
	method  string
	pattern string
	params  map[string]string
}

var apiRoutes = []apiRoute{
	{method: "GET", pattern: "/repos", matcher: urlpath.New("/repos")},
	{method: "POST", pattern: "/repos", matcher: urlpath.New("/repos")},
	{method: "GET", pattern: "/repos/:repoId", matcher: urlpath.New("/repos/:repoId")},
	{method: "GET", pattern: "/version", matcher: urlpath.New("/version")},
	{method: "GET", pattern: "/users/:username/repos", matcher: urlpath.New("/users/:username/repos")},
}

// apiCmd 通用 API 命令
var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Call Gitea API",
	Long: `Call Gitea API by RESTful path.
Examples:
	backstage-gitea api GET /repos
	backstage-gitea api GET /repos/123
  backstage-gitea api POST /repos`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, err := gitea.NewAdapter()
		if err != nil {
			return outputError(err)
		}

		method := strings.ToUpper(args[0])
		path := args[1]

		// 路由分发
		result, err := routeAPI(adapter, path, method)
		if err != nil {
			return outputError(err)
		}

		// 输出结果
		if jsonFlag {
			return json.NewEncoder(os.Stdout).Encode(result)
		}

		// 简单 JSON 输出
		printResult(result)
		return nil
	},
}

// routeAPI 先进行纯路径匹配，再分发到 adapter 方法
func routeAPI(adapter *gitea.Adapter, path, method string) (interface{}, error) {
	match, err := matchAPIPath(method, path)
	if err != nil {
		return nil, err
	}

	return executeHandler(adapter, path, match)
}

func matchAPIPath(method, path string) (apiMatch, error) {
	normalizedMethod := strings.ToUpper(method)
	if !strings.HasPrefix(path, "/") {
		return apiMatch{}, fmt.Errorf("path must start with /: %s", path)
	}

	for _, route := range apiRoutes {
		if route.method != normalizedMethod {
			continue
		}

		match, ok := route.matcher.Match(path)
		if !ok {
			continue
		}

		params := match.Params
		if params == nil {
			params = map[string]string{}
		}

		return apiMatch{
			method:  normalizedMethod,
			pattern: route.pattern,
			params:  params,
		}, nil
	}

	return apiMatch{}, fmt.Errorf("unsupported API: %s %s", normalizedMethod, path)
}

// executeHandler 根据匹配到的路由模板执行对应的 adapter 方法
func executeHandler(adapter *gitea.Adapter, path string, match apiMatch) (interface{}, error) {
	switch match.pattern {
	case "/repos":
		if match.method == "GET" {
			return adapter.ListRepos()
		}
		if match.method == "POST" {
			return nil, fmt.Errorf("POST /repos requires --name and --description flags")
		}
	case "/repos/:repoId":
		return map[string]string{"repoId": match.params["repoId"]}, nil
	case "/version":
		if match.method == "GET" {
			return adapter.GetGiteaVersion()
		}
	case "/users/:username/repos":
		if match.method == "GET" {
			return adapter.ListUserRepos(match.params["username"])
		}
	}
	return nil, fmt.Errorf("unsupported API: %s %s", match.method, path)
}

// printResult 简单 JSON 输出
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
