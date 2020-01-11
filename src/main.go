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

	fmt.Println("Opening Database...")
	hdb, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	defer hdb.Close()
	fmt.Println("Database Opened.")

	// boolean flag for whether to create the database and tables
	var setupFlag bool
	// assigns setupFlag if it was given in the argument, otherwise it defaults to false
	flag.BoolVar(&setupFlag, "setup", false, "bool value")
	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()
	//check setupFlag
	if setupFlag {
		fmt.Println("Creating database...")
		// Creates the database for all the tables
		hockeydb.CreateDb(hdb)
		fmt.Println("Database created.")
		fmt.Println("Creating tables...")
		// Creates the tables for the database
		hockeydb.CreateTables(hdb)
		fmt.Println("Tables Created.")
	}

	fmt.Println("Complete.")
}
