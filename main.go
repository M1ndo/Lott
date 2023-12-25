package main

// As i noticed in the the numbers that were picked on each datetime of the day.
// I've realized they are all select from low to high. Incrementating Numbers.
// What can we do to try to predict those numbers.
// are they fully random? -> Probably yes.
// though they have a pattern that the follow? YES.
// Is there an algorithm in play? Yes -> Then are not truly random

// Approaches Taken.
// Through the month of october from 1 to 24 (Oct 24 2023)
// Gathering all Numbers that were picked and running a based check to find
// Numbers that more picked often .
// Starting from 1 day to 2days - 3days - 5days - 10days - 20days

import (
	"github.com/m1ndo/Lott/utils"
)

func main() {
	app := utils.NewApp()
	// app.DB.AnalyzeEntry(2023)
	// app.DB.AnalyzeSDay(2023)
	// app.DB.AnalyzeDay(2023)
	// app.DB.AnalyzeM(2023)
	// app.DB.AnalyzeY(2023)
	// app.SdayTable(2023, 0, "Monday-morning")
	// app.SdayTable(2023, 0, "Monday-noon")
	// app.SdayTable(2023, 0, "Monday-night")
	app.Logger.Log.Info("Nothing to do")
}
