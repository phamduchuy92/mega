package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

// GetAllOrderDetails return all items
func GetAllOrderDetails(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	db.Preload("Order").Preload("Order.Customer").Preload("Product")

	items := []mega.OrderDetail{}
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

// GetOrderDetail return a single item with given ID
func GetOrderDetail(c *fiber.Ctx) error {
	app.DBConn.Preload("Order").Preload("Order.Customer").Preload("Product")

	id := c.Params("id")
	item := mega.OrderDetail{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewOrderDetail create a new item
func NewOrderDetail(c *fiber.Ctx) error {
	app.DBConn.Preload("Order").Preload("Order.Customer").Preload("Product")

	item := mega.OrderDetail{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateOrderDetail update item info with given ID
func UpdateOrderDetail(c *fiber.Ctx) error {
	app.DBConn.Preload("Order").Preload("Order.Customer").Preload("Product")

	item := mega.OrderDetail{}
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

// DeleteOrderDetail delete the item with given ID
func DeleteOrderDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	item := mega.OrderDetail{}
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
