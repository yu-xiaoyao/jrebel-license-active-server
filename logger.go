// @author: yu-xiaoyao
// @github: https://github.com/yu-xiaoyao/jrebel-license-active-server
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const defaultFlag = log.Ldate | log.Ltime | log.Lmsgprefix | log.Lshortfile

type Level int

func (l *Level) String() string {
	return strconv.Itoa(int(*l))
}
func (l *Level) Set(value string) error {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	i := int(i64)
	if i >= int(Trace) && i <= int(Off) {
		*l = Level(i) // 安全转换（注意范围）
		return nil
	}
	return nil
}

const (
	Trace Level = iota
	// Debug level
	Debug
	// Info level
	Info
	// Warn level
	Warn
	// Error level
	Error
	Off
)

type ILogger interface {
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool

	Tracef(format string, args ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, args ...interface{})

	Traceln(v ...interface{})
	Debugln(v ...interface{})
	Infoln(v ...interface{})
	Warnln(v ...interface{})
	Errorln(v ...interface{})
}

type BaseLogger struct {
	level Level
	l     *log.Logger
}

func (b *BaseLogger) IsTrace() bool {
	return b.level <= Trace
}
func (b *BaseLogger) IsDebug() bool {
	return b.level <= Debug
}
func (b *BaseLogger) IsInfo() bool {
	return b.level <= Info
}
func (b *BaseLogger) IsWarn() bool {
	return b.level <= Warn
}
func (b *BaseLogger) IsError() bool {
	return b.level <= Error
}
func (b *BaseLogger) Tracef(format string, v ...interface{}) {
	if b.level <= Trace {
		b.l.SetPrefix("[TRACE] ")
		b.l.Printf(format, v...)
	}
}
func (b *BaseLogger) Debugf(format string, v ...interface{}) {
	if b.level <= Debug {
		b.l.SetPrefix("[DEBUG] ")
		b.l.Printf(format, v...)
	}
}
func (b *BaseLogger) Infof(format string, v ...interface{}) {
	if b.level <= Info {
		b.l.SetPrefix("[ INFO] ")
		b.l.Printf(format, v...)
	}
}
func (b *BaseLogger) Warnf(format string, v ...interface{}) {
	if b.level <= Warn {
		b.l.SetPrefix("[ WARN] ")
		b.l.Printf(format, v...)
	}
}
func (b *BaseLogger) Errorf(format string, v ...interface{}) {
	if b.level <= Error {
		b.l.SetPrefix("[ERROR] ")
		b.l.Printf(format, v...)
	}
}
func (b *BaseLogger) Traceln(v ...interface{}) {
	if b.level <= Trace {
		b.l.SetPrefix("[TRACE] ")
		b.l.Println(v...)
	}
}
func (b *BaseLogger) Debugln(v ...interface{}) {
	if b.level <= Debug {
		b.l.SetPrefix("[DEBUG] ")
		b.l.Println(v...)
	}
}
func (b *BaseLogger) Infoln(v ...interface{}) {
	if b.level <= Info {
		b.l.SetPrefix("[ INFO] ")
		b.l.Println(v...)
	}
}
func (b *BaseLogger) Warnln(v ...interface{}) {
	if b.level <= Warn {
		b.l.SetPrefix("[ WARN] ")
		b.l.Println(v...)
	}
}
func (b *BaseLogger) Errorln(v ...interface{}) {
	if b.level <= Error {
		b.l.SetPrefix("[ERROR] ")
		b.l.Println(v...)
	}
}

func NewConsoleLogger(level Level, flag int) ILogger {
	if flag <= 0 {
		flag = defaultFlag
	}
	return &BaseLogger{
		l:     log.New(os.Stdout, "", flag),
		level: level,
	}
}

type FileLogger struct {
	mu          sync.Mutex
	logDir      string
	maxSize     int64
	currentSize int64
	file        *os.File
	currentDate string
	level       Level
}

// Write 实现 io.Writer 接口，处理文件分割逻辑
func (fw *FileLogger) Write(p []byte) (n int, err error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now().Format("2006-01-02")
	if fw.file == nil || now != fw.currentDate || (fw.currentSize+int64(len(p))) > fw.maxSize {
		if err := fw.rotate(); err != nil {
			fmt.Fprintf(os.Stderr, "Switch Log failed: %v\n", err)
		}
	}

	n, err = fw.file.Write(p)
	fw.currentSize += int64(n)
	return n, err
}

