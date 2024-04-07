package models

import "time"

type VideoData struct {
	Title           string    `json:"Title"`
	Description     string    `json:"Description"`
	PublishTime     time.Time `json:"PublishTime"`
	ThumbnailUrl    string    `json:"ThumbnailUrl"`
	ThumbnailHeight int32     `json:"ThumbnailHeight"`
	ThumbnailWidth  int32     `json:"ThumbnailWidth"`
	ChannelTitle    string    `json:"ChannelTitle"`
}

type Thumbnail struct {
	Height int32  `json:"height"`
	Width  int32  `json:"width"`
	URL    string `json:"url"`
}
