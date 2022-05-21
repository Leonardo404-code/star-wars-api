package database

import (
	"ame-challenge/pkg/models"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDb() *gorm.DB {
	dsn := viper.GetString("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatalf("error in open database connectin: %v", err)
	}

	var planet models.Planet

	db.AutoMigrate(&planet)

	return db
}
