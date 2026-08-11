package cmd

import (
	"fmt"

	"github.com/ucarion/urlpath"
)

// Router 路由管理器
type Router struct {
	routes []registeredRoute
}

// registeredRoute 已注册的路由
type registeredRoute struct {
	Method   string
	Pattern  string
	Matcher  urlpath.Path
	Callback Callback
}

// Callback 路由处理函数
type Callback func(method, pattern, pathname string, params map[string]string, args map[string]string) (interface{}, error)

// Verb 注册路由
func (r *Router) Verb(method, pattern string, callback Callback) {
	r.routes = append(r.routes, registeredRoute{
		Method:   method,
		Pattern:  pattern,
		Matcher:  urlpath.New(pattern),
		Callback: callback,
	})
}

// Invoke 调用路由
func (r *Router) Invoke(method, pathname string, args map[string]string) (interface{}, error) {
	for _, route := range r.routes {
		if route.Method != method {
			continue
		}

		match, ok := route.Matcher.Match(pathname)
		if !ok {
			continue
		}

		params := match.Params
		if params == nil {
			params = map[string]string{}
		}

		return route.Callback(method, route.Pattern, pathname, params, args)
	}

	return nil, fmt.Errorf("unsupported path: %s %s", method, pathname)
}
