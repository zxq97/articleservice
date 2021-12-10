package server

import (
	"context"
	"log"
	"time"
)

func dbGetArticle(ctx context.Context, articleID int64) (*Article, error) {
	articleMap, err := dbBatchGetArticles(ctx, []int64{articleID})
	if err != nil || articleMap == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func dbBatchGetArticles(ctx context.Context, articleIDs []int64) (map[int64]*Article, error) {
	articles := []*Article{}
	err := dbCli.Model(&Article{}).Where("article_id in (?)", articleIDs).Find(&articles).Error
	if err != nil {
		log.Printf("ctx %v dbBatchGetArticles article_ids %v err %v", ctx, articleIDs, err)
		return nil, err
	}
	articleMap := make(map[int64]*Article, len(articleIDs))
	for _, v := range articles {
		articleMap[v.ArticleID] = v
	}
	return articleMap, nil
}

func dbGetTopic(ctx context.Context, topicID int64) (*Topic, error) {
	topicMap, err := dbBatchGetTopics(ctx, []int64{topicID})
	if err != nil || topicMap == nil {
		return nil, err
	}
	return topicMap[topicID], nil
}

func dbBatchGetTopics(ctx context.Context, topicIDs []int64) (map[int64]*Topic, error) {
	topics := []*Topic{}
	err := dbCli.Model(&Topic{}).Where("topic_id in (?)", topicIDs).Find(&topics).Error
	if err != nil {
		log.Printf("ctx %v dbBatchGetTopics topic_ids %v err %v", ctx, topicIDs, err)
		return nil, err
	}
	topicMap := make(map[int64]*Topic, len(topicIDs))
	for _, v := range topics {
		topicMap[v.TopicID] = v
	}
	return topicMap, nil
}

func dbUpdateVisibleType(ctx context.Context, articleID int64, visibleType int32) error {
	article := new(Article)
	err := dbCli.Model(article).Where("article_id = ?", articleID).Update("visible_type", visibleType).Error
	if err != nil {
		log.Printf("ctx %v dbUpdateVisibleType article_id %v visible_type %v err %v", ctx, articleID, visibleType, err)
	}
	return err
}

func dbAddArticle(ctx context.Context, articleID, topicID, uid int64, content string, visibleType int32) error {
	article := &Article{
		ArticleID:   articleID,
		UID:         uid,
		TopicID:     topicID,
		Content:     content,
		VisibleType: visibleType,
		Ctime:       time.Now(),
		Mtime:       time.Now(),
	}
	err := dbCli.Create(article).Error
	if err != nil {
		log.Printf("ctx %v dbAddArticle article_id %v topic_id %v uid %v content %v visible_type %v err %v", ctx, articleID, topicID, uid, content, visibleType, err)
	}
	return err
}

func dbDeleteArticle(ctx context.Context, articleID int64) error {
	article := new(Article)
	err := dbCli.Model(article).Where("article_id = ?", articleID).Update("is_delete", 1).Error
	if err != nil {
		log.Printf("ctx %v dbDeleteArticle article_id %v err %v", ctx, articleID, err)
	}
	return err
}
