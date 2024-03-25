package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Utils struct {
	Logger Logger
}

// Pair represents a number and its count
type Pair struct {
	Number uint64
	Count  uint64
}

// Most Unused or Uncommon Numbers that are picked
// func FindMissingNumbers(list1, list2 []uint64) {

// }

func sortBased(numbers map[uint64]uint64) []Pair {
	pairs := make([]Pair, 0, len(numbers))
	for num, count := range numbers {
		pairs = append(pairs, Pair{Number: num, Count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
	return pairs
}

// Most Used Numbers in combination
func (util *Utils) FindDupNums(List []uint64, filtered bool) string {
	DuplicatedNums := make(map[uint64]uint64)
	for _, num := range List {
		_, exists := DuplicatedNums[num]
		if exists {
			DuplicatedNums[num]++
		} else {
			DuplicatedNums[num] = 1
		}
	}
	// fmt.Println(DuplicatedNums)
	if filtered {
		FilteredNums := make(map[uint64]uint64)
		for num, count := range DuplicatedNums {
			if count > 1 {
				FilteredNums[num] = count
			}
		}
		Sorted := sortBased(FilteredNums)
		Dups := util.toJson(Sorted)
		return Dups
	}
	Sorted := sortBased(DuplicatedNums)
	Dups := util.toJson(Sorted)
	return Dups
}

func ConvertToInt(str string) ([]uint64, error) {
	numbersStr := strings.Split(str, ",")
	numbers := make([]uint64, len(numbersStr))
	for i, numStr := range numbersStr {
		num, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return nil, err
		}
		numbers[i] = num
	}
	return numbers, nil
}

func stringToInt(num string) int {
	number, _ := strconv.Atoi(num)
	return number
}

func ConvertToStr(Nums []uint64) string {
	numbersStr := make([]string, len(Nums))
	for i, num := range Nums {
		numbersStr[i] = strconv.FormatUint(num, 10)
	}

	numbersJoined := strings.Join(numbersStr, ",")
	return numbersJoined
}

func (util *Utils) ReadNumbers(file string) []uint64 {
	numbers := make([]uint64, 0)
	Numbers, err := os.Open(file)
	if err != nil {
		util.Logger.Log.Error(err)
	}
	defer Numbers.Close()
	scanner := bufio.NewScanner(Numbers)
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			util.Logger.Log.Error(err)
		}
		numbers = append(numbers, number)
	}

	if err := scanner.Err(); err != nil {
		util.Logger.Log.Error(err)
	}
	return numbers
}

// Retrieve date from a string
func GetDate(str string) (int, int, error) {
	numbers := strings.Split(str, " ")
	if len(numbers) >= 2 {
		firstNumber, err := strconv.Atoi(numbers[0])
		if err != nil {
			return 0, 0, err
		}

		secondNumber, err := strconv.Atoi(numbers[1])
		if err != nil {
			return 0, 0, err
		}

		return firstNumber, secondNumber, nil
	}
	return 0, 0, nil
}

// Convert an interger to a string
func quickConv(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// Create a date string
func NumToDate(month, year, day int) string {
	return quickConv(year) + "/" + quickConv(month) + "/" + quickConv(day)
}

func NumToDate2(year, month, day int, entry string) string {
	return fmt.Sprintf("%d/%d/%d_%s", year, month, day, entry)
}


// Get the day of the week
func DateToDay(dat string) (string, error) {
	date, err := time.Parse("2006/1/2", dat)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", nil
	}
	dayOfWeek := date.Weekday()
	return dayOfWeek.String(), nil
}

// Convert a pair to json
func (util *Utils) toJson(pairs []Pair) string {
	util.Logger.Log.Debugf("Encoding Pair %v to Json", pairs)
	stringPair, err := json.Marshal(pairs)
	if err != nil {
		util.Logger.Log.Error(err)
		return ""
	}
	return string(stringPair)
}

// Convert from json to a pair
func FromJson(stringPair []byte) ([]Pair, error) {
	var pairs []Pair
	pairs = make([]Pair, 0)
	err := json.Unmarshal(stringPair, &pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

// Encode a myResult type to json
func (Util *Utils) myResultToJson(result myResult) {
	decoded := Util.jsonTomyResult()
	for k, r := range result {
		decoded[k] = r
	}
	encoded, err := json.Marshal(decoded)
	Util.HandleError(err)
	Util.SaveToFile(encoded)
}

// Decode a json byte to myResult
func (Util *Utils) jsonTomyResult() myResult {
	var decoded myResult
	result := Util.ReadfromFile()
	err := json.Unmarshal(result, &decoded)
	Util.HandleError(err)
	return decoded
}

// Save Json to a file
func (Util *Utils) SaveToFile(content []byte) {
	OutputFile, err := os.OpenFile("results/jackpot.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	Util.HandleError(err)
	defer OutputFile.Close()
	_, err = OutputFile.Write(content)
	Util.HandleError(err)
}

// Reads json from a file.
func (Util *Utils) ReadfromFile() []byte {
	ReadFile, err := os.ReadFile("results/jackpot.json")
	Util.HandleError(err)
	return ReadFile
}

// Split A string date
func (Util *Utils) SplitDate(c string) (string, string) {
	entry := strings.Split(c, "_")
	entry_date := entry[0]
	entry_time := entry[1]
	return entry_date, entry_time
}

// Split a string date v2
func (Util *Utils) SplitDateV2(c string) (string, string) {
	entry := strings.Split(c, "/")
	return entry[0], entry[1]
}

// Return a string date into numbers,
func (Util *Utils) dateToNum(d string) (int, int, int) {
	t, err := time.Parse("2006/1/2", d)
	Util.HandleError(err)
	year, month, day := t.Year(), int(t.Month()), t.Day()
	return year, month, day
}

// Format Date into a string
func (Util *Utils) FormatDate(month, day int, entry string) string {
	return fmt.Sprintf("%d/%d_%s", month, day, entry)
}

// Check if a list of maps contains a key
func ContainsKey(List []map[int][]uint64, key int) bool {
	for _, m := range List {
		if _, ok := m[key]; ok {
			return true
		}
	}
	return false
}

// Compare two lists and return duplicated Numbers
func FilterDups(ListA, ListB []uint64) []uint64 {
	common := []uint64{}
	count := map[uint64]bool{}
	for _, num := range ListA {
		count[num] = true
	}
	for _, num := range ListB {
		if _, found := count[num]; found {
			common = append(common, num)
		}
	}
	return common
}
