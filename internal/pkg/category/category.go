package category

import (
	"fmt"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gorm.io/gorm"
)

type Category struct {
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

func (r *Repository) Get(id int) (Category, error) {

	model := &Category{}

	r.db.Connection().First(model, "id = ?", id)

	if (*model).Id == 0 {
		return Category{}, fmt.Errorf("category with id %d not found", id)
	}

	return *model, nil
}

func (r *Repository) GetAll() []Category {

	var categories []Category
	r.db.Connection().Find(&categories)

	return categories
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

	model := Category{
		Title: dto.Title,
	}

	r.db.Connection().Create(&model)
}

func (r *Repository) Update(id int, dto DTO) error {

	model, err := r.Get(id)

	if err != nil {
		return err
	}

	model.Title = dto.Title

	r.db.Connection().Model(&Category{Id: id}).Updates(model)
	return nil
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
