package server

import (
	"articleservice/rpc/article/pb"
	"articleservice/util/concurrent"
	"context"
)

func getArticle(ctx context.Context, articleID int64) (*article_service.ArticleInfo, error) {
	val, err := cacheGetArticle(ctx, articleID)
	if err != nil || val == nil {
		val, err = dbGetArticle(ctx, articleID)
		if err != nil || val == nil {
			return nil, err
		}
		concurrent.Go(func() {
			cacheSetArticle(ctx, val)
		})
	}
	return val.toArticleInfo(), nil
}

func getBatchArticle(ctx context.Context, articleIDs []int64) (map[int64]*article_service.ArticleInfo, error) {
	cacheArticleMap, missed, err := cacheBatchGetArticle(ctx, articleIDs)
	articleMap := make(map[int64]*article_service.ArticleInfo, len(articleIDs))
	for k, v := range cacheArticleMap {
		articleMap[k] = v.toArticleInfo()
	}
	if err != nil || len(missed) != 0 {
		var dbArticleMap map[int64]*Article
		dbArticleMap, err = dbBatchGetArticles(ctx, missed)
		if err != nil {
			return articleMap, nil
		}
		concurrent.Go(func() {
			cacheBatchSetArticle(ctx, dbArticleMap)
		})
		for k, v := range dbArticleMap {
			articleMap[k] = v.toArticleInfo()
		}
	}
	return articleMap, nil
}

func getTopic(ctx context.Context, topicID int64) (*article_service.TopicInfo, error) {
	val, err := cacheGetTopic(ctx, topicID)
	if err != nil || val == nil {
		val, err = dbGetTopic(ctx, topicID)
		if err != nil || val == nil {
			return nil, err
		}
		concurrent.Go(func() {
			cacheSetTopic(ctx, val)
		})
	}
	return val.toTopicInfo(), nil
}

func getBatchTopic(ctx context.Context, topicIDs []int64) (map[int64]*article_service.TopicInfo, error) {
	cacheTopicMap, missed, err := cacheBatchGetTopic(ctx, topicIDs)
	topicMap := make(map[int64]*article_service.TopicInfo, len(topicIDs))
	for k, v := range cacheTopicMap {
		topicMap[k] = v.toTopicInfo()
	}
	if err != nil || len(missed) != 0 {
		var dbTopicMap map[int64]*Topic
		dbTopicMap, err = dbBatchGetTopics(ctx, missed)
		if err != nil {
			return topicMap, nil
		}
		concurrent.Go(func() {
			cacheBatchSetTopic(ctx, dbTopicMap)
		})
		for k, v := range dbTopicMap {
			topicMap[k] = v.toTopicInfo()
		}
	}
	return topicMap, nil
}

func changeVisibleType(ctx context.Context, articleID int64, visibleType int32) error {
	err := dbUpdateVisibleType(ctx, articleID, visibleType)
	if err != nil {
		return err
	}
	// fixme 还需调用social服务 调整关注人的可见列表

	return cacheDelArticle(ctx, articleID)
}

func publishArticle(ctx context.Context, articleID, topicID, uid int64, content string, visibleType int32) error {
	err := dbAddArticle(ctx, articleID, topicID, uid, content, visibleType)
	return err
}

func deleteArticle(ctx context.Context, articleID int64) error {
	err := dbDeleteArticle(ctx, articleID)
	if err != nil {
		return err
	}
	return cacheDelArticle(ctx, articleID)
}
