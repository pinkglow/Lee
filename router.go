package lee

import (
	"net/http"
	"strings"
)

// 定义路由的结构体
type router struct {
	// 不同方法(如 GET POST) 对应不同的树
	roots map[string]*node
	// 根据路由规则，找到对应的处理器
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// addRouter 注册路由
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	rule := getRouterRule(method, pattern)
	// 如果没有对应方法的树，则创建
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 在对应方法的树中插入路径
	r.roots[method].insertNode(pattern)
	// 在对应路由规则中，设置 handler
	r.handlers[rule] = handler
}

// getRouter 根据请求路径，找到对应的 node，并解析参数
func (r *router) getRoute(method, path string) (node *node, params map[string]string) {
	parts := parsePattern(path)
	params = make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		// 如果找不到该方法对应的树，直接返回
		return nil, nil
	}

	node = root.findNode(path)
	// 解析路径中的参数
	if node != nil {
		// 从路由规则 pattern 中解析, 例如 node.pattern = /author/:name
		paramParts := parsePattern(node.pattern)
		// 两个数组
		// paramParts 	["author", ":name"]
		// parts		["author", "lookcos"]
		for index, part := range paramParts {
			//
			if part[0] == ':' {
				// 将name=lookcos 写入字典中
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(parts) > 1 {
				// 假设这里规则为 /images/*filepath，请求路径为 /images/a.jpg 或者 /images/dir/b.jpg
				// 则 parts ["images", "a.jpg"] 或 ["images", "dir", "b.jpg"]
				// paramParts ["images", "*filepath"]
				// 则参数为 filepath = a.jpg 或 filepath = dir/b.jpg
				params[part[1:]] = strings.Join(parts[index:], "/")
			}
		}
	}
	return
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		rule := getRouterRule(c.Method, n.pattern)
		// 调用对应的 handler 去处理请求
		r.handlers[rule](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}

}

// 通过HTTP Method 和 pattern 得到路由规则
func getRouterRule(method, pattern string) string {
	var builder strings.Builder
	builder.WriteString(method)
	builder.WriteString("-")
	builder.WriteString(pattern)
	return builder.String()
}