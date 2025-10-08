package database

import (
	"context"
	"errors"
)

type Friendship struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	FirendID       int64  `json:"friend_id,omitempty"`
	FriendshipType string `json:"friendship_type"`    //one-on-one or group
	GroupID        int64  `json:"group_id,omitempty"` //if group; remove id to remove member from group
	CreatedAt      string `json:"created_at"`
}

type FriendRequest struct {
	ID        int64  `json:"id"`
	SentBy    int64  `json:"sent_by"`
	SentTo    int64  `json:"sent_to"`
	Status    string `json:"status"` //accepted, pendding, rejected
	CreatedAt string `json:"created_at"`
}

// ------------------------------ Friend Request ----------------------------------------------------------------------
func (r *DataRepository) InsertFriendRequest(ctx context.Context, sentTo, sentBy int64) error {

	query := `INSERT INTO friendRequest(sent_by,sent_to,status) VALUES($1,$2,$3)`

	_, err := r.db.ExecContext(ctx, query, sentBy, sentTo, "pending")

	return err
}

func (r *DataRepository) DeleteFriendRequest(ctx context.Context, id int64) error {

	query := `DELETE FROM friendRequest WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return err

}

func (r *DataRepository) HasSentMeRequest(ctx context.Context, firenId, userId int64)  bool {

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND sent_to = $2`

	row, err := r.db.Query(query, firenId,userId)

	if err != nil {
		return true
	}

	defer row.Close()

	
	return row.Next()
}


func (r *DataRepository) CheckDuplicateRequest(ctx context.Context, userId, firendId int64)  bool {

	query := `SELECT * FROM friendRequest WHERE sent_by = $1 AND sent_to = $2`

	row, err := r.db.Query(query, userId,firendId)

	if err != nil {
		return true
	}

	defer row.Close()

	
	return row.Next()
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
		Data:       request,
		TotalCount: totalCount,
		Page:       int(page),
		Limit:      int(limit),
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
		Data:       request,
		TotalCount: totalCount,
		Page:       int(page),
		Limit:      int(limit),
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

//------------------------------ Friendship ----------------------------------------------------------------------

func (d *DataRepository) InsertFriendship(ctx context.Context, userId, firendId int64) error {

	query := `INSERT INTO friendship(user_id,friend_id) VALUES($1,$2,$3)`

	_, err := d.db.ExecContext(ctx, query, userId, firendId, "one-on-one")

	return err

}

func (d *DataRepository) InsertFriendshipGroup(ctx context.Context, userId, groupId int64) error {

	query := `INSERT INTO friendship(user_id,group_id) VALUES($1,$2,$3)`

	_, err := d.db.ExecContext(ctx, query, userId, groupId, "group")

	return err

}

func (d *DataRepository) RemoveGroupFromFriendship(ctx context.Context, id int64) error {
	query := `UPDATE friendship SET friendship SET group_id = $1  WHERE id = $2`

	_, err := d.db.ExecContext(ctx, query, 0, id)

	return err
}
func (d *DataRepository) DeleteFriendship(ctx context.Context, id int64) error {
	query := `DELETE FROM friendship WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}

func (d *DataRepository) GetFriendshipByUserID(ctx context.Context, userId, page, limit int64) (*PaginatedResponse, error) {

	offset := (page - 1) * limit
	var totalCount int

	query := "SELECT * FROM friendship WHERE user_id = $1 LIMIT = $2 OFFSET = $3"
	queryCount := `SELECT COUNT(*) WHERE user_id`

	counrRow := d.db.QueryRowContext(ctx, queryCount, userId)

	err := counrRow.Scan(&totalCount)

	if err != nil {
		return nil, err
	}

	row, err := d.db.QueryContext(ctx, query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	var friendship []Friendship

	for row.Next() {

		item := Friendship{}

		err := row.Scan(&item.ID, &item.UserID, &item.FirendID, &item.FriendshipType, &item.GroupID, &item.CreatedAt)

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
