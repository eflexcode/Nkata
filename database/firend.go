package database

import (
	"context"
	"errors"
	"time"
)

type Friendship struct {
	ID             int64     `json:"id"`
	FriendShipId   string    `json:"firendship_id"`
	Username       string    `json:"username"`
	FirendUsername string    `json:"friend_username,omitempty"`
	FriendshipType string    `json:"friendship_type"`    //one-on-one or group
	GroupID        int64     `json:"group_id,omitempty"` //if group; remove id to remove member from group
	CreatedAt      time.Time `json:"created_at"`
}

type FriendRequest struct {
	ID         int64     `json:"id"`
	SentBy     string    `json:"sent_by"`
	SentTo     string    `json:"sent_to"`
	Status     string    `json:"status"` //accepted, pendding, rejected
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// ------------------------------ Friend Request ----------------------------------------------------------------------
func (r *DataRepository) InsertFriendRequest(ctx context.Context, sentTo, sentBy string) error {

	query := `INSERT INTO friendRequest(sent_by,sent_to,status,modified_at) VALUES($1,$2,$3,$4)`

	_, err := r.db.ExecContext(ctx, query, sentBy, sentTo, "pending", time.Now())

	return err
}

func (r *DataRepository) DeleteFriendRequest(ctx context.Context, id int64) error {

	query := `DELETE FROM friendRequest WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return err

}

func (r *DataRepository) HasSentMeRequest(ctx context.Context, friendUsername, username string) bool {

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND sent_to = $2`

	row, err := r.db.Query(query, friendUsername, username)

	if err != nil {
		return true
	}

	defer row.Close()

	return row.Next()
}

func (r *DataRepository) CheckDuplicateRequest(ctx context.Context, username, firendUsername string) bool {

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND sent_to = $2`

	row, err := r.db.Query(query, username, firendUsername)

	if err != nil {
		return true
	}

	defer row.Close()

	return row.Next()
}

// request i (client) sent out
func (r *DataRepository) GetFriendRequestSentBy(ctx context.Context, sentByUsername string, page, limit int64) (*PaginatedResponse, error) {

	var request []FriendRequest

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND status = $2 LIMIT = $3 OFFSET = $4`
	queryCount := `SELECT COUNT(*) FROM friendRequest WHERE sent_by = $1 AND status = $2`

	var totalCount int

	cRow := r.db.QueryRowContext(ctx, queryCount, sentByUsername, "pending")

	err := cRow.Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	row, err := r.db.Query(query, sentByUsername, "pending", limit, offset)

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
		Data:       request,
		TotalCount: totalCount,
		Page:       int(page),
		Limit:      int(limit),
	}

	return &p, nil
}

// request i (client) was sent
func (r *DataRepository) GetFriendRequestSentTo(ctx context.Context, sentToUsername string, page, limit int64) (*PaginatedResponse, error) {

	var request []FriendRequest

	query := `SELECT * FROM friendRequest WHERE sent_to = $1 AND status = $2 LIMIT = $3 OFFSET = $4`
	queryCount := `SELECT COUNT(*) FROM friendRequest WHERE sent_to = $1 AND status = $2`

	var totalCount int

	cRow := r.db.QueryRowContext(ctx, queryCount, sentToUsername, "pending")

	err := cRow.Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	row, err := r.db.Query(query, sentToUsername, "pending", limit, offset)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {

		item := FriendRequest{}

		err := row.Scan(&item.ID, &item.SentBy, &item.SentTo, &item.Status, &item.CreatedAt, &item.ModifiedAt)

		if err != nil {
			return nil, err
		}

		request = append(request, item)

	}

	p := PaginatedResponse{
		Data:       request,
		TotalCount: totalCount,
		Page:       int(page),
		Limit:      int(limit),
	}

	return &p, nil
}

func (r *DataRepository) GetFriendRequestById(ctx context.Context, id int64) (*FriendRequest, error) {

	query := `SELECT * FROM friendRequest WHERE id = $1`

	row, err := r.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	row.Next()
	item := FriendRequest{}

	err = row.Scan(&item.ID, &item.SentBy, &item.SentTo, &item.Status, &item.CreatedAt, &item.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (d *DataRepository) UpdateFriendRequestStatus(ctx context.Context, status string, request_id int64) error {

	if status != "accepted" && status != "rejected" {
		return errors.New("status can either be accepted or rejected only")
	}

	// f, err := d.GetFriendRequestById(ctx, request_id)

	// if err != nil {
	// 	return err
	// }

	// if f.SentTo != sentToId {
	// 	return errors.New("unauthorized: sentToId is not same")
	// }

	query := `UPDATE friendRequest SET status = $1,modified_at = $2 WHERE id = $3`

	modifiedAt := time.Now()

	_, err := d.db.ExecContext(ctx, query, status, modifiedAt, request_id)

	return err
}

//------------------------------ Friendship ----------------------------------------------------------------------

func (d *DataRepository) InsertFriendship(ctx context.Context, username, firendUsername, friendship_id string) error {

	query := `INSERT INTO friendship(friendship_id,username,friend_username,modified_at) VALUES($1,$2,$3,$4,$5)`

	_, err := d.db.ExecContext(ctx, query, friendship_id, username, firendUsername, "one-on-one", time.Now())

	return err

}

func (d *DataRepository) InsertFriendshipGroup(ctx context.Context, userId string, groupId int64) error {

	query := `INSERT INTO friendship(username,group_id,modified_at) VALUES($1,$2,$3,$4)`

	_, err := d.db.ExecContext(ctx, query, userId, groupId, "group", time.Now())

	return err

}

func (d *DataRepository) RemoveGroupFromFriendship(ctx context.Context, id int64) error {
	query := `UPDATE friendship SET friendship SET group_id = $1,modified_at = $2  WHERE id = $3`

	_, err := d.db.ExecContext(ctx, query, 0, time.Now(), id)

	return err
}

func (d *DataRepository) DeleteFriendship(ctx context.Context, id int64) error {
	query := `DELETE FROM friendship WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}

func (d *DataRepository) GetFriendshipByUserID(ctx context.Context, username, page, limit int64) (*PaginatedResponse, error) {

	offset := (page - 1) * limit
	var totalCount int

	query := "SELECT * FROM friendship WHERE username = $1 LIMIT = $2 OFFSET = $3"
	queryCount := `SELECT COUNT(*) WHERE username`

	counrRow := d.db.QueryRowContext(ctx, queryCount, username)

	err := counrRow.Scan(&totalCount)

	if err != nil {
		return nil, err
	}

	row, err := d.db.QueryContext(ctx, query, username, limit, offset)

	if err != nil {
		return nil, err
	}

	var friendship []Friendship

	for row.Next() {

		item := Friendship{}

		err := row.Scan(&item.ID, &item.Username, &item.FirendUsername, &item.FriendshipType, &item.GroupID, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		friendship = append(friendship, item)
	}

	p := PaginatedResponse{
		Data:       friendship,
		Page:       int(page),
		Limit:      int(limit),
		TotalCount: totalCount,
	}

	return &p, nil

}
