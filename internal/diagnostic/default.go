package diagnostic

var (
	_default = NewDiagnostic(NewOptions(
		WithName("Default"),
	))
)

func Trace(message string, fields ...Field) {
	_default.Trace(message, fields...)
}

func Debug(message string, fields ...Field) {
	_default.Debug(message, fields...)
}

func Info(message string, fields ...Field) {
	_default.Info(message, fields...)
}

func Warn(message string, fields ...Field) {
	_default.Warn(message, fields...)
}

func Error(message string, fields ...Field) {
	_default.Error(message, fields...)
}

func Fatal(message string, fields ...Field) {
	_default.Fatal(message, fields...)
}
