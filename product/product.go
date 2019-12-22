package product

import (
	"encoding/json"
	"io/ioutil"

	"github.com/emanpicar/minimart-api/db"
	"github.com/emanpicar/minimart-api/db/entities"
	"github.com/emanpicar/minimart-api/logger"
)

type (
	Manager interface {
		PopulateDefaultData()
		GetAllProducts() *[]ProductCollection
	}

	productHandler struct {
		dbManager db.Manager
	}

	ProductCollection struct {
		entities.ProductCollection
		Images     []string `json:"images,omitempty"`
		Image      string   `json:"image"`
		SalesPrice float32  `json:"sales_price"`
	}
)

func NewManager(dbManager db.Manager) Manager {
	return &productHandler{dbManager}
}

func (p *productHandler) PopulateDefaultData() {
	var products []ProductCollection

	bytesData, err := ioutil.ReadFile("./jsondata/products.json")
	if err != nil {
		logger.Log.Errorf("Unable to create default data due to: %v", err)
	}

	if err = json.Unmarshal(bytesData, &products); err != nil {
		logger.Log.Errorf("Unable to create default data due to: %v", err)
	}

	productsModel := p.populateCollectionForModel(&products)

	p.dbManager.BatchFirstOrCreate(productsModel)
}

func (p *productHandler) GetAllProducts() *[]ProductCollection {
	productList := p.dbManager.GetProductCollection()
	jsonReadyList := p.populateCollectionForJSON(productList)

	return jsonReadyList
}

func (p *productHandler) populateCollectionForModel(products *[]ProductCollection) *[]entities.ProductCollection {
	var dbEntity []entities.ProductCollection

	for _, product := range *products {
		dbEntity = append(dbEntity, entities.ProductCollection{
			ID:     product.ID,
			Images: p.populateArrayImgForModel(product.Images),
			Name:   product.Name,
			Offers: product.Offers,
			Slug:   product.Slug,
		})
	}

	return &dbEntity
}

func (p *productHandler) populateArrayImgForModel(images []string) []entities.ProductImages {
	var productImages []entities.ProductImages

	for _, image := range images {
		productImages = append(productImages, entities.ProductImages{
			Value: image,
		})
	}

	return productImages
}

func (p *productHandler) populateCollectionForJSON(products *[]entities.ProductCollection) *[]ProductCollection {
	var dbEntity []ProductCollection

	for _, product := range *products {
		var img string
		if len(product.Images) > 0 {
			img = product.Images[0].Value
		}

		var salesPrice float32
		if len(product.Offers) > 0 {
			salesPrice = product.Offers[0].Price
		}

		dbEntity = append(dbEntity, ProductCollection{
			ProductCollection: entities.ProductCollection{
				ID:   product.ID,
				Name: product.Name,
				Slug: product.Slug,
			},
			Image:      img,
			SalesPrice: salesPrice,
		})
	}

	return &dbEntity
}
