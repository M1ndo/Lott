package utils

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DB *sql.DB
	L *Utils
}

var Times = map[int]string{1:"morning", 2:"noon", 3:"night"}
var query string

func (m *DB) InitializeDB() chan error {
	chErrors := make(chan error)
	go func() {
		db, err := sql.Open("sqlite3", "luck.db")
		if err != nil {
			chErrors <- fmt.Errorf("failed to open database: %w", err)
			close(chErrors)
			return
		}
		m.DB = db
		if err := m.CreateTable(); err != nil {
			chErrors <- fmt.Errorf("failed to create table: %w", err)
		}
		if err := m.CJackpot(); err != nil {
			chErrors <- fmt.Errorf("failed to create table: %w", err)
		}
		if err := m.CreateTdaily(); err != nil {
			chErrors <- fmt.Errorf("failed to create Tdaily: %w", err)
		}
		close(chErrors)
	}()
	return chErrors
}

func (m *DB) CreateTdaily() error {
	query = `
		CREATE TABLE IF NOT EXISTS analyzedD (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER,
      month INTERGER,
      day INTEGER,
      sday TEXT,
			mostused TEXT
		);
  `
	_, err := m.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *DB) CreateTmonth() error {
	query = `
		CREATE TABLE IF NOT EXISTS analyzedD (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER,
			mostused TEXT
		);
  `
	return nil
}

func (m *DB) CreateTweeky() error {
	query = `
		CREATE TABLE IF NOT EXISTS analyzedD (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER,
			mostused TEXT
		);
  `
	return nil
}

func (m *DB) CreateTable() error {
	query = `
		CREATE TABLE IF NOT EXISTS data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER,
			month INTEGER,
			day INTEGER,
      sday TEXT,
			entry_time TEXT,
			numbers INTEGER,
			mostused TEXT
		);
	`
	_, err := m.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *DB) CJackpot() error {
	query = `
		CREATE TABLE IF NOT EXISTS jackpot (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER,
			month INTEGER,
			day INTEGER,
			entry_time TEXT,
			myear INTEGER,
			mmonth INTEGER,
			mday INTEGER,
			mentry_time TEXT,
			numbers INTEGER
		);
	`
	_, err := m.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}


func (m *DB) getNumbers(query string) ([]uint64, error) {
	var numbers []uint64
	m.L.Logger.Log.Debugf("Executing query getNumbers %s", query)
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var nums string
		err := rows.Scan(&nums)
		if err != nil {
			return nil, err
		}
    numrs, _ := ConvertToInt(nums)
		for _, n := range numrs {
			numbers = append(numbers, n)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return numbers, nil
}

// Get all numbers of a year
func (m *DB) NumGetYear(year int) ([]uint64, error) {
	query = fmt.Sprintf("SELECT numbers from data where year = %d;", year)
	return m.getNumbers(query)
}

// Get all numbers of a month (with entry_time) from a particular year
func (m *DB) NumGetMonth(year, month int, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d;", year, month)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and entry_time = '%s';", year, month, entry)
	}
	return m.getNumbers(query)
}

// Get numbers of a particular date and with an entry.
func (m *DB) NumGetDay(year, day int, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and day = %d", year, day)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and day = %d and entry_time = '%s';", year, day, entry)
	}
	return m.getNumbers(query)
}


// Get numbers of a particular day and a month and with an entry.
func (m *DB) NumGetDayTime(year, day, month int, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and day = %d and month = %d;", year, day, month)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and day = %d and month = %d and entry_time = '%s';", year, day, month, entry)
	}
	return m.getNumbers(query)
}

func (m *DB) NumGetSday(year, month int, sday, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and sday = '%s'", year, month, sday)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and sday = '%s' and entry_time = '%s';", year, month, sday, entry)
	}
	return m.getNumbers(query)
}

