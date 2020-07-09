package main

import (
	"awesomeProject/tinyweb"
	"log"
	"net/http"
)

func main() {
	r := tinyweb.NewEngine()
	r.POST("/index", SayHello, &Print{})
	log.Fatal(http.ListenAndServe(":8080", r))
}

type HelloRsp struct {
	Content string
}

type HelloReq struct {
	tinyweb.Request
	Content string
}

func SayHello(h *HelloReq) *HelloRsp {
	log.Println(h.Content)
	return &HelloRsp{
		Content: "HelloRsp",
	}
}

type Print struct {
}

func (l *Print) Before(c tinyweb.Context) {
	log.Println("before", c.Param("before"))
}

func (l *Print) After(c tinyweb.Context) {
	log.Println("after", c.Param("after"))
}
