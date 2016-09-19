package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Location", "some value")
		guid := uuid.New().String()
		w.Header().Add("Location", "/matches/"+guid)
		formatter.JSON(w, http.StatusCreated, struct{ Test string }{"This is a test"})
	}
}
