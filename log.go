package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type MixedJSON map[string]interface{}

type Log struct {
	name     string
	mutex    sync.Mutex
	out      io.Writer
	buffer   []byte
	logLevel int
}

type Formatted struct {
	time     time.Time
	hostname string
	pid      int
	level    int
	fields   *MixedJSON
}

var LogLevels = map[string]int{
	"fatal": 60,
	"error": 50,
	"warn":  40,
	"debug": 20,
	"trace": 10}

func New(n string) Log {
	rawEnv, existsLogLevel := os.LookupEnv("LOG_LEVEL")

	logger := Log{
		name:     n,
		mutex:    sync.Mutex{},
		out:      os.Stderr,
		logLevel: LogLevels["error"]}

	if existsLogLevel {
		env, err := LogLevels[strings.ToLower(rawEnv)]
		if !err {
			logger.Fatal("Bad LOG_LEVEL")
			os.Exit(1)
		}
		logger.logLevel = env
	}

	return logger
}

func (l *Log) shouldLog(level int) bool {
	if l.logLevel <= level {
		return true
	} else {
		return false
	}
}

func (l *Log) Output(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.buffer = l.buffer[:0]
	l.buffer = append(l.buffer, message...)
	l.out.Write(l.buffer)
}

func (l *Log) Fatal(message interface{}) {
	l.Fatalf(nil, message)
}

func (l *Log) Error(message interface{}) {
	l.Errorf(nil, message)
}

func (l *Log) Warn(message interface{}) {
	l.Warnf(nil, message)
}

func (l *Log) Debug(message interface{}) {
	l.Debugf(nil, message)
}

func (l *Log) Trace(message interface{}) {
	l.Tracef(nil, message)
}

func (l *Log) Fatalf(fields MixedJSON, message interface{}) {
	level := LogLevels["fatal"]
	if !l.shouldLog(level) {
		return
	}
	l.Output(l.toString(level, &fields, message))
	os.Exit(1)
}

func (l *Log) Errorf(fields MixedJSON, message interface{}) {
	level := LogLevels["error"]
	if !l.shouldLog(level) {
		return
	}
	l.Output(l.toString(level, &fields, message))
}

func (l *Log) Warnf(fields MixedJSON, message interface{}) {
	level := LogLevels["warn"]
	if !l.shouldLog(level) {
		return
	}
	l.Output(l.toString(level, &fields, message))
}

func (l *Log) Debugf(fields MixedJSON, message interface{}) {
	level := LogLevels["debug"]
	if !l.shouldLog(level) {
		return
	}
	l.Output(l.toString(level, &fields, message))
}

func (l *Log) Tracef(fields MixedJSON, message interface{}) {
	level := LogLevels["trace"]
	if !l.shouldLog(level) {
		return
	}
	l.Output(l.toString(level, &fields, message))
}

func (l *Log) toString(level int, fields *MixedJSON, message interface{}) string {
	formatted := l.format(level, fields, message)
	encoded, _ := json.Marshal(formatted)
	return fmt.Sprintf("%s\n", encoded)
}

func (l *Log) format(level int, fields *MixedJSON, message interface{}) *MixedJSON {
	hostname, _ := os.Hostname()
	formatted := &MixedJSON{
		"time":     time.Now(),
		"hostname": hostname,
		"pid":      os.Getpid(),
		"level":    level,
		"fields":   &fields,
		"v":        1,
		"name":     l.name,
		"msg":      message}
	return formatted
}
