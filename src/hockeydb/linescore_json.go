package hockeydb

type linescore struct {
	// Copyright            string           `json:"copyright"`
	CurrentPeriod uint8 `json:"currentPeriod"`
	// CurrentPeriodOrdinal string           `json:"currentPeriodOrdinal"`
	CurrentPeriodTimeRemaining string `json:"currentPeriodTimeRemaining"`
	// Periods              []periods        `json:"periods"`
	ShootoutInfo shootoutInfo `json:"shootoutInfo"`
	// Teams                homeAway         `json:"teams"`
	// PowerPlayStrength    string           `json:"powerPlayStrength"`
	// HasShootout          bool             `json:"hasShootout"`
	// IntermissionInfo     intermissionInfo `json:"intermissionInfo"`
	// PowerPlayInfo        powerPlayInfo    `json:"powerPlayInfo"`
}

/*type periods struct {
	PeriodType string      `json:"periodType"`
	StartTime  string      `json:"startTime"`
	EndTime    string      `json:"endTime"`
	Num        uint8       `json:"num"`
	OrdinalNum string      `json:"ordinalNum"`
	Home       periodStats `json:"home"`
	Away       periodStats `json:"away"`
}*/

/*type periodStats struct {
	Goals       uint8  `json:"goals"`
	ShotsOnGoal uint8  `json:"shotsOnGoal"`
	RinkSide    string `json:"rinkSide"`
}*/

type shootoutInfo struct {
	Away shootoutTeam `json:"away"`
	Home shootoutTeam `json:"home"`
}

type shootoutTeam struct {
	Scores   uint8 `json:"scores"`
	Attempts uint8 `json:"attempts"`
}

/*type homeAway struct {
	Home linescoreTeam `json:"home"`
	Away linescoreTeam `json:"away"`
}*/

/*type linescoreTeam struct {
	TeamInfo     linescoreTeamInfo `json:"team"`
	Goals        uint8             `json:"goals"`
	ShotsOnGoal  uint8             `json:"shotsOnGoal"`
	GoaliePulled bool              `json:"goaliePulled"`
	NumSkaters   uint8             `json:"numSkaters"`
	PowerPlay    bool              `json:"powerPlay"`
}*/

/*type linescoreTeamInfo struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}*/

/*type intermissionInfo struct {
	IntermissionTimeRemaining uint16 `json:"intermissionTimeRemaining"`
	IntermissionTimeElapsed   uint16 `json:"intermissionTimeElapsed"`
	InIntermission            bool   `json:"inIntermission"`
}*/

/*type powerPlayInfo struct {
	SituationTimeRemaining uint16 `json:"situationTimeRemaining"`
	SituationTimeElapsed   uint16 `json:"situationTimeElapsed"`
	InSituation            bool   `json:"inSituation"`
}*/
