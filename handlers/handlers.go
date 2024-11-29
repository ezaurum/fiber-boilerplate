package handlers

import (
	"boilerplate/database"
	"boilerplate/models"
	"encoding/gob"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"
	"time"

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

	log.Println("User created:", user.Name)

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

// UserGet returns a user
// @Summary Get a user
// @Description Get a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func UserGet(ctx *fiber.Ctx) error {
	var data = make(map[string]interface{})
	data["id"] = ctx.Params("id")
	data["name"] = "John Doe"
	return ctx.JSON(data)
}

func UserDelete(c *fiber.Ctx) error {
	return c.SendString("User deleted")
}

func init() {
	// User 타입 등록
	gob.Register(models.User{})
}

type LoginRequest struct {
	Email string `json:"email" form:"email"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); nil != err {
		return err
	}
	if req.Email == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	s := c.Locals("session").(*session.Session)
	if req.Email == "test" {
		user := models.User{
			Model: models.Model{
				ID: 0,
				UpdateModel: models.UpdateModel{
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: nil,
				},
			},
			Name:  "test",
			Posts: nil,
		}
		s.Set("user", user)
		if err := s.Save(); nil != err {
			return err
		}
		return c.JSON(user)
	}
	user := database.FindByName(req.Email)
	if nil != user && user.ID != 0 {
		s.Set("user", *user)
		_ = s.Save()
		return c.JSON(user)
	}
	return c.SendStatus(fiber.StatusNotFound)
}
