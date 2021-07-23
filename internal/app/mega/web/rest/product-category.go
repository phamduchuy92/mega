package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services/utils"
	"gorm.io/gorm/clause"
)

// GetAllProductCategories return all items
func GetAllProductCategories(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(utils.Paginate(c))

	items := []mega.ProductCategory{}
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

// GetProductCategory return a single item with given ID
func GetProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	item := mega.ProductCategory{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewProductCategory create a new item
func NewProductCategory(c *fiber.Ctx) error {
	item := mega.ProductCategory{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateProductCategory update item info with given ID
func UpdateProductCategory(c *fiber.Ctx) error {
	item := mega.ProductCategory{}
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
		return dbErr.Error
	}

	return c.JSON(item)
}

// DeleteProductCategory delete the item with given ID
func DeleteProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	item := mega.ProductCategory{}
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
