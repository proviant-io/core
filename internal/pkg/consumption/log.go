package consumption

import (
	"fmt"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	gorm.Model
	Id         int   `json:"id" gorm:"primaryKey;autoIncrement;"`
	ProductId  int   `json:"product_id"`
	Quantity   uint  `json:"quantity"`
	ConsumedAt int64 `json:"consumed_at"`
	AccountId  int   `json:"account_id" gorm:"default:0;index"`
	UserId     int   `json:"user_id" gorm:"default:0"`
}

func (Log) TableName() string {
	return "consumption_logs"
}

type DTO struct {
	Id         int   `json:"id"`
	ProductId  int   `json:"product_id"`
	Quantity   uint  `json:"quantity"`
	ConsumedAt int64 `json:"consumed_at"`
	UserId     int   `json:"user_id"`
	AccountId  int   `json:"account_id"`
}

type ConsumeDTO struct {
	ProductId int  `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type LogRepository struct {
	db db.DB
}

func (r *LogRepository) Get(id int, accountId int) (Log, *errors.CustomError) {

	model := Log{}
	r.db.Connection().First(&model, "id = ? and account_id = ?", id, accountId)

	if (model).Id == 0 {
		return Log{}, errors.NewErrNotFound(i18n.NewMessage("consumption log entry with id %d not found", id))
	}

	return model, nil
}

func (r *LogRepository) GetAllByProductId(id int, accountId int) []Log {

	var s []Log
	r.db.Connection().Where("product_id = ? and account_id = ?", id, accountId).Order("consumed_at DESC").Find(&s)

	return s
}

func (r *LogRepository) Delete(id int, accountId int) *errors.CustomError {

	model, err := r.Get(id, accountId)

	if err != nil {
		return errors.NewErrNotFound(i18n.NewMessage("consumption log entry with id %d not found", id))
	}

	r.db.Connection().Unscoped().Delete(model, id)
	return nil
}

func (r *LogRepository) Create(dto ConsumeDTO, accountId int, userId int) Log {

	model := Log{
		Quantity:   dto.Quantity,
		ProductId:  dto.ProductId,
		ConsumedAt: time.Now().Unix(),
		AccountId:  accountId,
		UserId:     userId,
	}

	r.db.Connection().Create(&model)

	return model
}

func (r *LogRepository) Migrate() error {
	// Migrate the schema
	err := r.db.Connection().AutoMigrate(&Log{})
	if err != nil {
		return fmt.Errorf("migration of Stock table failed: %v", err)
	}
	return nil
}

func ModelToDTO(m Log) DTO {
	return DTO{
		Id:         m.Id,
		Quantity:   m.Quantity,
		ProductId:  m.ProductId,
		ConsumedAt: m.ConsumedAt,
		UserId:     m.UserId,
		AccountId:  m.AccountId,
	}
}

func LogSetup(d db.DB) (*LogRepository, error) {

	repo := &LogRepository{}

	repo.db = d
	err := repo.Migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil

}
