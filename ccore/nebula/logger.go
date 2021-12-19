package nebula

var (
	_ Logger = noOpLogger{}
)

type (
	Logger interface {
		Info(...interface{})
		Infof(string, ...interface{})
		Warn(...interface{})
		Warnf(string, ...interface{})
		Error(...interface{})
		Errorf(string, ...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
	}
)

// noOpLogger is used as a placeholder for the default logger
type noOpLogger struct{}

func (noOpLogger) Info(...interface{})           {}
func (noOpLogger) Infof(string, ...interface{})  {}
func (noOpLogger) Warn(...interface{})           {}
func (noOpLogger) Warnf(string, ...interface{})  {}
func (noOpLogger) Error(...interface{})          {}
func (noOpLogger) Errorf(string, ...interface{}) {}
func (noOpLogger) Fatal(...interface{})          {}
func (noOpLogger) Fatalf(string, ...interface{}) {}
