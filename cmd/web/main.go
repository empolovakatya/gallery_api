package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", "./pkg/data.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `photos` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`image` VARCHAR(255) NOT NULL, `preview` VARCHAR(255) NOT NULL);")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/photos", getPhotos).Methods("GET")
	r.HandleFunc("/photos/{id}", getPhoto).Methods("GET")
	r.HandleFunc("/photos", uploadPhoto).Methods("POST")
	r.HandleFunc("/photos/{id}", deletePhoto).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
