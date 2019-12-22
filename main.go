package main

import (
	"fmt"

	"github.com/emanpicar/minimart-api/auth"
	"github.com/emanpicar/minimart-api/cart"
	"github.com/emanpicar/minimart-api/db"
	"github.com/emanpicar/minimart-api/logger"
	"github.com/emanpicar/minimart-api/product"
	"github.com/emanpicar/minimart-api/routes"
	"github.com/emanpicar/minimart-api/settings"

	"net/http"
)

func main() {
	logger.Init(settings.GetLogLevel())
	logger.Log.Infoln("Initializing Minimart API")

	dbManager := db.NewDBManager()
	productManager := product.NewManager(dbManager)
	cartManager := cart.NewManager(dbManager)
	authHandler := auth.NewManager()

	productManager.PopulateDefaultData()

	logger.Log.Fatal(http.ListenAndServeTLS(
		fmt.Sprintf("%v:%v", settings.GetServerHost(), settings.GetServerPort()),
		settings.GetServerPublicKey(),
		settings.GetServerPrivateKey(),
		routes.NewRouter(productManager, cartManager, authHandler),
	))
}
