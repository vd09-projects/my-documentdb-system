package parsers

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/utils"
)

type CSVParser struct {
	database db.Database
}

func (p *CSVParser) Parse(ctx context.Context, file io.Reader, w http.ResponseWriter, userID string) {
	fmt.Println("Parse in CSVParser")
	reader := csv.NewReader(bufio.NewReader(file))
	headers, err := reader.Read()
	if err != nil {
		http.Error(w, "Failed to read CSV headers", http.StatusBadRequest)
		return
	}

	var insertCount, quarantineCount int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.database.InsertToQuarantine(ctx, record, userID, "Error reading CSV")
			quarantineCount++
			continue
		}

		if len(record) != len(headers) {
			p.database.InsertToQuarantine(ctx, record, userID, "Mismatched header and record lengths")
			quarantineCount++
			continue
		}

		data := mapRowToKeyValue(headers, record)
		if !utils.IsValidRecord(data) {
			p.database.InsertToQuarantine(ctx, record, userID, "Failed validation")
			quarantineCount++
			continue
		}

		p.database.InsertToValid(ctx, data, userID)
		insertCount++
	}

	resultMsg := fmt.Sprintf("Successfully inserted %d records; quarantined %d records.", insertCount, quarantineCount)
	w.Write([]byte(resultMsg))
}

func mapRowToKeyValue(headers, record []string) map[string]string {
	data := make(map[string]string)
	for i, key := range headers {
		data[key] = record[i]
	}
	return data
}
