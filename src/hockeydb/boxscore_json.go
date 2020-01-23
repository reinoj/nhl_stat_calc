package hockeydb

type boxscore struct {
	Teams awayHomeBox `json:"teams"`
	// Officials []official  `json:"officials"`
}

type awayHomeBox struct {
	Away boxTeamInfo `json:"away"`
	Home boxTeamInfo `json:"home"`
}

type boxTeamInfo struct {
	// Team       teamBoxDescription `json:"team"`
	TeamStats teamStats `json:"teamStats"`
	// Players    []boxPlayer        `json:"players"`
	// Goalies    []uint32           `json:"goalies"`
	// Skaters    []uint32           `json:"skaters"`
	// OnIce      []uint32           `json:"onIce"`
	// OnIcePlus  []onIcePlus        `json:"onIcePlus"`
	// Scratches  []uint32           `json:"scratches"`
	// PenaltyBox []penaltyBox       `json:"penaltyBoxPlus"`
	// Coaches    []boxCoaches       `json:"coaches"`
}

/*type teamBoxDescription struct {
	ID           uint8  `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
	TriCode      string `json:"triCode"`
}*/

type teamStats struct {
	TeamSkaterStats teamSkaterStats `json:"teamSkaterStats"`
}

type teamSkaterStats struct {
	// Goals                  uint8  `json:"goals"`
	// PIM                    uint8  `json:"pim"`
	Shots uint8 `json:"shots"`
	// PowerPlayPercentage    string `json:"powerPlayPercentage"`
	// PowerPlayGoals         uint8  `json:"powerPlayGoals"`
	// PowerPlayOpportunities uint8  `json:"powerPlayOpportunities"`
	// FaceOffWinPercentage   string `json:"faceOffWinPercentage"`
	Blocked uint8 `json:"blocked"`
	// Takeaways              uint8  `json:"takeaways"`
	// Giveaways              uint8  `json:"giveaways"`
	// Hits                   uint8  `json:"hits"`
}

/*type boxPlayer struct {
	Person       boxPlayerInfo `json:"person"`
	JerseyNumber string        `json:"jerseyNumber"`
	Position     boxPosition   `json:"position"`
	Stats        boxStats      `json:"stats"`
}*/

/*type onIcePlus struct {
	PlayerID uint32 `json:"playerId"`
	ShiftDuration uint16 `json:"shiftDuration"`
	Stamina uint16 `json:"stamina"`
}*/

/*type boxPlayerInfo struct {
	ID            uint32 `json:"id"`
	FullName      string `json:"fullName"`
	Link          string `json:"link"`
	ShootsCatches string `json:"shootsCatches"`
	RosterStatus  string `json:"rosterStatus"`
}*/

/*type boxPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}*/

/*type boxStats struct {
	SkaterStats skaterStats `json:"skaterStats"`
	GoalieStats goalieStats `json:"goalieStats"`
}*/

/*type skaterStats struct {
	TimeOnIce            string `json:"timeOnIce"`
	Assists              uint8  `json:"assists"`
	Goals                uint8  `json:"goals"`
	Shots                uint8  `json:"shots"`
	Hits                 uint8  `json:"hits"`
	PowerPlayGoals       uint8  `json:"powerPlayGoals"`
	PowerPlayAssists     uint8  `json:"powerPlayAssists"`
	PenaltyMinutes       uint8  `json:"penaltyMinutes"`
	FaceOffWins          uint8  `json:"faceOffWins"`
	FaceOffTaken         uint8  `json:"faceoffTaken"`
	Takeaways            uint8  `json:"takeaways"`
	Giveaways            uint8  `json:"giveaways"`
	ShortHandedGoals     uint8  `json:"shortHandedGoals"`
	ShortHandedAssists   uint8  `json:"shortHandedAssists"`
	Blocked              uint8  `json:"blocked"`
	PlusMinus            int8   `json:"plusMinus"`
	EvenTimeOnIce        string `json:"evenTimeOnIce"`
	PowerPlayTimeOnIce   string `json:"powerPlayTimeOnIce"`
	ShortHandedTimeOnIce string `json:"shortHandedTimeOnIce"`
}*/

/*type goalieStats struct {
	TimeOnIce                  string  `json:"timeOnIce"`
	Assists                    uint8   `json:"assists"`
	Goals                      uint8   `json:"goals"`
	PIM                        uint8   `json:"pim"`
	Shots                      uint8   `json:"shots"`
	Saves                      uint8   `json:"saves"`
	PowerPlaySaves             uint8   `json:"powerPlaySaves"`
	ShortHandedSAves           uint8   `json:"shortHandedSaves"`
	EvenSaves                  uint8   `json:"evenSaves"`
	ShortHandedShotsAgainst    uint8   `json:"shortHandedShotsAgainst"`
	EvenShotsAgainst           uint8   `json:"evenShotsAgainst"`
	PowerPlayShotsAgainst      uint8   `json:"powerPlayShotsAgainst"`
	Decision                   string  `json:"decision"`
	SavePercentage             float64 `json:"savePercentage"`
	PowerPlaySavePercentage    float64 `json:"powerPlaySavePercentage"`
	ShortHandedSavePercentage  float64 `json:"shortHandedSavePercentage"`
	EvenStrengthSavePercentage float64 `json:"evenStrengthSavePercentage"`
}*/

/*type penaltyBox struct {
	ID            uint32 `json:"id"`
	TimeRemaining string `json:"timeRemaining"`
	Active        bool   `json:"active"`
}*/

/*type boxCoaches struct {
	Person   boxCoachDescription `json:"person"`
	Position boxPosition         `json:"position"`
}*/

/*type boxCoachDescription struct {
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}*/

/*type official struct {
	Official     officialInfo `json:"official"`
	OfficialType string       `json:"officialType"`
}*/

/*type officialInfo struct {
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}*/
