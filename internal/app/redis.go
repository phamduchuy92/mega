package app

import (
	redis "github.com/go-redis/redis/v8"
	"github.com/knadh/koanf/providers/confmap"
)

// RedisClient hold the active redis connection
var RedisClient *redis.Client

// RedisDefaults set up default configuration for redis client
func RedisDefaults() {
	Config.Load(confmap.Provider(map[string]interface{}{
		"redis.db":       0,
		"redis.address":  "localhost:6379",
		"redis.password": "",
		"redis.ttl":      24 * 60,
	}, "."), nil)
}

// RedisInit create the redis client based on koanf configuration
func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Config.String("redis.address"),
		Password: Config.String("redis.password"), // no password set
		DB:       Config.Int("redis.db"),          // use default DB
	})
}
