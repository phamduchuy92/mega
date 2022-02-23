package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/web/rest"
)

func setupRoutes(app *fiber.App) {
	// API Endpoint for Product
	app.Get("api/products", rest.GetAllProducts)
	app.Get("api/products/:id", rest.GetProduct)
	app.Post("api/products", rest.NewProduct)
	app.Put("api/products", rest.UpdateProduct)
	app.Delete("api/products/:id", rest.DeleteProduct)

	// API Endpoint for ProductCategory
	app.Get("api/product-categories", rest.GetAllProductCategories)
	app.Get("api/product-categories/:id", rest.GetProductCategory)
	app.Post("api/product-categories", rest.NewProductCategory)
	app.Put("api/product-categories", rest.UpdateProductCategory)
	app.Delete("api/product-categories/:id", rest.DeleteProductCategory)
}

// configure application runtime
func configure() {
	// koanf defautl values
	app.Config.Load(confmap.Provider(map[string]interface{}{
		"http.listen": ":3001",
		// + db settings
		"db.user":     "mega",
		"db.pass":     "mega",
		"db.host":     "localhost",
		"db.port":     5432,
		"db.name":     "mega",
		"db.sslmode":  "disable",
		"db.timezone": "Asia/Ho_Chi_Minh",
	}, "."), nil)

	// override configuration with YAML
	app.Config.Load(file.Provider("configs/product-api.yaml"), yaml.Parser())
}

// main function
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configure()
	srv := fiber.New(fiber.Config{
		ErrorHandler: app.ProblemJSONErrorHandle,
	})
	srv.Use(cors.New())
	srv.Use(logger.New())

	app.DatabaseInit()
	// specific tables
	app.DBConn.AutoMigrate(&mega.ProductCategory{})
	app.DBConn.AutoMigrate(&mega.Product{})
	fmt.Println("Database Migrated")

	setupRoutes(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}
