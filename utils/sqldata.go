package utils

import (
	"database/sql"
	"fmt"

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

func (m *DB) getNumbers(query string) ([]uint64, error) {
	var numbers []uint64
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

// Get all numbers of a month (with entry_time) from a particular year
func (m *DB) NumGetMonth(year, month int, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d;", year, month)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and entry_time = '%s';", year, month, entry)
	}
	return m.getNumbers(query)
}

// Get numbers of a particular date and month or an entry.
func (m *DB) NumGetDay(year, month, day int, entry string) ([]uint64, error) {
	if entry == "" {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and day = %d", year, month, day)
	} else {
		query = fmt.Sprintf("SELECT numbers from data where year = %d and month = %d and day = %d and entry_time = '%s';", year, month, day, entry)
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

// Insert into Other Ts
func (m *DB) InsertTo(table string, d []DataHold) error {
	for _, entry := range d {
		query := fmt.Sprintf("INSERT INTO %s (year, month, sday, mostused) values (?, ?, ?, ?)", table)
		_, err := m.DB.Exec(query, entry.Year, entry.Month, entry.DayT ,entry.Nums)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get Most Used Of Sdaily
func (m *DB) GetMostUsed(table, day_entry string, year, month int) (string, error) {
	var mostused string
	if day_entry == "" {
		query = fmt.Sprintf("SELECT mostused from %s where year = %d and month = %d;", table, year, month)
	} else {
		query = fmt.Sprintf("SELECT mostused from %s where year = %d and month = %d and sday = '%s';", table, year, month, day_entry)
	}
	rows, err := m.DB.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&mostused)
		if err != nil {
			return "", err
		}
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return mostused, nil
}
