package utils

import "flag"

type VarsUsed struct {
	File string
	HtmlFile string
	Month int
	Year int
	Day int
	Sday string
	Entry string

}

func Parse() *VarsUsed {
	Vars := &VarsUsed{}
	flag.StringVar(&Vars.File, "file", "", "Input file")
	flag.StringVar(&Vars.HtmlFile, "hfile", "", "HTML File to Parse")
	flag.IntVar(&Vars.Month, "month", 0, "Month In a number of the HTML file data")
	flag.IntVar(&Vars.Year, "year", 2023, "Year to set")
	flag.IntVar(&Vars.Day, "day", 1, "Day to set (1 to 31 - depending on month)")
	flag.StringVar(&Vars.Sday, "sday", "", "Day to set (Morning, Tuesday ...)")
	flag.StringVar(&Vars.Entry, "entry", "", "Entry to set (Morning, Noon, Night)")
	flag.Parse()
	return Vars
}
