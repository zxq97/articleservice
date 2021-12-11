package server

import (
	"articleservice/client/social"
	"articleservice/rpc/article/pb"
	"articleservice/util/cast"
	"articleservice/util/concurrent"
	"context"
	"time"
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
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		article := &Article{
			ArticleID:   articleID,
			TopicID:     topicID,
			UID:         uid,
			Content:     content,
			VisibleType: visibleType,
			Ctime:       time.Now(),
		}
		cacheSetArticle(ctx, article)
	})
	return nil
}

func deleteArticle(ctx context.Context, articleID int64) error {
	err := dbDeleteArticle(ctx, articleID)
	if err != nil {
		return err
	}
	return cacheDelArticle(ctx, articleID)
}

func pushFollowFeed(ctx context.Context, uid, articleID int64, uids []int64, ok bool) error {
	var err error
	if !ok {
		err = cachePushInBox(ctx, uid, articleID)
		if err != nil {
			return err
		}
	}
	err = cachePushOutBox(ctx, uid, articleID, uids)
	return err
}

func followAddOutBox(ctx context.Context, uid, toUID int64) (bool, error) {
	early, err := outBoxGetEarly(ctx, uid)
	if err != nil {
		return false, err
	}
	var (
		uids       []int64
		articleMap map[int64]int64
	)
	if early == "" {
		uids, err = social.GetFollowAll(ctx, toUID)
		if err != nil {
			return false, err
		}
		early = time.Now().AddDate(0, 0, 30).Format("2006-01-02 15:04:05")
		articleMap, err = dbGetArticlesEarly(ctx, uids, early)
		if err != nil {
			return false, err
		}
	} else {
		articleMap, err = getInBoxEarly(ctx, toUID, early)
		if err != nil {
			return false, err
		}
		if articleMap == nil {
			ctime := time.Unix(cast.ParseInt(early, 0), 0).Format("2006-01-02 15:04:05")
			articleMap, err = dbGetArticleEarly(ctx, toUID, ctime)
			if err != nil {
				return false, err
			}
		}
	}
	if articleMap == nil {
		return false, nil
	}
	err = addOutBox(ctx, articleMap, uid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func unfollowDeleteOutBox(ctx context.Context, uid, toUID int64) error {
	early, err := outBoxGetEarly(ctx, uid)
	if err != nil {
		return err
	}
	if early != "" {
		var articleMap map[int64]int64
		articleMap, err = getInBoxEarly(ctx, toUID, early)
		if err != nil {
			return err
		}
		if articleMap == nil {
			ctime := time.Unix(cast.ParseInt(early, 0), 0).Format("2006-01-02 15:04:05")
			articleMap, err = dbGetArticleEarly(ctx, toUID, ctime)
			if err != nil {
				return err
			}
		}
		if articleMap != nil {
			err = delOutBox(ctx, articleMap, uid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
