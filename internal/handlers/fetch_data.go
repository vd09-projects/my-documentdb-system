package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/utils"
)

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Extract claims from context (userID from token)
	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Missing user info in context", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	// Optional date parsing
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var fromTime, toTime *time.Time
	var err error

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "Invalid 'from' date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		fromTime = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "Invalid 'to' date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		toTime = &t
	}

	// Query the DB
	database := db.NewRecordDB(db.MongoClient, db.DatabaseName)
	results, err := database.GetUserData(ctx, userID, fromTime, toTime)
	if err != nil {
		http.Error(w, "Failed to query data", http.StatusInternalServerError)
		return
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func GetAllDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database := db.NewRecordDB(db.MongoClient, db.DatabaseName)
	results, err := database.GetAllValidData(ctx)
	if err != nil {
		http.Error(w, "Failed to query data", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
