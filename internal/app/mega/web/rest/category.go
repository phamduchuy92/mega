package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega-backend/internal/app"
	"gitlab.com/emi2/mega-backend/internal/app/mega"
	"gitlab.com/emi2/mega-backend/internal/app/mega/services/utils"
	"gorm.io/gorm/clause"
)

// GetAllCategories return all items
func GetAllCategories(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(utils.Paginate(c))

	items := []mega.Category{}
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

// GetCategory return a single item with given ID
func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.Category{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewCategory create a new item
func NewCategory(c *fiber.Ctx) error {
	item := mega.Category{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateCategory update item info with given ID
func UpdateCategory(c *fiber.Ctx) error {
	item := mega.Category{}
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

// DeleteCategory delete the item with given ID
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	item := mega.Category{}
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
