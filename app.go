package main

import (
	"boilerplate/auth"
	"boilerplate/configs"
	"boilerplate/database"
	"boilerplate/handlers"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/redis/v3"
	"log"
	"runtime"
	"time"
)

func main() {
	config := configs.New()

	// Connected with database
	db := database.Connect(database.Connection{
		Host:     config.GetString("DB_HOST"),
		Port:     config.GetInt("DB_PORT"),
		User:     config.GetString("DB_USER"),
		Password: config.GetString("DB_PASS"),
		Name:     config.GetString("DB_NAME"),
	})

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: config.GetBool("PREFORK"), // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(swagger.New(
		swagger.Config{
			BasePath: "/docs/",
			Path:     "v1",
			FilePath: "./docs/swagger.json",
		},
	))

	// Initialize redis config from .env
	storage := redis.New(redis.Config{
		Host:      config.GetString("REDIS_HOST"),
		Port:      config.GetInt("REDIS_PORT"),
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
	sessionStore := session.New(session.Config{
		Storage: storage,
	})
	app.Use(func(c *fiber.Ctx) error {
		get, err := sessionStore.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", get)
		return c.Next()
	})
	const HeaderName = "X-CSRF-Token"
	app.Use(csrf.New(
		csrf.Config{
			KeyLookup:         "header:" + HeaderName,
			CookieName:        "__Host-csrf_",
			CookieSameSite:    "Lax",
			CookieSecure:      true,
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			Expiration:        1 * time.Hour,
			KeyGenerator:      utils.UUIDv4,
			Extractor:         csrf.CsrfFromHeader(HeaderName),
			Session:           sessionStore,
			SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
		}))

	authz := auth.CasbinMiddleware(db)
	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users",
		authz.RequiresPermissions([]string{"user:create"}),
		handlers.UserCreate)
	v1.Delete("/users", authz.RequiresPermissions([]string{"user:delete"}), handlers.UserDelete)
	v1.Get("/users/:id", handlers.UserGet)

	// Setup static files
	app.Static("/", "./static/public")

	// websocket
	app.Get("/ws/:id", websocket.New(handlers.WebSocket))

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on
	log.Fatal(app.Listen(config.ListenString()))
}
