package constructors

import (
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"gorm.io/driver/postgres"
	"os"

	"gorm.io/gorm"
)

type Db struct {
	DB *gorm.DB
}

func NewDb(cfg *configs.Configs) *Db {
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDns()), &gorm.Config{})

	if err != nil {
		os.Exit(0)
	}

	return &Db{DB: db}
}
