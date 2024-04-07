CREATE DATABASE IF NOT EXISTS YoutubeDatabase;

USE YoutubeDatabase;

CREATE TABLE IF NOT EXISTS video_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title TEXT,
    description TEXT,
    publish_time DATETIME,
    thumbnail_url TEXT,
    thumbnail_height INT,
    thumbnail_width INT,
    channel_title VARCHAR(255),
    UNIQUE(title(255), description(255), channel_title)
);
