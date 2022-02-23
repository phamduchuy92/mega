package rest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

// GetAllProductCategories return all items
func GetAllProductCategories(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	db2 := app.DBConn.Table("product_categories")

	item := mega.ProductCategory{}
	if err := c.QueryParser(&item); err != nil {
		return fiber.ErrBadRequest
	}
	if item.Name != "" {
		db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(item.Name)+"%")
		db2.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(item.Name)+"%")
		item.Name = ""
	}
	if item.Description != "" {
		db.Where("LOWER(description) LIKE ?", "%"+strings.ToLower(item.Description)+"%")
		db2.Where("LOWER(description) LIKE ?", "%"+strings.ToLower(item.Description)+"%")
		item.Description = ""
	}
	db2.Where("deleted_at IS NULL")

	rows := make([]mega.ProductCategory, 0)
	db.Where(item).Find(&rows)

	var count int64
	db2.Count(&count)

	c.Set("X-Total-Count", fmt.Sprint(count))

	return c.JSON(rows)
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

	user_id := c.Get("X-User-Id")
	if user_id != "" {
		u64, err := strconv.ParseUint(user_id, 0, 0)
		if err != nil {
			return fiber.ErrBadGateway
		}
		item.UserID = uint(u64)
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

	oldItem := mega.ProductCategory{}
	result := app.DBConn.Find(&oldItem, item.ID)
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
