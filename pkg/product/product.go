package product

import (
	"fmt"
	"gorm.io/driver/mysql"
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
	db *gorm.DB
}

func (r *Repository) Get(id int) Product{

	log.Printf("id: %d\n", id)

	p := &Product{}
	r.db.First(p, "id = ?", id)

	return *p
}

func (r *Repository) Delete(id int){

	log.Printf("delete id: %d\n", id)

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Delete(&Product{}, id)
}

func (r *Repository) Create(dto DTO){

	p := Product{
		Title: dto.Title,
		Description: dto.Description,
		Link: dto.Link,
		Image: dto.Image,
		Barcode: dto.Barcode,
	}

	r.db.Create(&p)
}

func (r *Repository) Update(id int, dto DTO){

	p := Product{
		Title: dto.Title,
		Description: dto.Description,
		Link: dto.Link,
		Image: dto.Image,
		Barcode: dto.Barcode,
	}

	r.db.Model(&Product{Id: id}).Updates(p)
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
	err := r.db.AutoMigrate(&Product{})
	if err != nil{
		return fmt.Errorf("migration of Product table failed: %v", err)
	}
	return nil
}

func Setup() (*Repository, error){

	repo := &Repository{}

	dsn := "root:product@tcp(127.0.0.1:7777)/product?multiStatements=true&parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	repo.db = db
	err = repo.Migrate()
	if err != nil{
		return nil, err
	}

	return repo, nil

}