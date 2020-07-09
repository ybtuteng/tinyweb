package tinyweb

import "net/http"

type Context struct {
	Writer http.ResponseWriter
	Params Params
}

func (c *Context) Param(key string) string {
	s, _ := c.Params.Get(key)
	return s
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (p Params) Get(name string) (string, bool) {
	for _, entry := range p {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}
