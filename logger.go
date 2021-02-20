package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	Formatter := new(log.TextFormatter)

	/* Set Print Timestamp And Format */
	Formatter.FullTimestamp = true
	Formatter.TimestampFormat = "2006-01-02T15:04:05.000000"
	Formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
		_, filename := path.Split(f.File)
		retFile := fmt.Sprintf("[ %s : %d ]", filename, f.Line)
		return "", retFile
	}
	Formatter.DisableColors = true
	logger.SetFormatter(Formatter)
}

type logrusWriter struct {
	logger log.FieldLogger
}
