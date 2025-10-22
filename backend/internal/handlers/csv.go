package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
)

type CsvRow map[string]string

func GetCsvData(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data/dummy.csv")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to open csv"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read csv header"})
		return
	}

	var rows []CsvRow
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		row := CsvRow{}
		for i, h := range headers {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
