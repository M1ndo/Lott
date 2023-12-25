package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (app *Lott) SdayTable(Year, Month int, Sday string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Year", "Day", "Number", "Count"})
	app.L.Logger.Log.Info("Getting Data from Sql.")
	data, err := app.GetMostUsed("AnalyzedD", Sday, Year, 0, Month)
	if err != nil {
		app.L.Logger.Log.Error(err)
	}
	dataNums, err := FromJson([]byte(data))
	for _, d := range dataNums {
		t.AppendRows([]table.Row{{Year, Sday, d.Number, d.Count}})
	}
	// t.SetStyle(table.StyleColoredBright)
	t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	t.Render()
	// app.Logger.Log.Info("Data Return from sql")
	// app.Logger.Log.Info(data)
}

// This function gets date from analyze GetMostDate()
func (app *Lott) ShowAnalyzed(Year int) {
	data := app.DB.GetMostDate(Year)
	// app.Logger.Log.Infof("%v", data)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Year", "Original Date", "MMonth", "MDay", "MEntry", "Numbers", "Counts"})
	for i := 1; i < NumIter; i++ {
		day := data[int64(i)]
		for day, count := range day {
			daytime := strings.Split(day, "_")
			dayed := stringToInt(daytime[0])
			entry := daytime[1]
			if count > 2 {
				numbers := app.DB.getJackPotNums(Year, i, dayed, entry)
				originaldate := app.DB.OriginalNumbers(Year, i, dayed, entry)
				// app.Logger.Log.Infof("Original Date %v", originaldate)
				// for orday, ormonth := range originaldate {
				// 	date_time := strings.Split(orday, "_")
				// 	date_day := date_time[0]
				// 	date_entry := date_time[1]
				// 	app.Logger.Log.Infof("%s, %s // %d", date_day, date_entry, ormonth)
				// 	originnumbers, err := app.DB.getNumbers(fmt.Sprintf("SELECT numbers from jackpot where year = %d and month = %d and day = %s and entry_time = '%s'", Year, ormonth, date_day, date_entry))
				// 	app.L.HandleError(err)
				// 	app.Logger.Log.Infof("Dates Number %v", originnumbers)
				// 	numorigin := ConvertToStr(originnumbers)
				// 	t.AppendRows([]table.Row{{Year, fmt.Sprintf("%s/%d", orday, ormonth), i, dayed, entry, numorigin, count}})
				// }
				t.AppendRows([]table.Row{{Year, originaldate, i, dayed, entry, numbers, count}})
			}
		}
	}
	t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	t.Render()
}

// func (app *Lott) ShowWeird() {
// 	alldata := make(map[int][]int[][]int)
// 	for year := 2020; year < 2024; year++ {
// 		data := app.DB.GetMostDate(year)
// 		alldata[year] = data
// 	}
// 	app.Logger.Log.Infof("%v", alldata)
// }

func (app *Lott) ShowAnalyzed2() {
	compared := make(map[string]int)
	// founded := make(map[string]int)
	for year := 2020; year < 2024; year++ {
		data := app.DB.GetMostDate(year)
		// app.Logger.Log.Infof("%v", data)
		for month := 1; month < NumIter; month++ {
			days := data[int64(month)]
			for day := range days {
				daytime := strings.Split(day, "_")
				dayed := stringToInt(daytime[0])
				entry := daytime[1]
				date := app.L.FormatDate(month, dayed, entry)
					compared[date] += 1
				}
			}
		}
		app.Logger.Log.Infof("%v", compared)
		os.Exit(0)
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Date", "Counts"})
		for date, count := range compared {
			if count >= 4 {
				t.AppendRows([]table.Row{{date, count}})
			}
		}
		t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
		t.Render()
}

// Get everynumber from 2020/2023 and check if its repeated atleast 3times which should exactly be the same numbers.
// Since all the numbers would be in the same order like (1,2,3,4,5,6,7,8,9,10)
// Having a match would be insanely good luck
func (app *Lott) InsanelyLucky() {
	data := make(map[int][][]uint64)
	for year := 2020; year < 2024; year++ {
		dd := app.DB.JackPotNumbers(year)
		data[year] = append(dd[year])
	}
	allthem := make(map[string]int)
	for _, nums := range data {
		for _, num := range nums {
			allthem[fmt.Sprint(num)]++
		}
	}
	for num, count := range allthem {
		if count >= 3 {
			app.Logger.Log.Info("Jackpot")
			app.Logger.Log.Infof("%s", num)
		}
	}
}
