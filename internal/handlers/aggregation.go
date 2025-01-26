package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// e.g. GET /listRecordTypes
func ListRecordTypesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	metadataDB := db.NewRecordDB(db.MongoClient, db.DatabaseName)
	recordTypes, err := metadataDB.GetRecordTypesForUser(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recordTypes)
}

// e.g. GET /listFields?recordType=sales
func ListFieldsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	recordType := r.URL.Query().Get("recordType")
	if recordType == "" {
		http.Error(w, "Missing recordType", http.StatusBadRequest)
		return
	}

	metadataDB := db.NewRecordDB(db.MongoClient, db.DatabaseName)
	fields, err := metadataDB.GetFieldsForUserAndType(ctx, userID, recordType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fields) // e.g. ["price", "quantity", "timestamp"]
}

// e.g. GET /aggregate?recordType=sales&field=price&op=sum
func AggregateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	recordType := r.URL.Query().Get("recordType")
	field := r.URL.Query().Get("field")
	op := r.URL.Query().Get("op")

	if recordType == "" || field == "" || op == "" {
		http.Error(w, "Missing recordType/field/op", http.StatusBadRequest)
		return
	}

	// Use the RecordDB interface for database logic
	metadataDB := db.NewRecordDB(db.MongoClient, db.DatabaseName)
	result, err := metadataDB.AggregateData(ctx, userID, recordType, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aggregateResult, err := aggregateResultsFromDB(result, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the computed result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"result": aggregateResult})
}

func aggregateResultsFromDB(data []float64, op string) (float64, error) {
	// Perform the specified aggregation operation
	var result float64 = 0
	switch op {
	case "sum":
		for _, v := range data {
			result += v
		}
	case "average":
		if len(data) > 0 {
			for _, v := range data {
				result += v
			}
			result /= float64(len(data))
		}
	case "min":
		if len(data) > 0 {
			result = data[0]
			for _, v := range data {
				if v < result {
					result = v
				}
			}
		}
	case "max":
		if len(data) > 0 {
			result = data[0]
			for _, v := range data {
				if v > result {
					result = v
				}
			}
		}
	default:
		return 0, fmt.Errorf("invalid operation: %s", op)
	}

	return math.Round(result*100) / 100, nil
}
