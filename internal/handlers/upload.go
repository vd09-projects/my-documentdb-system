package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/parsers"
	"github.com/vd09-projects/my-documentdb-system/internal/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("datafile")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()
	database := db.NewMongoDB(db.MongoClient, db.DatabaseName)

	fileName := fileHeader.Filename
	parser := parsers.GetParser(fileName, database)
	if parser == nil {
		http.Error(w, "Unsupported file type. Upload CSV or JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve claims from context
	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	if !ok {
		http.Error(w, "Missing user info in context", http.StatusUnauthorized)
		return
	}

	parser.Parse(ctx, file, w, claims.UserID)
}

func GetAllDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database := db.NewMongoDB(db.MongoClient, db.DatabaseName)
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
