package stock

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	Id int `json:"id"`
	ProductId int `json:"product_id"`
	Quantity int `json:"quantity"`
	Expire int `json:"expire"`
}

type DTO struct {
	ProductId int `json:"product_id"`
	Quantity int `json:"quantity"`
	Expire int `json:"expire"`
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) Get(id int) *Stock{

	s := &Stock{}
	r.db.First(s, "id = ?", id)

	return s
}

func (r *Repository) GetAllByProductId(id int) []Stock{

	var s []Stock
	r.db.Where("product_id = ?", id).Find(&s)

	return s
}

func (r *Repository) Delete(id int){

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Delete(&Stock{}, id)
}

func (r *Repository) Create(dto DTO){

	s := Stock{
		Quantity: dto.Quantity,
		ProductId: dto.ProductId,
		Expire: dto.Expire,
	}

	r.db.Create(&s)
}

func (r *Repository) Update(id int, dto DTO){

	s := Stock{
		Quantity: dto.Quantity,
		ProductId: dto.ProductId,
		Expire: dto.Expire,
	}

	r.db.Model(&Stock{Id: id}).Updates(s)
}

func (r *Repository) Migrate() error{
	// Migrate the schema
	err := r.db.AutoMigrate(&Stock{})
	if err != nil{
		return fmt.Errorf("migration of Stock table failed: %v", err)
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