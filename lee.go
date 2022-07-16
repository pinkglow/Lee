package lee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)
type HandlersChain []HandlerFunc

type Engine struct {
	router *router
	groups []*RouteGroup // Engine 控制着所有分组
	// Engine 要有新建分组的能力, 比如
	//engine := lee.New()
	// engine.Group("/v1")
	*RouteGroup
}

// New 方法初始化了Engine，同时Engine也作为一个顶层的Group
// 应该报着设计方法来初始化
func New() *Engine {
	engine := &Engine{router: newRouter()}
	// engine 作为顶层分组，因此它的 prefix, middlewares, parent
	// 等属性都是空的，但所有分组共享一个 Engine 实例
	engine.RouteGroup = &RouteGroup{engine: engine, basePath: ""}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

// 添加路由规则
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodGet, pattern, handler)
}

// POST 方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPost, pattern, handler)
}

// PUT 方法
func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPut, pattern, handler)
}

// PATCH 方法
func (engine *Engine) PATCH(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPatch, pattern, handler)
}

// DELETE 方法
func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodDelete, pattern, handler)
}

// Run 启动 HTTP 服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP
// Engine 实现了 http.HandlerFunc 的 ServeHTTP 的方法
// 因此，Engine 是一个 HandlerFunc 类型的实例
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handlers HandlersChain
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.basePath) {
			handlers = append(handlers, group.handlers...)
		}
	}
	c := newContext(w, req)
	c.handlers = handlers
	engine.router.handle(c)
}
