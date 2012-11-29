package campfire

type Upload struct {
	ByteSize    int    `json:"byte_size"`
	ContentType string `json:"content-type"`
	CreatedAt   string `json:"created_at"`
	FullURL     string `json:"full_url"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	RoomID      int    `json:"room_id"`
	UserID      int    `json:"user_id"`
}
