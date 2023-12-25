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
	Aday bool
	Asday bool
	Amonth bool
	Ayear bool
	EntryAnalyze bool
	Debug bool
}

func Parse() *VarsUsed {
	Vars := &VarsUsed{}
	flag.StringVar(&Vars.File, "file", "", "Input file")
	flag.StringVar(&Vars.HtmlFile, "hfile", "", "HTML File to Parse")
	flag.IntVar(&Vars.Month, "month", 0, "Month In a number of the HTML file data")
	flag.IntVar(&Vars.Year, "year", 0, "Year to set")
	flag.IntVar(&Vars.Day, "day", 0, "Day to set (1 to 31 - depending on month)")
	flag.StringVar(&Vars.Sday, "sday", "", "Day to set (Morning, Tuesday ...)")
	flag.StringVar(&Vars.Entry, "entry", "", "Entry to set (Morning, Noon, Night)")
	flag.BoolVar(&Vars.Aday, "aday", false, "Return analytics of an day (1 to 31)")
	flag.BoolVar(&Vars.Asday, "asday", false, "Return analytics of an sday (Monday to Sunday)")
	flag.BoolVar(&Vars.Amonth, "amonth", false, "Return analytics of an month (Jan to Dec)")
	flag.BoolVar(&Vars.Ayear, "ayear", false, "Return analytics of an year (2021 to 2023)")
	flag.BoolVar(&Vars.EntryAnalyze, "AnalyzeEntry", false, "Analyze Entry (require a year or None for all years)")
	flag.BoolVar(&Vars.Debug, "verbose", false, "Enable Verbose/Debug")
	flag.Parse()
	return Vars
}
