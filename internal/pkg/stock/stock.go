package stock

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
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
	Id int `json:"id"`
	ProductId int `json:"product_id"`
	Quantity int `json:"quantity"`
	Expire int `json:"expire"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) *Stock{

	s := &Stock{}
	r.db.Connection().First(s, "id = ?", id)

	return s
}

func (r *Repository) GetAllByProductId(id int) []Stock{

	var s []Stock
	r.db.Connection().Where("product_id = ?", id).Find(&s)

	return s
}

func (r *Repository) Delete(id int){

	//db.Unscoped().Delete(&order) to delete permanently
	r.db.Connection().Delete(&Stock{}, id)
}

func (r *Repository) Consume(dto DTO){
	// do something smart here
}

func (r *Repository) Add(dto DTO){
	r.Create(dto)
}

func (r *Repository) Create(dto DTO){

	s := Stock{
		Quantity: dto.Quantity,
		ProductId: dto.ProductId,
		Expire: dto.Expire,
	}

	r.db.Connection().Create(&s)
}

func (r *Repository) Update(id int, dto DTO){

	s := Stock{
		Quantity: dto.Quantity,
		ProductId: dto.ProductId,
		Expire: dto.Expire,
	}

	r.db.Connection().Model(&Stock{Id: id}).Updates(s)
}

func (r *Repository) Migrate() error{
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Stock{})
	if err != nil{
		return fmt.Errorf("migration of Stock table failed: %v", err)
	}
	return nil
}

func ModelToDTO(m Stock) DTO {
	return DTO{
		Id: m.Id,
		Quantity: m.Quantity,
		ProductId: m.ProductId,
		Expire: m.Expire,
	}
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