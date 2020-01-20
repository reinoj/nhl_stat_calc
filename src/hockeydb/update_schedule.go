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

func getLinescore(hdb *sql.DB, gameNum string, gameLinescore *linescore) {
	// url to the individual games
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
	// unmarshal the json into the linescore pointer
	json.Unmarshal(byteLinescore, gameLinescore)
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
	start.Next()
	var gameNum sql.NullInt64
	if err = start.Scan(&gameNum); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting Schedule update at GameNum %d.\n", gameNum.Int64)
	if gameNum.Valid {
		for ; gameNum.Int64 <= 1271; gameNum.Int64++ {
			var gameLinescore linescore
			getLinescore(hdb, strconv.FormatUint(2019020000+uint64(gameNum.Int64), 10), &gameLinescore)

			fmt.Println(gameLinescore.CurrentPeriodTimeRemaining)
			if gameLinescore.CurrentPeriodTimeRemaining != "Final" {
				break
			}
			var resultPrefix string
			switch gameLinescore.CurrentPeriod {
			case 4:
				resultPrefix = "OT"
			case 5:
				resultPrefix = "SO"
			}
			var awayResult, homeResult string

			if gameLinescore.Teams.Away.Goals > gameLinescore.Teams.Home.Goals {
				awayResult, homeResult = "w", "L"
			} else {
				awayResult, homeResult = "L", "W"
			}
			/*
				UPDATE Schedule
				SET AwayResult = \"%s\", HomeResult = \"%s\"
				WHERE GameNum = %d
			*/
			updateString := fmt.Sprintf("UPDATE Schedule SET AwayResult = \"%s\", HomeResult = \"%s\" WHERE GameNum = %d", resultPrefix+awayResult, resultPrefix+homeResult, gameNum.Int64)
			fmt.Println(updateString)
			_, err := hdb.Exec(updateString)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Println("No rows to update in Schedule table.")
	}
	fmt.Println("Finished updating Schedule table.")
}
