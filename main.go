package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

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

	log.Println("Starting server on :8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}