package utils

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

var (
	Days           = []string{"Monday", "Tuesday", "Thursday", "Wednesday", "Friday", "Saturday", "Sunday"}
	YearSuccessive = map[int][]int{2020: {2021, 2022, 2023}, 2021: {2020, 2022, 2023}, 2022: {2020, 2021, 2023}, 2023: {2020, 2021, 2022}}
	rw             sync.RWMutex
	wg             sync.WaitGroup
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

func (m *DB) AnalyzeSuccess() {
	totalDays := make(map[int]int)
	// var result sync.Map
	// result := make(map[int]map[int][][]uint64)
	result := make(map[int]map[string][]uint64)
	for year := 2023; year < 2024; year++ { // Changed 2024 to 2021
		monthresult := make(map[string][]uint64)
		for d := 1; d < 13; d++ {
			if _, ok := totalDays[d]; !ok {
				totalDays[d] = m.Total_days(d)
			}
			fmt.Printf("Analyzing Jackpot year %d, month %d\n", year, d)
			for i := 1; i < totalDays[d]+1; i++ {
				for n := 1; n < 4; n++ {
					entry_numbers, err := m.NumGetDayTime(year, i, d, Times[n])
					m.L.HandleError(err)
					if len(entry_numbers) != 0 {
						date := fmt.Sprintf("%d/%d/%d_%s", year, d, i, Times[n])
						monthresult[date] = entry_numbers
					}
				}
			}
		}
		result[year] = monthresult
	}
	// fmt.Println(result)
	// m.startWorkers(numWorkers)
	// for year := 2023; year < 2024; year++ {
	// 	job := Job{
	// 		Year:   year,
	// 		Result: result[year],
	// 		AllResults: &allresults,
	// 	}
	// 	addJob(job)
	// 	break
	// }
	// closeAndWait()
	var wg sync.WaitGroup
	var mu sync.Mutex
	allresults := make(map[string]map[string][]uint64)
	sem := make(chan bool, 10)              // Buffered channel acting as a semaphore
	for year := 2023; year < 2024; year++ { // changed 2024 to 2021
		Years := YearSuccessive[year]
		for date, row := range result[year] {
			wg.Add(1)
			go func(row []uint64, date string) {
				defer wg.Done()
				sem <- true // Acquire a token
				dataa := m.Improvedcentry(row, date, Years)
				mu.Lock()
				allresults[date] = dataa
				mu.Unlock()
				<-sem // Release a token
			}(row, date)
		}
	}
	wg.Wait() // Wait for all goroutines to finish
	// m.L.Logger.Log.Infof("Found %v", allresults)
	//
	m.L.Logger.Log.Infof("Total matches found %d", len(allresults))
	for date, matches := range allresults {
		entry_date, entry_time := m.L.SplitDate(date)
		year, month, day := m.L.dateToNum(entry_date)
		for k, c := range matches { // replace c with _
			mentry_date, mentry_time := m.L.SplitDate(k)
			myear, mmonth, mday := m.L.dateToNum(mentry_date)
			m.ImportJack(year, month, day, myear, mmonth, mday, entry_time, mentry_time, c)
			// m.L.Logger.Log.Infof("Date %d/%d/%d %s Found matches %d/%d/%d %s ", year, month, day, entry_time, myear, mmonth, mday, mentry_time)
		}
	}
}

// Upgraded Version (Faster)
func (m *DB) Improvedcentry(row []uint64, date string, years []int) map[string][]uint64 {
	matchedNum := new(sync.Map)
	totalDays := make(map[int]int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	m.L.Logger.Log.Infof("Analyzing date %s, Total Years %v", date, years)
	for _, y := range years {
		// m.L.Logger.Log.Infof("Going Through Year %d", y)
		wg.Add(1)
		go func(date string, y2 int, row []uint64) {
			defer wg.Done()
			for i := 1; i < 13; i++ {
				mu.Lock()
				if _, ok := totalDays[i]; !ok {
					totalDays[i] = m.Total_days(i)
				}
				mu.Unlock()
				// fmt.Printf("Going through date %s month %d Sel Year %d\n", date, i, y2)
				for d := 1; d < totalDays[i]+1; d++ {
					for n := 1; n < 4; n++ {
						thisday_numbers, _ := m.NumGetDayTime(y2, d, i, Times[n])
						// m.L.HandleError(err)
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
							date := fmt.Sprintf("%s_%s", NumToDate(i, y2, d), Times[n])
							matchedNum.Store(date, matches)
						}
					}
				}
			}
		}(date, y, row)
	}
	wg.Wait()
	result := make(map[string][]uint64)
	matchedNum.Range(func(key, value interface{}) bool {
		if date, ok := key.(string); ok {
			if matches, ok := value.([]uint64); ok {
				result[date] = matches
			}
		}
		return true
	})
	return result
}

// defunc to analyze certain date
func (m *DB) quickTest() {
	query := "select year, day, mmonth, mday, mentry_time from jackpot where month = 12;"
	var results map[string][]map[int]string
	rows, err := m.DB.Query(query)
	m.L.HandleError(err)
	defer rows.Close()
	results = make(map[string][]map[int]string) // Initialize the results map
	for rows.Next() {
		var day int
		var mday int
		var mmonth int
		var mentry string
		var year int
		err := rows.Scan(&year, &day, &mmonth, &mday, &mentry)
		m.L.HandleError(err)
		key := fmt.Sprintf("%d_%d", year, day)
		innerMap := make(map[int]string)
		innerMap[mmonth] = fmt.Sprintf("%d_%s", day, mentry)
		results[key] = append(results[key], innerMap)
	}
	dd := make(map[string]int)
	for _, v := range results {
		for _, j := range v {
			for m, d := range j {
				dd[fmt.Sprintf("%d_%s", m, d)]++
			}
		}
	}
	cc := make(map[string][]string)
	for d, c := range dd {
		if c >= 2 {
			cd := strings.Split(d, "_")
			month := cd[0]
			day := cd[1]
			entry := cd[2]
			fmt.Println(month, day, entry)
			date := fmt.Sprintf("SELECT year, month, day, entry_time from jackpot where mmonth = %s and mday = %s and mentry_time = '%s' and month is 12;", month, day, entry)
			m.L.Logger.Log.Debugf("Executing Query %s", date)
			rows, err := m.DB.Query(date)
			m.L.HandleError(err)
			defer rows.Close()
			results := []string{}
			for rows.Next() {
				var day int
				var month int
				var entry string
				var year int
				err := rows.Scan(&year, &month, &day, &entry)
				m.L.HandleError(err)
				key := fmt.Sprintf("%d/%d/%d_%s", year, month, day, entry)
				results = append(results, key)
			}
			cc[d] = results
		}
	}
	fmt.Printf("%v\n", cc)
}

// Return Jackpot numbers either being shown 3times straight over span of 2020-2023
func (m *DB) AnalyzeJackPot() map[string]int {
	data := make(map[int][][]uint64)
	for year := 2020; year < 2024; year++ {
		dd := m.JackPotNumbers(year)
		data[year] = append(dd[year])
	}
	allthem := make(map[string]int)
	for _, nums := range data {
		for _, num := range nums {
			allthem[fmt.Sprint(num)]++
		}
	}
	return allthem
}

// Arrange Dates.
// Arranging Dates By Not Repeating new dates to a corresponding old dates.
// For example 2023/1/1 matching to 2021/1/1 Should be instead 2021/1/1 to 2023/1/1
func (m *DB) ArrangeDates(date, date2 string) bool {
	fmt.Println(date, date2)
	extend_date, _ := m.L.SplitDate(date)
	extended_date, _ := m.L.SplitDate(date2)
	year, month, day := m.L.dateToNum(extend_date)
	myear, mmonth, mday := m.L.dateToNum(extended_date)
	if myear > year {
		return true
	} else if myear == year {
		if mmonth >= month && mday >= day {
			return true
		}
	}
	return false
}

// Function to analyze dates, marking only repeated dates.
func (m *DB) MarkedDates(month int) map[string][]map[string]int {
	var allthem map[string]int
	allthem = m.AnalyzeJackPot()
	results := make(map[string]int)
	for num, count := range allthem {
		if count >= 2 {
			numed := strings.Trim(num, "[]")
			Num := strings.Replace(numed, " ", ",", 11)
			date, date2 := m.GetjDates(Num)
			for i := range date {
				if m.ArrangeDates(date[i], date2[i]) {
					dated, entry := m.L.SplitDate(date2[i])
					_, month, day := m.L.dateToNum(dated)
					formatted_date := m.L.FormatDate(month, day, entry)
					results[formatted_date]++
				}
			}
		}
	}
	scavengeddata := make(map[string][]map[string]int) // Well little complex but its works
	for date, _ := range results {
		dd, entry := m.L.SplitDate(date)
		month, day := m.L.SplitDateV2(dd)
		if stringToInt(month) == 1 {
			full_dates := m.GetJDates(month, day, entry)
			daddadido := make(map[string]int)
			for _, dd3 := range full_dates {
				// dd_1 := fmt.Sprintf("%s_%s", NumToDate(stringToInt(month), year, stringToInt(day)), entry)
				for _, dd4 := range dd3 {
					dd4_date, dd4_entry := m.L.SplitDate(dd4)
					_, dd4_month, dd4_day := m.L.dateToNum(dd4_date)
					formatted_luck := m.L.FormatDate(dd4_month, dd4_day, dd4_entry)
					daddadido[formatted_luck]++
				}
			}
			scavengeddata[date] = append(scavengeddata[date], daddadido)
		}
	}
	return scavengeddata
}

// Filter Numbers Return Only Used and Unused Based on Year.
// Analyze based on date returns a list of repeated  numbers at certain date based on other dates
func (m *DB) MarkedDates2(month, day int, entry string) {
	// var numbers []uint64
	MostUsed, UnusedNums := m.mDSortie(month, day, entry)
	dates := m.RetrieveDates(month, day, entry) // eg 2021/1/12-M [dates ...] (2021:[2022/....])
	// m.L.Logger.Log.Info(dates)
	// m.L.Logger.Log.Info("")
	jackPotResults := make(map[int]map[uint64][]int) // general numbers repeated across all matches.
	entry_fullresults := []map[string][]map[int][]uint64{}
	for y := range dates {
		sameYearNums := make(map[int][]uint64)
		for _, item := range dates[y] {
			for date, num := range item {
				nums, _ := ConvertToInt(num)   // Jackpot results numbers
				d, e := m.L.SplitDate(date)    // Date splited
				yy, mm, dd := m.L.dateToNum(d) // year, month, day
				for _, nu := range nums {      // Looping through numbers
					sameYearNums[yy] = append(sameYearNums[yy], nu) // adding to at temporal map
				}
				entry_results := m.mDFutureDates(yy, mm, dd, y, month, day, e, entry, nums)
				entry_fullresults = append(entry_fullresults, entry_results)
			}
		}
		jackPotResults[y] = m.mDSameYear(sameYearNums)
	}
	NumbersUnl := make(map[uint64][]uint64)
	for n, _ := range MostUsed {
		closest := closestNumbers(UnusedNums,n)
		if len(closest) > 0 {
			NumbersUnl[n] = closest
		}
	}
	// m.L.Logger.Log.Info(MostUsed)
	// m.L.Logger.Log.Info(UnusedNums)
	// for y,m := range jackPotResults {
	// 	if len(m) > 1 {
	// 		fmt.Println(y,m)
	// 	}
	// }
	//
	m.L.Logger.Log.Info("Most Used ", MostUsed)
	// m.L.Logger.Log.Info(UnusedNums)
	m.L.Logger.Log.Info("Closest Number", NumbersUnl)
	m.L.Logger.Log.Info("JackPot Results ", jackPotResults)
	m.L.Logger.Log.Info("Full Entry ", entry_fullresults)
}


// Find number that are either bigger or smaller by 1
func closestNumbers(slice []uint64, A uint64) []uint64 {
	closest := []uint64{}
	minDiff := A-1
	maxDiff := A+1
	if slices.Contains(slice, minDiff) {
		closest = append(closest, minDiff)
	}
	if slices.Contains(slice, maxDiff) {
		closest = append(closest, maxDiff)
	}
	return closest
}
// Continuation to MarkedDates2 reduce function size.
func (m DB) mDSortie(month, day int, entry string) (map[uint64][]int, []uint64) {
	var n uint64
	MostUsed := make(map[uint64][]int)
	var UnusedNums []uint64
	data := m.getDNumbers(month, day, entry)
	for y := range data {
		for n = 1; n < 86; n++ {
			if slices.Contains(data[y], n) {
				MostUsed[n] = append(MostUsed[n], y)
			}
		}
	}
	for n = 1; n < 86; n++ {
		if len(MostUsed[n]) == 0 {
			UnusedNums = append(UnusedNums, n)
		}
	}
	return MostUsed, UnusedNums
}
// Continuation to MkDate2 - find the most used numbers on date
func (m *DB) mDSameYear(sameYear map[int][]uint64) map[uint64][]int {
	var n uint64
	UsedInDate := make(map[uint64][]int)
	for y, nums := range sameYear {
		for n = 1; n < 86; n++ {
			if slices.Contains(nums, n) {
				UsedInDate[n] = append(UsedInDate[n], y)
			}
		}
	}
	return UsedInDate
}

// Continuation to MkDate2 - Find common numbers between date and its future dates
// return what matched number has been repeated on other years.
func (m *DB) mDFutureDates(year, month, day, OriginalYear, cmonth, cday int, entry, centry string, ornums []uint64) map[string][]map[int][]uint64 {
	// currentDateNums, _ := m.NumGetDayTime(year, month, day, entry)
	// yyear,cmonth,cday,entry - dates used to find common in.
	// year, month, day, enrry - date that matched yyear/cmonth/cday/centry
	years := []int{2020, 2021, 2022, 2023}
	results := make(map[string][]map[int][]uint64)
	date := NumToDate2(year, month, day, entry)
	results[date] = append(results[date], map[int][]uint64{OriginalYear: ornums})
	ListA, _ := m.NumGetDayTime(year, day, month, entry) // List oF numbers of matched secondary year
	for _, y := range years {
		if !ContainsKey(results[date], y) {
			data := make(map[int][]uint64)
			ListB, _ := m.NumGetDayTime(y, cday, cmonth, centry) // List oF numbers of selective other year
			commonNums := FilterDups(ListA, ListB)
			data[y] = commonNums
			// fmt.Println(y, year, ListB, commonNums)
			results[date] = append(results[date], data)
		}
	}
	return results
}

// Func to check

// Function to analyze numbers counting matched 5 or 6 7-8 numbers
// This way it show results more than it needs.
