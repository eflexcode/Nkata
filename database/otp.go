package database

type Otp struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Token      int32  `json:"token"`
	Email      string `json:"email"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}
