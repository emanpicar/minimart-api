package entities

import (
	"github.com/jinzhu/gorm"
)

type (
	ProductCollection struct {
		ID     uint            `gorm:"unique;primary_key" json:"id"`
		Images []ProductImages `gorm:"foreignkey:ProductID" json:"images"`
		Name   string          `gorm:"type:varchar(100)" json:"name"`
		Offers []ProductOffers `gorm:"foreignkey:ProductID" json:"offers,omitempty"`
		Slug   string          `gorm:"type:varchar(100)" json:"slug"`
	}

	ProductOffers struct {
		gorm.Model `json:"-"`
		Price      float32 `gorm:"type:decimal(10,2)" json:"price"`
		ProductID  uint    `json:"-"`
	}

	ProductImages struct {
		gorm.Model
		Value     string `gorm:"type:varchar(500)"`
		ProductID uint
	}

	Credential struct {
		gorm.Model
		Username string `gorm:"type:varchar(40)"`
		Password string `gorm:"type:varchar(40)"`
	}
)

func (ProductCollection) TableName() string {
	return "product_collections"
}

func (ProductOffers) TableName() string {
	return "product_offers"
}

func (ProductImages) TableName() string {
	return "product_images"
}

func (Credential) TableName() string {
	return "credentials"
}
