package rest

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetAllUser get all users
func GetAllUser(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	items := []mega.User{}
	result := db.Find(&items)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}
	c.Set("X-Total-Count", fmt.Sprint(result.RowsAffected))

	return c.JSON(items)
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.User{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewUser create new user
func NewUser(c *fiber.Ctx) error {
	item := mega.User{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}
	if strings.TrimSpace(item.Password) == "" {
		item.Password = RandStringBytesMaskImprSrcUnsafe(8)
		c.Response().Header.Add("X-Password", item.Password)
	}
	hash, err := services.HashPassword(item.Password)
	if err != nil {
		return fiber.ErrBadGateway
	}

	item.Password = hash

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	item := mega.User{}
	existsUser := mega.User{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	result := app.DBConn.Find(&existsUser, item.UserDTO.Id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	if strings.TrimSpace(item.Password) == "" {
		item.Password = existsUser.Password
	} else {
		hash, err := services.HashPassword(item.Password)
		if err != nil {
			return fiber.ErrBadGateway
		}

		item.Password = hash
	}

	if dbErr := app.DBConn.Omit(clause.Associations).Save(&item); dbErr.Error != nil {
		return fiber.ErrBadGateway
	}

	return c.JSON(item)
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.User{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	app.DBConn.Delete(&item)

	return c.SendStatus(fiber.StatusNoContent)
}

// GetAuthorities return list of authorities
func GetAuthorities(c *fiber.Ctx) error {
	return c.JSON(app.Config.Strings("security.authorities"))
}

// RandStringBytesMaskImprSrcUnsafe generate random string
func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
