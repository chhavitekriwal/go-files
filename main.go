package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	env_err := godotenv.Load()
  	if env_err != nil {
    	log.Fatal("Error loading .env file")
  	}

	db,db_err:= GetDB()

	if db_err != nil {
		log.Fatalf("Error connecting to the database: %v",db_err)
	}

	mig_err := MigrateModels(db)
	if mig_err != nil {
		log.Fatalf("Database migration error: %v",mig_err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Up and Running")
	})
	http.HandleFunc("/register",RegisterHandler(db))
	http.HandleFunc("/login",LoginHandler(db))
	
	http.Handle("/upload",AuthenticationMiddleware(http.HandlerFunc(UploadHandler(db))))
	http.Handle("/download",AuthenticationMiddleware(http.HandlerFunc(DownloadHandler(db))))
	http.Handle("/files",AuthenticationMiddleware(http.HandlerFunc(ListFilesHandler)))
	http.Handle("/delete",AuthenticationMiddleware(http.HandlerFunc(DeleteHandler(db))))
	http.Handle("/transactions",AuthenticationMiddleware(http.HandlerFunc(GetTransactionsHandler(db))))

	log.Println("Starting server on :8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}