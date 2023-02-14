package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	redis "github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// AppConfig contains the global app configuration
// Config is loaded via envconfig
type AppConfig struct {
	Port      int    `default:"8080"`
	RedisHost string `split_words:"true" default:"localhost:6379"`
	RedisDB   int    `split_words:"true" default:"0"`
}

// loadConfig loads the app configuration from environment variables into a struct
func (c *AppConfig) loadConfig() {
	if err := envconfig.Process("VISITORS", c); err != nil {
		log.Fatalf("failed to load config. Got error: %s", err)
		os.Exit(1)
	}
}

// print prints out the values of the config
func (c AppConfig) print() {
	format := "%s is set to %v"

	log.Infof(format, "web port", c.Port)
	log.Infof(format, "redis host", c.RedisHost)
	log.Infof(format, "redis database", c.RedisDB)
}

// redisConnectionTest tests connection to redis database
// with retry and delay 
func redisConnectionTest(c *redis.Client, maxRetries int, delay int8) error {

	for i := 0; i <= maxRetries; i++ {
		_, err := c.Ping().Result()
		if  err == nil {
			break
		}

		log.Warnf("failed to connect to server. Got error: %s", err)
		log.Warnf("retrying database connection in %d second(s). Attempt %d of %d", delay, i, maxRetries)
		time.Sleep(time.Duration(delay) * time.Second)

		if i == maxRetries {
			return fmt.Errorf("%s maximum retries reached", err)
		}
	}

	log.Info("database connection established")
	return nil
}

// initCounter initialize the redis value for the counter
func initCounter(c *redis.Client, ctx string) {
	val := "1"
	err := c.Set(ctx, val, 0).Err()
	if err != nil {
		log.Fatalf("failed to initialize redis value. Got error: %s",err)
	}
}

// incrementCounter adds 1 to the counter value and stores it to redis
func incrementCounter(c *redis.Client, k string, cv string) {
	// convert value to int and increment
	newVal, err := strconv.ParseInt(cv, 10, 32)
	if err != nil {
		log.Fatalf("failed to convert current value to type int. Got error: %s", err)
	}
	newVal++

	// store new value to redis
	err = c.Set(k, newVal, 0).Err()
	if err != nil {
		log.Fatalf("failed to set new value. Got error: %s", err)
	}
}

func main() {
	// setup echo web framework with middlewares 
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Load app configuration from environment
	var config AppConfig
	config.loadConfig()
	config.print()

	redisKey := "count"

	// create redis connection client
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisHost,
		DB: config.RedisDB,
	})
	defer client.Close()

	// test the redis connection or exists
	if err := redisConnectionTest(client, 4, 1); err != nil {
		log.Fatalf("failed to establish database connection. Got error: %s",err)
		os.Exit(1)
	}

	// Check if redis entry exists, if not initialize key with 0
	_, err := client.Get(redisKey).Result()
	if err == redis.Nil {
		log.Infof("initialize redis value of key %s", redisKey)
		initCounter(client, redisKey)
	} else if err != nil {
		log.Fatalf("failed to check redis key. Got error:  %s", err)
	}

	// Root route handler
	// Display number of visitors from redis.
	e.GET("/", func(c echo.Context) error {
		val, err := client.Get(redisKey).Result()
		if err != nil {
			log.Info("error while reading from redis.")
		}

		str := fmt.Sprintf("Hey there! You are the %s visitor of this website.", val)
		incrementCounter(client, redisKey, val)

		return c.HTML(http.StatusOK, str)
	})

	// Health check route
	// Check if redis connection can be established
	e.GET("/health", func(c echo.Context) error {
		_, err := client.Ping().Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, struct{ Status string }{ Status: "Failure" })
		}
		return c.JSON(http.StatusOK, struct{ Status string }{ Status: "OK" })
	})

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
