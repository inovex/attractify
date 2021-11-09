package responses

type AnalyzeEvents struct {
	Event   string `json:"event"`
	Total   int    `json:"total"`
	Year    int    `json:"year"`
	Month   int    `json:"month"`
	Day     int    `json:"day"`
	WeekDay int    `json:"weekDay"`
	Hour    int    `json:"hour"`
}

type AnalyzeRates struct {
	Delivered int `json:"delivered"`
	Shown     int `json:"shown"`
	Hidden    int `json:"hidden"`
	Declined  int `json:"declined"`
	Accepted  int `json:"accepted"`
}

type AnalyzeReach struct {
	Channel string `json:"channel"`
	Total   int    `json:"total"`
}
