package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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
	_, err := common.CreateFileWithPerm(absLogsPath+"/", "0720")

	if err != nil && os.IsNotExist(err) {
		log.Fatalf("create file %s with error: %s", absLogsPath, err.Error())
	}

	logFilePath := filepath.Join(
		absLogsPath,
		"test.log",
	)
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s","MaxSize":104857600,"perm":"0620"}`, logFilePath))
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
	_, err = common.CreateFileWithPerm(absUploadsPath+"/", "0720")

	if err != nil && os.IsNotExist(err) {
		log.Fatalf("create file %s with error: %s", absLogsPath, err.Error())
	}

	/*
		use channel to wait server quit.
	*/
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-quit

		beego.Info("server is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		beego.BeeApp.Server.SetKeepAlivesEnabled(false)
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			beego.Error(err.Error())
		}

		close(done)
	}()

	beego.Run()

	<-done
	beego.Info("server closed")
}
