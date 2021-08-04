package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"gitlab.com/emi2/mega/internal/app"
	"gitlab.com/emi2/mega/internal/app/mega"
	"gitlab.com/emi2/mega/internal/app/mega/middleware"
	"gitlab.com/emi2/mega/internal/app/mega/web/rest"
)

func setupRoutes(app *fiber.App) {
	// Jhipster endpoint for ROLE_USER
	app.Get("api/account", rest.GetAccount)                                     // getAccount
	app.Post("api/account", rest.SaveAccount)                                   // saveAccount
	app.Post("api/account/change-password", rest.ChangePassword)                // ChangePassword
	app.Post("​api​/account​/reset-password​/finish", rest.FinishPasswordReset) // finishPasswordReset
	app.Post("api​/account​/reset-password​/init", rest.RequestPasswordReset)   // requestPasswordReset

	// Account public endpoint
	app.Get("api/activate", rest.ActivateAccount)  // activateAccount
	app.Post("api/authenticate", rest.Login)       // isAuthenticated
	app.Post("api/register", rest.RegisterAccount) // registerAccount

	// User
	app.Get("api/authorities", middleware.HasAuthority("ROLE_ADMIN"), rest.GetAuthorities)
	app.Get("api/users", middleware.HasAuthority("ROLE_ADMIN"), rest.GetAllUser)
	app.Get("api/users/:id", middleware.HasAuthority("ROLE_ADMIN"), rest.GetUser)
	app.Post("api/users", middleware.HasAuthority("ROLE_ADMIN"), rest.NewUser)
	app.Put("api/users", middleware.HasAuthority("ROLE_ADMIN"), rest.UpdateUser)
	app.Delete("api/users/:id", middleware.HasAuthority("ROLE_ADMIN"), rest.DeleteUser)
}

// configure application runtime
func configure() {
	// koanf defautl values
	app.Config.Load(confmap.Provider(map[string]interface{}{
		"http.listen": ":3000",
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
	app.Config.Load(file.Provider("configs/admin-mgmt.yaml"), yaml.Parser())
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
	app.DBConn.AutoMigrate(&mega.User{})
	fmt.Println("Database Migrated")

	// JWT Middleware
	srv.Use(jwtware.New(jwtware.Config{
		// return true to skip middleware
		Filter: func(c *fiber.Ctx) bool {
			//log.Printf("Checking jwt on path %s", c.Path())
			return strings.HasPrefix(c.Path(), "/api/activate") ||
				strings.HasPrefix(c.Path(), "/api/authenticate") ||
				strings.HasPrefix(c.Path(), "/api/register") ||
				strings.HasPrefix(c.Path(), "/​api​/account​/reset-password​")
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			// declare locals:account to create audit log
			if c.Locals("user") != nil {
				token := c.Locals("user").(*jwt.Token)
				claims := token.Claims.(jwt.MapClaims)
				subject := claims["sub"].(string)
				c.Locals("account", subject)
			} else {
				c.Locals("account", "")
			}

			return c.Next()
		},
		SigningKey: []byte(app.Config.MustString("security.jwt-secret")),
	}))
	setupRoutes(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}
