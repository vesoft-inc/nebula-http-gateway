package common

import (
	"net/http"
	"runtime"

	"github.com/astaxie/beego/logs"
)

func LogPanic(r interface{}) {
	if r == http.ErrAbortHandler {
		return
	}

	// Same as stdlib http server code. Manually allocate stack trace buffer size
	// to prevent excessively large logs
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	if _, ok := r.(string); ok {
		logs.Warn("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		logs.Warn("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
