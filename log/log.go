package log

import (
	"os"
	"strings"

	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
)

const (
	ProdPath = "log/error.log"
	TestPath = "../log/error.log"
)

type Logger struct {
	*logrus.Logger
}

func (l *Logger) ClientError() {
	render := figlet4go.NewAsciiRender()

	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorRed,
	}

	str, err := render.RenderOpts("Server error", options)
	if err != nil {
		l.Errorln(err)
	}
	l.Fatalln(str)
}

func (l *Logger) OpenLogFile() *os.File {
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

func New() *Logger {
	log := new(Logger)
	log.Logger = logrus.New()
	log.Logger.SetOutput(log.OpenLogFile())
	log.Logger.SetLevel(logrus.FatalLevel)
	return log
}
