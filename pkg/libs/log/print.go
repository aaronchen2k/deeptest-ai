package _logs

import (
	"fmt"
	"runtime/debug"
)

func Info(str string) {
	Zap.Info(str)
}
func Infof(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Zap.Info(msg)
}

func Warn(str string) {
	Zap.Warn(str)
}
func Warnf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Zap.Warn(msg)
}

func Error(str string) {
	Zap.Error(str)
	s := string(debug.Stack())
	fmt.Printf("err=%v, stack=%s\n", str, s)
}
func Errorf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Zap.Error(msg)
	s := string(debug.Stack())
	fmt.Printf("err=%v, stack=%s\n", msg, s)
}
func Debug(str string) {
	Zap.Debug(str)
}
func Debugf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	Zap.Debug(msg)
}
