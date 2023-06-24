package api

import "github.com/gin-gonic/gin"

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (api *API) Route(r gin.IRouter) {}
