package mega

import (
	"time"

	"gorm.io/gorm"
)

// Customer
type Customer struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Name string `json:"name,omitempty" query:"name"`
}
