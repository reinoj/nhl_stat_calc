package hockeydb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// NumTeams is the number of teams, Seattle not in league yet
const NumTeams uint8 = 31

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
			Abbreviation VARCHAR(255) NOT NULL,
			DivisionName VARCHAR(255) NOT NULL,
			ConferenceName VARCHAR(255) NOT NULL,
			PRIMARY KEY (ID)
		)
	*/
	// executes the create statement to make Teams table
	_, err := hdb.Exec("CREATE TABLE Teams (ID INT NOT NULL, TeamName VARCHAR(255) NOT NULL, FullName VARCHAR(255) NOT NULL, Abbreviation VARCHAR(255) NOT NULL, DivisionName VARCHAR(255) NOT NULL, ConferenceName VARCHAR(255) NOT NULL, PRIMARY KEY (ID));")
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
			GameKey INT NOT NULL AUTO_INCREMENT,
			GameID CHAR(10) NOT NULL,
			Away CHAR(3) NOT NULL,
			Home CHAR(3) NOT NULL,
			PRIMARY KEY (GameKey)
		)
	*/
	_, err = hdb.Exec("CREATE TABLE Schedule (GameKey INT NOT NULL AUTO_INCREMENT, GameID CHAR(10) NOT NULL, Away CHAR(3) NOT NULL, Home CHAR(3) NOT NULL, PRIMARY KEY (GameKey));")
	if err != nil {
		if err.Error() != "Error 1050: Table 'Schedule' already exists" {
			log.Fatal(err)
		} else {
			fmt.Println("Schedule table already exists.")
		}
	}
	fmt.Println("Tables created.")
}

// populateTeamsTable populates the Teams database
func populateTeamsTable(hdb *sql.DB, allTeamInfo *nhlTeams) {
	fmt.Println("Populating Teams table...")

	for i := uint8(0); i < NumTeams; i++ {
		sqlStr := fmt.Sprintf("INSERT INTO Teams VALUES (%d, \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", allTeamInfo.Teams[i].ID, allTeamInfo.Teams[i].TeamName, allTeamInfo.Teams[i].Name, allTeamInfo.Teams[i].Abbreviation, allTeamInfo.Teams[i].Division.Name, allTeamInfo.Teams[i].Conference.Name)

		_, err := hdb.Exec(sqlStr)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Finished populating Teams table.")
}

// GetTeams retrieves the teams json from the api and stores relevant info
func GetTeams(hdb *sql.DB) {
	// url for the teams json
	url := "https://statsapi.web.nhl.com/api/v1/teams"
	fmt.Println("Getting teams json...")
	// gets the html from the url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved teams json.")
	// reads the response into []bytes
	teamData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// struct in the format of the json from nhl teams page
	var allTeamInfo nhlTeams
	// gets the json info from teamData ([]bytes) to allTeamInfo (nhlTeams struct)
	json.Unmarshal(teamData, &allTeamInfo)

	populateTeamsTable(hdb, &allTeamInfo)
}

func populateScheduleTable(hdb *sql.DB, fullSchedule *schedule) {
	fmt.Println("Populating Schedule table...")

	numDates := len(fullSchedule.Dates)
	for i := 0; i < numDates; i++ {
		numGames := len(fullSchedule.Dates[i].Games)
		fmt.Printf("Populating games from %s\n", fullSchedule.Dates[i].Date)
		for j := 0; j < numGames; j++ {
			if fullSchedule.Dates[i].Games[j].GameType != "PR" {
				var teamAbreviations [2]string

				awayAbbrev, err := hdb.Query("SELECT Abbreviation FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Away.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				awayAbbrev.Next()
				if err = awayAbbrev.Scan(&teamAbreviations[0]); err != nil {
					log.Fatal(err)
				}
				homeAbbrev, err := hdb.Query("SELECT Abbreviation FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Home.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				homeAbbrev.Next()
				if err = homeAbbrev.Scan(&teamAbreviations[1]); err != nil {
					log.Fatal(err)
				}

				sqlStr := fmt.Sprintf("INSERT INTO Schedule (GameID, Away, Home) VALUES (\"%s\", \"%s\", \"%s\")", strconv.FormatUint(uint64(fullSchedule.Dates[i].Games[j].GamePK), 10), teamAbreviations[0], teamAbreviations[1])

				_, err = hdb.Exec(sqlStr)
				if err != nil {
					log.Fatal(err)
				}
				awayAbbrev.Close()
				homeAbbrev.Close()
			}
		}
	}

	fmt.Println("Finished populating Schedule table.")
}

// GetSchedule retrieves the full schedule and puts the info in the Schedule table
func GetSchedule(hdb *sql.DB) {
	// url for the schedule json
	url := "https://statsapi.web.nhl.com/api/v1/schedule?season=20192020"
	fmt.Println("Getting NHL schedule...")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved schedule json.")
	// read in the json
	scheduleData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// variable to hold the json info
	var fullSchedule schedule
	// put the json info into the variable
	json.Unmarshal(scheduleData, &fullSchedule)

	populateScheduleTable(hdb, &fullSchedule)
}
