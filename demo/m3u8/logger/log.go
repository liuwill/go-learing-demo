package logger

import (
	"fmt"
	"time"
)

const (
	color_red = uint8(iota + 91)
	color_success
	erro = "[ERRO]"
	info = "[INFO]"
)

func Error(format string, a ...interface{}) {
	prefix := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, erro)
	logContent := time.Now().Format("2006/01/02 15:04:05") + " " + prefix + " "
	fmt.Println(logContent, fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
	prefix := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_success, info)
	logContent := time.Now().Format("2006/01/02 15:04:05") + " " + prefix + " "
	fmt.Println(logContent, fmt.Sprintf(format, a...))
}
