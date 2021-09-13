package shopping

import (
	"fmt"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Item struct {
	gorm.Model
	Id        int             `json:"id" gorm:"primaryKey;autoIncrement;"`
	ListId    int             `json:"list_id" gorm:"index"`
	Title     string          `json:"title"`
	Quantity  int             `json:"quantity"`
	Checked   bool            `json:"checked"`
	CheckedAt time.Time       `json:"checked_at"`
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(20,2);"`
	AccountId int             `json:"account_id" gorm:"default:0;index"`
}

func (Item) TableName() string {
	return "shopping_list_items"
}

type ItemDTO struct {
	Id        int             `json:"id"`
	ListId    int             `json:"list_id" gorm:"index"`
	Title     string          `json:"title"`
	Quantity  int             `json:"quantity"`
	Checked   bool            `json:"checked"`
	CheckedAt time.Time       `json:"checked_at"`
	Price     decimal.Decimal `json:"price"`
}

type ItemRepository struct {
	db db.DB
}

func (r *ItemRepository) Get(id int, accountId int) (Item, *errors.CustomError) {

	model := &Item{}

	r.db.Connection().First(model, "id = ? and account_id = ?", id, accountId)

	if (*model).Id == 0 {
		return Item{}, errors.NewErrNotFound(i18n.NewMessage("shopping list item with id %d not found", id))
	}

	return *model, nil
}

func (r *ItemRepository) GetAll(accountId int) []Item {

	var models []Item
	r.db.Connection().Where("account_id = ?", accountId).Find(&models)

	return models
}

func (r *ItemRepository) GetAllByList(listId int, accountId int) []Item {

	var models []Item
	r.db.Connection().Where("list_id = ? and account_id = ?", listId, accountId).Find(&models)

	return models
}

func (r *ItemRepository) Delete(id int, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return err
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *ItemRepository) Create(dto ItemDTO, accountId int) Item {

	model := Item{
		Title:     dto.Title,
		AccountId: accountId,
		ListId:    dto.ListId,
		Quantity:  dto.Quantity,
		Checked:   false,
		Price:     dto.Price,
	}

	r.db.Connection().Create(&model)
	return model
}

func (r *ItemRepository) Update(id int, dto ItemDTO, accountId int) (Item, *errors.CustomError) {

	model, err := r.Get(id, accountId)

	if err != nil {
		return Item{}, err
	}

	model.Title = dto.Title
	model.Quantity = dto.Quantity
	model.Price = dto.Price
	if dto.Checked {
		model.Checked = true
		model.CheckedAt = time.Now()
	} else {
		model.Checked = false
		// should be cleaning
		model.CheckedAt = time.Now()
	}

	r.db.Connection().Model(&Item{Id: id}).Updates(&model)

	if !model.Checked {
		r.db.Connection().Model(&model).Select("Checked").Updates(map[string]interface{}{"checked": false})
	}

	return model, nil
}

func (r *ItemRepository) updateChecked(id int, checked bool, accountId int) (Item, *errors.CustomError) {

	model, err := r.Get(id, accountId)

	if err != nil {
		return Item{}, err
	}

	model.Checked = checked

	r.db.Connection().Model(&model).Select("Checked").Updates(map[string]interface{}{"checked": model.Checked})
	r.db.Connection().Model(&model).Select("CheckedAt").Updates(map[string]interface{}{"checked_at": time.Now()})
	return model, nil
}

func (r *ItemRepository) Check(id int, accountId int) (Item, *errors.CustomError) {
	return r.updateChecked(id, true, accountId)
}

func (r *ItemRepository) Uncheck(id int, accountId int) (Item, *errors.CustomError) {
	return r.updateChecked(id, false, accountId)
}

func ItemToDTO(m Item) ItemDTO {
	return ItemDTO{
		Id:        m.Id,
		Title:     m.Title,
		ListId:    m.ListId,
		Quantity:  m.Quantity,
		Checked:   m.Checked,
		CheckedAt: m.CheckedAt,
		Price:     m.Price,
	}
}

func (r *ItemRepository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Item{})
	if err != nil {
		return fmt.Errorf("migration of ShoppingList table failed: %v", err)
	}
	return nil
}

func ItemSetup(d db.DB) (*ItemRepository, error) {

	repo := &ItemRepository{}

	repo.db = d

	err := repo.Migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil
}
