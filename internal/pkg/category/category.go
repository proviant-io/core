package category

import (
	"fmt"
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title     string `json:"title"`
	AccountId int    `json:"account_id" gorm:"default:0"`
}

type DTO struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id, accountId int) (Category, *errors.CustomError) {

	model := &Category{}

	r.db.Connection().First(model, "id = ? and account_id = ?", id, accountId)

	if (*model).Id == 0 {
		return Category{}, errors.NewErrNotFound(i18n.NewMessage("category with id %d not found", id))
	}

	return *model, nil
}

func (r *Repository) GetByIds(ids []int, accountId int) []Category {

	var categories []Category
	r.db.Connection().Where("id IN (?) and account_id = ?", ids, accountId).Find(&categories)

	return categories
}

func (r *Repository) GetAll(accountId int) []Category {

	var categories []Category
	r.db.Connection().Where("account_id = ?", accountId).Find(&categories)

	return categories
}

func (r *Repository) Delete(id, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return err
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto DTO, accountId int) Category {

	model := Category{
		Title: dto.Title,
		AccountId: accountId,
	}

	r.db.Connection().Create(&model)
	return model
}

func (r *Repository) Update(id int, dto DTO, accountId int) (Category, *errors.CustomError) {

	model, err := r.Get(id, accountId)

	if err != nil {
		return Category{}, err
	}

	model.Title = dto.Title

	r.db.Connection().Model(&Category{Id: id}).Updates(&model)
	return model, nil
}

func ModelToDTO(m Category) DTO {
	return DTO{
		Id:    m.Id,
		Title: m.Title,
	}
}

func (r *Repository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Category{})
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
