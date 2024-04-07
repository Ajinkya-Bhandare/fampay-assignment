package api

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/getSearch", handleGetSearch)
	router.HandleFunc("/getVideos", handleGetVideos)
	return router
}
