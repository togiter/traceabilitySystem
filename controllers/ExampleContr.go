package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type ExampleContr struct {
}

func (e *ExampleContr) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("Get", "/example/hello", "Hello")
	a.Handle("Get", "/example/goods", "Goods")
}

func (e *ExampleContr) Hello(ctx iris.Context) string {
	return ctx.Path()
}
func (e *ExampleContr) Goods() string {

	return "goods"
}

func (e *ExampleContr) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome<h1>",
	}
}

func (e *ExampleContr) GetName() string {
	return "yangminghai"
}

func (e *ExampleContr) GetMsg() interface{} {
	return map[string]string{
		"msg": "this is iris web framework！",
	}
}

func (e *ExampleContr) postLogin() interface{} {
	return "login"
}
