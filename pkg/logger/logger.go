package logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warning(msg string)
	Fatal(msg string)
	Error(msg string)
	SetOutputFile(fileName string) error
}
