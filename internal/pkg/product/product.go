package product

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`
	Image string `json:"image"`
	Barcode string `json:"barcode"`
}

type DTO struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`
	Image string `json:"image"`
	Barcode string `json:"barcode"`
	Categories []int `json:"categories"`
	ListId int `json:"list_id"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) Product{

	log.Printf("id: %d\n", id)

	p := &Product{}

	r.db.Connection().First(p, "id = ?", id)

	return *p
}

func (r *Repository) GetAll() []Product{

	var products []Product
	r.db.Connection().Find(&products)

	return products
}

func (r *Repository) Delete(id int){

	log.Printf("delete id: %d\n", id)

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Delete(&Product{}, id)
}

func (r *Repository) Create(dto DTO){

	p := Product{
		Title: dto.Title,
		Description: dto.Description,
		Link: dto.Link,
		Image: dto.Image,
		Barcode: dto.Barcode,
	}

	r.db.Connection().Create(&p)
}

func (r *Repository) Update(id int, dto DTO){

	p := Product{
		Title: dto.Title,
		Description: dto.Description,
		Link: dto.Link,
		Image: dto.Image,
		Barcode: dto.Barcode,
	}

	r.db.Connection().Model(&Product{Id: id}).Updates(p)
}

func ModelToDTO(m Product) DTO {
	return DTO{
		Id: m.Id,
		Title: m.Title,
		Description: m.Description,
		Link: m.Link,
		Image: m.Image,
		Barcode: m.Barcode,
	}
}


func (r *Repository) Migrate() error{
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Product{})
	if err != nil{
		return fmt.Errorf("migration of Product table failed: %v", err)
	}
	return nil
}

func Setup(d db.DB) (*Repository, error) {

	repo := &Repository{}

	repo.db = d

	err := repo.Migrate()
	if err != nil{
		return nil, err
	}

	return repo, nil

}