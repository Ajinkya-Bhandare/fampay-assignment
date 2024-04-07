package models

type GetVideoRequest struct {
	PageNumber int32 `json:"page_number"`
	PageSize   int32 `json:"page_size"`
}

type GetVideoResponse struct {
	NumberOfVideos int32       `json:"numVideos"`
	Videos         []VideoData `json:"videos"`
}

type GetSearchRequest struct {
	Query string `json:"query"`
}

type GetSearchResponse struct {
	NumberOfResults int32       `json:"numVideos"`
	Videos          []VideoData `json:"videos"`
}
