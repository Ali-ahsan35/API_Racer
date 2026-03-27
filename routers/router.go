package routers

import (
	"apiracer/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/benchmark", &controllers.BenchmarkController{}, "get:RunBenchmark")
}
