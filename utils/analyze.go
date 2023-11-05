package utils

import (
	"fmt"
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
