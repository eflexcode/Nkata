package database

import (
	"context"
	"time"
)

type Otp struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Token      int64  `json:"token"`
	Email      string `json:"email"`
	Purpose    string `json:"Purpose"`
	Exp        string `json:"exp"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

func (r *DataRepository) InsertOtp(ctx context.Context, username, email, purpose string, token int64) error {

	query := `INSERT INTO otp (username,email,purpose,exp,token,modified_at)`

	exp := time.Now().Add(time.Minute * 20)

	_, err := r.db.ExecContext(ctx, query, username, email, purpose, exp, token, time.Now())

	return err
}

func (r *DataRepository) GetOtp(ctx context.Context, token int64) (*Otp, error) {

	query := `SELECT * FROM otp WHERE token = $1`

	row, err := r.db.QueryContext(ctx, query, token)

	if err != nil {
		return nil, err
	}

	row.Next()

	var otp Otp

	err = row.Scan(&otp.ID, &otp.Username, &otp.Token, &otp.Purpose, &otp.Exp, &otp.CreatedAt, &otp.ModifiedAt)

	if err != nil {
		return nil, err
	}
	
	return &otp, nil
}
