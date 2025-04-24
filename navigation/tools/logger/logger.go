package logger

import (
	"encoding/json"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

const logPatternWithData = "{\"level\":%d,\"time\":\"%s\",\"pid\":%d,\"hostname\":\"%s\",\"name\":\"%s\",\"data\":\"%s\",\"msg\":\"%s\",\"v\":1}\n"
const logPatternWithoutData = "{\"level\":%d,\"time\":\"%s\",\"pid\":%d,\"hostname\":\"%s\",\"name\":\"%s\",\"msg\":\"%s\",\"v\":1}\n"

type Config struct {
	ServiceName string
	LogLevel    string
}

type TAPSIFormatter struct {
	pid      int
	name     string
	hostname string
}

func (formatter *TAPSIFormatter) Format(entry *log.Entry) ([]byte, error) {
	var loggableFields = log.Fields{}
	if errField, hasErrField := entry.Data[log.ErrorKey]; hasErrField {
		if loggable, ok := errField.(LoggableError); ok {
			loggableFields = loggable.GetFields()
		}
	}

	for k, v := range entry.Data {
		switch k {
		case "level", "time", "pid", "hostname", "name", "data", "msg", "v":
		default:
			loggableFields[k] = v
		}
	}

	noData := false
	if len(loggableFields) == 0 {
		noData = true
	}
	if errorField, hasErrorField := loggableFields[log.ErrorKey]; hasErrorField {
		loggableFields[log.ErrorKey] = fmt.Sprintf("%s", errorField)
	}
	data, jsonErr := json.Marshal(loggableFields)
	if jsonErr != nil {
		stdlog.Println("Error while logging: Can not create json from log's data fields")
		noData = true
	}

	level := int32((7 - entry.Level) * 10)
	logtime := entry.Time.Format(time.RFC3339)
	msg := entry.Message
	msg = strings.ReplaceAll(msg, "\\", "\\\\")
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	var str string
	if noData {
		str = fmt.Sprintf(
			logPatternWithoutData,
			level,
			logtime,
			formatter.pid,
			formatter.hostname,
			formatter.name,
			msg,
		)
	} else {
		datastr := string(data)
		datastr = strings.ReplaceAll(datastr, "\\", "\\\\")
		datastr = strings.ReplaceAll(datastr, "\"", "\\\"")
		str = fmt.Sprintf(
			logPatternWithData,
			level,
			logtime,
			formatter.pid,
			formatter.hostname,
			formatter.name,
			datastr,
			msg,
		)
	}

	out := []byte(str)
	return out, nil
}

func Init(config Config) {
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		stdlog.Printf("LOG_LEVEL environment variable is not correct: %s", err)
		level = log.DebugLevel
	}

	if level < log.ErrorLevel {
		level = log.ErrorLevel
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	pid := os.Getpid()

	log.SetFormatter(&TAPSIFormatter{
		pid:      pid,
		hostname: hostname,
		name:     config.ServiceName,
	})
	log.SetOutput(io.Discard)
	log.SetLevel(level)

	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	})
}

type LoggableError interface {
	GetFields() log.Fields
}
