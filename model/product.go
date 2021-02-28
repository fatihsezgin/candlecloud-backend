package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	UUID        uuid.UUID  `gorm:"type:uuid; type:varchar(100);"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Description string     `json:"description"`
	ImagePath   string     `json:"image_path"`
	Price       float64    `json:"price"`
	Quantity    int        `json:"quantity"`
}

type ProductDTO struct {
	ID          uint      `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Description string    `json:"description"`
	ImagePath   string    `json:"image_path"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}

func ToProduct(productDTO *ProductDTO) *Product {
	return &Product{
		ID:          productDTO.ID,
		UUID:        productDTO.UUID,
		Description: productDTO.Description,
		ImagePath:   productDTO.ImagePath,
		Price:       productDTO.Price,
		Quantity:    productDTO.Quantity,
	}
}

func ToProductDTO(product *Product) *ProductDTO {
	return &ProductDTO{
		ID:          product.ID,
		UUID:        product.UUID,
		Description: product.Description,
		ImagePath:   product.ImagePath,
		Price:       product.Price,
		Quantity:    product.Quantity,
	}
}
func ToProductDTOTable(product Product) ProductDTO {
	return ProductDTO{
		ID:          product.ID,
		UUID:        product.UUID,
		Description: product.Description,
		ImagePath:   product.ImagePath,
		Price:       product.Price,
		Quantity:    product.Quantity,
	}
}

func ToProductDTOs(products []Product) []ProductDTO {
	productDTOs := make([]ProductDTO, len(products))
	for i, item := range products {
		productDTOs[i] = ToProductDTOTable(item)
	}
	return productDTOs
}
