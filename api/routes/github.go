package routes

import (
	"net/http"

	"github.com/jwtly10/googl-bye/api"
	"github.com/jwtly10/googl-bye/api/handlers"
	"github.com/jwtly10/googl-bye/api/middleware"
	"github.com/jwtly10/googl-bye/internal/common"
)

type GithubRoutes struct {
	l common.Logger
	h handlers.GithubHandler
}

func NewGithubRoutes(router api.AppRouter, l common.Logger, h handlers.GithubHandler, mws ...middleware.Middleware) GithubRoutes {
	routes := GithubRoutes{
		l: l,
		h: h,
	}

	BASE_PATH := "/v1/api"

	searchHandler := http.HandlerFunc(routes.h.Search)
	router.Post(
		BASE_PATH+"/search",
		middleware.Chain(searchHandler, mws...),
	)

	return routes
}
