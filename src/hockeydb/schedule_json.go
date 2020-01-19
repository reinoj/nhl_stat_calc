package hockeydb

/*type state struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}*/

/*type leagueRecord struct {
	Wins   uint8  `json:"wins"`
	Losses uint8  `json:"losses"`
	OT     uint8  `json:"ot"`
	Type   string `json:"type"`
}*/

type teamDescription struct {
	// ID uint8 `json:"id"`
	Name string `json:"name"`
	// Link string `json:"link"`
}

type aHTeam struct {
	// Record leagueRecord    `json:"leagueRecord"`
	// Score  uint8           `json:"score"`
	Team teamDescription `json:"team"`
}

type awayHome struct {
	Away aHTeam `json:"away"`
	Home aHTeam `json:"home"`
}

/*type venueInfo struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}*/

/*type gameContent struct {
	Link string `json:"link"`
}*/

type game struct {
	GamePK uint32 `json:"gamePk"`
	// Link     string      `json:"link"`
	GameType string `json:"gameType"`
	// Season   string      `json:"season"`
	// GameDate string      `json:"gameDate"`
	// Status   state       `json:"status"`
	Teams awayHome `json:"teams"`
	// Venue    venueInfo   `json:"venue"`
	// Content  gameContent `json:"content"`
}

type date struct {
	Date string `json:"date"`
	// TotalItems   uint8    `json:"totalItems"`
	// TotalEvents  uint8    `json:"totalEvents"`
	TotalGames uint8 `json:"totalGames"`
	// TotalMatches uint8    `json:"totalMatches"`
	Games []game `json:"games"`
	// Events       []string `json:"events"`
	// Matches      []string `json:"matches"`
}

// Schedule is used for holding the json from NHL schedule statsapi
type Schedule struct {
	// Copyright    string `json:"copyright"`
	// TotalItems   uint16 `json:"totalItems"`
	// TotalEvents  uint8  `json:"totalEvents"`
	// TotalGames   uint16 `json:"totalGames"`
	// TotalMatches uint8  `json:"totalMatches"`
	// Wait         int8   `json:"wait"`
	Dates []date `json:"dates"`
}
