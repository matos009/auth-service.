package models

import (
	"auth-service/db"
	"time"
)

type RefreshToken struct {
	ID               int
	UserID           string
	RefreshTokenHash string
	IPAdress         string
	ExpiresAt        time.Time
}

func GetRefreshTokenByToken(refreshToken string) (*RefreshToken, error) {
	var token RefreshToken

	row := db.DB.QueryRow(
		`SELECT id, user_id, refresh_token_hash, ip_address, expires_at
		 FROM refresh_tokens WHERE refresh_token_hash = $1`, refreshToken,
	)

	err := row.Scan(&token.ID, &token.UserID, &token.RefreshTokenHash, &token.IPAdress, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshToken) Save() error {
	_, err := db.DB.Exec(
		`INSERT INTO refresh_tokens (user_id, refresh_token_hash, ip_address, expires_at)
		VALUES($1, $2, $3, $4)`,
		r.UserID, r.RefreshTokenHash, r.IPAdress, r.ExpiresAt)

	return err
}

func GetRefTokenById(user_id string) (*RefreshToken, error) {
	row := db.DB.QueryRow(
		`SELECT id, user_id, refresh_token_hash, ip_address, expires_at
		 FROM refresh_tokens WHERE user_id = $1`, user_id,
	)
	var token RefreshToken
	err := row.Scan(&token.ID, &token.UserID, &token.RefreshTokenHash, &token.IPAdress, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshToken) DeleteById() error {
	_, err := db.DB.Exec(`DELETE FROM refresh_tokens WHERE id = $1`, r.ID)
	return err

}
