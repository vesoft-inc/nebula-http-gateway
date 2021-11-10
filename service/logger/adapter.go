package logger

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

type HttpGatewayLogger struct{}

func (l HttpGatewayLogger) Info(msg string) {
	logs.Info(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Warn(msg string) {
	logs.Warn(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Error(msg string) {
	logs.Error(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Fatal(msg string) {
	logs.Emergency(fmt.Sprintf("[nebula-clients] %s", msg))
}
