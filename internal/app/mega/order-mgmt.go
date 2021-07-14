package mega

import (
	"time"

	"gorm.io/gorm"
)

// Order
type Order struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Title       string   `json:"title,omitempty" query:"title"`
	Price       uint     `json:"price,omitempty" query:"price"`
	Quantity    uint     `json:"quantity,omitempty" query:"quantity"`
	Description string   `json:"description,omitempty" query:"description"`
	Image       string   `json:"image,omitempty" query:"image"`
	Images      []string `json:"images,omitempty" query:"images"`
	Category    Category `json:"category,omitempty" query:"category"`
}

// OrderDetail
type OrderDetail struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Title string `json:"title,omitempty" query:"title"`
}
