package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FriendRequestPayload struct {
	FriendUsername string `json:"friend_username"`
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
// @Router /v1/firendship/request/send [post]
func (apiService *ApiService) SendFriendRequest(w http.ResponseWriter, r *http.Request) {

	var payload FriendRequestPayload

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	if username == payload.FriendUsername{
		forbidden(w,r,errors.New("user cannot send friend request to self"))
		return
	}

	userExist := apiService.database.CheackUsernameAvailability(ctx, payload.FriendUsername)

	if !userExist {
		notFound(w, r, errors.New("no user found with username: "+payload.FriendUsername))
		return
	}

	boolean := apiService.database.HasSentMeRequest(ctx, payload.FriendUsername, username)

	if boolean {
		conflict(w, r, errors.New("user already sent you a friend request"))
		return
	}

	duplicate := apiService.database.CheckDuplicateRequest(ctx, username, payload.FriendUsername)

	if duplicate {
		conflict(w, r, errors.New("you already a friend request to this user"))
		return
	}

	err = apiService.database.InsertFriendRequest(ctx, payload.FriendUsername, username)

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
// @Router /v1/firendship/request/responed [post]
func (api *ApiService) RespondFriendRequest(w http.ResponseWriter, r *http.Request) {

	var payload RespondFriendRequestPayload
	ctx := r.Context()

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	frendRequest, err := api.database.GetFriendRequestById(ctx, payload.Id)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			notFound(w, r, errors.New("no friend request found with id: "+strconv.Itoa(int(payload.Id))))
			return
		}

		internalServer(w, r, err)
		return
	}

	if payload.Status == "accepted" {

		err := api.database.UpdateFriendRequestStatus(ctx, payload.Status, payload.Id)

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

		err := api.database.UpdateFriendRequestStatus(ctx, payload.Status, payload.Id)

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

// @Summary Deleted Request
// @Description Responds with json
// @Tags Friendship
// @Param id  path string true "friend request id"
// @Produce json
// @Success 200 {file} file
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/firendship/request/delete/{id}  [delete]
func (api *ApiService) DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	ctx := r.Context()
	id := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		badRequest(w, r, errors.New("Id is not a number"))
		return
	}

	request, err := api.database.GetFriendRequestById(ctx, int64(idInt))

	if err != nil {

		if err.Error() == "sql: no rows in result set"  || err.Error() == "sql: Rows are closed"{
			notFound(w, r, errors.New("no request found with id: "+id))
			return
		}

		internalServer(w, r, err)
		return
	}

	if request.SentBy != username {
		unauthorized(w, r, errors.New("user does not have permision to perform this action"))
		return
	}

	err = api.database.DeleteFriendRequest(ctx, int64(idInt))
	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "firend request deleted successfully",
	}

	writeJson(w, 200, s)
}

// @Summary Get Friend Request Sent
// @Description Responds with json
// @Tags Friendship
// @Param page  query string true "page to get"
// @Param limit  query string true "page limit"
// @Produce json
// @Success 200 {object} database.PaginatedResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/firendship/request/get-sent  [get]
func (api *ApiService) GetFriendRequestSent(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	ctx := r.Context()

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	limitInt, errp := strconv.Atoi(limit)

	if err != nil || errp != nil {
		badRequest(w, r, errors.New("page or limitis not a number"))
		return
	}

	response, err := api.database.GetFriendRequestSentBy(ctx, username, int64(pageInt), int64(limitInt))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, 200, response)
}

// @Summary Get Friend Request Recieved
// @Description Responds with json
// @Tags Friendship
// @Param page  query string true "page to get"
// @Param limit  query string true "page limit"
// @Produce json
// @Success 200 {object} database.PaginatedResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/firendship/request/get-recieved  [get]
func (api *ApiService) GetFriendRequestRecieved(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	ctx := r.Context()

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	limitInt, errp := strconv.Atoi(limit)

	if err != nil || errp != nil {
		badRequest(w, r, errors.New("page or limitis not a number"))
		return
	}

	response, err := api.database.GetFriendRequestSentTo(ctx, username, int64(pageInt), int64(limitInt))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, 200, response)

}
