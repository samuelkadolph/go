package campfire

type User struct {
	Admin        bool   `json:""`
	AvatarURL    string `json:"avatar_url"`
	CreatedAt    string `json:"created_at"`
	EmailAddress string `json:"email_address"`
	ID           int    `json:""`
	Name         string `json:""`
	Type         string `json:""`
}
