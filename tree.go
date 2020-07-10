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
	wildcard
)

type nodes []*node

type node struct {
	path        string
	method      string
	children    nodes
	nType       nodeType
	handler     Handler
	middlewares []Middleware
}

func (ns nodes) find(path, method string) (*node, error) {
	for _, child := range ns {
		n, err := child.find(path, method)
		if err == nil {
			return n, err
		}
	}
	return nil, errors.New("node not found")
}

func (n *node) addRoute(path, method string, handler Handler, middlewares ...Middleware) {
	t := reflect.TypeOf(handler)
	if t.NumIn() != 2 {
		panic("handler method need exact two argument of context&request")
	}

	if len(n.path) == 0 && len(n.children) == 0 {
		n.nType = root
		n.insertChild(path, method, handler, middlewares...)
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
		child := &node{}
		n.children = append(n.children, child)
		child.insertChild(path[i:], method, handler, middlewares...)
	}

	return
}

func (n *node) find(path, method string) (*node, error) {
	if path == n.path && strings.ToUpper(method) == n.method {
		return n, nil
	}

	if strings.Contains(path, n.path) {
		n, err := n.children.find(path[len(n.path):], method)
		if err == nil {
			return n, err
		}
	}

	if len(n.path) != 0 && (n.path[0] == ':' || n.path[0] == '*') {
		i := strings.Index(path, "/")
		if i == -1 {
			return n, nil
		}

		n, err := n.children.find(path[i:], method)
		if err == nil {
			return n, err
		}
	}

	return nil, errors.New("not find")
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

func (n *node) insertChild(path, method string, handler Handler, middlewares ...Middleware) {
	for {
		// Find prefix until first wildcardSubStr
		wildcardSubStr, i, valid := findWildcard(path)
		if i < 0 { // No wildcardSubStr found
			break
		}

		// The wildcardSubStr name must not contain ':' and '*'
		if !valid {
			panic("only one wildcardSubStr per path segment is allowed, has: '" +
				wildcardSubStr + "' in path '" + path + "'")
		}

		// check if the wildcardSubStr has a name
		if len(wildcardSubStr) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + path + "'")
		}

		// Check if this node has existing children which would be
		// unreachable if we insert the wildcardSubStr here
		if len(n.children) > 0 {
			panic("wildcardSubStr segment '" + wildcardSubStr +
				"' conflicts with existing children in path '" + path + "'")
		}

		if wildcardSubStr[0] == ':' { // param
			if i > 0 {
				// Insert prefix before the current wildcardSubStr
				n.path = path[:i]
				path = path[i:]
			}

			child := &node{
				nType:  wildcard,
				path:   wildcardSubStr,
				method: method,
			}
			n.children = []*node{child}
			n = child

			// if the path doesn't end with the wildcardSubStr, then there
			// will be another non-wildcardSubStr subpath starting with '/'
			if len(wildcardSubStr) < len(path) {
				path = path[len(wildcardSubStr):]

				child := &node{}
				n.children = []*node{child}
				n = child
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			n.handler = handler
			n.middlewares = middlewares
			return
		}

		// catchAll
		if i+len(wildcardSubStr) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + path + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("catch-all conflicts with existing handle for the path segment root in path '" + path + "'")
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + path + "'")
		}

		n.path = path[:i]

		// First node: catchAll node with empty path
		child := &node{
			nType: catchAll,
		}

		n.children = []*node{child}
		n = child

		// second node: node holding the variable
		child = &node{
			path:        path[i:],
			nType:       catchAll,
			handler:     handler,
			middlewares: middlewares,
		}
		n.children = []*node{child}

		return
	}

	// If no wildcard was found, simply insert the path and handle
	n.path = path
	n.method = strings.ToUpper(method)
	n.handler = handler
	n.middlewares = middlewares
}

func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' && c != '*' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}
