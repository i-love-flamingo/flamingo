package infrastructure

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
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

func GetRedirectData() ([]CsvContent) {
	return csvContent
}

func readCSV() ([]CsvContent, error) {
	f, err := os.Open("/resources/redirects.csv")
	if err != nil {
		log.Printf("Error - CsvAssetRedirectAdapter %v", err)
		return nil, err
	}

	var CsvContents []CsvContent
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
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
