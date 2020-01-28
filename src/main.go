package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reinoj/nhl_stat_calc/src/hockeydb"
	"github.com/reinoj/nhl_stat_calc/src/statcalc"
)

func main() {
	fmt.Println("Starting...")

	//---------------FLAGS---------------
	// boolean flag for whether to create the database
	var initialSetupFlag bool
	// boolean flag for whether to create the tables
	var createTablesFlag bool
	// boolean flag for whether to update the Schedule table
	var updateTablesFlag bool
	// boolean flag for whether to calculate corsi
	var calculateCorsiFlag bool
	// boolean flag for whether to output corsi stats
	var outputCorsiFlag bool
	// string flag for the mysql user name
	var mysqlUserFlag string
	// string flag for the mysql password
	var mysqlPasswordFlag string

	flag.BoolVar(&initialSetupFlag, "initialSetup", false, "create the database.")
	flag.BoolVar(&initialSetupFlag, "createTables", false, "create the tables.")
	flag.BoolVar(&updateTablesFlag, "updateTables", false, "update the Schedule and ShotInfo tables.")
	flag.BoolVar(&calculateCorsiFlag, "calculateCorsi", false, "calculate corsi from ShotInfo table.")
	flag.BoolVar(&outputCorsiFlag, "outputCorsi", false, "output corsi info for the season.")
	flag.StringVar(&mysqlUserFlag, "mysqlUser", "root", "user name for mysql.")
	flag.StringVar(&mysqlPasswordFlag, "mysqlPassword", "root", "password for mysql user.")

	// must be called after all flags are defined and before flags are accessed by the program
	flag.Parse()
	//---------------------------------------------

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
		//---------------------------------------------
	}

	fmt.Println("Opening Hockey database...")
	hdb, err := sql.Open("mysql", mysqlSignIn+"@tcp(127.0.0.1:3306)/Hockey")
	if err != nil {
		log.Fatal(err)
	}
	defer hdb.Close()
	fmt.Println("Hockey database Opened.")

	if createTablesFlag {
		//---------------CREATE TABLES---------------
		// Creates the tables for the database
		hockeydb.CreateTables(hdb)
		// Populate the Teams table
		hockeydb.GetTeams(hdb)
		var fullSchedule hockeydb.Schedule
		hockeydb.GetSchedule(hdb, &fullSchedule)
		// Populate the Schedule table
		hockeydb.PopulateScheduleTable(hdb, &fullSchedule)
		//---------------------------------------------
	}

	if updateTablesFlag {
		//---------------UPDATE RESULTS---------------
		hockeydb.UpdateTables(hdb)
		//---------------------------------------------
	}

	if calculateCorsiFlag {
		statcalc.CorsiCalc(hdb)
	}

	if outputCorsiFlag {
		statcalc.GetCorsiWins(hdb)
	}

	fmt.Println("Hockey database closed.")
	fmt.Println("Complete.")
}
