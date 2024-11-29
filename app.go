package main

import (
	"boilerplate/auth"
	"boilerplate/configs"
	"boilerplate/database"
	"boilerplate/handlers"
	"boilerplate/models"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/template/html/v2"
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
		Views:   html.New("./static/private", ".html"),
	})

	// 쿠키 암호화 미들웨어 (보안을 위해 암호화 키 설정)
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: "my-secret-key-32bytes", // 32바이트 비밀 키
	}))

	sessionExpiration := 30 * time.Minute

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
		Storage:    storage,
		Expiration: sessionExpiration,
	})
	app.Use(func(c *fiber.Ctx) error {
		get, err := sessionStore.Get(c)
		if err != nil {
			return err
		}
		if get.Fresh() {
			get.Set("user", models.User{})
		} else {
			// 세션 만료 시간 연장
			get.SetExpiry(sessionExpiration)
		}
		c.Locals("session", get)
		_ = get.Save()
		return c.Next()
	})
	const HeaderName = "X-CSRF-Token"
	const CookieName = "__Host-csrf_"
	app.Use(csrf.New(
		csrf.Config{
			KeyLookup:         "header:" + HeaderName,
			CookieName:        CookieName,
			CookieSameSite:    "Lax",
			CookieSecure:      true,
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			Expiration:        1 * time.Hour,
			KeyGenerator:      utils.UUIDv4,
			Extractor: func(c *fiber.Ctx) (string, error) {
				if s, e := csrf.CsrfFromQuery("_csrf")(c); nil == e {
					return s, e
				}
				if s, e := csrf.CsrfFromParam("_csrf")(c); nil == e {
					return s, e
				}
				if s, e := csrf.CsrfFromHeader(HeaderName)(c); nil == e {
					return s, e
				}
				if s, e := csrf.CsrfFromForm("_csrf")(c); nil == e {
					return s, e
				}
				return csrf.CsrfFromCookie(CookieName)(c)
			},
			Session:           sessionStore,
			SessionKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
			ContextKey:        "fiber.csrf.token_string",
		}))

	_ = auth.CasbinMiddleware(db)
	// Create a /api/v1 endpoint
	// auto login
	v1 := app.Group("/api/v1", func(c *fiber.Ctx) error {
		username := c.Cookies("auto_login")
		if username != "" {
			user := database.FindByName(username)
			if user != nil {
				s := c.Locals("session").(*session.Session)
				s.Set("user", *user)
				_ = s.Save()
			}
		}
		return c.Next()
	})

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users", func(c *fiber.Ctx) error {
		s := c.Locals("session")
		if s == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		ses := s.(*session.Session)
		u := ses.Get("user")
		if u == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		user := u.(models.User)
		if user.Name != "test" {
			return c.SendStatus(fiber.StatusForbidden)
		}
		log.Println("user:", user)
		return c.Next()
	},
		handlers.UserCreate)
	v1.Delete("/users", handlers.UserDelete)
	v1.Get("/users/:id", handlers.UserGet)

	v1.Post("/login", handlers.Login)
	v1.Post("/logout", handlers.Logout)

	// Setup static files
	app.Static("/", "./static/public")

	app.Get("/login", func(c *fiber.Ctx) error {
		// 템플릿 렌더링 시 토큰 전달
		return c.Render("login", fiber.Map{
			"csrfToken": c.Locals("fiber.csrf.token_string"),
		})
	})

	// websocket
	app.Get("/ws/:id", websocket.New(handlers.WebSocket))

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on
	log.Fatal(app.Listen(config.ListenString()))
}
