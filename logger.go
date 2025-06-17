// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"io"
	"log"
	"os"
	"sync/atomic"
)

const (
	// Debug level
	Debug = 1
	// Info level
	Info = 2
	// Warn level
	Warn = 3
	// Error level
	Error = 4
)

type SimpleLogger struct {
	level int64
	w     io.Writer
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func NewLogger(w io.Writer, level int64, flag int) *SimpleLogger {
	if w == nil {
		w = os.Stderr
	}

	if flag <= 0 {
		flag = log.LstdFlags
	}

	return &SimpleLogger{
		w:     w,
		level: level,
		debug: log.New(w, "[DEBUG] ", flag|log.Lmsgprefix),
		info:  log.New(w, "[INFO ] ", flag|log.Lmsgprefix),
		warn:  log.New(w, "[WARN ] ", flag|log.Lmsgprefix),
		error: log.New(w, "[ERROR] ", flag|log.Lmsgprefix),
	}
}

func (l *SimpleLogger) SetLevel(level int64) {
	if level < Debug || level > Error {
		return
	}
	atomic.StoreInt64(&l.level, level)
}

func (l *SimpleLogger) Debugln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Debug {
		return
	}
	l.debug.Println(v...)
}
func (l *SimpleLogger) Debugf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Debug {
		return
	}
	l.debug.Printf(format, v...)
}

func (l *SimpleLogger) IsDebug() bool {
	return atomic.LoadInt64(&l.level) <= Debug
}

func (l *SimpleLogger) Infoln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Info {
		return
	}
	l.info.Println(v...)
}
func (l *SimpleLogger) Infof(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Info {
		return
	}
	l.info.Printf(format, v...)
}
func (l *SimpleLogger) IsInfo() bool {
	return atomic.LoadInt64(&l.level) <= Info
}

func (l *SimpleLogger) Warnln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Warn {
		return
	}
	l.warn.Println(v...)
}
func (l *SimpleLogger) Warnf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Warn {
		return
	}
	l.warn.Printf(format, v...)
}

func (l *SimpleLogger) IsWarn() bool {
	return atomic.LoadInt64(&l.level) <= Warn
}

func (l *SimpleLogger) Errorln(v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Error {
		return
	}
	l.error.Println(v...)
}
func (l *SimpleLogger) Errorf(format string, v ...interface{}) {
	if atomic.LoadInt64(&l.level) > Error {
		return
	}
	l.error.Printf(format, v...)
}
func (l *SimpleLogger) IsError() bool {
	return atomic.LoadInt64(&l.level) <= Error
}
