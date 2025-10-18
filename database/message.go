package database

import (
	"context"
	"errors"
	"time"
)

type Seen struct {
	ID         int64     `json:"message_id"`
	ReceiverID string    `json:"receiver_id"`
	SeenAt     time.Time `json:"seen_at"`
}

type MediaType int

const (
	NoMedia MediaType = iota
	Image
	Video
	Audio
	VoiceNote
	VideoNote
)

type Media struct {
	MediaUrl  string `json:"media_url"`
	MediaType string `json:"media_type"` // NoMedia,Image,Video,Audio,Doc
}

type MessageType int

const (
	MessageChat MessageType = iota
	MessageRaction
	MessageInfo
)

type Message struct {
	ID             int64  `json:"message_id"`
	FriendshipID   int64  `json:"friendship_id"` //put groupd id here if group
	SenderUsername string `json:"sender_username"`
	MessageType    string `json:"message_type"` //MessageChat,MessageRaction,MessageInfo
	TextContent    string `json:"text_content"`
	Media          Media  `json:"media"`
	CreatedAt      string `json:"created_at"`
	ModifiedAt     string `json:"modified_at"`
}

func (d *DataRepository) InsertMessage(cxt context.Context, FriendshipID, SenderUsername, MessageType, TextContent string) error {

	query := `INSERT INTO message(friendship_id,sender_username,message_type,text_content,media_url,media_type,modified_at)`

	if MessageType != "MessageChat" && MessageType != "MessageRaction" && MessageType != "MessageInfo" {
		return errors.New("MessageType is invalide")
	}

	_, err := d.db.ExecContext(cxt, query, FriendshipID, SenderUsername, MessageType, TextContent, "", "NoMedia", time.Now())

	return err
}

func (d *DataRepository) InsertMessageMedia(cxt context.Context, FriendshipID, SenderUsername, MessageType, TextContent, MediaUrl, MediaType string) error {

	query := `INSERT INTO message(friendship_id,sender_username,message_type,text_content,media_url,media_type,modified_at)`

	if MessageType != "MessageChat" && MessageType != "MessageRaction" && MessageType != "MessageInfo" {
		return errors.New("MessageType is invalide")
	}

	if MediaType != "NoMedia" && MediaType != "Image" && MediaType != "Audio" && MediaType != "Video" && MediaType != "Doc" {
		return errors.New("MediaType is invalide")
	}

	_, err := d.db.ExecContext(cxt, query, FriendshipID, SenderUsername, MessageType, TextContent, MediaUrl, MediaType, time.Now())

	return err
}

func (d *DataRepository) DeleteMessageById(cxt context.Context, id int) error {

	query := `DELETE FROM message WHERE id = $1`
	_, err := d.db.ExecContext(cxt, query, id)

	return err
}

func (d *DataRepository) GetMessageById(cxt context.Context, id int) (*Message, error) {

	query := `SELECT * FROM message WHERE id = $1`

	row, err := d.db.QueryContext(cxt, query, id)

	if err != nil {
		return nil, err
	}

	row.Next()

	var message Message

	row.Scan(&message.ID, &message.FriendshipID, &message.SenderUsername, &message.TextContent, &message.Media.MediaUrl, &message.Media.MediaType, &message.CreatedAt, &message.ModifiedAt)

	return &message, nil
}

func (d *DataRepository) GetMessages(cxt context.Context, FriendshipID string, page, limit int) (*PaginatedResponse, error) {

	var totalCount int
	offset := (page - 1) * limit
	query := `SELECT * FROM message WHERE friendship_id = $1 LIMIT = $2 OFFSET = $3 ORDER BY created_at DESC`
	queryCount := `SELECT COUNT(*) FROM message WHERE friendship_id = $1`

	counrRow := d.db.QueryRowContext(cxt, queryCount, FriendshipID)

	err := counrRow.Scan(&totalCount)

	if err != nil {
		return nil, err
	}

	row, err := d.db.QueryContext(cxt, query, FriendshipID, limit, offset)

	if err != nil {
		return nil, err
	}

	var messages []Message

	for row.Next() {

		var message Message

		row.Scan(&message.ID, &message.FriendshipID, &message.SenderUsername, &message.TextContent, &message.Media.MediaUrl, &message.Media.MediaType, &message.CreatedAt, &message.ModifiedAt)
		messages = append(messages, message)
	}

	s := PaginatedResponse{
		Data:       messages,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}
	return &s, nil

}

func (d *DataRepository) SearchMessages(cxt context.Context, FriendshipID, search, start_at, end_at string, page, limit int) (*PaginatedResponse, error) {

	query := `SELECT * FROM message WHERE friendship_id = $1 AND text_content LIKE =$2 AND created_at BETWEEN =$3 AND =$4  LIMIT = $5 OFFSET = $6 ORDER BY created_at DESC`
	queryCount := `SELECT COUNT(*) FROM message WHERE friendship_id = $1 AND text_content LIKE =$2 AND created_at BETWEEN =$3 AND =$4`

	var totalCount int
	offset := (page - 1) * limit

	counrRow := d.db.QueryRowContext(cxt, queryCount, FriendshipID, search, start_at, end_at)

	err := counrRow.Scan(&totalCount)

	if err != nil {
		return nil, err
	}

	row, err := d.db.QueryContext(cxt, query, FriendshipID, search, start_at, end_at, limit, offset)

	if err != nil {
		return nil, err
	}

	var messages []Message

	for row.Next() {

		var message Message

		row.Scan(&message.ID, &message.FriendshipID, &message.SenderUsername, &message.TextContent, &message.Media.MediaUrl, &message.Media.MediaType, &message.CreatedAt, &message.ModifiedAt)
		messages = append(messages, message)
	}

	s := PaginatedResponse{
		Data:       messages,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}
	
	return &s, nil
}

func (d *DataRepository) UpdateMessage(cxt context.Context, id int, updatedText string) error {

	query := `UPDATE message SET text_content = $1 WHERE id = $2`

	_, err := d.db.ExecContext(cxt, query, updatedText, id)

	return err
}
