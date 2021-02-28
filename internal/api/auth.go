package api

import (
	"encoding/json"
	"net/http"

	"github.com/fatihsezgin/candlecloud-backend/internal/app"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/fatihsezgin/candlecloud-backend/model"
	"github.com/go-playground/validator"
)

var (
	Success       = "Success"
	signupSuccess = "User created successfully"
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
