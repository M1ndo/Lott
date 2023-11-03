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
}

const (
	Year    = 2023
	NumIter = 11
)

// Analyzes Everyday of all months of a year
func (m *DB) AnalyzeSDay() {
	m.L.Logger.Log.Info("Analyzing All Stored 10 months")
	datahold := make([]DataHold, 0)
	for n := 1; n < NumIter; n++ {
		for _, day := range Days {
			for i := 1; i < 4; i++ {
				m.L.Logger.Log.Debugf("Collecting Month %d, day %s, Entry %d", n, day, i)

				nums, err := m.NumGetSday(Year, n, day, Times[i])
				if err != nil {
					m.L.Logger.Log.Error(err)
					continue
				}

				numbers := make([]uint64, len(nums))
				copy(numbers, nums)

				Unlisted := m.L.FindDupNums(numbers, true)
				rangedNums, err := FromJson([]byte(Unlisted), true)
				if err != nil {
					m.L.Logger.Log.Error(err)
					continue
				}

				key := fmt.Sprintf("%s-%s", day, Times[i])
				jsonNums := m.L.toJson(rangedNums)
				datahold = append(datahold, DataHold{Nums: jsonNums, DayT: key, Month: n, Year: Year})
			}
		}
	}
	err := m.InsertTo("AnalyzedD", datahold)
	if err != nil {
		m.L.Logger.Log.Error(err)
	}
}
