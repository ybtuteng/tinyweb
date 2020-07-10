package tinyweb

import (
	"errors"
	"reflect"
	"strings"
)

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type node struct {
	path        string
	method      string
	children    []*node
	nType       nodeType
	handler     Handler
	middlewares []Middleware
}

func (n *node) addRoute(path, method string, handler Handler, middlewares ...Middleware) {
	t := reflect.TypeOf(handler)
	if t.NumIn() != 2 {
		panic("invalid method without 2 argument")
	}

	if len(n.path) == 0 && len(n.children) == 0 {
		n.path = path
		n.method = strings.ToUpper(method)
		n.nType = root
		n.handler = handler
		n.middlewares = middlewares
		return
	}

	i := longestCommonPrefix(path, n.path)

	if i < len(n.path) {
		child := node{
			path:        n.path[i:],
			method:      n.method,
			children:    n.children,
			handler:     n.handler,
			middlewares: n.middlewares,
		}

		n.path = path[:i]
		n.children = []*node{&child}
		n.handler = nil
	}

	if i < len(path) {
		child := &node{
			path:        path[i:],
			method:      strings.ToUpper(method),
			handler:     handler,
			middlewares: middlewares,
		}
		n.children = append(n.children, child)
	}
	return
}

func (n *node) find(path, method string) (*node, error) {
	if path == n.path && strings.ToUpper(method) == n.method {
		return n, nil
	}
	if match(n.path, path) {
		for _, child := range n.children {
			return child.find(path[len(n.path):], method)
		}
	}
	return nil, errors.New("not find")
}

func match(prefix, fullPath string) bool {
	return strings.Contains(fullPath, prefix)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}
