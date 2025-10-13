package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatGroup struct {
	Name string `json:"name"`
}

func (api *ApiService) CreateGroup(w http.ResponseWriter, r *http.Request) {

	var group CreatGroup

	ctx := r.Context()

	username, err := getUsernameFromCtx(ctx)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	if err := readJson(w, r, &group); err != nil {
		badRequest(w, r, err)
		return
	}

	id, err := api.database.InsertGroup(ctx, group.Name)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	err = api.database.InsertGroupMember(ctx, username, id, "admin")

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "Group created succefully",
	}

	writeJson(w, 200, s)

}

func (api *ApiService) DeleteGroup(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	username, err := getUsernameFromCtx(ctx)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	member, err := api.database.GetGroupMemberByUsername(ctx, username, idInt)

	if err != nil {
		unauthorized(w, r, errors.New("user does not have permision to perform this action"))
		return
	}

	if member.Role != "admin" {
		unauthorized(w, r, errors.New("user does not have permision to perform this action"))
		return
	}

	err = api.database.DeleteGroup(ctx, int64(idInt))

	if err != nil {
		internalServer(w, r,err)
		return
	}

	s:= StandardResponse{
		Status: 200,
		Message: "group deleted succefully",
	}
//TODO delete all mebers
	writeJson(w,200,s)

}

func (api *ApiService) GetGroupMembers(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	idInt, err := strconv.Atoi(id)
	pageInt, err1 := strconv.Atoi(page)
	limitInt, err2 := strconv.Atoi(limit)

	if err != nil || err1 != nil || err2 != nil {
		badRequest(w, r, errors.New("group id, page or limit might not be a number"))
		return
	}

	ctx := r.Context()

	result, err := api.database.GetGroupMembersByGroupId(ctx, int64(idInt), int64(limitInt), int64(pageInt))

	writeJson(w, 200, result)

}
