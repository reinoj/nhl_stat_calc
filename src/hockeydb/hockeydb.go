package hockeydb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// Switch to the Hockey database
	_, err := hdb.Exec("USE Hockey")
	if err != nil {
		log.Fatal(err)
	}

	/*
		CREATE TABLE Teams (
			ID int NOT NULL,
			TeamName VARCHAR(255) NOT NULL,
			LocationName VARCHAR(255) NOT NULL,
			Abbreviation VARCHAR(255) NOT NULL,
			DivisionName VARCHAR(255) NOT NULL,
			ConferenceName VARCHAR(255) NOT NULL,
			PRIMARY KEY (ID)
		)
	*/
	// executes the create statement to make Teams table
	_, err = hdb.Exec("CREATE TABLE Teams (ID INT NOT NULL, TeamName VARCHAR(255) NOT NULL, LocationName VARCHAR(255) NOT NULL, Abbreviation VARCHAR(255) NOT NULL, DivisionName VARCHAR(255) NOT NULL, ConferenceName VARCHAR(255) NOT NULL, PRIMARY KEY (ID))")
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
	_, err = hdb.Exec("CREATE TABLE Schedule (GameKey INT NOT NULL AUTO_INCREMENT, GameID CHAR(10) NOT NULL, Away CHAR(3) NOT NULL, Home CHAR(3) NOT NULL, PRIMARY KEY (GameKey))")
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
func populateTeamsTable(hdb *sql.DB, teams [31]teamInfo) {
	fmt.Println("Populating Teams table...")

	for i := uint8(0); i < NumTeams; i++ {
		sqlStr := fmt.Sprintf("INSERT INTO Teams VALUES (%d, \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", teams[i].ID, teams[i].TeamName, teams[i].LocationName, teams[i].Abbreviation, teams[i].DivisionName, teams[i].ConferenceName)

		fmt.Println(sqlStr)

		// Switch to the Hockey database
		_, err := hdb.Exec("USE Hockey")
		if err != nil {
			log.Fatal(err)
		}
		_, err = hdb.Query(sqlStr)
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

	// NumTeams length teamInfo struct that contains only the needed info
	var teams [NumTeams]teamInfo
	// assigns the info from allTeamInfo to teams
	for i := uint8(0); i < NumTeams; i++ {
		teams[i] = teamInfo{allTeamInfo.Teams[i].ID, allTeamInfo.Teams[i].TeamName, allTeamInfo.Teams[i].LocationName, allTeamInfo.Teams[i].Abbreviation, allTeamInfo.Teams[i].Division.Name, allTeamInfo.Teams[i].Conference.Name}
	}

	populateTeamsTable(hdb, teams)
}

func populateScheduleTable(hdb *sql.DB, fullSchedule schedule) {
	fmt.Println("Populating Schedule table...")

	fmt.Println("Finished populating Schedule table.")
}

// GetSchedule retrieves the full schedule and puts the info in the Schedule table
func GetSchedule(hdb *sql.DB) {
	url := "https://statsapi.web.nhl.com/api/v1/schedule?season=20192020"
	fmt.Println("Getting NHL schedule...")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved schedule json.")
	scheduleData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var fullSchedule schedule
	json.Unmarshal(scheduleData, &fullSchedule)
	populateScheduleTable(hdb, fullSchedule)
}
