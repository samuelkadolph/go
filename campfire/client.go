package campfire

type Client struct {
	Subdomain string
	Token     string

	connection *connection
}

func NewClient(subdomain, token string, config Config) *Client {
	return &Client{Subdomain: subdomain, Token: token, connection: newConnection(subdomain, token, config)}
}

func (c *Client) Account() (*Account, error) {
	var wrapper struct {
		Account Account `json:"account"`
	}

	if err := c.connection.get("/account", &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.Account, nil
}

func (c *Client) Me() (*User, error) {
	var wrapper struct {
		User User `json:"user"`
	}

	if err := c.connection.get("/users/me", &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) Presence() ([]*Room, error) {
	return c.rooms("/presence")
}

func (c *Client) RoomByID(id int) (*Room, error) {
	var wrapper struct {
		Room Room `json:"room"`
	}

	if err := c.connection.get("/room/%d", id, &wrapper); err != nil {
		return nil, err
	}

	room := &wrapper.Room
	room.connection = c.connection

	return room, nil
}

func (c *Client) RoomByName(name string) (*Room, error) {
	if rooms, err := c.Rooms(); err != nil {
		return nil, err
	} else {
		for _, r := range rooms {
			if (string)(r.Name) == name {
				return r, nil
			}
		}
	}

	return nil, nil
}

func (c *Client) Rooms() ([]*Room, error) {
	return c.rooms("/rooms")
}

func (c *Client) Search(term string) ([]*Message, error) {
	var wrapper struct {
		Messages []*Message `json:"messages"`
	}

	if err := c.connection.get("/search", map[string]string{"q": term}, &wrapper); err != nil {
		return nil, err
	}

	msgs := make([]*Message, len(wrapper.Messages))

	for i, m := range wrapper.Messages {
		msgs[i] = m
		msgs[i].connection = c.connection
	}

	return msgs, nil
}

func (c *Client) UserByID(id int) (*User, error) {
	var wrapper struct {
		User User `json:"user"`
	}

	if err := c.connection.get("/users/%d", id, &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) rooms(path string) ([]*Room, error) {
	var wrapper struct {
		Rooms []*Room `json:"rooms"`
	}

	if err := c.connection.get("/rooms", &wrapper); err != nil {
		return nil, err
	}

	rooms := make([]*Room, len(wrapper.Rooms))

	for i, r := range wrapper.Rooms {
		rooms[i] = r
		rooms[i].connection = c.connection
	}

	return rooms, nil
}
