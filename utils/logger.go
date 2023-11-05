package utils

import (
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
