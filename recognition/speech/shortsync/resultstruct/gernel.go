package resultstruct

type Speech struct {
	Name    string   `json:"name"`
	Label   int      `json:"label"`
	Review  bool     `json:"review"`
	Rate    float32  `json:"rate,omitempty"`
	Details []Detail `json:details`
}

type Detail struct {
	StartTime float32 `json:"startTime"`
	EndTime   float32 `json:"endTime"`
	Label     int     `json:"label"`
	Rate      float64 `json:"rate"`
}
