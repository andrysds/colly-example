package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Record struct {
	Data map[string]string
}

const headersEnvKey = "CSV_HEADERS"

func NewCSV(file io.Reader) ([]Record, error) {
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return []Record{}, err
	}

	if len(records) < 1 {
		return []Record{}, fmt.Errorf("csv file is empty")
	}

	if len(records) == 1 {
		return []Record{}, fmt.Errorf("csv file has only one record (headers)")
	}

	headersEnv := os.Getenv(headersEnvKey)
	headers := strings.Split(headersEnv, ",")

	if !reflect.DeepEqual(headers, records[0]) {
		return []Record{}, fmt.Errorf("csv file has different format")
	}

	var res []Record
	for _, r := range records[1:] {
		record := Record{Data: make(map[string]string)}
		for i, h := range headers {
			record.Data[h] = r[i]
		}
		res = append(res, record)
	}

	return res, nil
}
