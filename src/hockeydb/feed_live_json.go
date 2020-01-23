package hockeydb

type feedLive struct {
	// Copyright string   `json:"copyright"`
	GamePK uint32 `json:"gamePk"`
	// Link      string   `json:"link"`
	// MetaData  metaData `json:"metaData"`
	GameData gameData `json:"gameData"`
	LiveData liveData `json:"liveData"`
}

/*type metaData struct {
	Wait      uint8  `json:"wait"`
	TimeStamp string `json:"timeStamp"`
}*/

type gameData struct {
	// Game     feedLiveGame  `json:"game"`
	// DateTime dateTime      `json:"datetime"`
	// Status   status        `json:"status"`
	Teams feedLiveTeams `json:"teams"`
	// Players  []player      `json:"players"`
	// Venue    feedLiveVenue `json:"venue"`
}

/*type feedLiveGame struct {
	PK     uint32 `json:"pk"`
	Season string `json:"season"`
	Type   string `json:"type"`
}*/

/*type dateTime struct {
	DateTime    string `json:"dateTime"`
	EndDateTime string `json:"endDateTime"`
}*/

/*type status struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}*/

type feedLiveTeams struct {
	Away team `json:"away"`
	Home team `json:"home"`
	// team struct in teams_json.go
}

/*type player struct {
	ID               uint32          `json:"id"`
	FullName         string          `json:"fullName"`
	Link             string          `json:"link"`
	FirstName        string          `json:"firstName"`
	LastName         string          `json:"lastName"`
	PrimaryNumber    string          `json:"primaryNumber"`
	BirthDate        string          `json:"birthDate"`
	CurrentAge       uint8           `json:"currentAge"`
	BirthCity        string          `json:"birthCity"`
	BirthCountry     string          `json:"birthCountry"`
	Nationality      string          `json:"nationality"`
	Height           string          `json:"height"`
	Weight           uint16          `json:"weight"`
	Active           bool            `json:"active"`
	AlternateCaptain bool            `json:"alternateCaptain"`
	Captain          bool            `json:"captain"`
	Rookie           bool            `json:"rookie"`
	ShootsCatches    string          `json:"shootsCatches"`
	RosterStatus     string          `json:"rosterStatus"`
	CurrentTeam      feedLiveTeam    `json:"currentTeam"`
	PrimaryPosition  primaryPosition `json:"primaryPosition"`
}*/

type feedLiveTeam struct {
	ID uint8 `json:"id"`
	// Name    string `json:"name"`
	// Link    string `json:"link"`
	// TriCode string `json:"triCode"`
}

/*type primaryPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}*/

/*type feedLiveVenue struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}*/

type liveData struct {
	Plays     plays     `json:"plays"`
	Linescore linescore `json:"linescore"`
	Boxscore  boxscore  `json:"boxscore"`
	// Decisions decisions `json:"decisions"`
}

type plays struct {
	AllPlays []play `json:"allPlays"`
	// ScoringPlays  []play            `json:"scoringPlays"`
	// PenaltyPlays  []play            `json:"penaltyPlays"`
	// PlaysByPeriod []periodPlaysInfo `json:"playsByPeriod"`
}

/*type periodPlaysInfo struct {
	StartIndex uint16 `json:"startIndex"`
	Plays      []play `json:"plays"`
	EndIndex   uint16 `json:"endIndex"`
}*/

type play struct {
	// Players     []playPlayerInfo `json:"players"`
	Result result `json:"result"`
	// About       about            `json:"about"`
	// Coordinates coordinates      `json:"coordinates"`
	Team feedLiveTeam `json:"team"`
}

/*type playPlayerInfo struct {
	Player     playPlayerDescription `json:"player"`
	PlayerType string                `json:"playerType"`
}*/

/*type playPlayerDescription struct {
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}*/

type result struct {
	Event string `json:"event"`
	// EventCode   string `json:"eventCode"`
	// EventTypeID string `json:"eventTypeId"`
	// Description string `json:"description"`
}

/*type about struct {
	EventIDx            uint16 `json:"eventIdx"`
	EventID             uint16 `json:"eventId"`
	Period              uint8  `json:"period"`
	PeriodType          string `json:"periodType"`
	OrdinalNum          string `json:"ordinalNum"`
	PeriodTime          string `json:"periodTime"`
	PeriodTimeRemaining string `json:"periodTimeRemaining"`
	DateTime            string `json:"dateTime"`
	Goals               goals  `json:"goals"`
}*/

/*type goals struct {
	Away uint8 `json:"away"`
	Home uint8 `json:"home"`
}*/

/*type coordinates struct {
	X int8 `json:"x"`
	Y int8 `json:"y"`
}*/

/*type decisions struct {
	Winner     playPlayerDescription `json:"winner"`
	Loser      playPlayerDescription `json:"loser"`
	FirstStar  playPlayerDescription `json:"firstStar"`
	SecondStar playPlayerDescription `json:"secondStar"`
	ThirdStar  playPlayerDescription `json:"thirdStar"`
}*/
