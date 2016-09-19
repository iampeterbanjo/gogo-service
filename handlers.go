package main

import (
	"net/http"

	"github.com/cloudnativego/gogo-engine"
	"github.com/google/uuid"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		newMatch := gogo.NewMatch(5)
		repo.addMatch(newMatch)
		guid := uuid.New().String()
		w.Header().Add("Location", "/matches/"+guid)
		formatter.JSON(w, http.StatusCreated, &newMatchResponse{ID: guid})
	}
}
