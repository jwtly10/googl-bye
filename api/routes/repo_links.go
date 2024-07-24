package routes

import (
	"net/http"

	"github.com/jwtly10/googl-bye/api"
	"github.com/jwtly10/googl-bye/api/handlers"
	"github.com/jwtly10/googl-bye/api/middleware"
	"github.com/jwtly10/googl-bye/internal/common"
)

type RepoLinkRoutes struct {
	l common.Logger
	h handlers.RepoLinkHandler
}

func NewRepoLinkRoutes(router api.AppRouter, l common.Logger, h handlers.RepoLinkHandler, mws ...middleware.Middleware) RepoLinkRoutes {
	routes := RepoLinkRoutes{
		l: l,
		h: h,
	}

	BASE_PATH := "/v1/api"

	repoLinkHandler := http.HandlerFunc(routes.h.GetRepoLinks)
	router.Get(
		BASE_PATH+"/repoLinks",
		middleware.Chain(repoLinkHandler, mws...),
	)

	return routes
}
