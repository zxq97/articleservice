package server

import (
	"articleservice/conf"
	"articleservice/rpc/article/pb"
	"context"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io"
	"log"
)

type ArticleService struct {
}

var (
	mcCli    *memcache.Client
	redisCli *redis.Client
	dbCli    *gorm.DB
)

func InitService(config *conf.Conf) error {
	var err error
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	mcCli = conf.GetMC(config.MC.Addr)
	redisCli = conf.GetRedis(config.Redis.Addr, config.Redis.DB)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	return err
}

func (as *ArticleService) GetArticle(ctx context.Context, req *article_service.ArticleRequest, res *article_service.ArticleResponse) error {
	articleInfo, err := getArticle(ctx, req.ArticleId)
	if err != nil {
		return err
	}
	res.ArticleInfo = articleInfo
	return nil
}

func (as *ArticleService) GetBatchArticle(ctx context.Context, req *article_service.ArticleBatchRequest, res *article_service.ArticleBatchResponse) error {
	articleInfoMap, err := getBatchArticle(ctx, req.ArticleIds)
	if err != nil {
		return err
	}
	res.ArticleInfos = articleInfoMap
	return nil
}

func (as *ArticleService) GetTopic(ctx context.Context, req *article_service.TopicRequest, res *article_service.TopicResponse) error {
	topicInfo, err := getTopic(ctx, req.TopicId)
	if err != nil {
		return err
	}
	res.TopicInfo = topicInfo
	return nil
}

func (as *ArticleService) GetBatchTopic(ctx context.Context, req *article_service.TopicBatchRequest, res *article_service.TopicBatchResponse) error {
	topicInfoMap, err := getBatchTopic(ctx, req.TopicIds)
	if err != nil {
		return err
	}
	res.TopicInfos = topicInfoMap
	return nil
}

func (as *ArticleService) ChangeVisibleType(ctx context.Context, req *article_service.VisibleTypeRequest, res *article_service.EmptyResponse) error {
	err := changeVisibleType(ctx, req.ArticleId, req.VisibleType)
	return err
}

func (as *ArticleService) PublishArticle(ctx context.Context, req *article_service.PublishArticleRequest, res *article_service.EmptyResponse) error {
	err := publishArticle(ctx, req.ArticleInfo.ArticleId, req.ArticleInfo.TopicId, req.ArticleInfo.Uid, req.ArticleInfo.Content, req.ArticleInfo.VisibleType)
	return err
}

func (as *ArticleService) DeleteArticle(ctx context.Context, req *article_service.ArticleRequest, res *article_service.EmptyResponse) error {
	err := deleteArticle(ctx, req.ArticleId)
	return err
}

func (as *ArticleService) PushFollowFeed(ctx context.Context, stream article_service.ArticleServer_PushFollowFeedStream) error {
	var flag bool
	defer stream.Close()
	for {
		req, err := stream.Recv()
		if err == nil {
			uid := req.Uid
			articleID := req.ArticleId
			uids := req.Uids
			err = pushFollowFeed(ctx, uid, articleID, uids, flag)
			flag = true
			if err != nil {
				return err
			}
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	// fixme 还需 更新提醒服务 缓存增加红点
	return nil
}

func (as *ArticleService) FollowAddOutBox(ctx context.Context, req *article_service.FollowRequest, res *article_service.EmptyResponse) error {
	ok, err := followAddOutBox(ctx, req.Uid, req.ToUid)
	if err != nil {
		return err
	}
	if ok {
		// fixme 还需 更新提醒服务 缓存增加红点
	}
	return nil
}

func (as *ArticleService) UnfollowDeleteOutBox(ctx context.Context, req *article_service.FollowRequest, res *article_service.EmptyResponse) error {
	err := unfollowDeleteOutBox(ctx, req.Uid, req.ToUid)
	return err
}
