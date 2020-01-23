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
	var allTeamInfo currentTeams
	// gets the json info from teamData ([]bytes) to allTeamInfo (nhlTeams struct)
	json.Unmarshal(teamData, &allTeamInfo)

	populateTeamsTable(hdb, &allTeamInfo)
}

// populateTeamsTable populates the Teams database
func populateTeamsTable(hdb *sql.DB, allTeamInfo *currentTeams) {
	fmt.Println("Populating Teams table...")

	for i := uint8(0); i < NumTeams; i++ {

		sqlStr := fmt.Sprintf("INSERT INTO Teams VALUES (%d, \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")",
			allTeamInfo.Teams[i].ID,
			allTeamInfo.Teams[i].TeamName,
			allTeamInfo.Teams[i].Name,
			allTeamInfo.Teams[i].Abbreviation,
			allTeamInfo.Teams[i].Division.Name,
			allTeamInfo.Teams[i].Conference.Name)
		_, err := hdb.Exec(sqlStr)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Finished populating Teams table.")
}
