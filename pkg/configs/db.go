package configs

import (
	"os"

	"github.com/prateeksonii/shutter-api-go/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() error {
	var err error

	Db, err = gorm.Open(sqlite.Open(os.Getenv("DATABASE_URI")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&models.User{}, &models.Invite{})

	return nil
}
