package routes

import (
	"encoding/json"
	"net/http"

	"github.com/emanpicar/minimart-api/logger"

	"github.com/emanpicar/minimart-api/product"

	"github.com/gorilla/mux"
)

type (
	Router interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}

	routeHandler struct {
		productManager product.Manager
		router         *mux.Router
	}

	JsonMessage struct {
		Message string `json:"message"`
	}
)

func NewRouter(productManager product.Manager) Router {
	routeHandler := &routeHandler{productManager: productManager}

	return routeHandler.newRouter()
}

func (rh *routeHandler) newRouter() *mux.Router {
	router := mux.NewRouter()
	rh.registerRoutes(router)

	return router
}

func (rh *routeHandler) registerRoutes(router *mux.Router) {
	router.HandleFunc("/api/products", rh.getAllProducts).Methods("GET")
	router.HandleFunc("/api/carts", rh.getAllProducts).Methods("GET")
	router.HandleFunc("/api/add-cart", rh.getAllProducts).Methods("POST")
	router.HandleFunc("/api/update-cart", rh.getAllProducts).Methods("PUT")
	router.HandleFunc("/api/remove-cart", rh.getAllProducts).Methods("DELETE")

	rh.router = router
}

func (rh *routeHandler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Getting all products")

	w.Header().Set("Content-Type", "application/json")
	data := rh.productManager.GetAllProducts()

	rh.encodeError(json.NewEncoder(w).Encode(data), w)
}

func (rh *routeHandler) encodeError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
