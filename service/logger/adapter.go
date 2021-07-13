package logger

import "fmt"

type HttpGatewayLogger struct{}

func (l HttpGatewayLogger) Info(msg string) {
	info(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Warn(msg string) {
	warn(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Error(msg string) {
	error(fmt.Sprintf("[nebula-clients] %s", msg))
}

func (l HttpGatewayLogger) Fatal(msg string) {
	fatal(fmt.Sprintf("[nebula-clients] %s", msg))
}