func (fw *FileLogger) rotate() error {
	if fw.file != nil {
		fw.file.Close()
	}

	_ = os.MkdirAll(fw.logDir, 0755)
	now := time.Now().Format("2006-01-02")
	fw.currentDate = now

	var newFilename string
	for i := 1; ; i++ {
		newFilename = filepath.Join(fw.logDir, fmt.Sprintf("%s.%d.log", now, i))
		if _, err := os.Stat(newFilename); os.IsNotExist(err) {
			break
		}
		if i > 999 {
			break
		} // 安全阈值
	}

	f, err := os.OpenFile(newFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, _ := f.Stat()
	fw.currentSize = info.Size()
	fw.file = f
	return nil
}

// StartCleanTask 启动后台清理协程
func (fw *FileLogger) StartCleanTask(maxDays int) {
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		for range ticker.C {
			files, _ := os.ReadDir(fw.logDir)
			cutoff := time.Now().AddDate(0, 0, -maxDays)
			for _, f := range files {
				if f.IsDir() || len(f.Name()) < 10 {
					continue
				}
				fileDate, err := time.Parse("2006-01-02", f.Name()[:10])
				if err == nil && fileDate.Before(cutoff) {
					_ = os.Remove(filepath.Join(fw.logDir, f.Name()))
				}
			}
		}
	}()
}

func NewFileLogger(dir string, maxSizeMB int, keepDays int, level Level, flag int) ILogger {
	if flag <= 0 {
		flag = defaultFlag
	}
	fWriter := &FileLogger{
		logDir:  dir,
		maxSize: int64(maxSizeMB) * 1024 * 1024,
	}
	if keepDays > 0 {
		fWriter.StartCleanTask(keepDays)
	}
	return &BaseLogger{
		l:     log.New(fWriter, "", flag),
		level: level,
	}
}

type MultiLogger struct {
	mainLogger ILogger
	loggers    []ILogger
}

func (m *MultiLogger) IsTrace() bool {
	return m.mainLogger.IsTrace()
}
func (m *MultiLogger) IsDebug() bool {
	return m.mainLogger.IsDebug()
}
func (m *MultiLogger) IsInfo() bool {
	return m.mainLogger.IsInfo()
}
func (m *MultiLogger) IsWarn() bool {
	return m.mainLogger.IsWarn()
}
func (m *MultiLogger) IsError() bool {
	return m.mainLogger.IsError()
}
func (m *MultiLogger) Tracef(format string, v ...interface{}) {
	m.mainLogger.Tracef(format, v...)
	for _, l := range m.loggers {
		l.Tracef(format, v...)
	}
}
func (m *MultiLogger) Debugf(format string, v ...interface{}) {
	m.mainLogger.Debugf(format, v...)
	for _, l := range m.loggers {
		l.Debugf(format, v...)
	}
}
func (m *MultiLogger) Infof(format string, v ...interface{}) {
	m.mainLogger.Infof(format, v...)
	for _, l := range m.loggers {
		l.Infof(format, v...)
	}
}
func (m *MultiLogger) Warnf(format string, v ...interface{}) {
	m.mainLogger.Warnf(format, v...)
	for _, l := range m.loggers {
		l.Warnf(format, v...)
	}
}
func (m *MultiLogger) Errorf(format string, v ...interface{}) {
	m.mainLogger.Errorf(format, v...)
	for _, l := range m.loggers {
		l.Errorf(format, v...)
	}
}
func (m *MultiLogger) Traceln(v ...interface{}) {
	m.mainLogger.Traceln(v...)
	for _, l := range m.loggers {
		l.Traceln(v...)
	}
}
func (m *MultiLogger) Debugln(v ...interface{}) {
	m.mainLogger.Debugln(v...)
	for _, l := range m.loggers {
		l.Debugln(v...)
	}
}
func (m *MultiLogger) Infoln(v ...interface{}) {
	m.mainLogger.Infoln(v...)
	for _, l := range m.loggers {
		l.Infoln(v...)
	}
}
func (m *MultiLogger) Warnln(v ...interface{}) {
	m.mainLogger.Warnln(v...)
	for _, l := range m.loggers {
		l.Warnln(v...)
	}
}
func (m *MultiLogger) Errorln(v ...interface{}) {
	m.mainLogger.Errorln(v...)
	for _, l := range m.loggers {
		l.Errorln(v...)
	}
}

func NewMultiLogger(mainLogger ILogger, loggers ...ILogger) ILogger {
	return &MultiLogger{
		mainLogger: mainLogger,
		loggers:    loggers}
}
