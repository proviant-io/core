package product_category

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gorm.io/gorm"
)

type ProductCategory struct {
	gorm.Model
	ProductId int `json:"product_id",gorm:"uniqueIndex:idx_member"`
	CategoryId int `json:"category_id",gorm:"uniqueIndex:idx_member"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) GetByProductId(id int) []ProductCategory {

	var models []ProductCategory

	r.db.Connection().Where("product_id = ?", id).Find(&models)

	return models
}

func (r *Repository) Link(productId int, categories []int) {

	r.db.Connection().Where("product_id = ?", productId).Delete(&ProductCategory{})

	for _, category := range categories{
		r.db.Connection().Create(&ProductCategory{ProductId: productId, CategoryId: category})
	}
}

func (r *Repository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&ProductCategory{})
	if err != nil {
		return fmt.Errorf("migration of Product table failed: %v", err)
	}
	return nil
}

func Setup(d db.DB) (*Repository, error) {

	repo := &Repository{}

	repo.db = d

	err := repo.Migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil

}