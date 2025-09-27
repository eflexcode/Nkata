package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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
	Role         string `json:"_"`
	Enabled      bool   `json:"_"`
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

	query := `INSERT INTO users (username,display_name,email,password,image_url,bio,is_online,friends_count,groups_count,role,enabled,modified_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("error hashing user password")
	}

	_, err = r.db.ExecContext(ctx, query, user.Username, user.DisplayName, "", string(hashedPassword), "", "", false, 0, 0, 0, true, time.Now())

	if err != nil {

		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {

	query := `SELECT id,username,display_name,email,password,image_url,bio,is_online,friends_count,groups_count,created_at,modified_at FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user User

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.Password, &user.ImageUrl, &user.Bio, &user.IsOnline, &user.FriendsCount, &user.GroupsCount, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {

	query := `SELECT id,username,display_name,email,password,image_url,bio,is_online,friends_count,groups_count,created_at,modified_at FROM users WHERE username = $1`

	row := r.db.QueryRowContext(ctx, query, username)

	var user User

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.Password, &user.ImageUrl, &user.Bio, &user.IsOnline, &user.FriendsCount, &user.GroupsCount, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) Update(ctx context.Context, username, displayName, bio string) error {

	queryBoth := `UPDATE users SET display_name = $1, bio =$2 WHERE username = $3?`
	queryBio := `UPDATE users SET bio = $1 WHERE username = $2`
	queryDisplay := `UPDATE users SET display_name = $1 WHERE username = $2`

	if displayName != "" && bio != "" {
		_, err := r.db.ExecContext(ctx, queryBoth, displayName, bio, username)
		if err != nil {
			return err
		}
		return nil

	} else if displayName != "" {
		_, err := r.db.ExecContext(ctx, queryDisplay, displayName, username)
		if err != nil {
			return err
		}
		return nil
	} else if bio != "" {
		log.Println(queryBio)
		_, err := r.db.ExecContext(ctx, queryBio, bio, username)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("display_name and bio cannot both be empty")

}

func (r *UserRepository) UpdateProfilePicUrl(ctx context.Context, username, imageUrl string) error {

	query := `UPDATE users SET image_url = $1 WHERE username = $2`

	_, err := r.db.ExecContext(ctx, query, imageUrl, username)
	return err

}

func (r *UserRepository) CheackUsernameAvailability(ctx context.Context, username string) bool {

	query := `SELECT username FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)

	if row.Err().Error() == sql.ErrNoRows.Error() {
		return false
	}

	return true

}
