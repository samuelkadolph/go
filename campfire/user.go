package campfire

type User struct {
	Admin        bool   `json:"admin"`
	APIAuthToken string `json:"api_auth_token"`
	AvatarURL    string `json:"avatar_url"`
	CreatedAt    string `json:"created_at"`
	EmailAddress string `json:"email_address"`
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

func Me(subdomain, username, password string, config Config) (*User, error) {
	c := newConnection(subdomain, "", config)
	c.Password = password
	c.Username = username

	var wrapper struct {
		User User `json:"user"`
	}

	if err := c.get("/users/me", &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}
