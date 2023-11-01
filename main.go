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

type Lott struct {
	utils.Utils
	utils.DB
}

var (
	a = &Lott{
		Utils: utils.Utils{},
		DB: utils.DB{},
	}
)

func main() {
	// Set Logger
	a.SetLogger()
	a.DB.L = &a.Utils
	// DB INit
	errCh := a.DB.InitializeDB()
	for err := range errCh {
		a.Logger.Log.Error(err)
	}
	// Parse Flags
	VarsUsed := utils.Parse()
	if VarsUsed.HtmlFile != "" {
		data := utils.GetDataFF(VarsUsed.HtmlFile)
		err := a.DB.ImportData(data, VarsUsed.Year, VarsUsed.Month)
		if err != nil {
			a.Logger.Log.Error(err)
		}
	}
	// if VarsUsed.File != "" {
	// 	AllNums := a.ReadNumbers(VarsUsed.File)
	// 	FindDupls := a.FindDupNums(AllNums)
	// 	a.Logger.Log.Infof("Numbers that are most picked (high to low)")
	// 	for _, num := range FindDupls {
	// 		a.Logger.Log.Infof("Number %d, Count %d", num.Number, num.Count)
	// 	}
	// }
	// // Get Numbers of a day
	// nums, err := a.DB.NumGetDay(VarsUsed.Year, VarsUsed.Month, VarsUsed.Day, VarsUsed.Entry)
	// if err != nil {
	// 	a.Logger.Log.Error(err)
	// }
	// a.Logger.Log.Infof("Day %d %v Numbers", VarsUsed.Day, nums)
	// // Get Numbers of a String day
	// nums, err = a.DB.NumGetSday(VarsUsed.Year, VarsUsed.Month, VarsUsed.Sday, VarsUsed.Entry)
	// if err != nil {
	// 	a.Logger.Log.Error(err)
	// }
	// a.Logger.Log.Infof("Day %s %v Numbers", VarsUsed.Sday, nums)

	// // Get Number of a month
	// nums, err = a.DB.NumGetMonth(VarsUsed.Year, VarsUsed.Month, VarsUsed.Entry)
	// if err != nil {
	// 	a.Logger.Log.Error(err)
	// }
	// a.Logger.Log.Infof("Month %d %v Numbers", VarsUsed.Month, nums)

	// Get Numbers of a sday
	a.Logger.Log.Info("Nothing to do")
}
