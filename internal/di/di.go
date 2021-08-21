package di

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/proviant-io/core/internal/apm"
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/db"
	"github.com/proviant-io/core/internal/pkg/image"
	"github.com/proviant-io/core/internal/pkg/shopping"
	"os"
)

type DI struct {
	Cfg          *config.Config
	Version      string
	ImageSaver   image.Saver
	Apm          apm.Apm
	ShoppingList *shopping.ListRepository
	ShoppingListItem *shopping.ItemRepository
}

func NewDI(d db.DB, cfg *config.Config, apm apm.Apm, version string) (*DI, error) {

	pool := &DI{}

	pool.Cfg = cfg
	pool.Version = version
	pool.Apm = apm

	shoppingListRepo, err := shopping.ListSetup(d)

	if err != nil {
		return nil, err
	}

	pool.ShoppingList = shoppingListRepo

	shoppingListItemRepo, err := shopping.ItemSetup(d)

	if err != nil {
		return nil, err
	}

	pool.ShoppingListItem = shoppingListItemRepo

	switch cfg.UserContent.Mode {
	case config.UserContentModeLocal:
		pool.ImageSaver = image.NewLocalSaver(cfg.UserContent.Location)
	case config.UserContentModeGCS:

		if cfg.API.GCS.JsonCredentialPath == "" {
			return nil, fmt.Errorf("credentials for GCS required")
		}

		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cfg.API.GCS.JsonCredentialPath)
		if err != nil {
			return nil, err
		}

		client, err := storage.NewClient(context.Background())
		if err != nil {
			return nil, err
		}
		pool.ImageSaver = image.NewGcsSaver(client, cfg.API.GCS.BucketName, cfg.API.GCS.ProjectId, cfg.UserContent.Location)

	default:
		return nil, fmt.Errorf("unsupported user content saver: %s", cfg.UserContent.Mode)
	}

	return pool, nil
}
