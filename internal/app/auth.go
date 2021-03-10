package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatihsezgin/candlecloud-backend/internal/cache"
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

	// access token uuid and expire date is initialized
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AtUUID = uuid.NewV4().String()

	// refresh token uuid and expire date is initialized
	td.RtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.RtUUID = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = false
	if user.Role == "Admin" {
		atClaims["authorized"] = true
	}

	atClaims["user_id"] = user.ID
	atClaims["user_uuid"] = user.UUID.String()
	atClaims["exp"] = td.AtExpires
	atClaims["auuid"] = td.AtUUID
	atClaims["ruuid"] = td.RtUUID

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = user.ID
	rtClaims["user_uuid"] = user.UUID.String()
	rtClaims["exp"] = td.RtExpires
	rtClaims["ruuid"] = td.RtUUID

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(token *jwt.Token) (*model.TokenDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		access_uuid, ok := claims["auuid"].(string)
		refresh_uuid, _ := claims["ruuid"].(string) // null geliyor
		if !ok {
			return nil, ErrUnauthorized
		} else {
			return &model.TokenDetails{
				AtUUID: access_uuid,
				RtUUID: refresh_uuid,
			}, nil
		}

	}
	return nil, errors.New("something went wrong")

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

func CreateAuth(userid uint, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	var ctx = context.Background()

	errAccess := cache.GetClient().Set(ctx, td.AtUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := cache.GetClient().Set(ctx, td.RtUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

//TODO: implement the DeleteAuth function
func DeleteAuth(token *jwt.Token) error {
	td, _ := ExtractTokenMetadata(token)
	var ctx = context.Background()
	deletedAt, err := cache.GetClient().Del(ctx, td.AtUUID).Result()
	if err != nil {
		return err
	}
	deletedRt, err := cache.GetClient().Del(ctx, td.RtUUID).Result()
	if err != nil {
		return err
	}
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("Someting went wrong while deleting auth token")
	}
	return nil
}
