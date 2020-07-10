package tinyweb

import (
	"log"
	"net/http"
	"reflect"
)

type Handler interface{}

type RouteGroup interface {
	Route
	Group(string, ...Middleware) *RouteGroup
}

type Route interface {
	Use(string, ...Middleware)

	Handle(string, string, ...Middleware)
	GET(string, Handler, ...Middleware)
	POST(string, Handler, ...Middleware)
	DELETE(string, Handler, ...Middleware)
	PATCH(string, Handler, ...Middleware)
	PUT(string, Handler, ...Middleware)
	OPTIONS(string, Handler, ...Middleware)
	HEAD(string, Handler, ...Middleware)

	StaticFile(string, string)
	Static(string, string)
	StaticFS(string, http.FileSystem)
	Find(string, string) *node
}

type MiddlewareChain []*Middleware

type routeInfo struct {
	path   string
	method string
}

type tRoute struct {
	handlers        map[string]Handler
	middlewareChain map[string]MiddlewareChain
	tree            *node
}

type Response struct {
}

func (r *Response) print() {
	log.Println("i run")
}

func NewTRoute() *tRoute {
	return &tRoute{
		handlers:        make(map[string]Handler),
		middlewareChain: make(map[string]MiddlewareChain),
		tree:            new(node),
	}
}

func (t tRoute) Use(path string, middleware ...Middleware) {
	_, ok := t.middlewareChain[path]
	if !ok {
		t.middlewareChain[path] = []*Middleware{}
	}
	for _, md := range middleware {
		t.middlewareChain[path] = append(t.middlewareChain[path], &md)
	}
}

func (t tRoute) Handle(s2 string, s22 string, middleware ...Middleware) {
	panic("implement me")
}

func (t tRoute) GET(path string, handler Handler, middleware ...Middleware) {
	t.tree.addRoute(path, "GET", handler, middleware...)
}

func (t tRoute) POST(path string, handler Handler, middleware ...Middleware) {
	t.tree.addRoute(path, "POST", handler, middleware...)
}

func (t tRoute) DELETE(path string, handler Handler, middleware ...Middleware) {
	t.tree.addRoute(path, "DELETE", handler, middleware...)
}

func (t tRoute) PATCH(path string, handler Handler, middleware ...Middleware) {
	t.tree.addRoute(path, "PATCH", handler, middleware...)
}

func (t tRoute) PUT(path string, handler Handler, middleware ...Middleware) {
	t.tree.addRoute(path, "PUT", handler, middleware...)
}

func (t tRoute) OPTIONS(s2 string, handler Handler, middleware ...Middleware) {
	panic("implement me")
}

func (t tRoute) HEAD(s2 string, handler Handler, middleware ...Middleware) {
	panic("implement me")
}

func (t tRoute) StaticFile(s2 string, s22 string) {
	panic("implement me")
}

func (t tRoute) Static(s2 string, s22 string) {
	panic("implement me")
}

func (t tRoute) StaticFS(s2 string, system http.FileSystem) {
	panic("implement me")
}

func (t tRoute) Find(path, method string) *node {
	if node, err := t.tree.find(path, method); err == nil {
		return node
	}
	return nil
}

func newReqInstance(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Ptr:
		return newReqInstance(t.Elem())
	case reflect.Interface:
		return nil
	default:
		return reflect.New(t).Interface()
	}
}
