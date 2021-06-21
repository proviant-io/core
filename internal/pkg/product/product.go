package product

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Barcode     string `json:"barcode"`
	ListId      int    `json:"list_id"`
}

type DTO struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Link        string        `json:"link"`
	Image       string        `json:"image"`
	Barcode     string        `json:"barcode"`
	CategoryIds []int         `json:"category_ids"`
	Categories  []interface{} `json:"categories"`
	ListId      int           `json:"list_id"`
	List        interface{}   `json:"list"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) (Product, error) {

	log.Printf("id: %d\n", id)

	p := &Product{}

	r.db.Connection().First(p, "id = ?", id)

	if (*p).Id == 0 {
		return Product{}, fmt.Errorf("product with id %d not found", id)
	}

	return *p, nil
}

func (r *Repository) GetAll() []Product {

	var products []Product
	r.db.Connection().Find(&products)

	return products
}

func (r *Repository) Delete(id int) error {

	model, err := r.Get(id)

	if err != nil {
		return err
	}

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto DTO) {

	p := Product{
		Title:       dto.Title,
		Description: dto.Description,
		Link:        dto.Link,
		Image:       dto.Image,
		Barcode:     dto.Barcode,
		ListId:      dto.ListId,
	}

	r.db.Connection().Create(&p)
}

func (r *Repository) Update(id int, dto DTO) error {

	model, err := r.Get(id)

	if err != nil {
		return err
	}

	model.Title = dto.Title
	model.Description = dto.Description
	model.Link = dto.Link
	model.Image = dto.Image
	model.Barcode = dto.Barcode
	model.ListId = dto.ListId

	r.db.Connection().Model(&Product{Id: id}).Updates(model)
	return nil
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
