package routers

import (
	"nebula-http-gateway/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.DatabaseController{}, "*:Home")
	beego.Router("/api/connect", &controllers.DatabaseController{}, "POST:Connect")
	beego.Router("/api/exec", &controllers.DatabaseController{}, "POST:Execute")
}
