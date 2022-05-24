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

	"github.com/astaxie/beego/logs"
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
	beego.BConfig.WebConfig.Session.SessionName = beego.AppConfig.String("sessionkey")
	beego.BConfig.WebConfig.Session.SessionOn = true

	/*
		logger config
	*/
	logsPath := beego.AppConfig.DefaultString("logspath", "./logs/")
	absLogsPath, _ := filepath.Abs(logsPath)
	_, err := common.CreateFileWithPerm(absLogsPath+"/", "0720")

	if err != nil && os.IsNotExist(err) {
		log.Fatalf("create file %s with error: %s", absLogsPath, err.Error())
	}

	logFileName := beego.AppConfig.DefaultString("appname", "nebula-http-gateway")
	logFileName += ".log"

	logFilePath := filepath.Join(
		absLogsPath,
		logFileName,
	)

	logLevelString := beego.AppConfig.String("logLevel")
	logLevel := logs.LevelWarning
	switch logLevelString {
	case "error":
		logLevel = logs.LevelError
	case "warning", "warn":
		logLevel = logs.LevelWarning
	case "notice":
		logLevel = logs.LevelNotice
	case "informational", "info":
		logLevel = logs.LevelInformational
	case "debug":
		logLevel = logs.LevelDebug
	}

	logs.SetLogger("file", fmt.Sprintf(`{"filename":"%s","MaxSize":104857600,"perm":"0620"}`, logFilePath))
	logs.GetBeeLogger().DelLogger("console")
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLevel(logLevel)
	defer func() {
		logs.GetBeeLogger().Flush()
	}()

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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if !signal.Ignored(syscall.SIGHUP) {
		signal.Notify(quit, syscall.SIGHUP)
	}

	go func() {
		<-quit

		logs.Info("server is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		beego.BeeApp.Server.SetKeepAlivesEnabled(false)
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			logs.Error(err.Error())
		}

		close(done)
	}()

	beego.Run()

	<-done
	logs.Info("server closed")
}
