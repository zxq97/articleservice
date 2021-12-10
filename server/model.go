package server

import (
	"articleservice/rpc/article/pb"
	"time"
)

type Article struct {
	ArticleID   int64     `json:"article_id"`
	UID         int64     `json:"uid"`
	Content     string    `json:"content"`
	TopicID     int64     `json:"topic_id"`
	VisibleType int32     `json:"visible_type"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

type Topic struct {
	TopicID   int64  `json:"topic_id" gorm:"column:topic_id"`
	TopicName string `json:"topic_name" gorm:"column:topic_name"`
}

func (a *Article) toArticleInfo() *article_service.ArticleInfo {
	return &article_service.ArticleInfo{
		ArticleId:   a.ArticleID,
		Uid:         a.UID,
		Content:     a.Content,
		TopicId:     a.TopicID,
		VisibleType: a.VisibleType,
		Ctime:       a.Ctime.Unix(),
	}
}

func (t *Topic) toTopicInfo() *article_service.TopicInfo {
	return &article_service.TopicInfo{
		TopicId:   t.TopicID,
		TopicName: t.TopicName,
	}
}

func (a *Article) TableName() string {
	return "article"
}

func (t *Topic) TableName() string {
	return "topic"
}
