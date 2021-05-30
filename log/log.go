package log

import (
	"os"
	"strings"

	"github.com/Kenshin/cprint"
	"github.com/sirupsen/logrus"
)

var (
	Logger = New()
)

const (
	ProdPath = "log/error.log"
	TestPath = "../log/error.log"
)

type Log struct {
	*logrus.Logger
}

func (l *Log) ClientError() {
	cprint.P(cprint.ERROR, "\nClient does not respond, try to restart ...")
}

func (l *Log) OpenLogFile() *os.File {
	path := ProdPath
	if f, err := os.Executable(); err == nil && strings.Contains(f, "test") {
		path = TestPath
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err := os.Create(path); err != nil {
			l.Fatalln(err)
		}
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		l.Fatalln(err)
	}
	return f
}

func New() *Log {
	log := new(Log)
	log.Logger = logrus.New()
	log.Logger.SetOutput(log.OpenLogFile())
	log.Logger.SetLevel(logrus.FatalLevel)
	return log
}
