package routes

import (
	"encoding/json"
	"net/http"

	"github.com/emanpicar/minimart-api/auth"

	"github.com/emanpicar/minimart-api/cart"
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
		cartManager    cart.Manager
		authManager    auth.Manager
		router         *mux.Router
	}

	JsonMessage struct {
		Message string `json:"message"`
	}
)

func NewRouter(productManager product.Manager, cartManager cart.Manager, authManager auth.Manager) Router {
	routeHandler := &routeHandler{
		productManager: productManager,
		cartManager:    cartManager,
		authManager:    authManager,
	}

	return routeHandler.newRouter()
}

func (rh *routeHandler) newRouter() *mux.Router {
	router := mux.NewRouter()
	rh.registerRoutes(router)

	return router
}

func (rh *routeHandler) registerRoutes(router *mux.Router) {
	router.HandleFunc("/api/authenticate", rh.authenticate).Methods("POST")
	router.HandleFunc("/api/products", rh.authMiddleware(rh.getAllProducts)).Methods("GET")
	router.HandleFunc("/api/carts", rh.authMiddleware(rh.getAllCarts)).Methods("GET")
	router.HandleFunc("/api/carts", rh.authMiddleware(rh.addToCart)).Methods("POST")
	router.HandleFunc("/api/carts/{productId}", rh.authMiddleware(rh.updateCart)).Methods("PUT")
	router.HandleFunc("/api/carts/{productId}", rh.authMiddleware(rh.deleteCart)).Methods("DELETE")

	rh.router = router
}

func (rh *routeHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Authenticating user")

	w.Header().Set("Content-Type", "application/json")
	data, err := rh.authManager.Authenticate(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{err.Error()}), w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(data), w)
}

func (rh *routeHandler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Getting all products")

	w.Header().Set("Content-Type", "application/json")
	data := rh.productManager.GetAllProducts()

	rh.encodeError(json.NewEncoder(w).Encode(data), w)
}

func (rh *routeHandler) getAllCarts(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Getting all carts")

	w.Header().Set("Content-Type", "application/json")
	data := rh.cartManager.GetAllCarts(r)

	rh.encodeError(json.NewEncoder(w).Encode(data), w)
}

func (rh *routeHandler) addToCart(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infoln("Adding to cart")

	w.Header().Set("Content-Type", "application/json")
	data, err := rh.cartManager.AddToCart(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{err.Error()}), w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{data}), w)
}

func (rh *routeHandler) updateCart(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infof("Updating cart by id:%v", mux.Vars(r)["productId"])

	w.Header().Set("Content-Type", "application/json")
	data, err := rh.cartManager.UpdateCart(r, mux.Vars(r)["productId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{err.Error()}), w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{data}), w)
}

func (rh *routeHandler) deleteCart(w http.ResponseWriter, r *http.Request) {
	logger.Log.Infof("Deleting cart by id:%v", mux.Vars(r)["productId"])

	w.Header().Set("Content-Type", "application/json")
	data, err := rh.cartManager.DeleteCart(r, mux.Vars(r)["productId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{err.Error()}), w)
		return
	}

	rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{data}), w)
}

func (rh *routeHandler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := rh.authManager.ValidateRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rh.encodeError(json.NewEncoder(w).Encode(&JsonMessage{err.Error()}), w)
			return
		}

		next(w, r)
	})
}

func (rh *routeHandler) encodeError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
