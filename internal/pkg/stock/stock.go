package stock

import (
	"fmt"
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	Id int `json:"id" gorm:"primaryKey;autoIncrement;"`
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

type ConsumeDTO struct {
	ProductId int `json:"product_id"`
	Quantity int `json:"quantity"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) (Stock, error){

	model := Stock{}
	r.db.Connection().First(&model, "id = ?", id)

	if (model).Id == 0 {
		return Stock{}, fmt.Errorf("stock with id %d not found", id)
	}

	return model, nil
}

func (r *Repository) GetAllByProductId(id int) []Stock{

	var s []Stock
	r.db.Connection().Where("product_id = ?", id).Order("expire ASC").Find(&s)

	return s
}

func (r *Repository) DeleteByProductId(id int) []Stock{

	var s []Stock
	r.db.Connection().Where("product_id = ?", id).Order("expire ASC").Unscoped().Delete(&Stock{})

	return s
}

func (r *Repository) Delete(id int) *errors.CustomError{

	model, err := r.Get(id)

	if err != nil {
		return errors.NewErrNotFound(i18n.NewMessage("stock with id %d not found", id))
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Consume(dto ConsumeDTO){
	// do something smart here

	quantityLeftToConsume := dto.Quantity

	models := r.GetAllByProductId(dto.ProductId)

	for _, model := range models {
		if model.Quantity < quantityLeftToConsume{
			quantityLeftToConsume -= model.Quantity
			r.Delete(model.Id)
		}else{
			model.Quantity -= quantityLeftToConsume
			r.Update(model.Id, DTO{
				ProductId: model.ProductId,
				Quantity:model.Quantity,
				Expire: model.Expire,
			})
			quantityLeftToConsume = 0
		}

		if quantityLeftToConsume == 0 {
			break
		}
	}
}

func (r *Repository) Add(dto DTO) Stock{
	return r.Create(dto)
}

func (r *Repository) Create(dto DTO) Stock{

	model := Stock{
		Quantity: dto.Quantity,
		ProductId: dto.ProductId,
		Expire: dto.Expire,
	}

	r.db.Connection().Create(&model)

	return model
}

func (r *Repository) Update(id int, dto DTO) (Stock, error){

	model, err := r.Get(id)

	if err != nil {
		return Stock{}, err
	}

	model.Quantity = dto.Quantity
	model.Expire = dto.Expire

	r.db.Connection().Model(&Stock{Id: id}).Updates(&model)
	return model, nil
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