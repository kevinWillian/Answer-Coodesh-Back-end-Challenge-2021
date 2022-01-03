package models

import "github.com/google/uuid"

type Launche struct {
	ID        uint      `json:"localid" gorm:"primaryKey"`
	SyncedID  uuid.UUID `json:"id"`
	Provider  string
	ArticleID uint
	Article   Article `gorm:"foreignKey:ArticleID"`
}

type Event struct {
	ID        uint   `json:"localid" gorm:"primaryKey"`
	SyncedID  uint32 `json:"id"`
	Provider  string
	ArticleID uint
	Article   Article `gorm:"foreignKey:ArticleID"`
}

type Article struct {
	ID          uint   `json:"localid" gorm:"primaryKey"`
	SyncedID    uint32 `json:"id"`
	Featured    bool
	Title       string
	Url         string
	ImageUrl    string
	NewsSite    string
	Summary     string
	PublishedAt string
	Launches    []Launche
	Events      []Event
}

type IgnoredArticle struct {
	ID          uint   `json:"localid" gorm:"primaryKey"`
	SyncedID    uint32 `json:"id"`
	Featured    bool
	Title       string
	Url         string
	ImageUrl    string
	NewsSite    string
	Summary     string
	PublishedAt string
}

type Config struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Vaule string `json:"vaule"`
}
