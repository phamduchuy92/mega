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

	CustomerID uint     `json:"customer_id" query:"customer_id"`
	Customer   Customer `json:"customer" query:"customer"`
}

// OrderDetail
type OrderDetail struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	OrderID   uint    `json:"order_id" query:"order_id"`
	Order     Order   `json:"order" query:"order"`
	ProductID uint    `json:"product_id" query:"product_id"`
	Product   Product `json:"product" query:"product"`
	Quantity  uint    `json:"quantity" query:"quantity"`
}
