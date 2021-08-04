package services

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserByUsername return user based on its login
func GetUserByUsername(login string, ctx context.Context) (*mega.User, error) {
	var user mega.User
	result := app.DBConn.Where("users.login = ?", login).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fiber.ErrNotFound
	}
	return &user, nil
}

// HashPassword hash the given password with bcrypt method
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}

// ValidToken check if current login match the given jwt subject
func ValidToken(t *jwt.Token, login string) bool {
	claims := t.Claims.(jwt.StandardClaims)
	return claims.Subject == login
}

// ValidUser validate one user retrieve from etcd
func ValidUser(login string, password string, ctx context.Context) bool {
	var user mega.User
	result := app.DBConn.Where("users.login = ?", login).Find(&user)
	if result.Error != nil {
		return false
	}
	if result.RowsAffected == 0 {
		return false
	}
	if user.Login == "" {
		return false
	}
	if !CheckPasswordHash(password, user.Password) {
		return false
	}
	return true
}
