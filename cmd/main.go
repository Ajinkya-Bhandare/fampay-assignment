package main

import (
	"fampay/internal/api"
	"fampay/internal/cron"
	"fmt"
	"net/http"
)

func main() {

	cron.StartCron()

	router := api.NewRouter()
	fmt.Println("Starting API server on port 8080...")
	http.ListenAndServe(":8080", router)
}
