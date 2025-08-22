package database

type Firend struct {
	ID        int64  `json:"id"`
	OwnerID   int64  `json:"owner_id"`
	FirendID  int64  `json:"friend_id"`
	CreatedAt string `json:"created_at"`
}
