package hockeydb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GetSchedule retrieves the full schedule and returns the schedule json in a schedule struct
func GetSchedule(hdb *sql.DB, fullSchedule *Schedule) {
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
	// put the json info into the variable
	json.Unmarshal(scheduleData, &fullSchedule)
}

// PopulateScheduleTable takes the schedule and inserts info from it to the Schedule table and adds the GameNum to the ShotInfo table
func PopulateScheduleTable(hdb *sql.DB, fullSchedule *Schedule) {
	fmt.Println("Populating Schedule table...")
	// length of the Dates array
	numDates := len(fullSchedule.Dates)
	for i := 0; i < numDates; i++ {
		// length of the TotalGames array in each Date
		numGames := fullSchedule.Dates[i].TotalGames
		fmt.Printf("Populating games from %s\n", fullSchedule.Dates[i].Date)
		for j := uint8(0); j < numGames; j++ {
			if fullSchedule.Dates[i].Games[j].GameType == "R" {
				// holds the returned IDs from the queries
				var teamIDs [2]uint8
				// away team id
				awayID, err := hdb.Query("SELECT ID FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Away.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				// Can now use Scan() to get the return of the query
				awayID.Next()
				if err = awayID.Scan(&teamIDs[0]); err != nil {
					log.Fatal(err)
				}
				// home team id
				homeID, err := hdb.Query("SELECT ID FROM Teams WHERE FullName = ?;", fullSchedule.Dates[i].Games[j].Teams.Home.Team.Name)
				if err != nil {
					log.Fatal(err)
				}

				// Can now use Scan() to get the return of the query
				homeID.Next()
				if err = homeID.Scan(&teamIDs[1]); err != nil {
					log.Fatal(err)
				}

				gameNum := fullSchedule.Dates[i].Games[j].GamePK - 2019020000
				/*sqlStr := fmt.Sprintf("INSERT INTO Schedule (GameNum, GameID, Away, Home) VALUES (%d, \"%s\", %d, %d)",
					gameNum,
					strconv.FormatUint(uint64(fullSchedule.Dates[i].Games[j].GamePK), 10),
					teamIDs[0],
					teamIDs[1])
				// execute insert command for Schedule Table
				_, err = hdb.Exec(sqlStr)
				if err != nil {
					log.Fatal(err)
				}*/

				_, err = hdb.Exec("INSERT INTO ShotInfo (GameNum) VALUES (?)", gameNum)

				awayID.Close()
				homeID.Close()
			}
		}
	}

	fmt.Println("Finished populating Schedule table.")
}
