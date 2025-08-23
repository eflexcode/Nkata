package database

import (
	"database/sql"
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
	Role         Role   `json:"role"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository{
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *User) error {

}

func (r *UserRepository) GetByID(id int64) error{
	
}

func (r *UserRepository) CheackUserNameAvailability(id int64) (bool,error){
	
}