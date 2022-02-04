package server

import (
	"articleservice/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MCKeyArticleInfoTTL = 5 * 60
	MCKeyTopicInfoTTL   = 10 * 60
	MCKeyArticleInfo    = "article_service_info_%v"       // article_id article
	MCKeyTopicInfo      = "article_service_topic_info_%v" // topic_id topic
)

func cacheGetArticle(ctx context.Context, articleID int64) (*Article, error) {
	articleMap, _, err := cacheBatchGetArticle(ctx, []int64{articleID})
	if err != nil || articleMap[articleID] == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func cacheBatchGetArticle(ctx context.Context, articleIDs []int64) (map[int64]*Article, []int64, error) {
	keys := make([]string, 0, len(articleIDs))
	for _, v := range articleIDs {
		keys = append(keys, fmt.Sprintf(MCKeyArticleInfo, v))
	}
	res, err := mcCli.GetMulti(keys)
	if err != nil {
		global.ExcLog.Printf("ctx %v cache get article_ids %v err %v", ctx, articleIDs, err)
		return nil, articleIDs, err
	}
	articleMap := make(map[int64]*Article, len(articleIDs))
	for _, v := range res {
		article := Article{}
		err = json.Unmarshal(v.Value, &article)
		if err != nil {
			global.ExcLog.Printf("ctx %v cache get article %v josn err %v", ctx, v.Value, err)
			continue
		}
		articleMap[article.ArticleID] = &article
	}
	missed := make([]int64, 0, len(articleIDs))
	for _, v := range articleIDs {
		if _, ok := articleMap[v]; !ok {
			missed = append(missed, v)
		}
	}
	return articleMap, missed, nil
}

func cacheSetArticle(ctx context.Context, article *Article) {
	val, err := json.Marshal(article)
	if err != nil {
		global.ExcLog.Printf("ctx %v cache set article_id %v json err %v", ctx, article.ArticleID, err)
		return
	}
	err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, article.ArticleID), Value: val, Expiration: MCKeyArticleInfoTTL})
	if err != nil {
		global.ExcLog.Printf("ctx %v cache set article_id %v mc err %v", ctx, article.ArticleID, err)
	}
}

func cacheBatchSetArticle(ctx context.Context, articleMap map[int64]*Article) {
	for k, v := range articleMap {
		val, err := json.Marshal(v)
		if err != nil {
			global.ExcLog.Printf("ctx %v cache set article_id %v json err %v", ctx, k, err)
			continue
		}
		err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, k), Value: val, Expiration: MCKeyArticleInfoTTL})
		if err != nil {
			global.ExcLog.Printf("ctx %v cache set article_id %v mc err %v", ctx, k, err)
		}
	}
}

func cacheDelArticle(ctx context.Context, articleID int64) error {
	err := mcCli.Delete(fmt.Sprintf(MCKeyArticleInfo, articleID))
	if err != nil {
		global.ExcLog.Printf("ctx %v cache del article_id %v err %v", ctx, articleID, err)
	}
	return err
}

func cacheGetTopic(ctx context.Context, topicID int64) (*Topic, error) {
	topicMap, _, err := cacheBatchGetTopic(ctx, []int64{topicID})
	if err != nil || topicMap == nil {
		return nil, err
	}
	return topicMap[topicID], nil
}

func cacheBatchGetTopic(ctx context.Context, topicIDs []int64) (map[int64]*Topic, []int64, error) {
	keys := make([]string, 0, len(topicIDs))
	for _, v := range topicIDs {
		keys = append(keys, fmt.Sprintf(MCKeyTopicInfo, v))
	}
	res, err := mcCli.GetMulti(keys)
	if err != nil {
		global.ExcLog.Printf("ctx %v cacheBatchGetTopic topic_ids %v err %v", ctx, topicIDs, err)
		return nil, topicIDs, err
	}
	topicMap := make(map[int64]*Topic, len(topicIDs))
	for _, v := range res {
		topic := Topic{}
		err = json.Unmarshal(v.Value, &topic)
		if err != nil {
			global.ExcLog.Printf("ctx %v cacheBatchGetTopic topic %v josn err %v", ctx, v.Value, err)
			continue
		}
		topicMap[topic.TopicID] = &topic
	}
	missed := make([]int64, 0, len(topicIDs))
	for _, v := range topicIDs {
		if _, ok := topicMap[v]; !ok {
			missed = append(missed, v)
		}
	}
	return topicMap, missed, nil
}

func cacheSetTopic(ctx context.Context, topic *Topic) {
	val, err := json.Marshal(topic)
	if err != nil {
		global.ExcLog.Printf("ctx %v cache set topic_id %v json err %v", ctx, topic.TopicID, err)
		return
	}
	err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyTopicInfo, topic.TopicID), Value: val, Expiration: MCKeyTopicInfoTTL})
	if err != nil {
		global.ExcLog.Printf("ctx %v cache set topic_id %v mc err %v", ctx, topic.TopicID, err)
	}
}

func cacheBatchSetTopic(ctx context.Context, topicMap map[int64]*Topic) {
	for k, v := range topicMap {
		val, err := json.Marshal(v)
		if err != nil {
			global.ExcLog.Printf("ctx %v cache set topic_id %v json err %v", ctx, k, err)
			continue
		}
		err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyTopicInfo, k), Value: val, Expiration: MCKeyArticleInfoTTL})
		if err != nil {
			global.ExcLog.Printf("ctx %v cache set topic_id %v mc err %v", ctx, k, err)
		}
	}
}
