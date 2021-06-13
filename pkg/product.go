package pkg

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`
	Image string `json:"image"`
	Barcode string `json:"barcode"`
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) Get(id int) *Product{

	log.Printf("id: %d\n", id)

	p := &Product{}
	r.db.First(p, "id = ?", id)

	return p
}

func (r *Repository) Migrate() error{
	// Migrate the schema
	err := r.db.AutoMigrate(&Product{})
	if err != nil{
		return fmt.Errorf("migration of Product table failed: %v", err)
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