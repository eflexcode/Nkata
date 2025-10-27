package api

import (
	"database/sql"
	"errors"
	"net/http"

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

func (api *ApiService) DeleteMessageByMessageId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "message_id")

	ctx := r.Context()

	err := api.database.DeleteMessageById(ctx, id)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status: 200,
		Message: "Message deleted successfully",
	}

	writeJson(w, http.StatusOK, s)

}
