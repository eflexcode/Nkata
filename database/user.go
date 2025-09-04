package database

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	NUser Role = iota
	Admin
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	ImageUrl     string `json:"image_url"`
	Bio          string `json:"bio"`
	IsOnline     bool   `json:"is_online"`
	FriendsCount int64  `json:"friends_count"`
	GroupsCount  int16  `json:"groups_count"`
	Role         Role   `json:"_"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {

	query := `INSERT INTO users (username,display_name,password,role) VALUES($1,$2,$3,$4)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("Error hashing user password")
	}

	_, err = r.db.ExecContext(ctx, query, user.Username, user.DisplayName, string(hashedPassword), 0)

	if err != nil {

		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return errors.New("User with username " + user.Username + " already exist")
		} else {
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {

	query := `SELECT id,username,display_name,email,image_url,bio,is_online,friends_count,groups_count,created_at,modified_at FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user User

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.ImageUrl, &user.Bio, &user.IsOnline, &user.FriendsCount, &user.GroupsCount, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) CheackUsernameAvailability(ctx context.Context, id int64) bool {

	query := `SELECT id FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	return row.Err().Error() == sql.ErrNoRows.Error()

}
