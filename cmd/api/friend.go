package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type FriendRequestPayload struct {
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
}

type RespondFriendRequestPayload struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}

// @Summary Send friend request
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body FriendRequestPayload true "id and friend id"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Router /v1/firendship/send-friend-request [post]
func (apiService *ApiService) SendFriendRequest(w http.ResponseWriter, r *http.Request) {

	var payload FriendRequestPayload

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	boolean := apiService.database.HasSentMeRequest(ctx, payload.FriendId, payload.UserId)

	if boolean {
		conflict(w, r, errors.New("user already sent you a friend request"))
		return
	}

	duplicate := apiService.database.CheckDuplicateRequest(ctx, payload.UserId, payload.FriendId)

	if duplicate {
		conflict(w, r, errors.New("you already a friend request to this user"))
		return
	}

	err := apiService.database.InsertFriendRequest(ctx, payload.FriendId, payload.UserId)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "Friend request sent successfully",
	}

	writeJson(w, 200, s)

}

// @Summary Responed to friend request
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body RespondFriendRequestPayload true "id and status"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Router /v1/firendship/responed-friend-request [post]
func (api *ApiService) RespondFriendRequest(w http.ResponseWriter, r *http.Request) {

	var payload RespondFriendRequestPayload
	ctx := r.Context()

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	frendRequest, err := api.database.GetFriendRequestById(ctx, payload.Id)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	if payload.Status == "accepted" {

		err := api.database.UpdateFriendRequestStatus(ctx, payload.Status, payload.Id, frendRequest.SentTo)

		if err != nil {
			internalServer(w, r, err)
			return
		}

		var friendship_id = uuid.New().String()

		err1 := api.database.InsertFriendship(ctx, frendRequest.SentBy, frendRequest.SentTo, friendship_id)

		err = api.database.InsertFriendship(ctx, frendRequest.SentTo, frendRequest.SentBy, friendship_id)

		if err != nil || err1 != nil {
			internalServer(w, r, err)
			return
		}

		s := StandardResponse{
			Status:  200,
			Message: "firend request accepted successfully",
		}

		writeJson(w, 200, s)
		return

	} else if payload.Status == "rejected" {

		err := api.database.UpdateFriendRequestStatus(ctx, payload.Status, payload.Id, frendRequest.SentTo)

		if err != nil {
			internalServer(w, r, err)
			return
		}

		s := StandardResponse{
			Status:  200,
			Message: "firend request rejected successfully",
		}

		writeJson(w, 200, s)
		return

	} else {
		badRequest(w, r, errors.New("invalid status type: status can either be accepted or rejected"))
		return
	}

}
