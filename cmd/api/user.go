package api

type Role int

const (
	NUser Role = iota
	Admin
)

type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	ImageUrl       string `json:"image_url"`
	Bio            string `json:"bio"`
	IsOnline       bool   `json:"is_online"`
	FollowersCount int64  `json:"followers_count"`
	GroupsCount    int16  `json:"groupsCount"`
	Role           Role   `json:"role"`
	CreatedAt      string `json:"created_at"`
	ModifiedAt     string `json:"modified_at"`
}
