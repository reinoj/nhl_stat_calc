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
			Abbreviation CHAR(3) NOT NULL,
			DivisionName VARCHAR(255) NOT NULL,
			ConferenceName VARCHAR(255) NOT NULL,
			PRIMARY KEY (ID)
		)
	*/
	// executes the create statement to make Teams table
	_, err := hdb.Exec("CREATE TABLE Teams (ID INT NOT NULL, TeamName VARCHAR(255) NOT NULL, FullName VARCHAR(255) NOT NULL, Abbreviation CHAR(3) NOT NULL, DivisionName VARCHAR(255) NOT NULL, ConferenceName VARCHAR(255) NOT NULL, PRIMARY KEY (ID));")
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
			GameNum INT NOT NULL,
			GameID CHAR(10) NOT NULL,
			Away INT NOT NULL,
			AwayResult VARCHAR(3),
			Home INT NOT NULL,
			HomeResult VARCHAR(3),
			PRIMARY KEY (GameNum)
		)
	*/
	_, err = hdb.Exec("CREATE TABLE Schedule (GameNum INT NOT NULL, GameID CHAR(10) NOT NULL, Away INT NOT NULL, AwayResult VARCHAR(3), Home INT NOT NULL, HomeResult VARCHAR(3), PRIMARY KEY (GameNum));")
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

// GetSchedule retrieves the full schedule and returns the schedule json in a schedule struct
func GetSchedule(hdb *sql.DB) Schedule {
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
	var fullSchedule Schedule
	// put the json info into the variable
	json.Unmarshal(scheduleData, &fullSchedule)
	return fullSchedule
}

// PopulateScheduleTable takes the schedule and inserts info from it to the Schedule table
func PopulateScheduleTable(hdb *sql.DB, fullSchedule *Schedule) {
	fmt.Println("Populating Schedule table...")

	numDates := len(fullSchedule.Dates)
	for i := 0; i < numDates; i++ {
		numGames := fullSchedule.Dates[i].TotalGames
		fmt.Printf("Populating games from %s\n", fullSchedule.Dates[i].Date)
		for j := uint8(0); j < numGames; j++ {
			if fullSchedule.Dates[i].Games[j].GameType == "R" {
				var teamIDs [2]uint8

				awayAbbrev, err := hdb.Query("SELECT ID FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Away.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				awayAbbrev.Next()
				if err = awayAbbrev.Scan(&teamIDs[0]); err != nil {
					log.Fatal(err)
				}
				homeAbbrev, err := hdb.Query("SELECT ID FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Home.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				homeAbbrev.Next()
				if err = homeAbbrev.Scan(&teamIDs[1]); err != nil {
					log.Fatal(err)
				}

				sqlStr := fmt.Sprintf("INSERT INTO Schedule (GameNum, GameID, Away, Home) VALUES (%d, \"%s\", %d, %d)",
					fullSchedule.Dates[i].Games[j].GamePK-2019020000,
					strconv.FormatUint(uint64(fullSchedule.Dates[i].Games[j].GamePK), 10),
					teamIDs[0],
					teamIDs[1])

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

func getLinescore(hdb *sql.DB, gameNum string, gameLinescore *linescore) {
	url := "http://statsapi.web.nhl.com/api/v1/game/" + gameNum + "/linescore"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	byteLinescore, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteLinescore, gameLinescore)
}

// UpdateScheduleResults updates the results in the Schedule table
func UpdateScheduleResults(hdb *sql.DB) {
	fmt.Println("Updating Schedule table with results...")
	var gameLinescore linescore
	for gameNum := uint64(2019020001); gameNum <= 2019021271; gameNum++ {
		getLinescore(hdb, strconv.FormatUint(gameNum, 10), &gameLinescore)

		if gameLinescore.CurrentPeriodTimeRemaining != "Final" {
			break
		}
	}
	fmt.Println("Finished updating Schedule table.")
}
