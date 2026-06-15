package modules

import (
	"github.com/caiolandgraf/grove-base/internal/modules/auth"
	"github.com/caiolandgraf/grove-base/internal/modules/users"
	"github.com/go-fuego/fuego"
)

var registry = []Factory{
	func(b Boot) Module { return users.Wire(b.DB) },
	func(b Boot) Module { return auth.Wire(b.DB, b.Session) },
}

// Mount wires and registers every module.
func Mount(api *fuego.Server, boot Boot) {
	for _, factory := range registry {
		factory(boot).Mount(api, boot.Session)
	}
}
