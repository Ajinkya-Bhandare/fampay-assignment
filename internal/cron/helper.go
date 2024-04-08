package cron

import (
	"encoding/json"
	"fampay/internal/db"
	"fampay/internal/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type YouTubeVideo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func fetchData() {

	config, err := getConfig()

	if err != nil {
		fmt.Println("Error fetching config:", err)
		return
	}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	q := req.URL.Query()
	q.Add("key", config.ApiKey)
	q.Add("q", config.Query)
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

func getConfig() (Config models.SearchConfig, err error) {
	file, err := os.Open("search.yaml")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	err = yaml.Unmarshal(contents, &Config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	return
}
