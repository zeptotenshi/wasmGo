//+build tinygo wasm,js

package web

import (
	"fmt"
	"syscall/js"
)

type Logger struct {
	Prefix string
	log    js.Value
}

func NewLogger() *Logger {
	return &Logger{log: js.Global().Get("console")}
}

func (l *Logger) Print(_msg string) {
	l.log.Call("log", fmt.Sprintf("{%s} %s", l.Prefix, _msg))
}

func (l *Logger) Error(_err error) {
	l.log.Call("log", fmt.Sprintf("{%s|ERROR} %v", l.Prefix, _err))
}

func (l *Logger) Debug(_msg string) {
	l.log.Call("log", fmt.Sprintf("{%s|DEBUG} %s", l.Prefix, _msg))
}

func (l *Logger) Info(_msg string) {
	l.log.Call("log", fmt.Sprintf("{%s|INFO} %s", l.Prefix, _msg))
}

func (l *Logger) LogElement(_e *Element) {
	l.log.Call("log", fmt.Sprintf("{%s|ELEMENT} ", l.Prefix), _e.Value)
}

func (l *Logger) LogValue(_v js.Value) {
	l.log.Call("log", _v)
}

func (l *Logger) Log(_v ...interface{}) {
	l.log.Call("log", _v)
}

func (l *Logger) Write(p []byte) (int, error) {
	l.Print(string(p))
	return len(p), nil
}
