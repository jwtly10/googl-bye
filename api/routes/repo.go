package routes

import (
	"net/http"

	"github.com/jwtly10/googl-bye/api"
	"github.com/jwtly10/googl-bye/api/handlers"
	"github.com/jwtly10/googl-bye/api/middleware"
	"github.com/jwtly10/googl-bye/internal/common"
)

type RepoRoutes struct {
	l common.Logger
	h handlers.RepoHandler
}

func NewRepoRoutes(router api.AppRouter, l common.Logger, h handlers.RepoHandler, mws ...middleware.Middleware) RepoRoutes {
	routes := RepoRoutes{
		l: l,
		h: h,
	}

	BASE_PATH := "/v1/api"

	saveHandler := http.HandlerFunc(routes.h.Save)
	router.Post(
		BASE_PATH+"/save",
		middleware.Chain(saveHandler, mws...),
	)

	return routes
}
