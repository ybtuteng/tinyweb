package main

import (
	"awesomeProject/tinyweb"
	"log"
	"net/http"
)

func main() {
	r := tinyweb.NewEngine()
	r.POST("/index", SayHello, &Print{})
	r.POST("/yes/:id", SayYes)
	log.Fatal(http.ListenAndServe(":8080", r))
}

type HelloRsp struct {
	Content string
}

type HelloReq struct {
	tinyweb.Request
	Content string
}

func SayHello(c tinyweb.Context, h *HelloReq) *HelloRsp {
	log.Println(h.Content)
	return &HelloRsp{
		Content: "HelloRsp",
	}
}

func SayYes(c tinyweb.Context, h *HelloReq) *HelloRsp {
	log.Println(h.Content)
	return &HelloRsp{
		Content: "yes",
	}
}

type Print struct {
}

func (l *Print) Before(c tinyweb.Context) {
	log.Println("before", c.Request.Header.Get("before"))
}

func (l *Print) After(c tinyweb.Context) {
	log.Println("after", c.Request.Header.Get("after"))
}
