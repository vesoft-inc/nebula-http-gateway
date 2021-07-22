package common

import (
	"log"
	"net/http"
	"runtime"
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
		log.Printf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		log.Printf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
