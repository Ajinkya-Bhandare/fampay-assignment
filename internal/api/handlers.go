package api

import (
	"encoding/json"
	"fampay/internal/db"
	"fampay/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func handleGetVideos(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	pageSize, err := strconv.Atoi(queryParams.Get("PageSize"))
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	pageNumber, err := strconv.Atoi(queryParams.Get("PageNumber"))
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	data := models.GetVideoRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
	}

	response, err := db.GetVideos(data)
	if err != nil {
		fmt.Println("Error getting videos")
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handleGetSearch(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	searchQuery := queryParams.Get("Query")
	if searchQuery == "" {
		http.Error(w, "Empty request", http.StatusBadRequest)
		return
	}

	data := models.GetSearchRequest{
		Query: searchQuery,
	}

	response, err := db.SearchVideos(data)
	if err != nil {
		fmt.Println("Error searching videos")
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
