package logger

import (
	"fmt"

	"github.com/astaxie/beego"
)

func Debug(v ...interface{}) {
	debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	info(fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	warn(fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	warn(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	error(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	fatal(fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	fatal(fmt.Sprintf(format, v...))
}

func debug(msg string) {
	beego.Debug(msg)
}

func info(msg string) {
	beego.Info(msg)
}

func warn(msg string) {
	beego.Warn(msg)
}

func error(msg string) {
	beego.Error(msg)
}

func fatal(msg string) {
	beego.Emergency(msg)
}
