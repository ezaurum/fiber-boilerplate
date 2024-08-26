package auth

import (
	_ "embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	casbinmiddle "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

//go:embed rbac.conf
var CasbinModel string

func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	m := model.NewModel()
	if err := m.LoadModelFromText(CasbinModel); nil != err {
		return nil, err
	}
	byDB, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("An error '%s' was not expected when creating a new casbin adapter", err)
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
	enforcer, err := NewEnforcer(db)
	if nil != err {
		log.Fatalf("An error '%s' was not expected when creating a new casbin enforcer", err)
	}
	return casbinmiddle.New(casbinmiddle.Config{
		Enforcer: enforcer,
	})
}
