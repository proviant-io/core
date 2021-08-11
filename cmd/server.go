package main

import (
	"fmt"
	"github.com/proviant-io/core/internal/apm"
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/di"
	"github.com/proviant-io/core/internal/http"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/proviant-io/core/internal/pkg/category"
	"github.com/proviant-io/core/internal/pkg/list"
	"github.com/proviant-io/core/internal/pkg/product"
	"github.com/proviant-io/core/internal/pkg/product_category"
	"github.com/proviant-io/core/internal/pkg/service"
	"github.com/proviant-io/core/internal/pkg/stock"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	Version = "dev"
)

func main() {

	err := viper.BindEnv("config")

	if err != nil {
		log.Fatalln(err)
	}

	configPath := viper.GetString("config")

	if configPath == "" {
		configPath = "/app/default-config.yml"
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	cfg, err := config.NewConfig(f)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)

	realApm := apm.NewApm(cfg.APM)

	var d db.DB

	switch cfg.Db.Driver {
	case config.DbDriverSqlite:
		d, err = db.NewSQLite(cfg.Db.Dsn)
		if err != nil {
			log.Fatalln(err)
		}
	case config.DbDriverMysql:
		d, err = db.NewMySQL(cfg.Db.Dsn)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalln(fmt.Sprintf("unsupported db driver: %s", cfg.Db.Driver))
	}

	productRepo, err := product.Setup(d)

	if err != nil {
		log.Fatalln(err)
	}

	stockRepo, err := stock.Setup(d)

	if err != nil {
		log.Fatalln(err)
	}

	categoryRepo, err := category.Setup(d)

	if err != nil {
		log.Fatalln(err)
	}

	listRepo, err := list.Setup(d)

	if err != nil {
		log.Fatalln(err)
	}

	productCategoryRepo, err := product_category.Setup(d)

	if err != nil {
		log.Fatalln(err)
	}

	i, err := di.NewDI(d, cfg, realApm, Version)

	if err != nil {
		log.Fatalln(err)
	}

	relationService := service.NewRelationService(productRepo, listRepo, categoryRepo, stockRepo, productCategoryRepo, i, *cfg)

	l := i18n.NewFileLocalizer()

	server := http.NewServer(productRepo, listRepo, categoryRepo, productCategoryRepo, stockRepo, relationService, l, i)

	hostPort := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Printf("starting server@%s\n", hostPort)

	err = server.Run(hostPort)

	if err != nil {
		log.Fatalln(err)
	}
}
