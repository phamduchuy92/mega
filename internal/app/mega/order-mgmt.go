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

	CustomerID uint     `json:"customer_id"`
	Customer   Customer `json:"customer"`
}

// OrderDetail
type OrderDetail struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	OrderID   uint    `json:"order_id"`
	Order     Order   `json:"order"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
}
