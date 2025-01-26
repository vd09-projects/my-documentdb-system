package main

import (
	"log"
	"net/http"

	"github.com/vd09-projects/my-documentdb-system/internal/db"
	"github.com/vd09-projects/my-documentdb-system/internal/handlers"
)

func main() {
	db.ConnectMongoDB("mongodb://db:27017")

	// Serve static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Public routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// Protected routes
	http.Handle("/upload", handlers.AuthMiddleware(http.HandlerFunc(handlers.UploadHandler)))
	// http.Handle("/data", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetAllDataHandler)))
	http.Handle("/userData", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetUserDataHandler)))

	http.Handle("/listRecordTypes", handlers.AuthMiddleware(http.HandlerFunc(handlers.ListRecordTypesHandler)))
	http.Handle("/listFields", handlers.AuthMiddleware(http.HandlerFunc(handlers.ListFieldsHandler)))
	http.Handle("/aggregate", handlers.AuthMiddleware(http.HandlerFunc(handlers.AggregateHandler)))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
