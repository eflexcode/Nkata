package database

import (
	"context"
	"errors"
)

type Group struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PicUrl      string `json:"pic_url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
}

type GroupMember struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"room_id"`
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"` // admin or member
	CreatedAt string `json:"created_at"`
} // once u add a user to a group they get added here and in friendship

type GroupMemberRemoved struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"room_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (d *DataRepository) InsertGroup(ctx context.Context, name string) error {

	query := `INSERT INTO group(name,pic_url,description) VALUES($1,$2,$3)`

	_, err := d.db.ExecContext(ctx, query, name, "", "")

	return err
}

func (d *DataRepository) GetGroupById(cxt context.Context, id int64) (*Group, error) {

	query := `SELECT * FROM group WHERE id = $1`

	row, err := d.db.QueryContext(cxt, query, id)

	if err != nil {
		return nil, err
	}

	row.Next()

	var group Group

	err = row.Scan(&group.ID, &group.Name, &group.PicUrl, &group.Description, &group.CreatedAt, &group.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (d *DataRepository) UpdateGroup(cxt context.Context, id int, name, description, picUrl string) error {

	queryAll := `UPDATE group SET name = $1, description = $2,pic_url = $3 WHERE id = $4?`
	queryPicUrl := `UPDATE group SET pic_url = $1 WHERE id = $2`
	queryName := `UPDATE group SET name = $1 WHERE id = $2`
	queryDescription := `UPDATE group SET description = $1 WHERE id = $2`

	if name != "" && description != "" && picUrl != "" {
		_, err := d.db.ExecContext(cxt, queryAll, name, description, picUrl)
		if err != nil {
			return err
		}
		return nil
	} else if name != "" {
		_, err := d.db.ExecContext(cxt, queryName, name)
		if err != nil {
			return err
		}
		return nil
	} else if description != "" {
		_, err := d.db.ExecContext(cxt, queryDescription, description)
		if err != nil {
			return err
		}
		return nil
	}else if picUrl != "" {
			_,err := d.db.ExecContext(cxt,queryPicUrl,picUrl)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("name,description,picUrl cannot all be empty")
}

func (d *DataRepository) DeleteGroup(ctx context.Context, id int64) error {
	query := `DELETE FROM group WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}