package database

import "time"

type Seen struct {
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
	MediaUrl  string    `json:"media_url"`
	MediaType MediaType `json:"media_type"`
}

type MessageType int

const (
	MessageChat MessageType = iota
	MessageRaction
)

type Message struct {
	ID          int64       `json:"message_id"`
	RoomID      int64       `json:"room_id"`
	SenderID    string      `json:"sender_id"`
	MessageChat string      `json:"message_chat"`
	MessageType MessageType `json:"message_type"`
	Media       Media       `json:"media"`
	Seen        Seen        `json:"seen"`
	CreatedAt   string      `json:"created_at"`
	ModifiedAt  string      `json:"modified_at"`
}
