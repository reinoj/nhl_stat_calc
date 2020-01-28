package statcalc

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	totalGames, wins := 0, 0
	for i := 1; i <= 1271; i++ {
		fmt.Println(i)
		awayResult, awayCorsi := getResultAndCorsi(hdb, &i)
		if awayResult.Valid && awayCorsi.Valid {
			// if both are not NULL then we increment the total number of games played
			totalGames++
			if (strings.HasSuffix(awayResult.String, "W") && awayCorsi.Int64 > 0) ||
				(!strings.HasSuffix(awayResult.String, "W") && awayCorsi.Int64 < 0) {
				// if the last character in the string is a W and the corsi is positive, then increment the number of wins
				wins++
			}
		} else {
			break
		}
	}
	fmt.Println(totalGames, "\t", wins)
	// write the new results to the file
	writeCorsiResult(totalGames, wins)
}

/*func getCurrentCorsi() (int, int) {
	corsiFile, err := os.Open("statcalc/corsi.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer corsiFile.Close()
	corsiCsv := csv.NewReader(corsiFile)
	csvArray, err := corsiCsv.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	tG, err := strconv.Atoi(csvArray[1][1])
	if err != nil {
		log.Fatal(err)
	}
	w, err := strconv.Atoi(csvArray[2][1])
	if err != nil {
		log.Fatal(err)
	}
	return tG + 1, w
}*/

func getResultAndCorsi(hdb *sql.DB, gameNum *int) (sql.NullString, sql.NullInt64) {
	result, err := hdb.Query("SELECT AwayResult FROM Schedule WHERE GameNum=?", gameNum)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	var awayResult sql.NullString
	result.Next()
	err = result.Scan(&awayResult)
	if err != nil {
		log.Fatal(err)
	}

	corsi, err := hdb.Query("SELECT AwayCorsi FROM ShotInfo WHERE GameNum=?", gameNum)
	if err != nil {
		log.Fatal(err)
	}
	defer corsi.Close()

	var awayCorsi sql.NullInt64
	corsi.Next()
	err = corsi.Scan(&awayCorsi)
	if err != nil {
		log.Fatal(err)
	}
	return awayResult, awayCorsi
}

func writeCorsiResult(totalGames, wins int) {
	output := [][]string{
		{"description", "numberOfGames"},
		{"totalGames", strconv.Itoa(int(totalGames))},
		{"wins", strconv.Itoa(int(wins))},
	}
	file, err := os.OpenFile("statcalc/corsi.csv", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range output {
		err = writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func madeIt(location string) {
	fmt.Printf("Made it %s.\n", location)
}
