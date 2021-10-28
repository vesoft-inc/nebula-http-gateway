package logger

import (
	"fmt"

	"github.com/astaxie/beego"
)

type HttpGatewayLogger struct{}

func (l HttpGatewayLogger) Info(msg string) {
	beego.Info(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Warn(msg string) {
	beego.Warn(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Error(msg string) {
	beego.Error(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Fatal(msg string) {
	beego.Emergency(fmt.Sprintf("[nebula-clients] %s", msg))
}
