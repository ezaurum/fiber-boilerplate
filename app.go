package main

import (
	"boilerplate/auth"
	"boilerplate/configs"
	"boilerplate/database"
	"boilerplate/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/swagger"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
        "github.com/gofiber/fiber/v2/middleware/cors"

	_ "boilerplate/docs"
)

func main() {
	config := configs.New()

	// Connected with database
	db := database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: config.GetBool("PREFORK"), // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	authz := auth.CasbinMiddleware(db)
	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users", handlers.UserCreate)
	v1.Delete("/users", authz.RequiresPermissions([]string{"user:create"}), handlers.UserCreate)
	v1.Get("/users/:id", handlers.UserGet)

	// Setup static files
	app.Static("/", "./static/public")

	// websocket
	app.Get("/ws/:id", websocket.New(handlers.WebSocket))

	app.Get("/swagger/*", swagger.HandlerDefault) // default
	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(config.ListenString())) // go run app.go -port=:3000
}
