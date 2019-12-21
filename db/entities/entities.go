package entities

import (
	"github.com/jinzhu/gorm"
)

type (
	ProductCollection struct {
		ID     uint            `gorm:"unique;primary_key" json:"id"`
		Images []ProductImages `gorm:"foreignkey:ProductID" json:"images"`
		Name   string          `gorm:"type:varchar(100)" json:"name"`
		Offers []ProductOffers `gorm:"foreignkey:ProductID" json:"offers"`
		Slug   string          `gorm:"type:varchar(100)" json:"slug"`
	}

	ProductOffers struct {
		gorm.Model
		Price     float32 `gorm:"type:decimal(10,2)" json:"price"`
		ProductID uint
	}

	ProductImages struct {
		gorm.Model
		Value     string `gorm:"type:varchar(500)"`
		ProductID uint
	}
)

func (ProductCollection) TableName() string {
	return "product_collection"
}

func (ProductOffers) TableName() string {
	return "product_offers"
}

func (ProductImages) TableName() string {
	return "product_images"
}
