package models

type User struct {
	ID        string `json:"id"`
	Plan      string `json:"plan"`
	Quota     int    `json:"quota"`
	RateLimit int    `json:"rate_limit"`
	Referrer  string `json:"referrer"`
	Keys      []string  `json:"keys"`
	Spend     float64  `json:"spend"`
	Loyalty   string    `json:"loyalty"`
	Affiliations []string `json:"affiliations"`
}


