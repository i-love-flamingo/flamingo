package infrastructure

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"

	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// CsvContent definition / dto
	CsvContent struct {
		HTTPStatusCode int
		OriginalPath   string
		RedirectTarget string
	}

	// RedirectData loads redirect data from csv
	RedirectData struct {
		filePath   string
		logger     flamingo.Logger
		csvContent []CsvContent
	}
)

// NewRedirectData provides a new dataset with autoload from csv
func NewRedirectData(
	cfg *struct {
		RedirectsCsv string `inject:"config:redirects.csv"`
	},
	logger flamingo.Logger,
) *RedirectData {
	rd := &RedirectData{
		filePath: cfg.RedirectsCsv,
		logger:   logger,
	}
	rd.csvContent, _ = rd.readCSV()

	return rd
}

// Get returns the Redirect data
func (rd *RedirectData) Get() []CsvContent {
	return rd.csvContent
}

func (rd *RedirectData) readCSV() ([]CsvContent, error) {
	f, err := os.Open(rd.filePath)
	if err != nil {
		rd.logger.Error("Error - CsvAssetRedirectAdapter ", err)
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
			rd.logger.Error("Redirect load - Error - Invalid Row (wrong amount) %v in Row: %v", record, rowCount)
			continue
		}

		if isAlpha(record[0]) {
			// invalid row
			rd.logger.Error("Redirect load - Error - Invalid Row (no int status) %v in Row: %v", record, rowCount)
			continue
		}

		statusCode, err := strconv.Atoi(record[0])

		if err != nil {
			rd.logger.Error("Error - CsvAssetRedirectAdapter %v", err)
			continue
		}

		CsvContents = append(CsvContents, CsvContent{
			HTTPStatusCode: statusCode,
			OriginalPath:   record[1],
			RedirectTarget: record[2],
		})
	}
	return CsvContents, nil
}
