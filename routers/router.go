package routers

import (
	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-http-gateway/controllers"
)

func init() {
	beego.Router("/", &controllers.DatabaseController{}, "*:Home")
	beego.Router("/api/db/connect", &controllers.DatabaseController{}, "POST:Connect")
	beego.Router("/api/db/exec", &controllers.DatabaseController{}, "POST:Execute")
	beego.Router("/api/db/disconnect", &controllers.DatabaseController{}, "POST:Disconnect")

	beego.Router("/api/task/import", &controllers.TaskController{}, "POST:Import")
	beego.Router("/api/task/action", &controllers.TaskController{}, "POST:Action")
}
