package server

import (
	"articleservice/global"
	"articleservice/util/cast"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	MCKeyArticleInfoTTL = 5 * 60
	MCKeyTopicInfoTTL   = 10 * 60
	MCKeyArticleInfo    = "article_service_info_%v"       // article_id article
	MCKeyTopicInfo      = "article_service_topic_info_%v" // topic_id topic

	RedisKeyZInBox = "article_service_in_box_%v" // uid article_id ctime
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

func cachePushInBox(ctx context.Context, uid, articleID int64) error {
	key := fmt.Sprintf(RedisKeyZInBox, uid)
	luaScript := redis.NewScript(`
		local val = redis.call("exist", KEYS[1])
		if ( val == 1 ) then
			return redis.call("zadd", KEYS[1], KEYS[2], KEYS[3])
		else
			return -1
		end
	`)
	_, err := luaScript.Run(ctx, redisCli, []string{key, cast.FormatInt(time.Now().Unix()), cast.FormatInt(articleID)}).Result()
	if err != nil {
		global.ExcLog.Printf("ctx %v cachePushInBox uid %v articleid %v err %v", ctx, uid, articleID, err)
	}
	return err
}

func cacheSetInBox(ctx context.Context, uid int64, articleMap map[int64]int64) error {
	key := fmt.Sprintf(RedisKeyZInBox, uid)
	zs := make([]*redis.Z, 0, len(articleMap))
	for k, v := range articleMap {
		zs = append(zs, &redis.Z{
			Member: k,
			Score:  float64(v),
		})
	}
	err := redisCli.ZAdd(ctx, key, zs...).Err()
	if err != nil {
		global.ExcLog.Printf("ctx %v cacheSetInBox uid %v articlemap %v err %v", ctx, uid, articleMap, err)
	}
	return err
}

func cacheGetInBox(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, bool, error) {
	var (
		val     []string
		nextCur int64
		hasMore bool
		err     error
	)
	key := fmt.Sprintf(RedisKeyZInBox, uid)
	val, err = redisCli.ZRevRange(ctx, key, cursor, cursor+offset).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, 0, false, nil
		}
		global.ExcLog.Printf("ctx %v cacheGetInBox uid %v cursor %v err %v", ctx, uid, cursor, err)
		return nil, 0, false, err
	}
	if len(val) > int(offset) {
		hasMore = true
	}
	ids := make([]int64, 0, offset)
	for _, v := range val {
		nextCur = cast.ParseInt(v, 0)
		if len(ids) == int(offset) {
			break
		}
		ids = append(ids, nextCur)
	}
	return ids, nextCur, hasMore, nil
}

//func cachePushInBox(ctx context.Context, uid, articleID int64) error {
//	key := fmt.Sprintf(RedisKeyZInBox, uid)
//	pipe := redisCli.Pipeline()
//	pipe.ZAdd(key, redis.Z{Member: articleID, Score: float64(time.Now().Unix())})
//	pipe.Expire(key, RedisKeyBoxTTL)
//	_, err := pipe.Exec()
//	if err != nil {
//		global.ExcLog.Printf("ctx %v cachePushInBox uid %v articleid %v err %v", ctx, uid, articleID, err)
//	}
//	return err
//}
//
//func cachePushOutBox(ctx context.Context, uid, articleID int64, uids []int64) error {
//	now := float64(time.Now().Unix())
//	pipe := redisCli.Pipeline()
//	var (
//		key  string
//		aKey string
//	)
//	for _, v := range uids {
//		key = fmt.Sprintf(RedisKeyZOutBox, v)
//		aKey = fmt.Sprintf(RedisKeySAlready, v)
//		pipe.ZAdd(key, redis.Z{Member: articleID, Score: now})
//		pipe.SAdd(aKey, uid)
//		pipe.Expire(key, RedisKeyBoxTTL)
//		pipe.Expire(aKey, RedisKeyBoxTTL)
//	}
//	_, err := pipe.Exec()
//	if err != nil {
//		global.ExcLog.Printf("ctx %v cachePushInBox articleid %v uids %v err %v", ctx, articleID, uids, err)
//	}
//	return err
//}
//
//func outBoxGetEarly(ctx context.Context, uid int64) (string, error) {
//	key := fmt.Sprintf(RedisKeyZOutBox, uid)
//	z, err := redisCli.ZRangeWithScores(key, 0, 0).Result()
//	if err != nil && err != redis.Nil {
//		global.ExcLog.Printf("ctx %v outBoxGetEarly uid %v err %v", ctx, uid, err)
//		return "", err
//	} else if err == redis.Nil || len(z) == 0 {
//		return "", nil
//	}
//	return cast.FormatInt(int64(z[0].Score)), nil
//}
//
//func getInBoxEarly(ctx context.Context, uid int64, ctime string) (map[int64]int64, error) {
//	key := fmt.Sprintf(RedisKeyZInBox, uid)
//	z, err := redisCli.ZRevRangeByScoreWithScores(key, redis.ZRangeBy{Max: ctime}).Result()
//	if err != nil && err != redis.Nil {
//		global.ExcLog.Printf("ctx %v getInBoxEarly uid %v ctime %v err %v", ctx, uid, ctx, err)
//		return nil, err
//	} else if err == redis.Nil || len(z) == 0 {
//		return nil, nil
//	}
//	articleMap := make(map[int64]int64, len(z))
//	for _, v := range z {
//		articleMap[cast.ParseInt(v.Member.(string), 0)] = int64(v.Score)
//	}
//	return articleMap, nil
//}
//
//func addOutBox(ctx context.Context, articleMap map[int64]int64, uid int64) error {
//	key := fmt.Sprintf(RedisKeyZInBox, uid)
//	z := make([]redis.Z, 0, len(articleMap))
//	for k, v := range articleMap {
//		z = append(z, redis.Z{
//			Member: cast.FormatInt(k),
//			Score:  float64(v),
//		})
//	}
//	err := redisCli.ZAdd(key, z...).Err()
//	if err != nil {
//		global.ExcLog.Printf("ctx %v addOutBox articlemap %v uid %v err %v", ctx, articleMap, uid, err)
//	}
//	return err
//}
//
//func delOutBox(ctx context.Context, articleMap map[int64]int64, uid int64) error {
//	key := fmt.Sprintf(RedisKeyZInBox, uid)
//	articleIDs := make([]int64, 0, len(articleMap))
//	for k := range articleMap {
//		articleIDs = append(articleIDs, k)
//	}
//	err := redisCli.ZRem(key, articleIDs).Err()
//	if err != nil {
//		global.ExcLog.Printf("ctx %v delOutBox articleids %v uid %v err %v", ctx, articleIDs, uid, err)
//	}
//	return err
//}
