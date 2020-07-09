package tinyweb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

type Request struct {
}

type Middleware interface {
	Before(Context)
	After(Context)
}

type Engine struct {
	Route
}

func NewEngine() *Engine {
	return &Engine{Route: NewTRoute()}
}

var DefaultEngine Engine = Engine{}

//AutoRoute
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Handler crashed with error %v", err)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				log.Println(file, line)
			}
		}
	}()

	node := e.Find(r.URL.Path)
	if node == nil {
		return
	}

	context := Context{}
	e.before(node, context)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ServeHTTP read body error %v", err)
		return
	}

	ft := reflect.TypeOf(node.handler)
	req := newReqInstance(ft.In(0))
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Printf("ServeHTTP unmarshal request error %v", err)
		e.serveError(w, 400, "invalid req")
		return
	}

	injector := NewInjector()
	injector.Map(req)
	ret, err := injector.Invoke(node.handler)
	if err != nil {
		log.Printf("ServeHTTP inject invoke error %v", err)
		e.serveError(w, 500, "internal error")
		return
	}

	i := ret[0].Interface()
	b, err := json.Marshal(i)
	if err != nil {
		log.Printf("ServeHTTP marshal json error %v", err)
		e.serveError(w, 500, "internal error")
		return
	}

	e.after(node, context)

	w.Write(b)
}

func (e *Engine) serveError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func (e *Engine) before(n *node, c Context) {
	for _, md := range n.middlewares {
		md.Before(c)
	}
}

func (e *Engine) after(n *node, c Context) {
	for _, md := range n.middlewares {
		md.After(c)
	}
}
