package db

import (
	"fmt"

	"github.com/emanpicar/minimart-api/settings"

	"github.com/emanpicar/minimart-api/db/entities"
	"github.com/emanpicar/minimart-api/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Manager interface {
		BatchFirstOrCreate(prodCollection *[]entities.ProductCollection)
		GetProductCollection() *[]entities.ProductCollection
	}

	dbHandler struct {
		database *gorm.DB
	}
)

func NewDBManager() Manager {
	dbHandler := &dbHandler{}
	dbHandler.connect(gorm.Open)
	dbHandler.migrateTables()

	return dbHandler
}

func (dbHandler *dbHandler) connect(openConnection func(dialect string, args ...interface{}) (db *gorm.DB, err error)) {
	logger.Log.Infoln("Establishing connection to DB")

	var err error
	dbHandler.database, err = openConnection("postgres", fmt.Sprintf("host=%v port=%v user=%v dbname=minimart_db password=%v sslmode=disable",
		settings.GetDBHost(), settings.GetDBPort(), settings.GetDBUser(), settings.GetDBPass(),
	))

	if err != nil {
		logger.Log.Fatalln(err)
	}

	logger.Log.Infoln("Successfully connected to DB")
}

func (dbHandler *dbHandler) migrateTables() {
	dbHandler.database.AutoMigrate(&entities.ProductCollection{})
	dbHandler.database.AutoMigrate(&entities.ProductOffers{}).AddForeignKey("product_id", "product_collection(id)", "CASCADE", "CASCADE")
	dbHandler.database.AutoMigrate(&entities.ProductImages{}).AddForeignKey("product_id", "product_collection(id)", "CASCADE", "CASCADE")
}

func (dbHandler *dbHandler) BatchFirstOrCreate(prodCollection *[]entities.ProductCollection) {
	for _, product := range *prodCollection {
		dbHandler.database.FirstOrCreate(&product, entities.ProductCollection{})
	}
}

func (dbHandler *dbHandler) GetProductCollection() *[]entities.ProductCollection {
	var data []entities.ProductCollection
	dbHandler.database.Set("gorm:auto_preload", true).Find(&data)

	return &data
}
