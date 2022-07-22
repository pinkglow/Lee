package lee

import "net/http"

type RouteGroup struct {
	basePath string        // 分组的公众前缀
	handlers HandlersChain // 提供中间件支持
	parent   *RouteGroup   // 使 Group 支持嵌套
	engine   *Engine       // 所有分组共享此 Engine 实例
}

// Group 创建一个新的分组
func (group *RouteGroup) Group(basePath string) *RouteGroup {
	// 所有的分组共享一个Engine实例
	engine := group.engine
	newGroup := &RouteGroup{
		basePath: group.basePath + basePath,
		parent:   group,
		engine:   engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouteGroup) addRoute(method, relativePath string, handler HandlerFunc) {
	pattern := group.basePath + relativePath
	group.engine.addRoute(method, pattern, handler)
}

func (group *RouteGroup) Use(middleware ...HandlerFunc) {
	group.handlers = append(group.handlers, middleware...)
}

func (group *RouteGroup) POST(relativePath string, handlerFunc HandlerFunc) {
	group.addRoute(http.MethodPost, relativePath, handlerFunc)
}

func (group *RouteGroup) GET(relativePath string, handlerFunc HandlerFunc) {
	group.addRoute(http.MethodGet, relativePath, handlerFunc)
}
