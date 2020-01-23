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
	// unmarshal the json into the linescore pointer
	json.Unmarshal(byteFeedLive, &gameFeedLive)
}

func updateShotInfo(hdb *sql.DB, gameFeedLive *feedLive, gameNum uint16) {
	awayMissed, homeMissed := uint8(0), uint8(0)
	away := gameFeedLive.GameData.Teams.Away.ID
	for i := range gameFeedLive.LiveData.Plays.AllPlays {
		if gameFeedLive.LiveData.Plays.AllPlays[i].Result.Event == "Missed Shot" {
			if gameFeedLive.LiveData.Plays.AllPlays[i].Team.ID == away {
				awayMissed++
			} else {
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
}

// UpdateScheduleResults updates the results in the Schedule table
func UpdateScheduleResults(hdb *sql.DB) {
	fmt.Println("Updating Schedule table with results...")

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

	var gameNum sql.NullInt64
	start.Next()
	if err = start.Scan(&gameNum); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting update at GameNum: %d.\n", gameNum.Int64)
	if gameNum.Valid {
		for ; gameNum.Int64 <= 1271; gameNum.Int64++ {
			var gameFeedLive feedLive
			getFeedLive(hdb, strconv.FormatUint(2019020000+uint64(gameNum.Int64), 10), &gameFeedLive)

			if gameFeedLive.LiveData.Linescore.CurrentPeriodTimeRemaining != "Final" {
				break
			}
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
			sqlStr := fmt.Sprintf("UPDATE Schedule SET AwayResult = \"%s\", HomeResult = \"%s\" WHERE GameNum = %d", resultPrefix+awayResult, resultPrefix+homeResult, gameNum.Int64)
			_, err := hdb.Exec(sqlStr)
			if err != nil {
				log.Fatal(err)
			}
			//---------------------------------------------

			//---------------SHOT INFO UPDATE---------------
			updateShotInfo(hdb, &gameFeedLive, uint16(gameNum.Int64))
			//---------------------------------------------
		}
	} else {
		fmt.Println("No rows to update in Schedule table.")
	}
	fmt.Println("Finished updating Schedule table.")
}
