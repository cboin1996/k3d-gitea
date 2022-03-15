package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("CREATE ROLE gitea WITH LOGIN PASSWORD 'gitea'", age)
	fmt.Printf("rows: \n%v\nerr: %v", rows, err)
	// rows, err := db.Query("CREATE ROLE gitea WITH LOGIN PASSWORD 'gitea'", age)
}