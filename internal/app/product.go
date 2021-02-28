package app

import (
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/fatihsezgin/candlecloud-backend/model"
	uuid "github.com/satori/go.uuid"
)

// TODO for being able to add product, the user should login
func CreateProduct(s storage.Store, productDTO *model.ProductDTO) (*model.Product, error) {
	// Check the payload is valid
	err := PayloadValidator(productDTO)
	if err != nil {
		return nil, err
	}
	productDTO.UUID = uuid.NewV4()
	return s.Products().Save(model.ToProduct(productDTO))
}
