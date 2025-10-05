package database

import (
	"context"
	"errors"
)

type Friendship struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	FirendID  int64  `json:"friend_id"`
	Status    string `json:"status"` //blocked, unblocked
	CreatedAt string `json:"created_at"`
}

type FriendRequest struct {
	ID        int64  `json:"id"`
	SentBy    int64  `json:"sent_by"`
	SentTo    int64  `json:"sent_to"`
	Status    string `json:"status"` // accepted, pendding, rejected
	CreatedAt string `json:"created_at"`
}

func (r *DataRepository) InsertFriendRequest(ctx context.Context, sentTo, sentBy int64) error {

	query := `INSERT INTO friendRequest(sent_by,sent_to,status) VALUES($1,$2,$3)`

	_, err := r.db.ExecContext(ctx, query, sentBy, sentTo, "pending")

	return err
}

// request i (client) sent out
func (r *DataRepository) GetFriendRequestSentBy(ctx context.Context, sentBy, page, limit int64) (*PaginatedResponse, error) {

	var request []FriendRequest

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND status = $2 LIMIT = $3 OFFSET = $4`
	queryCount := `SELECT COUNT(*) FROM friendRequest WHERE sent_by = $1 AND status = $2`

	var totalCount int

	cRow := r.db.QueryRowContext(ctx, queryCount, sentBy, "pending")

	err := cRow.Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	row, err := r.db.Query(query, sentBy, "pending", limit, offset)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {

		item := FriendRequest{}

		err := row.Scan(&item.ID, &item.SentBy, &item.SentTo, &item.Status, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		request = append(request, item)

	}

	p := PaginatedResponse{
		Data: request,
		TotalCount: totalCount,
		Page: int(page),
		Limit: int(limit),
	}

	return &p, nil
}

// request i (client) was sent
func (r *DataRepository) GetFriendRequestSentTo(ctx context.Context, sentTo, page, limit int64) (*PaginatedResponse, error) {

	var request []FriendRequest

	query := `SELECT * FROM friendRequest WHERE sent_to = $1 AND status = $2 LIMIT = $3 OFFSET = $4`
	queryCount := `SELECT COUNT(*) FROM friendRequest WHERE sent_to = $1 AND status = $2`

	var totalCount int

	cRow := r.db.QueryRowContext(ctx, queryCount, sentTo, "pending")

	err := cRow.Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	row, err := r.db.Query(query, sentTo, "pending", limit, offset)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {

		item := FriendRequest{}

		err := row.Scan(&item.ID, &item.SentBy, &item.SentTo, &item.Status, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		request = append(request, item)

	}

	p := PaginatedResponse{
		Data: request,
		TotalCount: totalCount,
		Page: int(page),
		Limit: int(limit),
	}

	return &p, nil
}

func (r *DataRepository) GetFriendRequestById(ctx context.Context, id int64) (*FriendRequest, error) {

	query := `SELECT * FROM friendRequest WHERE ido = $1`

	row, err := r.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	row.Next()
	item := FriendRequest{}

	err = row.Scan(&item.ID, &item.SentBy, &item.SentTo, &item.Status, &item.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (d *DataRepository) UpdateFriendRequestStatus(ctx context.Context, status string, sentToId, id int64) error {

	if status != "accepted" && status != "rejected" {
		return errors.New("status can either be accepted or rejected only")
	}

	f, err := d.GetFriendRequestById(ctx, id)

	if err != nil {
		return err
	}

	if f.SentTo != sentToId {
		return errors.New("unauthorized: sentToId is not same")
	}

	query := `UPDATE friendRequest SET status = $1 WHERE id = $2`

	_, err = d.db.ExecContext(ctx, query, status, id)

	return err
}
