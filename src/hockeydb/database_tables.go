package hockeydb

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateDb creates the hockey database
func CreateDb(db *sql.DB) {
	fmt.Println("Creating Hockey databse...")
	// executes the create statement to make the database
	_, err := db.Exec("CREATE DATABASE Hockey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully created Hockey database")
}

// CreateTables creates the tables inside the database
func CreateTables(hdb *sql.DB) {
	fmt.Println("Creating tables...")

	/*
		CREATE TABLE Teams (
			ID int NOT NULL,
			TeamName VARCHAR(255) NOT NULL,
			FullName VARCHAR(255) NOT NULL,
			Abbreviation CHAR(3) NOT NULL,
			DivisionName VARCHAR(255) NOT NULL,
			ConferenceName VARCHAR(255) NOT NULL,
			PRIMARY KEY (ID)
		)
	*/
	// executes the create statement to make Teams table
	_, err := hdb.Exec("CREATE TABLE Teams (ID INT NOT NULL, TeamName VARCHAR(255) NOT NULL, FullName VARCHAR(255) NOT NULL, Abbreviation CHAR(3) NOT NULL, DivisionName VARCHAR(255) NOT NULL, ConferenceName VARCHAR(255) NOT NULL, PRIMARY KEY (ID));")
	if err != nil {
		// I've read that checking the output of the .Error() function is bad practice, but it works
		if err.Error() != "Error 1050: Table 'Teams' already exists" {
			log.Fatal(err)
		} else {
			fmt.Println("Teams table already exists.")
		}
		// If the error is that the table already exists, just ignore it, otherwise it will print the error to the screen
	}

	/*
		CREATE TABLE Schedule (
			GameNum INT NOT NULL,
			GameID CHAR(10) NOT NULL,
			Away INT NOT NULL,
			AwayResult VARCHAR(3),
			Home INT NOT NULL,
			HomeResult VARCHAR(3),
			PRIMARY KEY (GameNum)
		)
	*/
	_, err = hdb.Exec("CREATE TABLE Schedule (GameNum INT NOT NULL, GameID CHAR(10) NOT NULL, Away INT NOT NULL, AwayResult VARCHAR(3), Home INT NOT NULL, HomeResult VARCHAR(3), PRIMARY KEY (GameNum));")
	if err != nil {
		if err.Error() != "Error 1050: Table 'Schedule' already exists" {
			log.Fatal(err)
		} else {
			fmt.Println("Schedule table already exists.")
		}
	}
	fmt.Println("Tables created.")
}
