package hockeydb

type state struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}

type leagueRecord struct {
	Wins   int8   `json:"wins"`
	Losses int8   `json:"losses"`
	OT     int8   `json:"ot"`
	Type   string `json:"type"`
}

type teamDescription struct {
	ID   int8   `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type aHTeam struct {
	Record leagueRecord    `json:"leagueRecord"`
	Score  int8            `json:"score"`
	Team   teamDescription `json:"team"`
}

type awayHome struct {
	Away aHTeam `json:"away"`
	Home aHTeam `json:"home"`
}

type venueInfo struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type gameContent struct {
	Link string `json:"link"`
}

type game struct {
	GamePK   string      `json:"gamePk"`
	Link     string      `json:"link"`
	GameType string      `json:"gameType"`
	Season   string      `json:"season"`
	GameDate string      `json:"gameDate"`
	Status   state       `json:"status"`
	Teams    awayHome    `json:"teams"`
	Venue    venueInfo   `json:"venue"`
	Content  gameContent `json:"content"`
}

type date struct {
	Date         string   `json:"date"`
	TotalItems   int8     `json:"totalItems"`
	TotalEvents  int8     `json:"totalEvents"`
	TotalGames   int8     `json:"totalGames"`
	TotalMatches int8     `json:"totalMatches"`
	Games        []game   `json:"games"`
	Events       []string `json:"events"`
	Matches      []string `json:"matches"`
}

type schedule struct {
	Copyright    string `json:"copyright"`
	TotalItems   int16  `json:"totalItems"`
	TotalEvents  int8   `json:"totalEvents"`
	TotalGames   int16  `json:"totalGames"`
	TotalMatches int8   `json:"totalMatches"`
	Wait         int8   `json:"wait"`
	Dates        []date `json:"dates"`
}
