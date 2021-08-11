package helpful

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogError
	LogNone
)

func LoggerFromPrinter(p Printer) defaultLogger {
	return defaultLogger{
		p: p,
	}
}

type defaultPrinter struct {
}

func (d defaultPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

var DefaultLogger = defaultLogger{
	p: defaultPrinter{},
}

type defaultLogger struct {
	p     Printer
	level LogLevel
}

func (d defaultLogger) Errorf(format string, args ...interface{}) {
	switch d.level {
	case LogInfo, LogError:
		d.p.Printf(d.time()+": "+format+"\r\n", args...)
		return
	default:
		return
	}
}

func (d defaultLogger) Infof(format string, args ...interface{}) {
	switch d.level {
	case LogInfo:
		d.p.Printf(d.time()+": "+format+"\r\n", args...)
		return
	default:
		return
	}
}

func (d defaultLogger) time() string {
	n := time.Now()
	return n.Format(time.Stamp)
}

func (d defaultLogger) WithLevel(level LogLevel) defaultLogger {
	res := d
	res.level = level
	return res
}
