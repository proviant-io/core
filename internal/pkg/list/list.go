package list

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gitlab.com/behind-the-fridge/product/internal/errors"
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type DTO struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Repository struct {
	db db.DB
}

func (r *Repository) Get(id int) (List, *errors.CustomError) {

	model := &List{}

	r.db.Connection().First(model, "id = ?", id)

	if (*model).Id == 0 {
		return List{}, errors.NewErrNotFound(fmt.Sprintf("list with id %d not found", id))
	}

	return *model, nil
}

func (r *Repository) GetAll() []List {

	var models []List
	r.db.Connection().Find(&models)

	return models
}

func (r *Repository) Delete(id int) *errors.CustomError {

	model, err := r.Get(id)

	if err != nil {
		return err
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *Repository) Create(dto DTO) List {

	model := List{
		Title: dto.Title,
	}

	r.db.Connection().Create(&model)
	return model
}

func (r *Repository) Update(id int, dto DTO) (List, *errors.CustomError) {

	model, err := r.Get(id)

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
