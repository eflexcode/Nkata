package api

type Room struct {
	ID int64 `json:"id"`
	FriendID        int64  `json:"friend_id"`
}
