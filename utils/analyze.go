package utils

import (
	"fmt"
	"strings"
)

var (
	Days = []string{"Monday", "Tuesday", "Thursday", "Wednesday", "Friday", "Saturday", "Sunday"}
)

type DataHold struct {
	Nums  string
	DayT  string
	Month int
	Year  int
	Day   int
}

type myRes map[string][]uint64
type myResult map[string]myRes
type day map[string]int64
type month map[int64]day

const (
	NumIter = 13
)

// Analyze everyday of all monthday of a year.
func (m *DB) AnalyzeDay(year int) {
	m.L.Logger.Log.Infof("Analyzing All Every day (1 to 30) Stored 12 months of %d", year)
	datahold := make([]DataHold, 0)
	for day := 1; day < 32; day++ {
		for i := 1; i < 4; i++ {
			m.L.Logger.Log.Debugf("Collecting Year %d, day %d, Entry %d", year, day, i)
			nums, err := m.NumGetDay(year, day, Times[i])
			if err != nil {
				m.L.Logger.Log.Error(err)
				return
			}
			numbers := make([]uint64, len(nums))
			copy(numbers, nums)
			Unlisted := m.L.FindDupNums(numbers, true)
			key := fmt.Sprintf("%d-%s", day, Times[i])
			datahold = append(datahold, DataHold{Nums: Unlisted, DayT: key, Day: day, Year: year})
		}
	}
	err := m.InsertTo("AnalyzedD", datahold)
	if err != nil {
		m.L.Logger.Log.Error(err)
	}
}

// Analyzes Everysday of all months of a year
func (m *DB) AnalyzeSDay(year int) {
	m.L.Logger.Log.Infof("Analyzing All Days \"Monday to Sunday\" Stored 12 months of %d", year)
	datahold := make([]DataHold, 0)
	for n := 1; n < NumIter; n++ {
		for _, day := range Days {
			for i := 1; i < 4; i++ {
				m.L.Logger.Log.Debugf("Collecting Month %d, day %s, Entry %d", n, day, i)
				nums, err := m.NumGetSday(year, n, day, Times[i])
				if err != nil || len(nums) == 0 {
					m.L.Logger.Log.Error(err)
					continue
				}
				numbers := make([]uint64, len(nums))
				copy(numbers, nums)
				Unlisted := m.L.FindDupNums(numbers, true)
				key := fmt.Sprintf("%s-%s", day, Times[i])
				datahold = append(datahold, DataHold{Nums: Unlisted, DayT: key, Month: n, Year: year})
			}
		}
	}
	err := m.InsertTo("AnalyzedD", datahold)
	if err != nil {
		m.L.Logger.Log.Error(err)
	}
}

// Most Used Numbers for a month
func (m *DB) AnalyzeM(year int) {
	m.L.Logger.Log.Info("Analyzing all numbers of all months of %d", year)
	datahold := make([]DataHold, 0)
	for i := 1; i < NumIter; i++ {
		nums, err := m.NumGetMonth(year, i, "")
		if err != nil || len(nums) == 0 {
			m.L.Logger.Log.Error(err)
		}
		numbers := make([]uint64, len(nums))
		copy(numbers, nums)
		Unlisted := m.L.FindDupNums(numbers, true)
		datahold = append(datahold, DataHold{Nums: Unlisted, Month: i, Year: year})
	}
	err := m.InsertTo("AnalyzedD", datahold)
	if err != nil {
		m.L.Logger.Log.Error(err)
	}
}

// Most Used Numbers for a month
func (m *DB) AnalyzeY(year int) {
	m.L.Logger.Log.Infof("Analyzing all numbers of %d", year)
	datahold := make([]DataHold, 0)
	nums, err := m.NumGetYear(year)
	if err != nil || len(nums) == 0 {
		m.L.Logger.Log.Error(err)
	}
	numbers := make([]uint64, len(nums))
	copy(numbers, nums)
	Unlisted := m.L.FindDupNums(numbers, true)
	datahold = append(datahold, DataHold{Nums: Unlisted, Year: year})
	err = m.InsertTo("AnalyzedD", datahold)
	if err != nil {
		m.L.Logger.Log.Error(err)
	}
}

// Analyzing each entry in a whole year.
// Getting an entry of the morning and add it run a match on it.
// Could be missed by 1 number.
func (m *DB) AnalyzeEntry(year int) {
	var result myResult
	m.L.Logger.Log.Infof("Analyzing Jackpot year %d", year)
	result = make(myResult)
	totalDays := make(map[int]int)

	for d := 1; d < 13; d++ {
		if _, ok := totalDays[d]; !ok {
			totalDays[d] = m.Total_days(d)
		}
		m.L.Logger.Log.Infof("Analyzing Jackpot year %d, month %d", year, d)
		for i := 1; i < totalDays[d]+1; i++ {
			for n := 1; n < 4; n++ {
				entry_numbers, err := m.NumGetDayTime(year, i, d, Times[n])
				m.L.HandleError(err)
				returned := m.centry(entry_numbers, year)
				if len(returned) != 0 {
					date := fmt.Sprintf("%s_%s", NumToDate(d, year, i), Times[n])
					// m.L.Logger.Log.Debugf("Found and adding %s", date)
					result[date] = returned
				}
			}
		}
	}
	for date, matches := range result {
		entry_date, entry_time := m.L.SplitDate(date)
		year, month, day := m.L.dateToNum(entry_date)
		for k, c := range matches {
			mentry_date, mentry_time := m.L.SplitDate(k)
			myear, mmonth, mday := m.L.dateToNum(mentry_date)
			m.ImportJack(year, month, day, myear, mmonth, mday, entry_time, mentry_time, c)
		}
	}
	// m.L.myResultToJson(result)
	// data := m.L.jsonTomyResult()
	// // fmt.Println(data)
}

// Continuation of Function Above
func (m *DB) centry(row []uint64, year int) map[string][]uint64 {
	matchedNum := make(map[string][]uint64)
	totalDays := make(map[int]int)

	for i := 1; i < 13; i++ {
		if _, ok := totalDays[i]; !ok {
			totalDays[i] = m.Total_days(i)
		}
		for d := 1; d < totalDays[i]+1; d++ {
			for n := 1; n < 4; n++ {
				thisday_numbers, err := m.NumGetDayTime(year, d, i, Times[n])
				m.L.HandleError(err)
				tempcomp := make(map[uint64]struct{})
				for _, num := range thisday_numbers {
					tempcomp[num] = struct{}{}
				}
				count := 0
				var matches []uint64
				for _, num := range row {
					if _, exists := tempcomp[num]; exists {
						count++
						matches = append(matches, num)
					}
				}
				if count >= 11 && count != 20 {
					date := fmt.Sprintf("%s_%s", NumToDate(i, year, d), Times[n])
					matchedNum[date] = matches
				}
			}
		}
	}
	return matchedNum
}

// Sort Jackpot Data
func (m *DB) GetMostDate(year int) month {
	data := m.AnalyzedJackPot(year)
	listings := make(month)
	for d, m := range data {
		dayformat := strings.Split(d, "_")
		dayd := fmt.Sprintf("%s_%s", dayformat[0], dayformat[1])
		if listings[m] == nil {
			listings[m] = make(day)
		}
		daye := listings[m]
		if daye[dayd] >= 0 {
			daye[dayd] += 1
		}
		listings[m] = daye
	}
	return listings
}
