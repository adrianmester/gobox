package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func GetLogger(name string, debug bool) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatMessage = func(i interface{}) string {
		if name == "" {
			return fmt.Sprintf("%s", i)
		}
		return fmt.Sprintf("[%s] %s", name, i)
	}
	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}
	return log.Output(output).Level(logLevel)
}
