package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/ldhk/tonton-be/pkg/openapi/internal/api"
)

type (
	Config struct{}
)

type Module struct {
	api *api.API
}

func InitModule(c Config) (*Module, error) {
	return &Module{
		api: api.NewAPI(),
	}, nil
}

func (m *Module) Route(r gin.IRouter) {
	m.api.Route(r)
}
