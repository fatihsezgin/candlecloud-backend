package api

import (
	"encoding/json"
	"net/http"

	"github.com/fatihsezgin/candlecloud-backend/internal/app"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/fatihsezgin/candlecloud-backend/model"
	"github.com/go-playground/validator"
)

func CreateProduct(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productDTO := new(model.ProductDTO)

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&productDTO); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		validate := validator.New()
		validateError := validate.Struct(productDTO)
		if validateError != nil {
			errs := GetErrors(validateError.(validator.ValidationErrors))
			RespondWithErrors(w, http.StatusBadRequest, InvalidRequestPayload, errs)
			return
		}

		createdProduct, err := app.CreateProduct(s, productDTO)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, model.ToProductDTO(createdProduct))
	}
}

// TODO get the all products from app
func GetAllProducts(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := s.Products().All()

		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		productDTOs := model.ToProductDTOs(products)
		// users = app.DecryptUserPasswords(users)
		RespondWithJSON(w, http.StatusOK, productDTOs)
	}
}
