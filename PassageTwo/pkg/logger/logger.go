package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota // 0
	LevelInfo               // 1
	LevelWarn               // 2
	LevelError              // 3
	LevelFatal              // 4
	LevelPanic              // 5
)

type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

// 将日志级别转化为字符串
func (l Level) LevelToString() string {
	switch l {
	case LevelDebug:
		return "LevelDebug"
	case LevelInfo:
		return "LevelInfo"
	case LevelWarn:
		return "LevelWarn"
	case LevelError:
		return "LevelError"
	case LevelFatal:
		return "LevelFatal"
	case LevelPanic:
		return "LevelPanic"
	}
	return ""
}

// 在log.New()的基础上封装一层
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

// clone函数
func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// init fields
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}
	return ll
}

func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = level.LevelToString()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			} // 如果data中没有就添加进来
		}
	}
	return data
}

func (l *Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Print(content)
	case LevelPanic:
		l.newLogger.Print(content)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}
func (l *Logger) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}
func (l *Logger) Warn(v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprintf(format, v...))
}
func (l *Logger) Error(v ...interface{}) {
	l.Output(LevelError, fmt.Sprint(v...))
}
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
}
func (l *Logger) Panic(v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprint(v...))
}
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprintf(format, v...))
}
