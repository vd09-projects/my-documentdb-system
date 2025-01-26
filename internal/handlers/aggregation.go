package handlers

import (
	"context"
	"encoding/json"
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
	result, err := metadataDB.AggregateData(ctx, userID, recordType, field, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the computed result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"result": result})
}
