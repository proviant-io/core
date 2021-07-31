package di

import (
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/db"
)

type DI struct {
	Cfg      *config.Config
	Version string
}

func NewDI(d db.DB, cfg *config.Config, version string) (*DI, error) {

	pool := &DI{}

	pool.Cfg = cfg
	pool.Version = version

	return pool, nil
}
