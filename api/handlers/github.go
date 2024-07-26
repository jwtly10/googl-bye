package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/service"
	"github.com/jwtly10/googl-bye/internal/utils"
)

type GithubHandler struct {
	log     common.Logger
	service service.GithubService
}

func NewGithubHandler(l common.Logger, g service.GithubService) *GithubHandler {
	return &GithubHandler{
		log:     l,
		service: g,
	}
}

func (gh *GithubHandler) SearchRepos(w http.ResponseWriter, r *http.Request) {
	res, err := gh.service.GithubSearchRepos(r)
	if err != nil {
		gh.log.Error("github search failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		gh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (gh *GithubHandler) SearchReposForUser(w http.ResponseWriter, r *http.Request) {
	res, err := gh.service.GithubSearchReposForUser(r)
	if err != nil {
		gh.log.Error("github search failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		gh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (gh *GithubHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	res, err := gh.service.GithubCreateIssue(r)
	if err != nil {
		gh.log.Error("github issue creation failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		gh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func (gh *GithubHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	res, err := gh.service.GithubSearchUsers(r)
	if err != nil {
		gh.log.Error("github search failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		gh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
