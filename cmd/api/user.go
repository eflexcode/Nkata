package api

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	ImageUrl   string `json:"image_url"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}