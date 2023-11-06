package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Date    string
	Numbers []uint64
}


func saveNumbersToFile(numbers []uint64, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, number := range numbers {
		strNumber := strconv.FormatUint(number, 10) + "\n"
		_, err := file.WriteString(strNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeJunk(str string) string {
	re := regexp.MustCompile(`(\d+)\s+de\s+\w+\.\w+\s+(\d+)`)
	matches := re.FindStringSubmatch(str)

	if len(matches) >= 3 {
		dateTime := matches[1] + " " + matches[2]
		return dateTime
	}
	return ""
}

func parseHTMLFile(filePath string) []Result {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
if err != nil {
			panic(err)
		}

		var results []Result

		doc.Find("div.SUP").Each(func(i int, supDiv *goquery.Selection) {
		pElements := supDiv.Find("p")
		if pElements.Length() >= 2 {
			dateElement := pElements.Eq(1)
			sortedElement := pElements.Eq(2)
			sorted := sortedElement.Text()
			date := dateElement.Text()
			dateAndWhen := removeJunk(date + sorted)
			ulElement := supDiv.Find("ul[aria-label='Números del Super 11']")
			if ulElement.Length() > 0 {
				numbers := []uint64{}
				ulElement.Find("li").Each(func(i int, li *goquery.Selection) {
					number, _ := strconv.ParseUint(li.Text(), 10, 64)
					numbers = append(numbers, number)
				})
				result := Result{
					Date:    dateAndWhen,
					Numbers: numbers,
				}
				results = append(results, result)
				}
		}
	})

	return results
}

func main() {
	filePath := flag.String("file", "", "Path to the HTML file")
	// fileOut := flag.String("outfile", "", "Output file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide the path to the HTML file using the -file flag")
		return
	}

	results := parseHTMLFile(*filePath)
	// numbers := []uint64{}
	for _, result := range results {
		// numbers = append(numbers, result.Numbers...)
		fmt.Printf("%s: %v\n", result.Date, result.Numbers)
	}
	// saveNumbersToFile(numbers, *fileOut)
}
