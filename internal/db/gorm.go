package db

import (
	"fampay/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateVideo(video *models.VideoData) error {
	gormDB, err := getDBConn()
	if err != nil {
		return err
	}
	result := gormDB.Create(video)
	if result.Error != nil {
		return result.Error
	}
	err = closeDBConn(gormDB)

	return err
}

func GetVideos(req models.GetVideoRequest) (resp models.GetVideoResponse, err error) {

	gormDB, err := getDBConn()
	if err != nil {
		return
	}
	offset := int((req.PageNumber - 1) * req.PageSize)

	var Videos []models.VideoData
	result := gormDB.Order("publish_time desc").Offset(offset).Limit(int(req.PageSize)).Find(&Videos)
	if result.Error != nil {
		return resp, result.Error
	}
	err = closeDBConn(gormDB)

	resp.Videos = Videos
	resp.NumberOfVideos = int32(len(Videos))

	return
}

func SearchVideos(req models.GetSearchRequest) (resp models.GetSearchResponse, err error) {

	gormDB, err := getDBConn()
	if err != nil {
		return
	}

	result := gormDB.Table("video_data").Order("publish_time desc").Limit(10).Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Query)).Find(&resp.Videos)
	if result.Error != nil {
		return resp, result.Error
	}
	err = closeDBConn(gormDB)

	resp.NumberOfResults = int32(len(resp.Videos))
	return
}

func getDBConn() (db *gorm.DB, err error) {
	dsn := "root:root@tcp(host.docker.internal:3306)/YoutubeDatabase?parseTime=true"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	return gormDB.Begin(), nil
}

func closeDBConn(db *gorm.DB) (err error) {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Failed to close database connection: %v\n", err)
	}
	sqlDB.Close()
	return nil
}
