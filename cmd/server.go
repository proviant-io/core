package main

import (
	"fmt"
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/di"
	"github.com/proviant-io/core/internal/http"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/proviant-io/core/internal/pkg/category"
	"github.com/proviant-io/core/internal/pkg/image"
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
		panic(err)
	}

	configPath := viper.GetString("config")

	if configPath == "" {
		configPath = "/app/default-config.yml"
	}

	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg, err := config.NewConfig(f)
	if err != nil {
		panic(err)
	}

	log.Println(cfg)

	var d db.DB

	switch cfg.Db.Driver {
	case config.DbDriverSqlite:
		d, err = db.NewSQLite(cfg.Db.Dsn)
		if err != nil {
			panic(err)
		}
	case config.DbDriverMysql:
		d, err = db.NewMySQL(cfg.Db.Dsn)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("unsupported db driver: %s", cfg.Db.Driver))
	}

	var imageSaver image.Saver
	switch cfg.UserContent.Mode {
	case config.UserContentModeLocal:
		imageSaver = image.NewLocalSaver(cfg.UserContent.Location)
	default:
		panic(fmt.Sprintf("unsupported user content saver: %s", cfg.UserContent.Mode))
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

	i, err := di.NewDI(d, cfg, Version)

	if err != nil {
		panic(err)
	}

	relationService := service.NewRelationService(productRepo, listRepo, categoryRepo, stockRepo, productCategoryRepo, imageSaver, *cfg)

	l := i18n.NewFileLocalizer()

	server := http.NewServer(productRepo, listRepo, categoryRepo, productCategoryRepo, stockRepo, relationService, l, i)

	hostPort := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Printf("starting server@%s\n", hostPort)

	err = server.Run(hostPort)

	if err != nil {
		panic(err)
	}
}
