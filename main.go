package main

import (
	"net"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/traceability-system/controllers"
)

/*golang iris mvc框架的服务端加载过程
整个iris框架共三层结构：
https://www.jianshu.com/p/e8ae82234391
应用的配置和注册信息，如路由、中间件、日志。
中间的服务端实例，从iris实例拿配置信息进行配置。
底层net/http包，负责TCP连接建立、监听接受，请求收取并解析，缓冲区管理，写入响应。
*/
//在写入 context.ResponseWriter() 之后可能无
// 法从 Context.Request().Body 读取内容。
// 严谨的处理程序应该首先读取 Context.Request().Body ，然后再响应。
func main() {

	//所有的配置项都是有默认值的，所有配置都会在当使用  iris.New() 时发挥功效。
	app := iris.New()
	//配置,toml格式用iris.TOML
	app.Configure(iris.WithConfiguration(iris.YAML("./configs/iris.yml")))

	app.Use(recover.New())
	app.Use(logger.New())
	//静态
	app.StaticWeb("/statics", "./static/html")
	//注册视图
	app.RegisterView(iris.HTML("./static/templates", ".html").Reload(true))

	aMvc := mvc.New(app)
	aMvc.Handle(new(controllers.ExampleContr))
	// aMvc.Handle(new(controllers.FabricContr))

	//使用分组路由配置控制器,mvc.Configure来配置路由组和控制器的设置
	mvc.Configure(app.Party("/fabric"), func(mvc *mvc.Application) {
		mvc.Handle(new(controllers.FabricContr))
	})

	// 使用自定义 net.Listener
	listen, err := net.Listen("tcp4", ":8088")
	if err != nil {
		panic(err)
	}
	// 在 Tcp 上监听网络地址 0.0.0.0:8080
	app.Run(iris.Listener(listen))
	// app.Run(iris.Addr(":8088"))

}
