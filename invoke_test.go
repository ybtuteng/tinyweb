package tinyweb

import (
    "log"
    "net/http"
    "testing"
)

func TestInvoke(t *testing.T) {
    injector := NewInjector()
    a := A{name: "doudou"}
    injector.Map(a)
    f := test
    injector.Invoke(f)
}

type A struct {
    name string
}

func test(a A) {
    log.Printf("name: %s", a.name)
}

func TestWeb(t *testing.T) {
    r := NewTRoute()
    r.GET("/index", SayHello)
    mux := http.NewServeMux()
    mux.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func SayHello() {
    log.Println("hello world!")
}
