package api

import (
	"errors"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type CreatGroup struct {
	Name string `json:"name"`
}

type UpdateGroupPayload struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddGroup struct {
	Username string `json:"username"`
	Id       int64  `json:"id"`
}

// @Summary Create a new group
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body CreatGroup true "name"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Router /v1/firendship/group/create [post]
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

	err = api.database.InsertFriendshipGroup(ctx, username, id)

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

// @Summary Delete group
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Failure 402  {object} errorslope
// @Router /v1/firendship/group/delete/{id} [delete]
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
		internalServer(w, r, err)
		return
	}

	err = api.database.DeleteAllGroupMembers(ctx, int64(idInt))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "group deleted succefully",
	}

	writeJson(w, 200, s)

}

// @Summary Get group members
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Failure 402  {object} errorslope
// @Router /v1/firendship/group/get-members/{id} [get]
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

// @Summary Add group member
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body AddGroup true "group id and username"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Failure 402  {object} errorslope
// @Router /v1/firendship/group/add-member [post]
func (api *ApiService) AddGroupMember(w http.ResponseWriter, r *http.Request) {

	var newMember AddGroup

	if err := readJson(w, r, &newMember); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	err := api.database.InsertGroupMember(ctx, newMember.Username, newMember.Id, "member")

	if err != nil {
		internalServer(w, r, err)
		return
	}

	err = api.database.InsertFriendshipGroup(ctx, newMember.Username, newMember.Id)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "member added to group",
	}

	writeJson(w, 200, s)

}

// @Summary Remove group member
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body AddGroup true "group id and username"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Failure 402  {object} errorslope
// @Router /v1/firendship/group/remove-member [delete]
func (api *ApiService) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {

	var newMember AddGroup

	if err := readJson(w, r, &newMember); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	username, err := getUsernameFromCtx(ctx)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	member, err := api.database.GetGroupMemberByUsername(ctx, username, int(newMember.Id))

	if err != nil {
		unauthorized(w, r, errors.New("user does not have permision to perform this action"))
		return
	}

	if member.Role != "admin" {
		unauthorized(w, r, errors.New("user does not have permision to perform this action"))
		return
	}

	err = api.database.DeleteGroupMember(ctx, newMember.Username, newMember.Id)
	if err != nil {
		internalServer(w, r, err)
		return
	}

	err = api.database.UpdateFriendshipGroupId(ctx, newMember.Username, 000)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  200,
		Message: "member removed from group",
	}

	writeJson(w, 200, s)
}

// @Summary Remove group member
// @Description Responds with json
// @Tags Friendship
// @Accept json
// @Produce json
// @Param payload body UpdateGroupPayload true "group details"
// @Success 200 {object} StandardResponse
// @Failure 400  {object} errorslope
// @Failure 500  {object} errorslope
// @Failure 402  {object} errorslope
// @Router /v1/firendship/group/update [delete]
func (api *ApiService) UpdateGroup(w http.ResponseWriter, r *http.Request) {

	var group UpdateGroupPayload

	if err := readJson(w, r, &group); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	err := api.database.UpdateGroup(ctx, int(group.Id), group.Name, group.Description, "")

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "group profile updated successfuly",
	}

	writeJson(w, http.StatusOK, s)
}

// @Summary Upload Group Pic
// @Description Responds with json
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param img formData file true "upload send png jpeg and gif"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Security ApiKeyAuth
// @Router /v1/firendship/group/upload-group-pic [post]
func (api *ApiService) UploadGroupPic(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	r.ParseMultipartForm(30 << 20)

	id := r.FormValue("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		badRequest(w, r, err)
		return
	}

	file, fileHeader, err := r.FormFile("img")

	if err != nil {
		badRequest(w, r, err)
		return
	}

	_, _, err = image.Decode(file)

	if err != nil {
		badRequest(w, r, err)
		return
	}

	defer file.Close()

	currentTime := time.Now().UnixMilli()

	currentTimeString := strconv.Itoa(int(currentTime)) + filepath.Ext(fileHeader.Filename)

	destinationFile, err := os.Create("/home/ifeanyi/nkata_storage/group/" + currentTimeString)

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

	url := "localhost:5557/v1/media/groups/" + currentTimeString

	err = api.database.UpdateGroup(ctx, idInt, "", "", url)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "user profile picture updated successfuly",
	}

	writeJson(w, http.StatusOK, s)
}

// @Summary Download Group Pic
// @Description Responds with json
// @Tags Media
// @Param img_name path string true "file name"
// @Produce octet-stream
// @Success 200 {file} file
// @Failure 404 {object} errorslope
// @Router /v1/media/groups/{img_name} [get]
func (api *ApiService) LoadGroupPic(w http.ResponseWriter, r *http.Request) {

	filename := chi.URLParam(r, "img_name")
	url := "/home/ifeanyi/nkata_storage/group/" + filename
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
