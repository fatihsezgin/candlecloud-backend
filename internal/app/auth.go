package app

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatihsezgin/candlecloud-backend/model"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

var (
	//ErrExpiredToken represents message for expired token
	ErrExpiredToken = fmt.Errorf("Token expired or invalid")
	//ErrUnauthorized represents message for unauthorized
	ErrUnauthorized = fmt.Errorf("Unauthorized")
)

func CreateToken(user *model.User) (*model.TokenDetails, error) {
	var err error
	secret := viper.GetString("server.secret")
	td := &model.TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AtUUID = uuid.NewV4()

	td.RtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.RtUUID = uuid.NewV4()

	atClaims := jwt.MapClaims{}

	atClaims["authorized"] = false
	if user.Role == "Admin" {
		atClaims["authorized"] = true
	}

	atClaims["user_uuid"] = user.UUID.String()
	atClaims["exp"] = td.AtExpires
	atClaims["uuid"] = td.AtUUID.String()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["user_uuid"] = user.UUID.String()
	rtClaims["exp"] = td.RtExpires
	rtClaims["uuid"] = td.RtUUID.String()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

//TokenValid ...
func TokenValid(bearerToken string) (*jwt.Token, error) {
	token, err := verifyToken(bearerToken)
	if err != nil {
		if token != nil {
			return token, err
		}
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, ErrUnauthorized
	}
	return token, nil
}

//verifyToken verify token
func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("server.secret")), nil
	})
	if err != nil {
		return token, ErrExpiredToken
	}
	return token, nil
}
