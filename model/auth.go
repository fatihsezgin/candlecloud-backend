package model

//AuthLoginDTO ...
type AuthLoginDTO struct {
	Email          string `validate:"required" json:"email"`
	MasterPassword string `validate:"required" json:"master_password"`
}

//AuthLoginResponse ...
type AuthLoginResponse struct {
	AccessToken string `json:"access_token"`
	*UserDTO
}
