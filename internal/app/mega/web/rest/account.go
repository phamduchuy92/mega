package rest

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

// GetAccount implement api endpoint
func GetAccount(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	subject := claims["sub"].(string)
	account, err := services.GetUserByUsername(subject, c.Context())
	if err != nil {
		return err
	}

	return c.JSON(account.UserDTO)
}

// SaveAccount implement api endpoint
func SaveAccount(c *fiber.Ctx) error {
	var userDTO mega.UserDTO
	if err := c.BodyParser(&userDTO); err != nil {
		return fiber.ErrBadRequest
	}
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	subject := claims["sub"].(string)
	account, err := services.GetUserByUsername(subject, c.Context())
	if err != nil {
		return err
	}

	account.UserDTO = userDTO
	fmt.Printf("account %v", account)
	if dbErr := app.DBConn.Omit(clause.Associations).Save(&account); dbErr.Error != nil {
		return fiber.ErrBadGateway
	}
	return c.JSON(account.UserDTO)
}

// ChangePassword implement api endpoint
func ChangePassword(c *fiber.Ctx) error {
	var input mega.PasswordChange
	if err := c.BodyParser(&input); err != nil {
		return fiber.ErrBadRequest
	}
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	subject := claims["sub"].(string)
	account, err := services.GetUserByUsername(subject, c.Context())
	if err != nil {
		return err
	}

	if !services.CheckPasswordHash(input.CurrentPassword, account.Password) {
		return fiber.NewError(fiber.StatusExpectationFailed, "Current password does not match")
	}

	hash, err := services.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	account.Password = hash
	if dbErr := app.DBConn.Omit(clause.Associations).Save(&account); dbErr.Error != nil {
		return fiber.ErrBadGateway
	}
	return c.JSON(account.UserDTO)
}

// FinishPasswordReset implement api endpoint
func FinishPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// RequestPasswordReset implement api endpoint
func RequestPasswordReset(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// ActivateAccount implement api endpoint
func ActivateAccount(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

// IsAuthenticated implement api endpoint
func IsAuthenticated(c *fiber.Ctx) error {
	return Login(c)
}

// RegisterAccount implement api endpoint
func RegisterAccount(c *fiber.Ctx) error {
	var user mega.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest

	}
	result := app.DBConn.Find(&user, user.Id)
	if result.Error != nil {
		return fiber.ErrBadGateway
	}
	if result.RowsAffected > 0 {
		return fiber.ErrConflict
	}

	hash, err := services.HashPassword(user.Password)
	if err != nil {
		return fiber.ErrBadGateway

	}

	user.Password = hash
	if dbErr := app.DBConn.Omit(clause.Associations).Save(&user); dbErr.Error != nil {
		return fiber.ErrBadGateway
	}
	return c.JSON(user.UserDTO)
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	var input mega.Login

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	ud, err := services.GetUserByUsername(input.Username, c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid login")
	}
	if !ud.Activated {
		return fiber.NewError(fiber.StatusExpectationFailed, "Account is not activated")
	}
	if !services.CheckPasswordHash(input.Password, ud.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = ud.Login
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["authorities"] = ud.Authorities

	t, err := token.SignedString([]byte(app.Config.String("security.jwt-secret")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to generate token")
	}
	return c.JSON(fiber.Map{"id_token": t, "login": ud.Login, "email": ud.Email, "id": ud.Id})
}
