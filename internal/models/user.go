package models

type User struct {
	Affiliations  []string `json:"affiliations"`
	Email         string   `json:"email"`
	ForwardedIP   string   `json:"forwarded_ip"`
	ForwardedHost string   `json:"forwarded_host"`
	ID            string   `json:"id"`
	Keys          []string `json:"keys"`
	LoyaltyScore  string   `json:"loyalty_score"`
	Platform      string   `json:"platform"`
	Quota         int      `json:"quota"`
	RateLimit     int      `json:"rate_limit"`
	RealIP        string   `json:"real_ip"`
	Spend         float64  `json:"spend"`
	Subscription  string   `json:"subscription"`
	Username      string   `json:"username"`
	Volume        int      `json:"volume"`
}

// NewUser creates a new user.
// func NewUser(userRealIP, username, key, forwaredIP, forwaredHost, subscription, affiliation, email, platform, quota, rateLimit) *User {
// 	return &User{
// 		Affiliations:  []string{},
// 		Email:         "",
// 		ForwardedIP:   "",
// 		ForwardedHost: "",
// 		ID:            "",
// 		Keys:          []string{},
// 		LoyaltyScore:  "",
// 		Platform:      "",
// 		Quota:         0,
// 		RateLimit:     0,
// 		RealIP:        "",
// 		Spend:         0,
// 		Subscription:  "",
// 		Username:      "",
// 		Volume:        0,
// 	}
// }
