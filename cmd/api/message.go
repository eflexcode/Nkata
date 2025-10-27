package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
