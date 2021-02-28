package product

import (
	"github.com/fatihsezgin/candlecloud-backend/model"
	"github.com/jinzhu/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (p *Repository) All() ([]model.Product, error) {
	products := []model.Product{}
	err := p.db.Find(&products).Error
	return products, err
}

func (p *Repository) FindByID(id uint) (*model.Product, error) {
	product := new(model.Product)
	err := p.db.Where(`id = ?`, id).First(&product).Error
	return product, err
}

func (p *Repository) FindByUUID(id string) (*model.Product, error) {
	product := new(model.Product)
	err := p.db.Where(`uuid = ?`, id).First(&product).Error
	return product, err
}

// Save ...
func (p *Repository) Save(product *model.Product) (*model.Product, error) {
	err := p.db.Save(&product).Error
	return product, err
}

// TODO might not work, this shoul be DELETE query not DROP SCHEMA
// Delete ...
// func (p *Repository) Delete(id uint, schema string) error {

// 	err := p.db.Exec("DROP SCHEMA " + schema + " CASCADE").Error
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = p.db.Delete(&model.Product{ID: id}).Error
// 	return err
// }

// Migrate ...
func (p *Repository) Migrate() error {
	return p.db.AutoMigrate(&model.Product{}).Error
}
