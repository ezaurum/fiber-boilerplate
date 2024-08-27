package handlers

import (
	"boilerplate/database"
	"boilerplate/models"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

func WebSocket(c *websocket.Conn) {
	// c.Locals property is added to the *websocket.Conn
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("id"))       // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
