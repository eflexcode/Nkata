package database

import (
	"context"
	"time"
)

type NotificationType int

const (
	MessageChatType NotificationType = iota
	InfoType
)

type Notification struct {
	ID       int64  `json:"id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	UrlImg   string `json:"img_url"`
	// NotificationType NotificationType `json:"notification_type"`
	CreatedAt string `json:"created_at"`
}

func (r *DataRepository) InsertNotification(ctx context.Context, userId, username, title, message, url string) (Notification, error) {

	query := `INSERT INTO notification(user_id,username,title,message,img_url,created_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,username,title,message,img_url,created_at`
	var notification Notification
	err := r.db.QueryRowContext(ctx, query, userId, username, title, message, url, time.Now().String()).Scan(&notification.ID, &notification.UserID, &notification.Username, &notification.Title, &notification.Message, &notification.UrlImg, &notification.CreatedAt)

	return notification, err
}

func (r *DataRepository) GetNotifications(ctx context.Context, userId string, page, limit int64) (*PaginatedResponse, error) {

	var request []Notification

	query := `SELECT * FROM notification WHERE user_id = $1 LIMIT $2 OFFSET $3`
	queryCount := `SELECT COUNT(*) FROM notification WHERE user_id = $1`

	var totalCount int

	cRow := r.db.QueryRowContext(ctx, queryCount, userId)

	err := cRow.Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	row, err := r.db.Query(query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {

		notification := Notification{}
		err := row.Scan(&notification.ID, &notification.UserID, &notification.Username, &notification.Title, &notification.Message, &notification.UrlImg, &notification.CreatedAt)

		if err != nil {
			return nil, err
		}

		request = append(request, notification)

	}

	p := PaginatedResponse{
		Data:       request,
		TotalCount: totalCount,
		Page:       int(page),
		Limit:      int(limit),
	}

	return &p, nil
}
