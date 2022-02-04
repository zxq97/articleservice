package server

import (
	"articleservice/conf"
	"articleservice/rpc/article/pb"
	"context"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ArticleService struct {
}

var (
	mcCli    *memcache.Client
	redisCli redis.Cmdable
	dbCli    *gorm.DB
	slaveCli *gorm.DB
)

func InitService(config *conf.Conf) error {
	var err error
	mcCli = conf.GetMC(config.MC.Addr)
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
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
	return changeVisibleType(ctx, req.ArticleId, req.VisibleType)
}

func (as *ArticleService) PublishArticle(ctx context.Context, req *article_service.PublishArticleRequest, res *article_service.EmptyResponse) error {
	return publishArticle(ctx, req.ArticleInfo.ArticleId, req.ArticleInfo.TopicId, req.ArticleInfo.Uid, req.ArticleInfo.Content, req.ArticleInfo.VisibleType)
}

func (as *ArticleService) DeleteArticle(ctx context.Context, req *article_service.ArticleRequest, res *article_service.EmptyResponse) error {
	return deleteArticle(ctx, req.ArticleId)
}
