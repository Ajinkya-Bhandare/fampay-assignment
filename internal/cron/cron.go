package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func fetchDataJob() {
	fmt.Println("Fetching data from youtube")
	fetchData()
}

func StartCron() {
	c := cron.New()

	// Schedule cron every minute
	_, err := c.AddFunc("* * * * *", fetchDataJob)

	// Start the cron
	c.Start()

	if err != nil {
		fmt.Println("Error adding cron job:", err)
		c.Stop()
		return
	}
}
