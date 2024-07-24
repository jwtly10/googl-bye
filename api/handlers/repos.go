package handlers

import (
	"net/http"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/service"
	"github.com/jwtly10/googl-bye/internal/utils"
)

type RepoHandler struct {
	log     common.Logger
	service service.RepoService
}

func NewRepoHandler(l common.Logger, r service.RepoService) *RepoHandler {
	return &RepoHandler{
		log:     l,
		service: r,
	}
}

type response struct {
	status  string
	message string
}

func (rh *RepoHandler) Save(w http.ResponseWriter, r *http.Request) {
	err := rh.service.SaveRepo(r)
	if err != nil {
		rh.log.Error("saving repos to databased failed with error: ", err)
		utils.HandleCustomErrors(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