// Import data
func (m *DB) ImportData(d map[string][]uint64, year, month int) error {
	m.L.Logger.Log.Infof("Importing And Analyzing Year %d, Month %d", year, month)
	for date, numbers := range d {
		if day, time, _ := GetDate(date); day != 0 {
			sday := NumToDate(month, year, day)
			sday, _ = DateToDay(sday)
			mostUsed := m.L.FindDupNums(numbers, false)
			query := fmt.Sprintf("INSERT into data (year, month, day, sday, entry_time, numbers, mostused) values (?, ?, ?, ?, ?, ?, ?)")
			_, err := m.DB.Exec(query, year, month, day, sday, Times[time], ConvertToStr(numbers), mostUsed)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Import jackpot data
func (m *DB) ImportJack(year, month, day, myear, mmonth, mday int, entry, mentry string, numbers []uint64) error {
	m.L.Logger.Log.Debugf("Importing And Analyzed Year %d, Month %d", year, month)
	query := fmt.Sprintf("INSERT into jackpot (year, month, day, myear, mmonth, mday, entry_time, mentry_time, numbers) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := m.DB.Exec(query, year, month, day, myear, mmonth, mday, entry, mentry, ConvertToStr(numbers))
	if err != nil {
		return err
	}
	return nil
}

// Insert into Other Ts
func (m *DB) InsertTo(table string, d []DataHold) error {
	for _, entry := range d {
		query := fmt.Sprintf("INSERT INTO %s (year, month, day, sday, mostused) values (?, ?, ?, ?, ?)", table)
		_, err := m.DB.Exec(query, entry.Year, entry.Month, entry.Day, entry.DayT ,entry.Nums)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get Most Used Of Sdaily
// GetMostUsed retrieves the most used value based on the given parameters.
func (m *DB) GetMostUsed(table, dayEntry string, year, day, month int) (string, error) {
	if year == 0 {
		return "", fmt.Errorf("Year cannot be empty")
	}
	if dayEntry == "" {
		if day != 0 {
			if month != 0 {
				query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d AND month = %d AND day = %d;", table, year, month, day)
			} else {
				query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d AND day = %d;", table, year, day)
			}
		} else if month != 0 {
			query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d AND month = %d;", table, year, month)
		} else if year != 0 {
			query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d ", table, year)
		}
	} else {
		if month != 0 {
			query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d AND month = %d AND sday = '%s';", table, year, month, dayEntry)
		} else {
			query = fmt.Sprintf("SELECT mostused FROM %s WHERE year = %d AND sday = '%s';", table, year, dayEntry)
		}
	}

	return m.getmostUsed(query)
}


func (m *DB) getmostUsed(query string) (string, error) {
	var mostUsed string
	m.L.Logger.Log.Debugf("Executing %s", query)
	rows, err := m.DB.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&mostUsed)
		if err != nil {
			return "", err
		}
	}

	if err = rows.Err(); err != nil {
		return "", err
	}
	return mostUsed, nil
}

// Get total days in a month
func (m *DB) Total_days(month int) int {
	query := fmt.Sprintf("SELECT COUNT(distinct day) as days_in_month from data where month = %d;", month)
	days_count, err := m.getmostUsed(query)
	m.L.HandleError(err)
	total_days, err := strconv.Atoi(days_count)
	m.L.HandleError(err)
	return total_days
}

// Get JackPot Data Of a Year
func (m *DB) AnalyzedJackPot(year int) map[string]int64 {
	query := fmt.Sprintf("select mmonth, mday, mentry_time from jackpot where year = %d;", year)
	rows, err := m.DB.Query(query)
	m.L.HandleError(err)
	defer rows.Close()
	date := make(map[string]int64)
	count := int64(0)
	for rows.Next() {
		var day int64
		var month int64
		var entry string
		err := rows.Scan(&month, &day, &entry)
		m.L.HandleError(err)
		count++
		monthe := fmt.Sprintf("%d_%s_%d", day, entry, count)
		date[monthe] = month
	}
	m.L.HandleError(rows.Err())
	return date
}

// Get original matched date
func (m *DB) OriginalNumbers(year, month, day int, entry_time string) map[string]int64 {
	query := fmt.Sprintf("select month, day, entry_time from jackpot where myear = %d and mmonth = %d and mday = %d and mentry_time = '%s';", year, month, day, entry_time)
	rows, err := m.DB.Query(query)
	m.L.HandleError(err)
	defer rows.Close()
	date := make(map[string]int64)
	for rows.Next() {
		var day int64
		var month int64
		var entry string
		err := rows.Scan(&month, &day, &entry)
		m.L.HandleError(err)
		monthe := fmt.Sprintf("%d_%s", day, entry)
		date[monthe] = month
	}
	m.L.HandleError(rows.Err())
	return date
}

func (m *DB) getJackPotNums(year, month, day int, entry string) string {
	query := fmt.Sprintf("SELECT numbers from jackpot where year = %d and mmonth = %d and mday = %d and mentry_time = '%s';", year, month, day, entry)
	numbers, err := m.getmostUsed(query)
	m.L.HandleError(err)
	return numbers
}

func (m *DB) JackPotNumbers(year int) map[int][][]uint64 {
	query := fmt.Sprintf("SELECT numbers from jackpot where year = %d", year)
	data := make(map[int][][]uint64)
	rows, err := m.DB.Query(query)
	m.L.HandleError(err)
	defer rows.Close()
	for rows.Next() {
		var nums string
		err := rows.Scan(&nums)
		m.L.HandleError(err)
    numrs, _ := ConvertToInt(nums)
		data[year] = append(data[year], numrs)
	}
	m.L.HandleError(rows.Err())
	return data
}
