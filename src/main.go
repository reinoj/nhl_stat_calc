package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reinoj/go_corsi_calc/src/hockeydb"
	//"github.com/reinoj/go_corsi_calc/src/hockeydb"
)

func main() {
	fmt.Println("Starting...")

	// boolean flag for whether to create the database and tables
	var createDatabaseFlag bool
	//
	var createTablesFlag bool
	// assigns createDatabase if it was given in the argument, otherwise it defaults to false
	flag.BoolVar(&createDatabaseFlag, "createDatabase", false, "bool value")
	// assigns creatTables if it was given in the argument, otherwise it defaults to false
	flag.BoolVar(&createTablesFlag, "createTables", false, "bool value")
	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()

	// check setupFlag
	if createDatabaseFlag {
		fmt.Println("Opening initial database...")
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Initial database Opened.")
		// Creates the database for all the tables
		hockeydb.CreateDb(db)
		db.Close()
	}
	fmt.Println("Opening Hockey database...")
	hdb, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Hockey")
	if err != nil {
		log.Fatal(err)
	}
	defer hdb.Close()
	fmt.Println("Hockey database Opened.")
	// check createTablesFlag
	if createTablesFlag {
		// Creates the tables for the database
		hockeydb.CreateTables(hdb)
		// Populate the Teams table
		hockeydb.GetTeams(hdb)
	}
	//hockeydb.GetTeams(hdb)
	fmt.Println("Hockey database closed.")
	fmt.Println("Complete.")
}
