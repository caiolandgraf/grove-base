package modules

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

// Module is a self-contained HTTP domain: it wires its own dependencies and
// registers its routes. Split modules into separate packages as the app grows.
type Module interface {
	Mount(api *fuego.Server, session *scs.SessionManager)
}

// Boot carries infra dependencies available at application startup.
type Boot struct {
	DB      *gorm.DB
	Session *scs.SessionManager
}

// Factory builds a Module from runtime infra. Register new domains in register.go.
type Factory func(boot Boot) Module
