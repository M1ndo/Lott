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
func FindMissingNumbers(list1, list2 []uint64) {

}

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
func (util *Utils) FindDupNums(List []uint64) string {
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
	// Disabling Filtration of 1 used digits or shall we ?
	// FilteredNums := make(map[uint64]uint64)
	// for num, count := range DuplicatedNums {
	// 	if count > 1 {
	// 		FilteredNums[num] = count
	// 	}
	// }
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


func quickConv(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func NumToDate(month, year, day int) string {
	return  quickConv(year) + "/" +quickConv(month) + "/" + quickConv(day)
}

func DateToDay(dat string) (string, error) {
	date, err := time.Parse("2006/1/2", dat)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", nil
	}

	// Get the day of the week
	dayOfWeek := date.Weekday()
	return dayOfWeek.String(), nil
}

func (util *Utils) toJson(pairs []Pair) (string) {
	util.Logger.Log.Infof("Encoding Pair %v to Json", pairs)
	stringPair, err := json.Marshal(pairs)
	if err != nil {
		util.Logger.Log.Error(err)
		return ""
	}
	return string(stringPair)
}

func FromJson(stringPair []byte) ([]Pair, error) {
	var pairs []Pair
	err := json.Unmarshal([]byte(stringPair), &pairs)
	if err != nil {

		return nil, err
	}
	return pairs, nil
}
