package main

import (
	"articleservice/client/remind"
	"articleservice/client/social"
	"articleservice/conf"
	"articleservice/global"
	"articleservice/rpc/article/pb"
	"articleservice/server"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

var (
	articleConf *conf.Conf
	remindConf  *conf.Conf
	socialConf  *conf.Conf
	err         error
)

func main() {
	articleConf, err = conf.LoadYaml(conf.ArticleConfPath)
	if err != nil {
		panic(err)
	}
	remindConf, err = conf.LoadYaml(conf.RemindConfPath)
	if err != nil {
		panic(err)
	}
	socialConf, err = conf.LoadYaml(conf.SocialConfPath)
	if err != nil {
		panic(err)
	}

	global.InfoLog, err = conf.InitLog(articleConf.InfoLog.Path)
	if err != nil {
		panic(err)
	}
	global.ExcLog, err = conf.InitLog(articleConf.ExcLog.Path)
	if err != nil {
		panic(err)
	}

	remind.InitClient(remindConf)
	social.InitClient(socialConf)

	err = server.InitService(articleConf)
	if err != nil {
		panic(err)
	}

	etcdRegistry := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = articleConf.Etcd.Addr
	})

	service := micro.NewService(
		micro.Name(articleConf.Grpc.Name),
		micro.Address(articleConf.Grpc.Addr),
		micro.Registry(etcdRegistry),
	)
	service.Init()
	err = article_service.RegisterArticleServerHandler(
		service.Server(),
		new(server.ArticleService),
	)
	if err != nil {
		panic(err)
	}
	err = service.Run()
	if err != nil {
		panic(err)
	}
}
