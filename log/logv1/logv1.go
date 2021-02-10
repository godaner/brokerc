package logv1

import (
	"fmt"
	"github.com/Shopify/sarama"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var n = fmt.Sprintln()

type LoggerV1 struct {
	DebugWriter, InfoWriter, WarnWriter, ErrorWriter io.Writer
}

func (l *LoggerV1) SetDebug(debug bool) {
	if debug {
		l.DebugWriter = os.Stdout
	} else {
		l.DebugWriter = ioutil.Discard
	}
	sarama.Logger = log.New(l.DebugWriter, "KAFKA ", 0)
}

func (l *LoggerV1) Debugf(fms string, arg ...interface{}) {
	l.DebugWriter.Write([]byte(fmt.Sprintf(fms, arg...) + n))
}

func (l *LoggerV1) Debug(arg ...interface{}) {
	l.DebugWriter.Write([]byte(fmt.Sprint(arg...) + n))
}

func (l *LoggerV1) Infof(fms string, arg ...interface{}) {
	l.InfoWriter.Write([]byte(fmt.Sprintf(fms, arg...) + n))
}

func (l *LoggerV1) Info(arg ...interface{}) {
	l.InfoWriter.Write([]byte(fmt.Sprint(arg...) + n))
}

func (l *LoggerV1) Warningf(fms string, arg ...interface{}) {
	l.WarnWriter.Write([]byte(fmt.Sprintf(fms, arg...) + n))
}

func (l *LoggerV1) Warning(arg ...interface{}) {
	l.WarnWriter.Write([]byte(fmt.Sprint(arg...) + n))
}

func (l *LoggerV1) Errorf(fms string, arg ...interface{}) {
	l.ErrorWriter.Write([]byte(fmt.Sprintf(fms, arg...) + n))
}

func (l *LoggerV1) Error(arg ...interface{}) {
	l.ErrorWriter.Write([]byte(fmt.Sprint(arg...) + n))
}
