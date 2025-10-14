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
	GroupID   int64  `json:"group_id"`
	Username  int64  `json:"username"`
	Role      string `json:"role"` // admin or member
	CreatedAt string `json:"created_at"`
} // once u add a user to a group they get added here and in friendship

// type GroupMemberRemoved struct {
// 	ID        int64  `json:"id"`
// 	GroupID   int64  `json:"group_id"`
// 	UserID    int64  `json:"user_id"`
// 	CreatedAt string `json:"created_at"`
// }

func (d *DataRepository) InsertGroup(ctx context.Context, name string) (int64, error) {

	query := `INSERT INTO group(name,pic_url,description) VALUES($1,$2,$3) RETURNING id `

	var id int64

	err := d.db.QueryRowContext(ctx, query, name, "", "").Scan(&id)

	return id, err
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
	} else if picUrl != "" {
		_, err := d.db.ExecContext(cxt, queryPicUrl, picUrl)
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

//------------------------------ GroupMemeber ----------------------------------------------------------------------

func (d *DataRepository) InsertGroupMember(ctx context.Context, username string, groupId int64, role string) error {

	query := `INSERT INTO group_member(username,group_id,role) VALUES($1,$2,$3)`

	_, err := d.db.ExecContext(ctx, query, username, groupId, role)

	return err
}

func (d *DataRepository) GetGroupMemberByUsername(cxt context.Context, username string, id int) (*GroupMember, error) {

	query := `SELECT * FROM group_member WHERE groud_id = $1 AND username = $2`

	row, err := d.db.QueryContext(cxt, query, id, username)

	if err != nil {
		return nil, err
	}

	var member GroupMember

	err = row.Scan(&member.ID, &member.GroupID, &member.Username, &member.Role, &member.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &member,err
}

func (d *DataRepository) GetGroupMembersByGroupId(cxt context.Context, id, limit, page int64) (*PaginatedResponse, error) {

	offset := (page - 1) * limit

	var totalCount int64

	query := `SELECT * FROM group_member WHERE groud_id = $1 LIMIT = $2 OFFSET = $3`
	queryCount := `SELECT COUNT(*) FROM group_member WHERE groud_id = $1`

	if err := d.db.QueryRowContext(cxt, queryCount, id).Scan(&totalCount); err != nil {
		return nil, err
	}

	row, err := d.db.QueryContext(cxt, query, id, limit, offset)

	if err != nil {
		return nil, err
	}

	var members []GroupMember

	for row.Next() {

		var member GroupMember

		err := row.Scan(&member.ID, &member.GroupID, &member.Username, &member.Role, &member.CreatedAt)

		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	p := PaginatedResponse{
		Data:       members,
		TotalCount: int(totalCount),
		Page:       int(page),
		Limit:      int(limit),
	}

	return &p, nil
}

//remove from group
func (d *DataRepository) DeleteGroupMember(ctx context.Context, username string,id int64) error {
	query := `DELETE FROM group_member WHERE username = $1 AND group_id = $2`

	_, err := d.db.ExecContext(ctx, query, username,id)

	return err
}

func (d *DataRepository) DeleteAllGroupMembers(ctx context.Context,id int64) error{
		query := `DELETE FROM group_member WHERE group_id = $1`

	_, err := d.db.ExecContext(ctx, query,id)

	return err
}