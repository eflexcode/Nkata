package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"main/database"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
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
	FriendshipID   string `json:"friendship_id"` //put grouped id here if group
	SenderUsername string `json:"sender_username"`
	MessageType    string `json:"message_type"` //MessageChat,MessageRaction,MessageInfo
	TextContent    string `json:"text_content"`
	Media          Media  `json:"media"`
}

type MessageNotInDb struct {
	MessageId string `json:"message_id"`
	Info      string `json:"info"`
}

type WsConnectionManager struct {
	sync.RWMutex
	Connections map[string]*websocket.Conn //the sting is the user_id
}

// var g = map[string] map[string]

// var connections map[string]*websocket.Conn
// var conns = WsConnectionManager{
// 	Connections: connections,
// }

type WsMessage struct {
	FriendshipID string `json:"friendship_id"` //put grouped id here if group
	MessageType  string `json:"message_type"`  //MessageChat,MessageRaction,MessageInfo
	TextContent  string `json:"text_content"`
	Media        Media  `json:"media"`
}

type WsNotification struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}

type ProfileUpdate struct {
	DisplayName      string `json:"display_name"`
	ProfilePicBinary string `json:"profile_pic_binary"`
	Online           string `json:"online"` //note i send date if date is greater than 1 min off line
}

type ProfileUpdateReturnPayload struct {
	DisplayName           string `json:"display_name"`
	ProfilePicDownloadUrl string `json:"profile_pic_download_url"`
	Online                string `json:"online"` //note i send date if date is greater than 1 min off line
}

type WsFriendRequestPayload struct {
	Type           string                      `json:"type"` //request,response
	FriendUsername string                      `json:"friend_username"`
	FriendUserId   string                      `json:"friend_user_id"`
	Respond        RespondFriendRequestPayload `json:"respond"`
	//in another
}

type WsPayload struct {
	UserId                 string                 `json:"user_id"`
	PayloadType            string                 `json:"payload_type"` //notification,message,friendRequest,myFriends,profileUpdate
	MessagePayload         MessagePayload         `json:"message"`
	WsNotification         WsNotification         `json:"notification"`
	ProfileUpdate          ProfileUpdate          `json:"profile_update"`
	WsFriendRequestPayload WsFriendRequestPayload `json:"friend_request"`
}

type WsStructure struct {
	sync.RWMutex
	UserId      string          `json:"user_id"`
	Conn        *websocket.Conn `json:"conn"`
	FriendShips []string        `json:"friendships"`
}

type WsHandShackPayload struct {
	FriendShips []string `json:"friendships"`
}

var conns []*WsStructure

func addConnection(userId string, conn *websocket.Conn, friendShips []string) {

	w := WsStructure{
		UserId:      userId,
		Conn:        conn,
		FriendShips: friendShips,
	}
	_ = append(conns, &w)
	// conns.Lock()
	// conns.Connections[userId] = conn
	// conns.Unlock()
}

func deleteConnection(userId string) {

	for conPosition := range conns {
		var con = conns[conPosition]
		con.Conn.Close()
		conns = slices.Delete(conns, conPosition, conPosition+1) //starts delete at current position and end at current position plus 1 the plus one would not get deleted: if it is plus 2 the 2 would not get deleted
		return
	}

	// conns.Lock()
	// defer conns.Unlock()

	// if conn, exist := conns.Connections[userId]; exist {
	// 	conn.Close()
	// 	delete(conns.Connections, userId)
	// }
}

