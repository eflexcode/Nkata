package database

type NotificationType int

const (
	MessageChatType NotificationType = iota
	InfoType
)

type Notification struct {
	ID        int64  `json:"id"`
	UserID int64 `json:"user_id"`
	Title string `json:"title"`
	Message string `json:"message"`
	UrlImg string `json:"img_url"`
	NotificationType NotificationType `json:"notification_type"`
	CreatedAt string `json:"created_at"`
}
