package statcalc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// (shots + missed + blocked against) - (shots against + missed shots against + blocked)

// CorsiCalc will populate a table with the corsi results
func CorsiCalc(hdb *sql.DB) {
	start, err := hdb.Query("SELECT MIN(GameNum) FROM ShotInfo WHERE AwayCorsi IS NULL;")
	if err != nil {
		log.Fatal(err)
	}
	defer start.Close()

	var gameStart sql.NullInt64
	start.Next()
	if err = start.Scan(&gameStart); err != nil {
		log.Fatal(err)
	}

	gameStart.Int64 = 1

	if gameStart.Valid {
		fmt.Printf("Starting Corsi calculation at game %d.\n", gameStart.Int64)
		for gameNum := uint16(gameStart.Int64); gameNum <= 1271; gameNum++ {
			var gameShotInfo [6]sql.NullInt64
			getGameCorsi(hdb, gameNum, &gameShotInfo)
			fmt.Println(gameShotInfo, gameShotInfo[0].Valid)
			if gameShotInfo[0].Valid {
				fmt.Printf("Calculating Corsi for game #%d\n", gameNum)

				awayCorsi := int16((gameShotInfo[0].Int64 + gameShotInfo[2].Int64 + gameShotInfo[4].Int64) - (gameShotInfo[3].Int64 + gameShotInfo[5].Int64 + gameShotInfo[1].Int64))
				fmt.Printf("awayCorsi = (%d + %d + %d) - (%d + %d + %d) = %d", gameShotInfo[0].Int64, gameShotInfo[2].Int64, gameShotInfo[4].Int64, gameShotInfo[3].Int64, gameShotInfo[5].Int64, gameShotInfo[1].Int64, awayCorsi)
				sqlStr := fmt.Sprintf("UPDATE ShotInfo SET AwayCorsi=%d, HomeCorsi=%d WHERE GameNum=%d", awayCorsi, -awayCorsi, gameNum)

				_, err = hdb.Exec(sqlStr)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				break
			}
		}
	} else {
		// if the return is NULL then there's no rows left to update
		fmt.Println("No rows to calculate in ShotInfo.")
	}
}

func getGameCorsi(hdb *sql.DB, gameNum uint16, gameShotInfo *[6]sql.NullInt64) {
	row, err := hdb.Query("SELECT AwayShots, AwayBlocked, AwayMissed, HomeShots, HomeBlocked, HomeMissed FROM ShotInfo WHERE GameNum=?;", gameNum)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	row.Next()
	if err = row.Scan(&gameShotInfo[0], &gameShotInfo[1], &gameShotInfo[2], &gameShotInfo[3], &gameShotInfo[4], &gameShotInfo[5]); err != nil {
		log.Fatal(err)
	}

}

// GetCorsiWins will output the number of games where the team with higher corsi won or loss
func GetCorsiWins(hdb *sql.DB) {
	totalGames, totalWins, otWins, soWins := 0, 0, 0, 0
	for i := 1; i <= 1271; i++ {
		awayGameResult := getResultAndCorsi(hdb, &i)
		if awayGameResult.AwayResult.Valid {
			// if both are not NULL then we increment the total number of games played
			totalGames++
			if (awayGameResult.AwayResult.String == "W" && awayGameResult.AwayCorsi.Int64 > 0) ||
				(awayGameResult.AwayResult.String == "L" && awayGameResult.AwayCorsi.Int64 < 0) {
				// if the string is a W and the corsi is positive or if the string is an L and the corsi is less than 0, increment the number of wins
				totalWins++

				if awayGameResult.OT.Bool {
					otWins++
				} else if awayGameResult.SO.Bool {
					soWins++
				}
			}
		} else {
			break
		}
	}
	fmt.Println("Total games: ", totalGames)
	// write the new results to the file
	writeCorsiResult(totalGames, totalWins, otWins, soWins)
}

type result struct {
	AwayResult sql.NullString
	OT         sql.NullBool
	SO         sql.NullBool
	AwayCorsi  sql.NullInt64
}

func getResultAndCorsi(hdb *sql.DB, gameNum *int) result {
	queryResult, err := hdb.Query("SELECT AwayResult, OT, SO FROM Schedule WHERE GameNum=?", gameNum)
	if err != nil {
		log.Fatal(err)
	}
	defer queryResult.Close()

	var gameResult result
	queryResult.Next()
	err = queryResult.Scan(&gameResult.AwayResult, &gameResult.OT, &gameResult.SO)
	if err != nil {
		log.Fatal(err)
	}

	corsi, err := hdb.Query("SELECT AwayCorsi FROM ShotInfo WHERE GameNum=?", gameNum)
	if err != nil {
		log.Fatal(err)
	}
	defer corsi.Close()

	corsi.Next()
	err = corsi.Scan(&gameResult.AwayCorsi)
	if err != nil {
		log.Fatal(err)
	}
	return gameResult
}

type corsiOutput struct {
	Description       string `json:"description"`
	TotalGames        int    `json:"totalGames"`
	AllCorsiWins      int    `json:"allCorsiWins"`
	NoSOCorsiWins     int    `json:"noSOCorsiWins"`
	NoOTOrSOCorsiWins int    `json:"noOTOrSOCorsiWins"`
}

func writeCorsiResult(totalGames, totalWins, otWins, soWins int) {
	outputPreMarshal := corsiOutput{"numberOfGames", totalGames, totalWins, totalWins - soWins, totalWins - otWins - soWins}
	output, err := json.Marshal(outputPreMarshal)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("statcalc/corsi.json", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(output)
}
