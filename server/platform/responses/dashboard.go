package responses

import "time"

type DashboardBucket struct {
	Bucket time.Time `json:"bucket"`
	Count  int       `json:"count"`
}

type Dashboard struct {
	Reactions []DashboardBucket `json:"reactions"`
	Profiles  []DashboardBucket `json:"profiles"`
	Events    int               `json:"events"`
	Actions   int               `json:"actions"`
}
