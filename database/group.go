package database

type Group struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PicUrl      string `json:"pic_url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
}

type GroupMember struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"room_id"`
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"` // admin or member
	CreatedAt string `json:"created_at"`
} // once u add a user to a group they get added here and in friendship

type GroupMemberRemoved struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"room_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
