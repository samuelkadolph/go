package campfire

type Upload struct {
	ByteSize    int    `json:"byte_size"`
	ContentType string `json:"content-type"`
	CreatedAt   string `json:"created_at"`
	FullURL     string `json:"full_url"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	RoomID      int64  `json:"room_id"`
	UserID      int64  `json:"user_id"`
}
