package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CreateDatabase(dataSourceName string) {
	newDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		CreateBatchSize: 1000,
	})

	if err != nil {
		panic(err)
	}

	DB = newDB
}
