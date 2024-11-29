package auth

import (
	"boilerplate/models"
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	casbinmiddle "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

//go:embed rbac.casbin
var casbinModel string

func enforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	m := model.NewModel()
	if err := m.LoadModelFromText(casbinModel); nil != err {
		return nil, err
	}
	byDB, err0 := gormadapter.NewAdapterByDB(db)
	if err0 != nil {
		log.Fatalf("An error '%s' was not expected when creating a new casbin adapter", err0)
	}
	if e, err := casbin.NewEnforcer(); nil != err {
		return nil, err
	} else if err = e.InitWithModelAndAdapter(m, byDB); nil != err {
		return e, err
	} else {
		return e, nil
	}
}

func CasbinMiddleware(db *gorm.DB) *casbinmiddle.Middleware {
	if enf, err := enforcer(db); nil != err {
		log.Fatalf("An error '%s' was not expected when creating a new casbin enforcer", err)
		return nil
	} else {
		return casbinmiddle.New(casbinmiddle.Config{
			Enforcer: enf,
			Lookup: func(c *fiber.Ctx) string {
				s := c.Locals("session").(*session.Session)
				if nil == s {
					return ""
				}
				u := s.Get("user")
				if nil == u {
					return ""
				}
				user := u.(*models.User)
				if user.Name == "test" {
					return ""
				}
				return "create:user"
			},
		})
	}
}
