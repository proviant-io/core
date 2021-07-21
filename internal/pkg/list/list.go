package list

import (
	"fmt"
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title     string `json:"title"`
	AccountId int    `json:"account_id" gorm:"default:0"`
}

type DTO struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	AccountId int    `json:"account_id"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int, accountId int) (List, *errors.CustomError) {

	model := &List{}

	r.db.Connection().First(model, "id = ? and account_id = ?", id, accountId)

	if (*model).Id == 0 {
		return List{}, errors.NewErrNotFound(i18n.NewMessage("list with id %d not found", id))
	}

	return *model, nil
}

func (r *Repository) GetAll(accountId int) []List {

	var models []List
	r.db.Connection().Where("account_id = ?", accountId).Find(&models)

	return models
}

func (r *Repository) Delete(id int, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return err
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto DTO, accountId int) List {

	model := List{
		Title: dto.Title,
		AccountId: accountId,
	}

	r.db.Connection().Create(&model)
	return model
}

func (r *Repository) Update(id int, dto DTO, accountId int) (List, *errors.CustomError) {

	model, err := r.Get(id, accountId)

	if err != nil {
		return List{}, err
	}

	model.Title = dto.Title

	r.db.Connection().Model(&List{Id: id}).Updates(&model)
	return model, nil
}

func ModelToDTO(m List) DTO {
	return DTO{
		Id:    m.Id,
		Title: m.Title,
	}
}

func (r *Repository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&List{})
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
