package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

// GetAllOrders return all items
func GetAllOrders(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	category := c.Query("category")
	if category != "" {
	}
	items := []mega.Order{}
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

// GetOrder return a single item with given ID
func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.Order{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewOrder create a new item
func NewOrder(c *fiber.Ctx) error {
	item := mega.Order{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateOrder update item info with given ID
func UpdateOrder(c *fiber.Ctx) error {
	item := mega.Order{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	result := app.DBConn.Find(&item, item.ID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	if dbErr := app.DBConn.Omit(clause.Associations).Save(&item); dbErr.Error != nil {
		return fiber.ErrBadGateway
	}

	return c.JSON(item)
}

// DeleteOrder delete the item with given ID
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.Order{}
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
