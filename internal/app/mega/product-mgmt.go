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

	Title             string          `json:"title" query:"title"`
	Image             string          `json:"image" query:"image"`
	Images            pq.StringArray  `json:"images" gorm:"type:text[]" query:"images"`
	Description       string          `json:"description" query:"description"`
	Price             uint            `json:"price" query:"price"`
	Quantity          uint            `json:"quantity" query:"quantity"`
	ShortDescription  string          `json:"short_description" query:"short_description"`
	ProductCategoryID uint            `json:"product_category_id" query:"product_category_id"`
	ProductCategory   ProductCategory `json:"product_category" query:"product_category"`
}

// ProductCategory
type ProductCategory struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Title string `json:"title" query:"title"`
}
