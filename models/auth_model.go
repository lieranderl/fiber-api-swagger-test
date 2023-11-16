package models

import "time"

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type AccessTokenRecord struct {
	User           string    `json:"user" bson:"user"`
	AcessTokenHash []byte    `json:"access_token_hash" bson:"access_token_hash"`
	Created_time   time.Time `json:"created_time" bson:"created_time"`
}
