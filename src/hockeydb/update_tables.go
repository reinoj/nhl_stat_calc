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

// UpdateTables will update the Schedule and ShotInfo tables
func UpdateTables(hdb *sql.DB) {
	fmt.Println("Updating tables...")

	/*
		SELECT MIN(GameNum)
		FROM Schedule
		WHERE AwayResult IS NULL
	*/
	// returns the GameNum of the first game without a result filled in
	start, err := hdb.Query("SELECT MIN(GameNum) FROM Schedule WHERE AwayResult IS NULL")
	if err != nil {
		log.Fatal(err)
	}
	defer start.Close()

	// will hold the result of the query, could return NULL so we use NullInt64 to take in the value instead of an int from the standard library
	var gameNum sql.NullInt64
	start.Next()
	if err = start.Scan(&gameNum); err != nil {
		log.Fatal(err)
	}

	// if the value is not NULL
	if gameNum.Valid {
		fmt.Printf("Starting update at GameNum: %d.\n", gameNum.Int64)
		for ; gameNum.Int64 <= 1271; gameNum.Int64++ {
			// holds the json from the feed/live page for the game
			var gameFeedLive feedLive
			getFeedLive(hdb, strconv.FormatUint(2019020000+uint64(gameNum.Int64), 10), &gameFeedLive)
			// if the game isn't finished exit the loop
			if gameFeedLive.LiveData.Linescore.CurrentPeriodTimeRemaining != "Final" {
				break
			}
			// ---------------SCHEDULE UPDATE---------------
			updateSchedule(hdb, &gameFeedLive, uint16(gameNum.Int64))

			//---------------SHOT INFO UPDATE---------------
			updateShotInfo(hdb, &gameFeedLive, uint16(gameNum.Int64))
		}
	} else {
		// if the return is NULL then there's no rows left to update
		fmt.Println("No rows to update in tables.")
	}
	fmt.Println("Finished updating tables.")
}

func getFeedLive(hdb *sql.DB, gameNum string, gameFeedLive *feedLive) {
	// url to the individual games
	url := "http://statsapi.web.nhl.com/api/v1/game/" + gameNum + "/feed/live"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	byteFeedLive, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal the json into the feedLive pointer
	json.Unmarshal(byteFeedLive, &gameFeedLive)
}

// UpdateShotInfo updates the ShotInfo table
func updateShotInfo(hdb *sql.DB, gameFeedLive *feedLive, gameNum uint16) {
	fmt.Println("Updating ShotInfo table...")
	// count away and home missed shots
	awayMissed, homeMissed := uint8(0), uint8(0)
	// only need the away team id for checking
	away := gameFeedLive.GameData.Teams.Away.ID
	// 'i' will be length of AllPlays
	for i := range gameFeedLive.LiveData.Plays.AllPlays {
		// check if the play is a missed shot
		if gameFeedLive.LiveData.Plays.AllPlays[i].Result.Event == "Missed Shot" {
			// check if it was the away team that missed
			if gameFeedLive.LiveData.Plays.AllPlays[i].Team.ID == away {
				awayMissed++
			} else { // if it wasn't the away team, then it was the home team
				homeMissed++
			}
		}
	}
	sqlStr := fmt.Sprintf("UPDATE ShotInfo SET AwayShots=%d, AwayBlocked=%d, AwayMissed=%d, HomeShots=%d, HomeBlocked=%d, HomeMissed=%d WHERE GameNum=%d",
		gameFeedLive.LiveData.Boxscore.Teams.Away.TeamStats.TeamSkaterStats.Shots,
		gameFeedLive.LiveData.Boxscore.Teams.Away.TeamStats.TeamSkaterStats.Blocked,
		awayMissed,
		gameFeedLive.LiveData.Boxscore.Teams.Home.TeamStats.TeamSkaterStats.Shots,
		gameFeedLive.LiveData.Boxscore.Teams.Home.TeamStats.TeamSkaterStats.Blocked,
		homeMissed,
		gameNum)
	_, err := hdb.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ShotInfo table updated.")
}

func updateSchedule(hdb *sql.DB, gameFeedLive *feedLive, gameNum uint16) {
	fmt.Println("Updating Schedule table with results...")

	// ---------------SCHEDULE UPDATE---------------
	var resultPrefix string
	switch gameFeedLive.LiveData.Linescore.CurrentPeriod {
	case 4:
		resultPrefix = "OT"
	case 5:
		resultPrefix = "SO"
	}
	var awayResult, homeResult string

	if gameFeedLive.LiveData.Linescore.Teams.Away.Goals > gameFeedLive.LiveData.Linescore.Teams.Home.Goals {
		awayResult, homeResult = "w", "L"
	} else {
		awayResult, homeResult = "L", "W"
	}
	/*
		UPDATE Schedule
		SET AwayResult = \"%s\", HomeResult = \"%s\"
		WHERE GameNum = %d
	*/
	sqlStr := fmt.Sprintf("UPDATE Schedule SET AwayResult = \"%s\", HomeResult = \"%s\" WHERE GameNum = %d", resultPrefix+awayResult, resultPrefix+homeResult, gameNum)
	_, err := hdb.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	//---------------------------------------------

	fmt.Println("Finished updating Schedule table.")
}
