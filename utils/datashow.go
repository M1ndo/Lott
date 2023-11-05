package utils

import (
	"os"

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
