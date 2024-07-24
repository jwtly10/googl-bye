package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/service"
	"github.com/jwtly10/googl-bye/internal/utils"
)

type RepoLinkHandler struct {
	log     common.Logger
	service service.RepoLinkService
}

func NewRepoLinkHandler(l common.Logger, r service.RepoLinkService) *RepoLinkHandler {
	return &RepoLinkHandler{
		log:     l,
		service: r,
	}
}

// TODO: Cache this for a few seconds/minute
func (rlh *RepoLinkHandler) GetRepoLinks(w http.ResponseWriter, r *http.Request) {
	repos, err := rlh.service.GetRepoLinks(r)
	if err != nil {
		rlh.log.Error("getting repo links from db failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(repos)
	if err != nil {
		rlh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (rlh *RepoLinkHandler) GetUserRepoLinks(w http.ResponseWriter, r *http.Request) {
	repos, err := rlh.service.GetUserRepoLinks(r)
	if err != nil {
		rlh.log.Error("getting repo links from db failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	jsonResponse, err := json.Marshal(repos)
	if err != nil {
		rlh.log.Error("marshaling response failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
