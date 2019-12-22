package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emanpicar/minimart-api/db/entities"

	"github.com/emanpicar/minimart-api/product"

	"github.com/emanpicar/minimart-api/auth"
	"github.com/emanpicar/minimart-api/db"

	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"

	jwt "github.com/dgrijalva/jwt-go"
	gocache "github.com/patrickmn/go-cache"
)

type (
	Manager interface {
		GetAllCarts(r *http.Request) *[]CartCollection
		AddToCart(r *http.Request) (string, error)
		UpdateCart(r *http.Request, productID string) (string, error)
		DeleteCart(r *http.Request, productID string) (string, error)
	}

	cartHandler struct {
		cache     *gocache.Cache
		dbManager db.Manager
	}

	CartCollection struct {
		product.ProductCollection
		Quantity int `json:"quantity"`
	}

	CartReqBody struct {
		ID       uint `json:"id"`
		Quantity int  `json:"quantity"`
	}
)

func NewManager(dbManager db.Manager) Manager {
	return &cartHandler{
		cache:     gocache.New(time.Hour*1, time.Minute*10),
		dbManager: dbManager,
	}
}

func (c *cartHandler) GetAllCarts(r *http.Request) *[]CartCollection {
	user := c.getUserInContext(r)

	if data, ok := c.cache.Get(user.Username); ok {
		myCart := data.(*[]CartCollection)

		return myCart
	}

	return &[]CartCollection{}
}

func (c *cartHandler) AddToCart(r *http.Request) (string, error) {
	var reqData CartReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return "", err
	}

	product, err := c.dbManager.GetProductByID(reqData.ID)
	if err != nil {
		return "", err
	}

	user := c.getUserInContext(r)
	cachedData, ok := c.cache.Get(user.Username)
	if ok {
		cachedList := cachedData.(*[]CartCollection)
		if c.isProductIDInCache(cachedList, reqData.ID) {
			return "", errors.New("Product already in cart instead use PUT to update cart")
		}

		*cachedList = append(*cachedList, c.populateToCartCollection(product, reqData.Quantity))
		c.cache.Set(user.Username, cachedList, gocache.DefaultExpiration)
	} else {
		cartCol := &[]CartCollection{c.populateToCartCollection(product, reqData.Quantity)}
		c.cache.Set(user.Username, cartCol, gocache.DefaultExpiration)
	}

	return "Successfully added to cart", nil
}

func (c *cartHandler) UpdateCart(r *http.Request, productID string) (string, error) {
	var reqData CartReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return "", err
	}

	pID, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("Unable to parse productID:%v", productID)
	}

	user := c.getUserInContext(r)
	cachedData, ok := c.cache.Get(user.Username)
	if !ok || !c.isProductIDInCache(cachedData.(*[]CartCollection), uint(pID)) {
		return "", errors.New("Product does not exist in cart instead use POST to add in cart")
	}

	cartCol := c.updateCartCollection(reqData, uint(pID), cachedData.(*[]CartCollection))
	c.cache.Set(user.Username, cartCol, gocache.DefaultExpiration)

	return "Successfully updated in cart", nil
}

func (c *cartHandler) DeleteCart(r *http.Request, productID string) (string, error) {
	pID, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("Unable to parse productID:%v", productID)
	}

	user := c.getUserInContext(r)
	cachedData, ok := c.cache.Get(user.Username)

	if !ok || !c.isProductIDInCache(cachedData.(*[]CartCollection), uint(pID)) {
		return "", errors.New("Product does not exist in cart")
	}

	cartCol := c.deleteInCartCollection(uint(pID), cachedData.(*[]CartCollection))
	if len(*cartCol) > 0 {
		c.cache.Set(user.Username, cartCol, gocache.DefaultExpiration)
	} else {
		c.cache.Delete(user.Username)
	}

	return "Successfully deleted in cart", nil
}

func (c *cartHandler) getUserInContext(r *http.Request) auth.User {
	decoded := context.Get(r, "tokenClaims")
	var user auth.User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	return user
}

func (c *cartHandler) updateCartCollection(reqData CartReqBody, pID uint, cachedCol *[]CartCollection) *[]CartCollection {
	var cartCol []CartCollection

	for _, data := range *cachedCol {
		if data.ID == reqData.ID {
			data.Quantity = reqData.Quantity
		}

		cartCol = append(cartCol, data)
	}

	return &cartCol
}

func (c *cartHandler) deleteInCartCollection(productID uint, cachedCol *[]CartCollection) *[]CartCollection {
	var cartCol []CartCollection

	for _, data := range *cachedCol {
		if data.ID != productID {
			cartCol = append(cartCol, data)
		}
	}

	return &cartCol
}

func (c *cartHandler) populateToCartCollection(product *entities.ProductCollection, quantity int) CartCollection {
	data := CartCollection{Quantity: quantity}
	data.ID = product.ID
	data.Name = product.Name
	data.Slug = product.Slug

	if len(product.Images) > 0 {
		data.Image = product.Images[0].Value
	}

	if len(product.Offers) > 0 {
		data.SalesPrice = product.Offers[0].Price
	}

	return data
}

func (c *cartHandler) isProductIDInCache(cachedList *[]CartCollection, pID uint) bool {
	for _, c := range *cachedList {
		if c.ID == pID {
			return true
		}
	}

	return false
}
