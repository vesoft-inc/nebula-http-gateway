package logger

import (
	"log"
)

/*
	@dev:
	v1: no logger
	v2.0.0 -- v2.6.0:
	type Logger interface {
		Info(msg string)
		Warn(msg string)
		Error(msg string)
		Fatal(msg string)
	}
*/

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type DefaultLogger struct {
}

func (l DefaultLogger) Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func (l DefaultLogger) Warn(msg string) {
	log.Printf("[WARNING] %s\n", msg)
}

func (l DefaultLogger) Error(msg string) {
	log.Printf("[ERROR] %s\n", msg)
}

func (l DefaultLogger) Fatal(msg string) {
	log.Fatalf("[FATAL] %s\n", msg)
}
