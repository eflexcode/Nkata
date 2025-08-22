package database

type RoomType int

const (
	Private RoomType = iota
	Group
)

type Room struct {
	ID              int64    `json:"id"`
	RoomType        RoomType `json:"room_type"`
	RoomName        string   `json:"room_name"` //for group
	RoomPicUrl      string   `json:"room_pic_url"`
	RoomDescription string   `json:"room_description"`
}

type RoomMembers struct {
	ID     int64 `json:"id"`
	RoomID int64 `json:"room_id"`
	UserID int64 `json:"user_id"`
}
