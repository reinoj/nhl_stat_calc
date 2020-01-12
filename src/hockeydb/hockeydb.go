package hockeydb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// Switch to the Hockey database
	_, err := hdb.Exec("USE Hockey")
	if err != nil {
		log.Fatal(err)
	}
	/*
		"CREATE TABLE Teams (
			ID int NOT NULL,
			TeamName VARCHAR(255) NOT NULL,
			LocationName VARCHAR(255) NOT NULL,
			Abbreviation VARCHAR(255) NOT NULL,
			DivisionName VARCHAR(255) NOT NULL,
			ConferenceName VARCHAR(255) NOT NULL,
			PRIMARY KEY (ID)
		)"
	*/
	// executes the create statement to make the table
	_, err = hdb.Exec("CREATE TABLE Teams (ID int NOT NULL, TeamName VARCHAR(255) NOT NULL, LocationName VARCHAR(255) NOT NULL, Abbreviation VARCHAR(255) NOT NULL, DivisionName VARCHAR(255) NOT NULL, ConferenceName VARCHAR(255) NOT NULL, PRIMARY KEY (ID))")
	if err != nil {
		log.Fatal(err)
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

	fmt.Println("Teams table populated.")
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
