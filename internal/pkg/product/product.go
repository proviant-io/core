package product

import (
	"fmt"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id          int             `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Link        string          `json:"link"`
	Image       string          `json:"image"`
	Barcode     string          `json:"barcode"`
	ListId      int             `json:"list_id"`
	Stock       uint            `json:"stock" gorm:"type:UINT"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(20,2);"`
	AccountId   int             `json:"account_id" gorm:"default:0;index"`
}

type CreateDTO struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Link        string          `json:"link"`
	Image       string          `json:"image"`
	ImageBase64 string          `json:"image_base64"`
	Barcode     string          `json:"barcode"`
	CategoryIds []int           `json:"category_ids"`
	ListId      int             `json:"list_id"`
	Stock       uint            `json:"stock"`
	Price       decimal.Decimal `json:"price"`
}

type UpdateDTO struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Link        string          `json:"link"`
	Image       string          `json:"image"`
	ImageBase64 string          `json:"image_base64"`
	Barcode     string          `json:"barcode"`
	CategoryIds []int           `json:"category_ids"`
	ListId      int             `json:"list_id"`
	Stock       uint            `json:"stock"`
	Price       decimal.Decimal `json:"price"`
}

type DTO struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Link        string          `json:"link"`
	Image       string          `json:"image"`
	Barcode     string          `json:"barcode"`
	CategoryIds []int           `json:"category_ids"`
	Categories  interface{}     `json:"categories"`
	ListId      int             `json:"list_id"`
	List        interface{}     `json:"list"`
	Stock       uint            `json:"stock"`
	Price       decimal.Decimal `json:"price"`
}

type Repository struct {
	db db.DB
}

type Query struct {
	Category int
	List     int
}

func (r *Repository) Get(id int, accountId int) (Product, *errors.CustomError) {

	p := &Product{}

	r.db.Connection().First(p, "id = ? and account_id = ?", id, accountId)

	if (*p).Id == 0 {
		return Product{}, errors.NewErrNotFound(i18n.NewMessage("product with id %d not found", id))
	}

	return *p, nil
}

func (r *Repository) GetAll(query *Query, accountId int) []Product {

	var products []Product

	if query == nil {
		r.db.Connection().Where("account_id = ?", accountId).Find(&products)
	} else {
		queryBuilder := &Product{}
		queryBuilder.AccountId = accountId

		if query.List != 0 {
			queryBuilder.ListId = query.List
		}
		r.db.Connection().Where(queryBuilder).Find(&products)
	}

	return products
}

func (r *Repository) Delete(id int, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return err
	}

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto CreateDTO, accountId int) Product {

	p := &Product{
		Title:       dto.Title,
		Description: dto.Description,
		Link:        dto.Link,
		Image:       dto.Image,
		Barcode:     dto.Barcode,
		ListId:      dto.ListId,
		Stock:       0,
		AccountId:   accountId,
		Price:       dto.Price,
	}

	r.db.Connection().Create(p)
	return *p
}

func (r *Repository) Save(model Product, accountId int) (Product, *errors.CustomError) {

	// sanity check
	_, err := r.Get(model.Id, accountId)

	if err != nil {
		return Product{}, err
	}

	r.db.Connection().Model(&Product{Id: model.Id}).Updates(model)

	if model.Stock == 0 {
		r.db.Connection().Model(&Product{Id: model.Id}).Select("Stock").Updates(Product{Stock: 0})
	}

	return model, nil
}

func (r *Repository) UpdateFromDTO(dto UpdateDTO, accountId int) (Product, *errors.CustomError) {

	model, err := r.Get(dto.Id, accountId)

	if err != nil {
		return Product{}, err
	}

	model.Title = dto.Title
	model.Description = dto.Description
	model.Link = dto.Link
	model.Image = dto.Image
	model.Barcode = dto.Barcode
	model.ListId = dto.ListId
	model.Stock = dto.Stock
	model.Price = dto.Price

	r.db.Connection().Model(&Product{Id: dto.Id}).Updates(model)

	if model.Stock == 0 {
		r.db.Connection().Model(&Product{Id: dto.Id}).Select("Stock").Updates(Product{Stock: 0})
	}

	return model, nil
}

func ModelToDTO(m Product) DTO {
	return DTO{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Link:        m.Link,
		Image:       m.Image,
		Barcode:     m.Barcode,
		ListId:      m.ListId,
		Stock:       m.Stock,
		Price:       m.Price,
	}
}

func (r *Repository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Product{})
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
