package mega

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Product
type Product struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Name              string          `json:"name"`
	Image             string          `json:"image"`
	Images            pq.StringArray  `json:"images" gorm:"type:text[]"`
	Description       string          `json:"description"`
	Price             uint            `json:"price"`
	Quantity          uint            `json:"quantity"`
	ShortDescription  string          `json:"shortDescription"`
	ProductCategoryID uint            `json:"productCategoryId"`
	ProductCategory   ProductCategory `json:"productCategory"`
	OnSale            bool            `json:"onSale"`
	SalePrice         uint            `json:"salePrice"`
	UserID            uint            `json:"userId"`
}

// ProductCategory
type ProductCategory struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"userId"`
}
