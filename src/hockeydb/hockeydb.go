package hockeydb

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateDb creates the hockey database
func CreateDb(hdb *sql.DB) {
	_, err := hdb.Exec("CREATE DATABASE hockey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Created Database")
}

// CreateTables creates the tables inside the database
func CreateTables(hdb *sql.DB) {
	_, err := hdb.Exec("CREATE TABLE teams")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("teams table created.")
}
