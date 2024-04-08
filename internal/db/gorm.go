package db

import (
	"fampay/internal/models"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"

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

	config, err := getConfig()

	if err != nil {
		fmt.Printf("Failed to get db config: %v\n", err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Username, config.Password, config.Host, config.Port, config.Database)

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

func getConfig() (Config models.SQLConfig, err error) {
	file, err := os.Open("mysql.yaml")
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
