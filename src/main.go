package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reinoj/corsi_calc/src/hockeydb"
)

func main() {
	fmt.Println("Starting...")

	//---------------FLAGS---------------
	// boolean flag for whether to create the database and tables
	var initialSetupFlag bool
	// string flag for the mysql user name
	var mysqlUserFlag string
	// string flag for the mysql password
	var mysqlPasswordFlag string

	flag.BoolVar(&initialSetupFlag, "initialSetup", false, "create the database and base tables.")
	flag.StringVar(&mysqlUserFlag, "mysqlUser", "root", "user name for mysql")
	flag.StringVar(&mysqlPasswordFlag, "mysqlPassword", "root", "mysql for mysql user")

	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()
	//---------------FLAGS---------------

	// Assigns the user name and password to the beginning of the string
	mysqlSignIn := fmt.Sprintf("%s:%s", mysqlUserFlag, mysqlPasswordFlag)

	// check initialSetupFlag
	if initialSetupFlag {
		//---------------CREATE INITIAL DATABASE---------------
		fmt.Println("Opening initial database...")
		db, err := sql.Open("mysql", mysqlSignIn+"@tcp(127.0.0.1:3306)/")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Initial database Opened.")
		// Creates the database for all the tables
		hockeydb.CreateDb(db)

		db.Close()
		//---------------CREATE INITIAL DATABASE---------------

		//---------------CREATE TABLES---------------
		fmt.Println("Opening Hockey database...")
		hdb, err := sql.Open("mysql", mysqlSignIn+"@tcp(127.0.0.1:3306)/Hockey")
		if err != nil {
			log.Fatal(err)
		}
		defer hdb.Close()
		fmt.Println("Hockey database Opened.")

		// Creates the tables for the database
		hockeydb.CreateTables(hdb)
		// Populate the Teams table
		hockeydb.GetTeams(hdb)
		// Populate the Schedule table
		hockeydb.GetSchedule(hdb)
		//---------------CREATE TABLES---------------
	}

	fmt.Println("Hockey database closed.")
	fmt.Println("Complete.")
}