// @Summary General ws handler
// @Description Responds with json
// @Tags ws
// @Produce json
// @Success 200 {object} database.Message
// @Produce octet-stream
// @Success 200 {file} file
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/general-authenticated/ws/{user_id} [post]
func (api *ApiService) GeneralWsHandler(w http.ResponseWriter, r *http.Request) {

	var wsHandShackPayload WsHandShackPayload
	if err := readJson(w, r, &wsHandShackPayload); err != nil {
		badRequest(w, r, errors.New("json payload cannot be decoded"))
		return
	}

	conn, err := upgradeConn.Upgrade(w, r, nil)
	if err != nil {
		internalServer(w, r, errors.New("failed to upgrade connection to ws"))
		return
	}

	user_id := chi.URLParam(r, "user_id")

	addConnection(user_id, conn, wsHandShackPayload.FriendShips)
	defer deleteConnection(user_id)

	username, err := getUsernameFromCtx(r.Context())
	if err != nil {
		internalServer(w, r, err)
		return
	}

	ctx := r.Context()

	for {

		messageType, data, err := conn.ReadMessage()

		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch messageType {
		//Note todo send files like audio and img
		case websocket.TextMessage:
			var payload WsPayload
			if err := json.Unmarshal(data, &payload); err != nil {
				log.Printf("Invalid payload sent: %v", err)
				return
			}

			var userId = payload.UserId
			var payloadType = payload.PayloadType

			if payloadType == "message" {

				var messageP = payload.MessagePayload
				var messageId = uuid.New().String()

				var now = time.Now()

				message := database.Message{
					MessageID:      messageId,
					FriendshipID:   messageP.FriendshipID,
					SenderUsername: messageP.SenderUsername,
					MessageType:    messageP.MessageType,
					TextContent:    messageP.TextContent,
					Media:          database.Media(messageP.Media),
					CreatedAt:      now.String(),
					ModifiedAt:     now.String(),
				}

				err = api.database.InsertMessage(r.Context(), messageId, message.FriendshipID, message.SenderUsername, message.MessageType, message.TextContent, now)

				if err != nil {

					info := MessageNotInDb{
						MessageId: messageId,
						Info:      "failed to insert with this id in db please remove",
					}

					log.Printf("%v", info)

					byteResponse, err := json.Marshal(message)

					if err != nil {
						log.Printf("failed to parse response to byte: %v", err)
						return
					}

					// search all  connections/ connected user
					for conPosition := range conns {
						wsConn := conns[conPosition]
						if wsConn.UserId == userId {
							//my user id return/continue
							// note once sent it enters user(sender chat first sqlite in sender device) so no need to send back to the person yet might need to for seen update
							continue
						}
						//indididual user friendships
						friendships := wsConn.FriendShips
						for fri := range friendships {
							//each friend
							friend := friendships[fri]
							//check if friendship id equals so therefore we are friends or in same group
							if messageP.FriendshipID == friend {
								//publish in that friend connection
								if err := wsConn.Conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
									log.Printf("socket publish failed: %g", err)
								}
								//TODO might send back to self to for seen
								var splitedFriendshipId = strings.Split(friend, "_")
								if splitedFriendshipId[0] == "chat" {
									continue
								}
							}
						}
					}
					//not done yet
					// if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
					// 	log.Panicf("socket publish failed: %g", err)
					// }

				}

			} else if payloadType == "notification" {

			} else if payloadType == "friendRequest" {

				var requestProcesses = payload.WsFriendRequestPayload

				if requestProcesses.Type == "request" {

					if username == requestProcesses.FriendUsername {
						forbidden(w, r, errors.New("user cannot send friend request to self"))
						return
					}

					userExist := api.database.CheackUsernameAvailability(ctx, requestProcesses.FriendUsername)

					if !userExist {
						notFound(w, r, errors.New("no user found with username: "+requestProcesses.FriendUsername))
						continue
					}

					boolean := api.database.HasSentMeRequest(ctx, requestProcesses.FriendUsername, username)

					if boolean {
						conflict(w, r, errors.New("user already sent you a friend request"))
						continue
					}

					duplicate := api.database.CheckDuplicateRequest(ctx, username, requestProcesses.FriendUsername)

					if duplicate {
						conflict(w, r, errors.New("you already a friend request to this user"))
						continue
					}

					fRequest, err := api.database.InsertFriendRequest(ctx, requestProcesses.FriendUsername, username)

					if err != nil {
						internalServer(w, r, err)
						continue
					}

					byteResponse, err := json.Marshal(fRequest)
					if err != nil {
						log.Printf("failed to parse response to byte: %v", err)
						return
					}

					//write to user name
					for conPosition := range conns {
						wsConn := conns[conPosition]
						if wsConn.UserId == userId {
							//my user id return/continue
							// note once sent it enters user(sender chat first sqlite in sender device) so no need to send back to the person yet might need to for seen update
							continue
						}

						if wsConn.UserId == requestProcesses.FriendUserId {
							if err := wsConn.Conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
								log.Printf("socket publish failed: %g", err)
							}
						}

					}
				} else if requestProcesses.Type == "response" {

					respond := requestProcesses.Respond

					friendRequest, err := api.database.GetFriendRequestById(ctx, respond.Id)

					if err != nil {

						if err.Error() == "sql: no rows in result set" {
							// notFound(w, r, errors.New("no friend request found with id: "+strconv.Itoa(int(respond.Id))))
							continue
						}

						// internalServer(w, r, err)
						continue
					}

					if respond.Status == "accepted" {

						err := api.database.UpdateFriendRequestStatus(ctx, respond.Status, respond.Id)

						if err != nil {
							// internalServer(w, r, err)
							continue
						}

						var friendship_id = "chat_" + uuid.New().String()

						err1 := api.database.InsertFriendship(ctx, friendRequest.SentBy, friendRequest.SentTo, friendship_id)

						err = api.database.InsertFriendship(ctx, friendRequest.SentTo, friendRequest.SentBy, friendship_id)

						if err != nil || err1 != nil {
							// internalServer(w, r, err)
							continue
						}

						// s := StandardResponse{
						// 	Status:  200,
						// 	Message: "friend request accepted successfully",
						// }

						// writeJson(w, 200, s)
						continue

					} else if respond.Status == "rejected" {

						err := api.database.UpdateFriendRequestStatus(ctx, respond.Status, respond.Id)

						if err != nil {
							// internalServer(w, r, err)
							continue
						}
						//send notification maybe
						err = api.database.DeleteFriendRequest(ctx, respond.Id)
						if err != nil {
							// internalServer(w, r, err)
							continue
						}

						// s := StandardResponse{
						// 	Status:  200,
						// 	Message: "friend request rejected successfully",
						// }

						// writeJson(w, 200, s)
						continue

					}

				}

			} else if payloadType == "myFriends" {

				response, err := api.database.GetFriends(ctx, username)

				if err != nil {
					internalServer(w, r, err)
					continue
				}

				friendsResponse, err := json.Marshal(response)

				//send to me
				if err := conn.WriteMessage(websocket.TextMessage, friendsResponse); err != nil {
					log.Printf("socket publish failed: %g", err)
				}
			} else if payloadType == "profileUpdate" {
				profileUpdate := payload.ProfileUpdate
				var profileUpdateR ProfileUpdateReturnPayload
				if profileUpdate.ProfilePicBinary != "" {
					var extention = ".png"
					var fileBinary []byte

					src := []byte(profileUpdate.ProfilePicBinary)
					fileBinary = make([]byte, len(src))
					// for dbytes, char := range profileUpdate.ProfilePicBinary {
					// 	var g = string(char)
					// 	fileBinary = append(fileBinary, byte(g))
					// }

					currentTime := time.Now().UnixMilli()

					currentTimeString := strconv.Itoa(int(currentTime)) + extention

					destinationFile, err := os.Create("/home/ifeanyi/nkata_storage/chat_storage/" + currentTimeString)

					if err != nil {
						internalServer(w, r, err)
						continue
					}

					defer destinationFile.Close()

					i, err := destinationFile.Write([]byte(fileBinary))
					if err != nil {
						internalServer(w, r, err)
						continue
					}

					if i == 0 {
						internalServer(w, r, errors.New("failed to write file sent"))
						continue
					}

					url := "http://localhost:5557/v1/media/profiles/" + currentTimeString
					err = api.database.UpdateProfilePicUrl(ctx, username, url)

					if err != nil {
						internalServer(w, r, err)
						continue
					}
					profileUpdateR = ProfileUpdateReturnPayload{
						ProfilePicDownloadUrl: url,
						DisplayName:           profileUpdate.DisplayName,
						Online:                profileUpdate.Online,
					}

				} else {
					profileUpdateR = ProfileUpdateReturnPayload{
						ProfilePicDownloadUrl: "url",
						DisplayName:           profileUpdate.DisplayName,
						Online:                profileUpdate.Online,
					}
				}

				byteResponse, err := json.Marshal(profileUpdateR)
				if err != nil {
					log.Printf("failed to parse response to byte: %v", err)
					return
				}

				// search all  connections/ connected user
				for conPosition := range conns {
					wsConn := conns[conPosition]
					if wsConn.UserId == userId {
						//my user id return/continue
						// note once sent it enters user(sender chat first sqlite in sender device) so no need to send back to the person yet might need to for seen update
						continue
					}
					//indididual user friendships
					friendships := wsConn.FriendShips
					for fri := range friendships {
						//each friend
						_ = friendships[fri]
						//check if friendship id equals so therefore we are friends or in same group
						// if messageP.FriendshipID == friend {
						//publish in that friend connection
						if err := wsConn.Conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
							log.Printf("socket publish failed: %g", err)
						}

						// }
					}
				}
			}

		case websocket.BinaryMessage:
			// var payload WsPayload
			// if err := json.Unmarshal(data, &payload); err != nil {
			// 	log.Printf("Invalid payload sent: %v", err)
			// 	return
			// }

			// var userId = payload.UserId
			// var payloadType = payload.PayloadType
			// if payloadType == "message" {

			// 	fileTypeHttp := http.DetectContentType(data)

			// 	var fileTypeHttpSplit = strings.Split(fileTypeHttp, "/")

			// 	fileExtention := "." + fileTypeHttpSplit[1]

			// 	currentTime := time.Now().UnixMilli()

			// 	currentTimeString := strconv.Itoa(int(currentTime)) + fileExtention

			// 	destinationFile, err := os.Create("/home/ifeanyi/nkata_storage/chat_storage/" + currentTimeString)

			// 	if err != nil {
			// 		internalServer(w, r, err)
			// 		return
			// 	}

			// 	defer destinationFile.Close()

			// 	i, err := destinationFile.Write(data)
			// 	if err != nil {
			// 		internalServer(w, r, err)
			// 		return
			// 	}

			// 	if i == 0 {
			// 		internalServer(w, r, errors.New("failed to write file sent"))
			// 		return
			// 	}
			// 	ctx := r.Context()
			// 	username, err := getUsernameFromCtx(ctx)
			// 	if err != nil {
			// 		internalServer(w, r, err)
			// 		return
			// 	}
			// 	now := time.Now()
			// 		var messageId = uuid.New().String()
			// 	message := database.Message{
			// 		MessageID:      messageId,
			// 		FriendshipID:   friendshipId,
			// 		SenderUsername: username,
			// 		MessageType:    "MessageChat",
			// 		Media:          database.Media{MediaUrl: url, MediaType: fileExtention},
			// 		CreatedAt:      now.String(),
			// 		ModifiedAt:     now.String(),
			// 	}

			// 	byteResponse, err := json.Marshal(message)

			// 	if err != nil {
			// 		log.Printf("failed to parse response to byte: %v", err)
			// 		return
			// 	}

			// 	if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
			// 		log.Panicf("socket publish failed: %t", err)
			// 	}

			// 	err = api.database.InsertMessageMedia(ctx, messageId, friendshipId, username, "MessageChat", url, fileExtention, now)

			// 	if err != nil {

			// 		info := MessageNotInDb{
			// 			MessageId: messageId,
			// 			Info:      "failed to insert with this id in db please remove",
			// 		}

			// 		byteResponse, err := json.Marshal(info)

			// 		if err != nil {
			// 			log.Printf("failed to parse response 2 to byte: %v", err)
			// 			return
			// 		}

			// 		if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
			// 			log.Panicf("socket publish failed: %g", err)
			// 		}

			// }
			// } else {

		}
	}
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// @Summary Message ws connection
// @Description Responds with json
// @Tags Message
// @Produce json
// @Success 200 {object} database.Message
// @Produce octet-stream
// @Success 200 {file} file
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/message/ws/{friendship_id} [get]
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
				log.Printf("%v", info)
				// byteResponse, err := json.Marshal(info)

				// if err != nil {
				// 	log.Printf("failed to parse response 2 to byte: %v", err)
				// 	return
				// }

				// if err := conn.WriteMessage(websocket.TextMessage, byteResponse); err != nil {
				// 	log.Panicf("socket publish failed: %g", err)
				// }

			}

		case websocket.BinaryMessage:

			fileTypeHttp := http.DetectContentType(data)

			var fileTypeHttpSplit = strings.Split(fileTypeHttp, "/")

			fileExtention := "." + fileTypeHttpSplit[1]

			currentTime := time.Now().UnixMilli()

			currentTimeString := strconv.Itoa(int(currentTime)) + fileExtention

			destinationFile, err := os.Create("/home/ifeanyi/nkata_storage/chat_storage/" + currentTimeString)

			if err != nil {
				internalServer(w, r, err)
				return
			}

			defer destinationFile.Close()

			i, err := destinationFile.Write(data)
			if err != nil {
				internalServer(w, r, err)
				return
			}

			if i == 0 {
				internalServer(w, r, errors.New("failed to write file sent"))
				return
			}
			ctx := r.Context()
			username, err := getUsernameFromCtx(ctx)
			if err != nil {
				internalServer(w, r, err)
				return
			}
			now := time.Now()
			friendshipId := chi.URLParam(r, "friendship_id")

			var messageId = uuid.New().String()

			//broadcast message before insert for latency

			url := "localhost:5557/v1/media/chat/" + currentTimeString

			message := database.Message{
				MessageID:      messageId,
				FriendshipID:   friendshipId,
				SenderUsername: username,
				MessageType:    "MessageChat",
				Media:          database.Media{MediaUrl: url, MediaType: fileExtention},
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

			err = api.database.InsertMessageMedia(ctx, messageId, friendshipId, username, "MessageChat", url, fileExtention, now)

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

		default:
			log.Printf("cannot determine incoming socket data type: %v", err)
		}

	}

}

