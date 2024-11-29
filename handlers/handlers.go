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

		// 쿠키 설정 (암호화된 사용자 ID 저장)
		c.Cookie(&fiber.Cookie{
			Name:     "auto_login",                       // 쿠키 이름
			Value:    "test",                             // 사용자 ID (암호화됨)
			Expires:  time.Now().Add(7 * 24 * time.Hour), // 7일 동안 유지
			HTTPOnly: true,                               // JavaScript에서 접근 금지
			Secure:   false,                              // HTTPS 환경에서는 true로 설정
		})

		if err := s.Save(); nil != err {
			return err
		}
		return c.JSON(user)
	}
	user := database.FindByName(req.Email)
	if nil != user && user.ID != 0 {
		s.Set("user", *user)
		// 쿠키 설정 (암호화된 사용자 ID 저장)
		c.Cookie(&fiber.Cookie{
			Name:     "auto_login",                       // 쿠키 이름
			Value:    user.Name,                          // 사용자 ID (암호화됨)
			Expires:  time.Now().Add(7 * 24 * time.Hour), // 7일 동안 유지
			HTTPOnly: true,                               // JavaScript에서 접근 금지
			Secure:   false,                              // HTTPS 환경에서는 true로 설정
		})

		_ = s.Save()
		return c.JSON(user)
	}
	return c.SendStatus(fiber.StatusNotFound)
}

// Logout 로그아웃
func Logout(c *fiber.Ctx) error {
	// 쿠키 삭제
	c.Cookie(&fiber.Cookie{
		Name:     "auto_login",
		Expires:  time.Now().Add(-time.Hour), // 과거로 설정하여 삭제
		HTTPOnly: true,
		Secure:   false,
	})

	c.Locals("session").(*session.Session).Destroy()

	return c.JSON(fiber.Map{"message": "Logged out successfully!"})
}
