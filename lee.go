package lee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
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
	c := newContext(w, req)
	engine.router.handle(c)
}