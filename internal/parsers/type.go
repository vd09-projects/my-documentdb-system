package parsers

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
)

type Parser interface {
	Parse(ctx context.Context, file io.Reader, w http.ResponseWriter, userID string, recordType string)
}

func GetParser(fileName string, database db.Database) Parser {
	if strings.HasSuffix(fileName, ".csv") {
		return &CSVParser{
			database: database,
		}
	} else if strings.HasSuffix(fileName, ".json") {
		return &JSONParser{
			database: database,
		}
	}
	return nil
}
