package utils

import (
	"fmt"
	"runtime"
	"strings"

	loggdb "github.com/m1ndo/LogGdb"
	"github.com/muesli/termenv"
)

type Logger struct {
	*loggdb.Logger
}

var Debug = loggdb.Debug

func (util *Utils) SetLogger() {
	CustomOpt := &loggdb.CustomOpt{
		Prefix:          "LoTT 👾 ",
		ReportTimestamp: true,
	}
		util.Logger = Logger{Logger: &loggdb.Logger{}}
		util.Logger.LogOptions = CustomOpt
		util.Logger.NewLogger()
		util.Logger.Log.SetColorProfile(termenv.TrueColor)
	}


func (util *Utils) HandleError(err error) {
	if err != nil {
		stackTrace := make([]uintptr, 50)
		length := runtime.Callers(1, stackTrace)
		stackTrace = stackTrace[:length]
		frames := runtime.CallersFrames(stackTrace)

		var stackLines []string
		for {
			frame, more := frames.Next()
			stackLines = append(stackLines, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
			if !more {
				break
			}
		}

		util.Logger.Log.Errorf("Error %s\nStack trace:\n%s", err, strings.Join(stackLines, "\n"))
	}
}
