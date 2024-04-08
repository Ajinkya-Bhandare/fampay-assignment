package cron

import (
	"encoding/json"
	"fampay/internal/db"
	"fampay/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type YouTubeVideo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func fetchData() {

	apiKey := "AIzaSyBTzJCmPaHuI7Uv9PyNL0weoDli1qaYVts"

	// query for yt video
	searchQuery := "cat videos"

	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("q", searchQuery)
	q.Add("part", "snippet")
	q.Add("maxResults", "10")
	q.Add("order", "date")

	req.URL.RawQuery = q.Encode()
	apiURL := req.URL.String()

	// Send HTTP request
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Extract data into json
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = parseData(body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func parseData(b []byte) (data []models.VideoData, err error) {

	var videoData map[string]interface{}
	err = json.Unmarshal(b, &videoData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract video data from the response
	if items, ok := videoData["items"].([]interface{}); ok && len(items) > 0 {

		for i := 0; i < len(items); i++ {
			video := items[i].(map[string]interface{})["snippet"].(map[string]interface{})

			defaultThumbnail := video["thumbnails"].(map[string]interface{})["default"].(map[string]interface{})

			thumbnail := models.Thumbnail{
				URL:    defaultThumbnail["url"].(string),
				Height: int32(defaultThumbnail["height"].(float64)),
				Width:  int32(defaultThumbnail["width"].(float64)),
			}
			parsedTime, err := time.Parse(time.RFC3339, video["publishedAt"].(string))
			if err != nil {
				fmt.Println("Error parsing time")
				return data, err
			}
			ytVideo := models.VideoData{
				Title:           video["title"].(string),
				Description:     video["description"].(string),
				PublishTime:     parsedTime,
				ChannelTitle:    video["channelTitle"].(string),
				ThumbnailUrl:    thumbnail.URL,
				ThumbnailHeight: thumbnail.Height,
				ThumbnailWidth:  thumbnail.Width,
			}
			data = append(data, ytVideo)
		}

	} else {
		fmt.Println("API response does not contain expected data.")
		return
	}

	err = addToDB(data)
	if err != nil {
		fmt.Println("Error adding videos to db")
		return
	}
	return
}

func addToDB(data []models.VideoData) (err error) {

	for _, video := range data {

		err = db.CreateVideo(&video)
		if err != nil {
			if isDuplicateErr(err) {
				// Skip if duplicate entry
				fmt.Printf("Skipping duplicate video: %v\n", video.Title)
				continue
			}
			fmt.Printf("Failed to create video: %v\n", err)
			return err
		}
	}
	err = nil
	return
}

func isDuplicateErr(err error) bool {
	// Check for duplicate entry
	return strings.Contains(err.Error(), "Duplicate entry")
}
