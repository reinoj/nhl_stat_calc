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

	var allTeams nhlTeams
	json.Unmarshal(teamData, &allTeams)

	fmt.Printf("\n\nname: %s\nabbrev: %s\ncity: %s\n\n", allTeams.Teams[0].Name, allTeams.Teams[0].Abbreviation, allTeams.Teams[0].Venue.City)
}
