package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"main/database"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgradeConn = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

type MediaType int

const (
	NoMedia MediaType = iota
	Image
	Video
	Audio
)

type Media struct {
	MediaUrl  string `json:"media_url"`
	MediaType string `json:"media_type"` // NoMedia,Image,Video,Audio,Doc
}

type MessagePayload struct {
	FriendshipID   string `json:"friendship_id"` //put groupd id here if group
	SenderUsername string `json:"sender_username"`
	MessageType    string `json:"message_type"` //MessageChat,MessageRaction,MessageInfo
	TextContent    string `json:"text_content"`
	Media          Media  `json:"media"`
}

type MessageNotInDb struct {
	MessageId string `json:"message_id"`
	Info      string `json:"info"`
}

func (api *ApiService) MessageWsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgradeConn.Upgrade(w, r, nil)

	if err != nil {
		internalServer(w, r, errors.New("failed to upgrade connection to ws"))
		return
	}

	defer conn.Close()

	for {

		messageType, data, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch messageType {

		case websocket.TextMessage:

			var messagePayload MessagePayload

			if err := json.Unmarshal(data, &messagePayload); err != nil {
				log.Printf("Invalid payload sent: %v", err)
				return
			}

			var messageId = uuid.New().String()

			//broadcast message before insert for latency

			now := time.Now()

			message := database.Message{
				MessageID:      messageId,
				FriendshipID:   messagePayload.FriendshipID,
				SenderUsername: messagePayload.SenderUsername,
				MessageType:    messagePayload.MessageType,
				TextContent:    messagePayload.TextContent,
				Media:          database.Media(messagePayload.Media),
				CreatedAt:      now.String(),
				ModifiedAt:     now.String(),
			}

			byteResponse, err := json.Marshal(message)

			if err != nil {
				log.Printf("failed to parse response to byte: %v", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
				log.Panicf("socket publish failed: %t", err)
			}

			err = api.database.InsertMessage(r.Context(), messageId, messagePayload.FriendshipID, messagePayload.SenderUsername, message.MessageType, message.TextContent, now)

			if err != nil {

				info := MessageNotInDb{
					MessageId: messageId,
					Info:      "failed to insert with this id in db please remove",
				}

				byteResponse, err := json.Marshal(info)

				if err != nil {
					log.Printf("failed to parse response 2 to byte: %v", err)
					return
				}

				if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
					log.Panicf("socket publish failed: %g", err)
				}

			}

		case websocket.BinaryMessage:

			strconv.Itoa(messageType).

			io.Reader.Read(data)

			currentTime := time.Now().UnixMilli()

			currentTimeString := strconv.Itoa(int(currentTime)) + filepath.Ext(fileHeader.Filename)

			destinationFile, err := os.Create("/home/ifeanyi/nkata_storage/chat_storage/" + currentTimeString)

			if err != nil {
				internalServer(w, r, err)
				return
			}

			defer destinationFile.Close()

			_, err = file.Seek(0, io.SeekStart)

			if err != nil {
				internalServer(w, r, err)
				return
			}

			_, err = io.Copy(destinationFile, file)

			if err != nil {
				internalServer(w, r, err)
				return
			}

		default:
			log.Printf("cannot determin incoming socket data type: %v", err)
		}

	}

}

func (api *ApiService) GetMessageByMessageId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "message_id")

	ctx := r.Context()

	message, err := api.database.GetMessageById(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			notFound(w, r, errors.New("no message found with message_id: "+id))
			return
		}
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, message)

}

func (api *ApiService) GetMessages(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "friendship_id")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	ctx := r.Context()

	pageInt, err1 := strconv.Atoi(page)
	limitInt, err2 := strconv.Atoi(limit)

	if err1 != nil || err2 != nil {
		badRequest(w, r, errors.New("page or limit might not be a number"))
		return
	}

	result, err := api.database.GetMessages(ctx, id, pageInt, limitInt)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, result)
}

func (api *ApiService) DeleteMessageByMessageId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "message_id")

	ctx := r.Context()

	err := api.database.DeleteMessageById(ctx, id)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "Message deleted successfully",
	}

	writeJson(w, http.StatusOK, s)

}

func (api *ApiService) SearchMessages(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "friendship_id")
	searchText := r.URL.Query().Get("q")
	startDate := r.URL.Query().Get("start_at")
	endDate := r.URL.Query().Get("end_at")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	ctx := r.Context()

	pageInt, err1 := strconv.Atoi(page)
	limitInt, err2 := strconv.Atoi(limit)

	if err1 != nil || err2 != nil {
		badRequest(w, r, errors.New("page or limit might not be a number"))
		return
	}

	result, err := api.database.SearchMessages(ctx, id, searchText, startDate, endDate, pageInt, limitInt)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, result)

}
