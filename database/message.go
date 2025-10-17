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
	MediaType string `json:"media_type"` // NoMedia,Image,Video,Audio
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

	_, err := d.db.ExecContext(cxt, query, FriendshipID, SenderUsername, MessageType, TextContent,"","NoMedia", time.Now())

	return err
}

func (d *DataRepository) InsertMessageMedia(cxt context.Context, FriendshipID, SenderUsername, MessageType, TextContent,MediaUrl,MediaType string) error {

	query := `INSERT INTO message(friendship_id,sender_username,message_type,text_content,media_url,media_type,modified_at)`

	if MessageType != "MessageChat" && MessageType != "MessageRaction" && MessageType != "MessageInfo" {
		return errors.New("MessageType is invalide")
	}

	if MediaType != "NoMedia" && MediaType != "Image" && MediaType != "Audio" && MediaType != "Video"{
		return errors.New("MediaType is invalide")
	}

	_, err := d.db.ExecContext(cxt, query, FriendshipID, SenderUsername, MessageType, TextContent,MediaUrl,MediaType, time.Now())

	return err
}