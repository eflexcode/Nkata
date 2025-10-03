package database

type Firends struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	FirendID  int64  `json:"friend_id"`
	CreatedAt string `json:"created_at"`
}

type FriendRequest struct {
	ID        int64  `json:"id"`
	SentBy    int64  `json:"sent_by"`
	SentTo    int64  `json:"sent_to"`
	Accepted  bool   `json:"accepted"`
	CreatedAt string `json:"created_at"`
}