// @Summary Get Messages with message_id
// @Description Responds with json
// @Tags Message
// @Param message_id path string true "message_id"
// @Produce json
// @Success 200 {object} database.Message
// @Failure 404 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/message/get/{message_id} [get]
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

// @Summary Download chat Pic
// @Description Responds with json
// @Tags Media
// @Param img_name path string true "file name"
// @Produce octet-stream
// @Success 200 {file} file
// @Failure 404 {object} errorslope
// @Router /v1/media/chat/{img_name} [get]
func (api *ApiService) LoadMessagefile(w http.ResponseWriter, r *http.Request) {

	filename := chi.URLParam(r, "img_name")
	url := "/home/ifeanyi/nkata_storage/chat_storage/" + filename
	file, err := os.Open(url)

	if err != nil {
		notFound(w, r, errors.New("the system cannot find the file specified"))
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename= "+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeContent(w, r, filename, time.Time{}, file)
}

// @Summary Get Messages with friendship id
// @Description Responds with json
// @Tags Message
// @Param friendship_id path string true "friendship id"
// @Param page query string true "current page if any"
// @Param limit query string true "page max lenght if any"
// @Produce json
// @Success 200 {object} database.PaginatedResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/message/get-messages/{friendship_id} [get]
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

// @Summary Delete Messages with message_id
// @Description Responds with json
// @Tags Message
// @Param message_id path string true "message_id"
// @Produce json
// @Success 200 {object} StandardResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/message/delete/{message_id} [delete]
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

// @Summary Search Messages with friendship_id
// @Description Responds with json
// @Tags Message
// @Param friendship_id path string true "friendship id"
// @Param page query string true "current page if any"
// @Param limit query string true "page max lenght if any"
// @Param q query string true "query text"
// @Param start_at query string true "start date"
// @Param end_at query string true "end date"
// @Produce json
// @Success 200 {object} database.PaginatedResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/message/search-messages/{friendship_id} [get]
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
