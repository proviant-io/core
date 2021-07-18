package product

import (
	"fmt"
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Barcode     string `json:"barcode"`
	ListId      int    `json:"list_id"`
	Stock       uint   `json:"stock",gorm:"type:UINT(10)"`
}

type CreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	ImageBase64 string `json:"image_base64"`
	Barcode     string `json:"barcode"`
	CategoryIds []int  `json:"category_ids"`
	ListId      int    `json:"list_id"`
}

type DTO struct {
	Id          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Link        string      `json:"link"`
	Image       string      `json:"image"`
	ImageBase64 string      `json:"image_base64"`
	Barcode     string      `json:"barcode"`
	CategoryIds []int       `json:"category_ids"`
	Categories  interface{} `json:"categories"`
	ListId      int         `json:"list_id"`
	List        interface{} `json:"list"`
	Stock       uint        `json:"stock"`
}

type Repository struct {
	db db.DB
}

type Query struct {
	Category int
	List     int
}

func (r *Repository) Get(id int) (Product, *errors.CustomError) {

	p := &Product{}

	r.db.Connection().First(p, "id = ?", id)

	if (*p).Id == 0 {
		return Product{}, errors.NewErrNotFound(i18n.NewMessage("product with id %d not found", id))
	}

	return *p, nil
}

func (r *Repository) GetAll(query *Query) []Product {

	var products []Product

	if query == nil {
		r.db.Connection().Find(&products)
	} else {
		queryBuilder := &Product{}

		if query.List != 0 {
			queryBuilder.ListId = query.List
		}
		r.db.Connection().Where(queryBuilder).Find(&products)
	}

	return products
}

func (r *Repository) Delete(id int) *errors.CustomError {

	model, err := r.Get(id)

	if err != nil {
		return err
	}

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto CreateDTO) Product {

	p := &Product{
		Title:       dto.Title,
		Description: dto.Description,
		Link:        dto.Link,
		Image:       dto.Image,
		Barcode:     dto.Barcode,
		ListId:      dto.ListId,
	}

	r.db.Connection().Create(p)
	return *p
}

func (r *Repository) Update(dto DTO) (Product, *errors.CustomError) {

	model, err := r.Get(dto.Id)

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
