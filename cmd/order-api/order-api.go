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
	// API Endpoint for Order
	app.Get("api/order", rest.GetAllOrders)
	app.Get("api/order/:id", rest.GetOrder)
	app.Post("api/order", rest.NewOrder)
	app.Put("api/order", rest.UpdateOrder)
	app.Delete("api/order/:id", rest.DeleteOrder)

	// API Endpoint for OrderDetail
	app.Get("api/order-detail", rest.GetAllOrderDetails)
	app.Get("api/order-detail/:id", rest.GetOrderDetail)
	app.Post("api/order-detail", rest.NewOrderDetail)
	app.Put("api/order-detail", rest.UpdateOrderDetail)
	app.Delete("api/order-detail/:id", rest.DeleteOrderDetail)
}

// configure application runtime
func configure() {
	// koanf defautl values
	app.Config.Load(confmap.Provider(map[string]interface{}{
		"http.listen": ":3002",
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
	app.Config.Load(file.Provider("configs/order-api.yaml"), yaml.Parser())
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
	app.DBConn.AutoMigrate(&mega.Order{})
	app.DBConn.AutoMigrate(&mega.OrderDetail{})
	fmt.Println("Database Migrated")

	setupRoutes(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}
