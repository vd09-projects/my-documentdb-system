package parsers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/utils"
)

type JSONParser struct {
	database db.Database
}

func (p *JSONParser) Parse(ctx context.Context, file io.Reader, w http.ResponseWriter, userID string, recordType string) {
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	var jsonRecords []map[string]interface{}
	var insertCount, quarantineCount int

	if err := json.Unmarshal(dataBytes, &jsonRecords); err == nil {
		for _, rec := range jsonRecords {
			if !utils.IsValidRecord(rec) {
				p.database.InsertToQuarantine(ctx, rec, userID, recordType, "Failed validation")
				quarantineCount++
				continue
			}
			p.database.InsertToValid(ctx, rec, userID, recordType)
			insertCount++
		}
	} else {
		var singleRec map[string]interface{}
		if err := json.Unmarshal(dataBytes, &singleRec); err != nil {
			p.database.InsertToQuarantine(ctx, string(dataBytes), userID, recordType, "Invalid JSON structure")
			quarantineCount++
		} else {
			if !utils.IsValidRecord(singleRec) {
				p.database.InsertToQuarantine(ctx, singleRec, userID, recordType, "Failed validation")
				quarantineCount++
			} else {
				p.database.InsertToValid(ctx, singleRec, userID, recordType)
				insertCount++
			}
		}
	}

	resultMsg := fmt.Sprintf("Successfully inserted %d records; quarantined %d records.", insertCount, quarantineCount)
	w.Write([]byte(resultMsg))
}
