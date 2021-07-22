package product_category

import (
	"fmt"
	"github.com/proviant-io/core/internal/db"
	"gorm.io/gorm"
)

type ProductCategory struct {
	gorm.Model
	ProductId  int `json:"product_id" gorm:"uniqueIndex:idx_member"`
	CategoryId int `json:"category_id" gorm:"uniqueIndex:idx_member"`
	AccountId  int `json:"account_id" gorm:"default:0"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) GetByProductId(id int, accountId int) []ProductCategory {

	var models []ProductCategory

	r.db.Connection().Where("product_id = ? and account_id = ?", id, accountId).Find(&models)

	return models
}

func (r *Repository) DeleteByProductId(id int, accountId int) {
	r.db.Connection().Where("product_id = ? and account_id = ?", id, accountId).Unscoped().Delete(&ProductCategory{})
}

func (r *Repository) DeleteByCategory(id int, accountId int) {
	r.db.Connection().Where("category_id = ? and account_id = ?", id, accountId).Unscoped().Delete(&ProductCategory{})
}

func (r *Repository) Link(productId int, categories []int, accountId int) {

	r.DeleteByProductId(productId, accountId)

	for _, category := range categories {
		r.db.Connection().Create(&ProductCategory{ProductId: productId, CategoryId: category, AccountId: accountId})
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
