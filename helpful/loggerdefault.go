package helpful

import (
	"fmt"
	"time"
)

type LogLevel int

const(
	LogInfo LogLevel = iota
	LogError
	LogNone
)

var DefaultLogger = defaultLogger{}

type defaultLogger struct {
	level LogLevel
}

func (d defaultLogger) Errorf(format string, args ...interface{}) {
	switch d.level {
	case LogInfo, LogError:
		fmt.Printf(d.time() + ": " + format + "\r\n", args...)
		return
	default:
		return
	}
}

func (d defaultLogger) Infof(format string, args ...interface{}) {
	switch d.level {
	case LogInfo:
		fmt.Printf(d.time()+": "+format+"\r\n", args...)
		return
	default:
		return
	}
}

func (d defaultLogger) time() string {
	n := time.Now()
	return n.Format(time.Stamp)
}

func (d defaultLogger) WithLevel(level LogLevel) defaultLogger{
	res := d
	res.level = level
	return res
}
