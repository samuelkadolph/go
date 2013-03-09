package campfire

type Account struct {
	CreatedAt   string `json:"created_at"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	OwnerID     int64  `json:"owner_id"`
	Plan        string `json:"plan"`
	StorageUsed int64  `json:"storage"`
	Subdomain   string `json:"subdomain"`
	Timezone    string `json:"time_zone"`
	UpdatedAt   string `json:"updated_at"`
}
