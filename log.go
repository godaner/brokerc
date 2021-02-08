package main

import (
	"github.com/godaner/brokerc/log"
	"github.com/godaner/brokerc/log/logv1"
	"io/ioutil"
	"os"
)

var logger log.Logger

func initLogger(debug bool) {
	logger = &logv1.LoggerV1{
		DebugWriter: ioutil.Discard,
		InfoWriter:  os.Stdout,
		WarnWriter:  os.Stdout,
		ErrorWriter: os.Stderr,
	}
	logger.SetDebug(debug)
}
