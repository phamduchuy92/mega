package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/knadh/koanf/providers/confmap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DBConn hold the connection to database
	DBConn *gorm.DB
	// QueryParam is the name of the query string key.
	QueryParam = "query"
)

// DatabaseDefaults set up default configuration for redis client
func DatabaseDefaults() {
	Config.Load(confmap.Provider(map[string]interface{}{
		"db.user":     "mega",
		"db.pass":     "mega",
		"db.host":     "localhost",
		"db.port":     5432,
		"db.name":     "mega",
		"db.sslmode":  "disable",
		"db.timezone": "Asia/Ho_Chi_Minh",
	}, "."), nil)

}

// DatabaseInit create the redis client based on koanf configuration
func DatabaseInit() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	var err error
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s TimeZone=%s",
		Config.String("db.user"),
		Config.String("db.pass"),
		Config.String("db.host"),
		Config.Int("db.port"),
		Config.String("db.name"),
		Config.String("db.sslmode"),
		Config.String("db.timezone"))
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected database")
}
