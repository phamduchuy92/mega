package mega

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// UserDTO store information about user
type UserDTO struct {
	Id          uint           `json:"id" gorm:"primarykey"`
	Login       string         `json:"login"`
	Email       string         `json:"email"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	ImageUrl    string         `json:"imageUrl"`
	LangKey     string         `json:"langKey"`
	Activated   bool           `json:"activated"`
	Authorities pq.StringArray `json:"authorities" gorm:"type:text[]"`
}

// User store information about user from admin point of view
type User struct {
	UserDTO
	Password  string         `json:"password"`
	CreatedBy string         `json:"createdBy"`
	UpdatedBy string         `json:"updatedBy"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// Login is a value model for Login request
type Login struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

// PasswordChange is a model for password change request
type PasswordChange struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
