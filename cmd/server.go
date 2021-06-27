package main

import (
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/http"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/service"
	"github.com/brushknight/proviant/internal/pkg/stock"
)

var SqliteLocation = "pantry.db"

func main() {

	d, err := db.NewSQLite(SqliteLocation)

	if err != nil {
		panic(err)
	}

	productRepo, err := product.Setup(d)

	if err != nil {
		panic(err)
	}

	stockRepo, err := stock.Setup(d)

	if err != nil {
		panic(err)
	}

	categoryRepo, err := category.Setup(d)

	if err != nil {
		panic(err)
	}

	listRepo, err := list.Setup(d)

	if err != nil {
		panic(err)
	}

	productCategoryRepo, err := product_category.Setup(d)

	if err != nil {
		panic(err)
	}

	relationService := service.NewRelationService(productRepo, listRepo, categoryRepo, stockRepo, productCategoryRepo)

	server := http.NewServer(productRepo, listRepo, categoryRepo, productCategoryRepo, stockRepo, relationService)

	err = server.Run("0.0.0.0:80")

	if err != nil {
		panic(err)
	}
}
