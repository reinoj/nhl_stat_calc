package hockeydb

type timezone struct {
	ID     string `json:"id"`
	Offset int8   `json:"offset"`
	TZ     string `json:"tz"`
}

type venue struct {
	Name     string   `json:"name"`
	Link     string   `json:"link"`
	City     string   `json:"city"`
	TimeZone timezone `json:"timeZone"`
}

type division struct {
	ID           int8   `json:"id"`
	Name         string `json:"name"`
	NameShort    string `json:"nameShort"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
}

type conference struct {
	ID   int8   `json:"id"`
	Name string `json:"teamName"`
	Link string `json:"link"`
}

type franchise struct {
	FranchiseID int8   `json:"franchiseId"`
	TeamName    string `json:"teamName"`
	Link        string `json:"link"`
}

type team struct {
	ID               int8       `json:"id"`
	Name             string     `json:"name"`
	Link             string     `json:"link"`
	Venue            venue      `json:"venue"`
	Abbreviation     string     `json:"abbreviation"`
	TeamName         string     `json:"teamName"`
	LocationName     string     `json:"locationName"`
	FirstYearOfPlay  string     `json:"firstYearOfPlay"`
	Division         division   `json:"division"`
	Conference       conference `json:"conference"`
	Franchise        franchise  `json:"franchise"`
	ShortName        string     `json:"shortName"`
	OfficialySiteURL string     `json:"officialSiteUrl"`
	FranchiseID      int8       `json:"franchiseId"`
	Active           bool       `json:"active"`
}

type nhlTeams struct {
	Copyright string `json:"copyright"`
	Teams     []team `json:"teams"`
}
