package utils


type Lott struct {
	Utils
	DB
}

func NewApp() *Lott {
	app := &Lott{
		Utils{},
		DB{},
	}
	// Logger
	app.SetLogger()
	app.DB.L = &app.Utils
	// DB INit
	errCh := app.DB.InitializeDB()
	for err := range errCh {
		app.Logger.Log.Error(err)
	}
	// Parse Flags
	app.ParserFlags()
	return app
}


func (app *Lott) ParserFlags() {
	VarsUsed := Parse()
	if VarsUsed.HtmlFile != "" {
		data := GetDataFF(VarsUsed.HtmlFile)
		err := app.DB.ImportData(data, VarsUsed.Year, VarsUsed.Month)
		if err != nil {
			app.Logger.Log.Error(err)
		}
	}
	if VarsUsed.File != "" {
		// AllNums := app.ReadNumbers(VarsUsed.File)
		// FindDupls := app.FindDupNums(AllNums, false)
		app.Logger.Log.Infof("Numbers that are most picked (high to low)")
		// for _, num := range FindDupls {
		// 	app.Logger.Log.Infof("Number %d, Count %d", num.Number, num.Count)
		// }
	}
	if VarsUsed.Debug {
		app.Logger.Log.SetLevel(Debug)
	}
	switch {
	case VarsUsed.Aday:
		app.Logger.Log.Infof("Numbers that are most picked of day %d", VarsUsed.Day)
	case VarsUsed.Asday:
		app.Logger.Log.Infof("Numbers that are most picked of day %s", VarsUsed.Sday)
	case VarsUsed.Amonth:
		app.Logger.Log.Infof("Numbers that are most picked of day %d", VarsUsed.Month)
	case VarsUsed.Ayear:
		app.Logger.Log.Infof("Numbers that are most picked of day %d", VarsUsed.Year)
	case VarsUsed.EntryAnalyze:
		// app.DB.GetMostDate()
		app.DB.MarkedDates2(VarsUsed.Month, VarsUsed.Day, VarsUsed.Entry) // Marked Dates
		// app.ShowMarked(VarsUsed.Month)
		// app.ShowAnalyzed(VarsUsed.Year)
		// app.InsanelyLucky() // Jackpot numbers
		// app.FuckingBeIt()
		// app.AnotherFuck()
		// app.MonthTwelve()
		// app.ShowAnalyzed2()
		// app.ShowWeird()
		// app.DB.AnalyzeSuccess() // All Analyzsis
		// app.DB.quickTest()
		// Analyze dates
		// if VarsUsed.Year != 0 {
		// 	app.DB.AnalyzeEntry(VarsUsed.Year)
		// } else {
		// 	for i := 2020; i < 2024; i++ {
		// 		app.DB.AnalyzeEntry(i)
		// 	}
		// }
	}
}

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
