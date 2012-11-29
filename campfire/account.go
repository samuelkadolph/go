package campfire

type Account struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"owner_id"`
	Plan      string `json:"plan"`
	Storage   int64  `json:"stoage"`
	Subdomain string `json:"subdomain"`
	TimeZone  string `json:"time_zone"`
	UpdatedAt string `json:"updated_at"`
}
