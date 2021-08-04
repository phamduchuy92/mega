package rest

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/services"
	"gorm.io/gorm/clause"
)

// GetAllProducts return all items
func GetAllProducts(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	db.Preload("ProductCategory")

	category := c.Query("category")
	if category != "" {
		db.Joins("LEFT JOIN product_categories ON product_categories.id = products.product_category_id")
		db.Where("LOWER(product_categories.title) LIKE ?", "%"+strings.ToLower(category)+"%")
	}

	items := []mega.Product{}
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

// GetProduct return a single item with given ID
func GetProduct(c *fiber.Ctx) error {
	app.DBConn.Preload("ProductCategory")

	id := c.Params("id")
	item := mega.Product{}
	result := app.DBConn.Find(&item, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.ErrNotFound
	}

	return c.JSON(item)
}

// NewProduct create a new item
func NewProduct(c *fiber.Ctx) error {
	app.DBConn.Preload("ProductCategory")

	item := mega.Product{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	app.DBConn.Create(&item)

	return c.JSON(item)
}

// UpdateProduct update item info with given ID
func UpdateProduct(c *fiber.Ctx) error {
	app.DBConn.Preload("ProductCategory")

	item := mega.Product{}
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

// DeleteProduct delete the item with given ID
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	item := mega.Product{}
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
