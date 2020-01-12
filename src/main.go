package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reinoj/go_corsi_calc/src/hockeydb"
)

func main() {
	fmt.Println("Starting...")

	// boolean flag for whether to create the database and tables
	var createDatabaseFlag bool
	//
	var createTablesFlag bool
	//
	var mysqlUserFlag string
	//
	var mysqlPasswordFlag string

	flag.BoolVar(&createDatabaseFlag, "createDatabase", false, "create the database")
	flag.BoolVar(&createTablesFlag, "createTables", false, "create the tables")
	flag.StringVar(&mysqlUserFlag, "mysqlUser", "root", "user name for mysql")
	flag.StringVar(&mysqlPasswordFlag, "mysqlPassword", "root", "mysql for mysql user")

	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()

	mysqlSignIn := fmt.Sprintf("%s:%s", mysqlUserFlag, mysqlPasswordFlag)
	// check setupFlag
	if createDatabaseFlag {
		fmt.Println("Opening initial database...")
		db, err := sql.Open("mysql", mysqlSignIn+"@tcp(127.0.0.1:3306)/")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Initial database Opened.")
		// Creates the database for all the tables
		hockeydb.CreateDb(db)
		db.Close()
	}
	fmt.Println("Opening Hockey database...")
	hdb, err := sql.Open("mysql", mysqlSignIn+"@tcp(127.0.0.1:3306)/Hockey")
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
	fmt.Println(mysqlSignIn)
	fmt.Println("Hockey database closed.")
	fmt.Println("Complete.")
}
