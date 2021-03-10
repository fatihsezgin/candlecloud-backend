package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatihsezgin/candlecloud-backend/internal/app"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/fatihsezgin/candlecloud-backend/model"
	"github.com/go-playground/validator"
)

const (
	//InvalidJSON represents a message for invalid json
	InvalidJSON = "Invalid json provided"
)

var (
	Success           = "Success"
	signupSuccess     = "User created successfully"
	userLoginErr      = "User email or master password is wrong."
	userVerifyErr     = "Please verify your email first."
	invalidUser       = "Invalid user"
	validToken        = "Token is valid"
	invalidToken      = "Token is expired or not valid!"
	noToken           = "Token could not found! "
	tokenCreateErr    = "Token could not be created"
	tokenDeleteErr    = "Token could not be deleted"
	userLogoutSuccess = "User logged out successfully"
)

func Signup(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSignup := new(model.UserSignup)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userSignup); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		// Run the validator to model.UserSignup validator tags
		err := app.PayloadValidator(userSignup)
		if err != nil {
			errs := GetErrors(err.(validator.ValidationErrors))
			RespondWithErrors(w, http.StatusBadRequest, InvalidRequestPayload, errs)
			return
		}

		// Check if user exists in the database
		userDTO := model.ConvertUserDTO(userSignup)
		_, err = s.Users().FindByEmail(userDTO.Email)
		if err == nil {
			RespondWithError(w, http.StatusBadRequest, "User couldn't created!")
			return
		}

		createdUser, err := app.CreateUser(s, userDTO)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = app.GenerateSchema(s, createdUser)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		//app.MigrateUserTables(s, updatedUser.Schema)

		response := model.Response{
			Code:    http.StatusOK,
			Status:  Success,
			Message: signupSuccess,
		}
		RespondWithJSON(w, http.StatusOK, response)
	}
}

func Signin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDTO model.AuthLoginDTO

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginDTO); err != nil {
			RespondWithError(w, http.StatusInternalServerError, InvalidJSON)
			return
		}
		defer r.Body.Close()

		err := app.PayloadValidator(loginDTO)
		if err != nil {
			errs := GetErrors(err.(validator.ValidationErrors))
			RespondWithErrors(w, http.StatusBadRequest, InvalidRequestPayload, errs)
			return
		}

		user, err := s.Users().FindByCredentials(loginDTO.Email, loginDTO.Password)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, userLoginErr)
			return
		}

		token, err := app.CreateToken(user)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, tokenCreateErr)
			return
		}
		saveErr := app.CreateAuth(user.ID, token)
		if saveErr != nil {
			log.Fatal("error while saving redis", saveErr)
			return
		}

		authLoginResponse := model.AuthLoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			UserDTO:      model.ToUserDTO(user),
		}
		RespondWithJSON(w, 200, authLoginResponse)
	}
}

func Singout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		au := app.ExtractToken(r)
		if au == "" {
			RespondWithError(w, http.StatusUnauthorized, noToken)
			return
		}
		token, err := app.TokenValid(au)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, invalidToken)
			return
		}
		deleteErr := app.DeleteAuth(token)
		if deleteErr != nil {
			RespondWithError(w, http.StatusInternalServerError, deleteErr.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, userLogoutSuccess)
	}
}
