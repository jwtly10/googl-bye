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

	searchRepoHandler := http.HandlerFunc(routes.h.SearchRepos)
	router.Post(
		BASE_PATH+"/search",
		middleware.Chain(searchRepoHandler, mws...),
	)

	searchUserRepoHandler := http.HandlerFunc(routes.h.SearchReposForUser)
	router.Get(
		BASE_PATH+"/search",
		middleware.Chain(searchUserRepoHandler, mws...),
	)

	searchUserHandler := http.HandlerFunc(routes.h.SearchUsers)
	router.Get(
		BASE_PATH+"/search-user",
		middleware.Chain(searchUserHandler, mws...),
	)

	createIssueHandler := http.HandlerFunc(routes.h.CreateIssue)
	router.Get(
		BASE_PATH+"/create-issue",
		middleware.Chain(createIssueHandler, mws...),
	)

	return routes
}
