package shopping

import (
	"fmt"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title     string `json:"title"`
	AccountId int    `json:"account_id" gorm:"default:0;index"`
}

func (List) TableName() string {
	return "shopping_lists"
}

type ListDTO struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type ListFilledDTO struct {
	Id    int       `json:"id"`
	Title string    `json:"title"`
	Items []ItemDTO `json:"items"`
}

type ListRepository struct {
	db db.DB
}

func (r *ListRepository) Get(id int, accountId int) (List, *errors.CustomError) {

	model := &List{}

	r.db.Connection().First(model, "id = ? and account_id = ?", id, accountId)

	if (*model).Id == 0 {
		return List{}, errors.NewErrNotFound(i18n.NewMessage("shopping list with id %d not found", id))
	}

	return *model, nil
}

func (r *ListRepository) GetAll(accountId int) []List {

	var models []List
	r.db.Connection().Where("account_id = ?", accountId).Find(&models)

	return models
}

func (r *ListRepository) Delete(id int, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return err
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *ListRepository) Create(dto ListDTO, accountId int) List {

	model := List{
		Title:     dto.Title,
		AccountId: accountId,
	}

	r.db.Connection().Create(&model)
	return model
}

func (r *ListRepository) Update(id int, dto ListDTO, accountId int) (List, *errors.CustomError) {

	model, err := r.Get(id, accountId)

	if err != nil {
		return List{}, err
	}

	model.Title = dto.Title

	r.db.Connection().Model(&List{Id: id}).Updates(&model)
	return model, nil
}

func ListToDTO(m List) ListDTO {
	return ListDTO{
		Id:    m.Id,
		Title: m.Title,
	}
}

func (r *ListRepository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&List{})
	if err != nil {
		return fmt.Errorf("migration of ShoppingList table failed: %v", err)
	}
	return nil
}

func ListSetup(d db.DB) (*ListRepository, error) {

	repo := &ListRepository{}

	repo.db = d

	err := repo.Migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil
}
