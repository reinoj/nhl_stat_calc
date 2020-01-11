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
func CreateDb(hdb *sql.DB) {
	// executes the create statement
	_, err := hdb.Exec("CREATE DATABASE hockey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Created Database")
}

// CreateTables creates the tables inside the database
func CreateTables(hdb *sql.DB) {
	// executes the create statement
	_, err := hdb.Exec("CREATE TABLE teams")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("teams table created.")
}

// MOVE THESE STRUCTS TO A DIFFERENT FILE

// GetTeams words words words
func GetTeams(hdb *sql.DB) {
	url := "https://statsapi.web.nhl.com/api/v1/teams"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	teamData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var allTeamInfo nhlTeams
	json.Unmarshal(teamData, &allTeamInfo)

	var teams [NumTeams]teamInfo
	for i := uint8(0); i < NumTeams; i++ {
		teams[i] = teamInfo{allTeamInfo.Teams[i].ID, allTeamInfo.Teams[i].TeamName, allTeamInfo.Teams[i].LocationName, allTeamInfo.Teams[i].Abbreviation, allTeamInfo.Teams[i].Division.Name, allTeamInfo.Teams[i].Conference.Name}
	}

	for i := uint8(0); i < NumTeams; i++ {
		fmt.Printf("name: %s\n", teams[i].teamName)
	}
}
