package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/vesoft-inc/nebula-http-gateway/routers"

	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-http-gateway/common"
)

func main() {

	/*
		session config
	*/
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 0
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 60 * 60 * 24
	beego.BConfig.WebConfig.Session.SessionName = "nsid"
	beego.BConfig.WebConfig.Session.SessionOn = true

	/*
		logger config
	*/
	logsPath := beego.AppConfig.String("logspath")
	absLogsPath, _ := filepath.Abs(logsPath)
	_, err := common.CreateFileWithPerm(absLogsPath+"/", "0660")

	if err != nil && os.IsNotExist(err) {
		log.Fatalf("create file %s with error: %s", absLogsPath, err.Error())
	}

	logFilePath := filepath.Join(
		absLogsPath,
		"test.log",
	)
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s","MaxSize":104857600,"perm":"0660"}`, logFilePath))
	beego.BeeLogger.DelLogger("console")
	beego.SetLogFuncCall(true)
	beego.BeeLogger.SetLogFuncCallDepth(3)
	// beego.SetLevel(beego.LevelInformational)
	beego.SetLevel(beego.LevelDebug)

	/*
		importer file uploads config
	*/
	uploadsPath := beego.AppConfig.String("uploadspath")
	absUploadsPath, _ := filepath.Abs(uploadsPath)
	_, err = common.CreateFileWithPerm(absUploadsPath+"/", "0640")

	if err != nil && os.IsNotExist(err) {
		log.Fatalf("create file %s with error: %s", absLogsPath, err.Error())
	}

	beego.Run()
}
