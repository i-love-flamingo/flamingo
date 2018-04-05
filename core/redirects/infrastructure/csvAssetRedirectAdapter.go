package infrastructure

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"

	"log"
)

type (
	CsvContent struct {
		HttpStatusCode int
		OriginalPath   string
		RedirectTarget string
	}
)

var csvContent []CsvContent

func init() {
	csvContent, _ = readCSV()
}

func GetRedirectData() []CsvContent {
	return csvContent
}

func readCSV() ([]CsvContent, error) {
	f, err := os.Open("resources/redirects.csv")
	if err != nil {
		log.Printf("Error - CsvAssetRedirectAdapter %v", err)
		return nil, err
	}

	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

	var CsvContents []CsvContent
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	rowCount := 0
	isFirstRow := true
	for {
		rowCount++
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if isFirstRow {
			isFirstRow = false
			continue
		}

		if len(record) != 3 {
			// invalid row
			log.Printf("Redirect load - Error - Invalid Row (wrong amount) %v in Row: %v", record, rowCount)
			continue
		}

		if isAlpha(record[0]) {
			// invalid row
			log.Printf("Redirect load - Error - Invalid Row (no int status) %v in Row: %v", record, rowCount)
			continue
		}

		statusCode, err := strconv.Atoi(record[0])

		if err != nil {
			log.Printf("Error - CsvAssetRedirectAdapter %v", err)
			continue
		}

		CsvContents = append(CsvContents, CsvContent{
			HttpStatusCode: statusCode,
			OriginalPath:   record[1],
			RedirectTarget: record[2],
		})
	}
	return CsvContents, nil
}
