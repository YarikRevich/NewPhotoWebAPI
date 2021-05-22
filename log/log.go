package log

import (
	"os"

	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
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
	if _, err := os.Stat("error.log"); os.IsNotExist(err) {
		if _, err := os.Create("error.log"); err != nil {
			l.Fatalln(err)
		}
	}
	f, err := os.OpenFile("error.log", os.O_WRONLY|os.O_APPEND, os.ModeAppend)
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
