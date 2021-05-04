package log

import (
	"os"

	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func (l *Logger)ClientError() {
	render := figlet4go.NewAsciiRender()

	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorRed,
	}

	str, err := render.RenderOpts("Server error", options)
	if err != nil{
		l.Errorln(err)
	}
	l.Fatalln(str)
}

func New() *Logger {
	log := new(Logger)
	log.Logger = logrus.New()
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(logrus.FatalLevel)
	return log
}
