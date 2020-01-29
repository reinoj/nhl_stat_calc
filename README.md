# nhl_stat_calc

## Running

### Dependencies

* MySQL
* Golang (I ran this on 1.13.5)
  * > go get github.com/reinoj/nhl_stat_calc

### Flags

* -initialSetup
  * creating the database, only do first time running
* -createTables
  * creating the tables in the database
* -updateTables
  * updating the Schedule and ShotInfo tables
* -calculateCorsi
  * calculate corsi for each game and store in ShotInfo table
* -outputCorsi
  * output the corsi results to .json file
* -mysqlUser
  * username for your mysql
* -mysqlPassword
  * password for your mysql

## To-Do

* Multiple years, instead of just 2019-20
