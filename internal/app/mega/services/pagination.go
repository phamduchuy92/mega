package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Paginate provide a GORM pagination
func Paginate(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}

		pageSize, err := strconv.Atoi(c.Query("size"))
		if (err != nil) || (pageSize <= 0) {
			pageSize = 20
		}

		offset := page * pageSize

		sort := "id ASC"
		sortOpts := strings.Split(c.Query("sort", "id,asc"), ",")
		if len(sortOpts) == 2 {
			sortOpts[0] = db.NamingStrategy.ColumnName("", sortOpts[0])
			if sortOpts[1] != "asc" {
				sortOpts[1] = "desc"
			}
			sort = fmt.Sprintf("%s %s", sortOpts[0], sortOpts[1])
		}

		return db.Offset(offset).Limit(pageSize).Order(sort)
	}
}
