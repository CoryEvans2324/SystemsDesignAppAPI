package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDatabase(dataSourceName string) *gorm.DB {
	newDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		CreateBatchSize: 1000,
	})

	if err != nil {
		panic(err)
	}

	return newDB
}
