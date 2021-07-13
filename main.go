package main

import (
	"fmt"

	_ "nebula-http-gateway/routers"
	common "nebula-http-gateway/utils"

	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 0
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 60 * 60 * 24
	beego.BConfig.WebConfig.Session.SessionName = "nsid"
	beego.BConfig.WebConfig.Session.SessionOn = true

	// logger config
	// logFilepath, _ := filepath.Abs("logs/test.log")
	logFilepath := "logs/test.log"
	permcode := "0660"
	common.CreateFileWithPerm(logFilepath, permcode)
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s","MaxSize":104857600,"perm":"%s"}`, logFilepath, permcode))
	// beego.BeeLogger.DelLogger("console")
	// beego.SetLevel(beego.LevelInformational)
	beego.SetLevel(beego.LevelDebug)
	beego.SetLogFuncCall(true)
	beego.BeeLogger.SetLogFuncCallDepth(3)

	beego.Run()
}
