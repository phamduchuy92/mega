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

// GetAllProducts return all items
func GetAllProducts(c *fiber.Ctx) error {
	db := app.DBConn.Scopes(services.Paginate(c))
	db2 := app.DBConn.Table("products")

	item := mega.Product{}
	if err := c.QueryParser(&item); err != nil {
		return fiber.ErrBadRequest
	}
	priceFrom := c.Query("priceFrom")
	priceTo := c.Query("priceTo")
	if priceFrom != "" && priceTo != "" {
		db.Where("price BETWEEN ? AND ?", priceFrom, priceTo)
		db2.Where("price BETWEEN ? AND ?", priceFrom, priceTo)
	}
	if item.Name != "" {
		db.Where("LOWER(name) LIKE ?", strings.ToLower("%"+item.Name+"%"))
		db2.Where("LOWER(name) LIKE ?", strings.ToLower("%"+item.Name+"%"))
		item.Name = ""
	}
	if item.ProductCategoryID != 0 {
		db.Where("product_category_id = ?", item.ProductCategoryID)
		db2.Where("product_category_id = ?", item.ProductCategoryID)
		item.ProductCategoryID = 0
	}
	db2.Where("deleted_at IS NULL")

	rows := make([]mega.Product, 0)
	db.Where(item).Find(&rows)

	var count int64
	db2.Count(&count)

	c.Set("X-Total-Count", fmt.Sprint(count))

	return c.JSON(rows)
}

// GetProduct return a single item with given ID
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	item := mega.Product{}
	result := app.DBConn.Preload("ProductCategory").Find(&item, id)
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
	item := mega.Product{}
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

// UpdateProduct update item info with given ID
func UpdateProduct(c *fiber.Ctx) error {
	item := mega.Product{}
	if err := c.BodyParser(&item); err != nil {
		return fiber.ErrBadRequest
	}

	oldItem := mega.Product{}
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
