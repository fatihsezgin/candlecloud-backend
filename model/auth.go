package model

import (
	uuid "github.com/satori/go.uuid"
)

//AuthLoginDTO ...
type AuthLoginDTO struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}

//AuthLoginResponse ...
type AuthLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	*UserDTO
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AtExpires    int64
	RtExpires    int64
	AtUUID       uuid.UUID
	RtUUID       uuid.UUID
}
