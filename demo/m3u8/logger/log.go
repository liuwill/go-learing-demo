package logger

import (
	"fmt"
	"time"
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	log_erro    = "[ERRO]"
	log_success = "[SUCCESS]"
	log_warning = "[WARNING]"
	log_info    = "[INFO]"
)

type Logger struct{}

func (logger *Logger) print(color uint8, tip string, format string, a ...interface{}) {
	prefix := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, tip)
	logContent := time.Now().Format("2006/01/02 15:04:05") + " " + prefix + " "
	fmt.Println(logContent, fmt.Sprintf(format, a...))
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.print(color_red, log_erro, format, a...)
}

func (logger *Logger) Success(format string, a ...interface{}) {
	logger.print(color_green, log_success, format, a...)
}

func (logger *Logger) Warning(format string, a ...interface{}) {
	logger.print(color_yellow, log_warning, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.print(color_blue, log_info, format, a...)
}
