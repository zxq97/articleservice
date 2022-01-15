package social

import (
	"articleservice/conf"
	"articleservice/global"
	"articleservice/rpc/social/pb"
	"context"
	"github.com/micro/go-micro"
	"io"
)

var (
	client social_service.SocialServerService
)

func InitClient(config *conf.Conf) {
	service := micro.NewService(micro.Name(
		config.Grpc.Name),
	)
	client = social_service.NewSocialServerService(
		config.Grpc.Name,
		service.Client(),
	)
}

func Follow(ctx context.Context, uid, toUID int64, fType int32) error {
	_, err := client.Follow(ctx, toFollowRequest(uid, toUID, fType))
	if err != nil {
		global.ExcLog.Printf("ctx %v Follow uid %v touid %v ftype %v err %v", ctx, uid, toUID, fType, err)
	}
	return err
}

func Unfollow(ctx context.Context, uid, toUID int64, fType int32) error {
	_, err := client.Unfollow(ctx, toFollowRequest(uid, toUID, fType))
	if err != nil {
		global.ExcLog.Printf("ctx %v Unfollow uid %v touid %v ftype %v err %v", ctx, uid, toUID, fType, err)
	}
	return err
}

func GetFollow(ctx context.Context, uid, lastID, offset int64, fType int32) ([]int64, bool, error) {
	res, err := client.GetFollow(ctx, toListRequest(uid, lastID, offset, fType))
	if err != nil {
		global.ExcLog.Printf("ctx %v GetFollow uid %v lastid %v offset %v ftype %v err %v", ctx, uid, lastID, offset, fType, err)
		return nil, false, err
	}
	return res.Uids, res.HasMore, nil
}

func GetFollower(ctx context.Context, uid, lastID, offset int64, fType int32) ([]int64, bool, error) {
	res, err := client.GetFollower(ctx, toListRequest(uid, lastID, offset, fType))
	if err != nil {
		global.ExcLog.Printf("ctx %v GetFollower uid %v lastid %v offset %v ftype %v err %v", ctx, uid, lastID, offset, fType, err)
		return nil, false, err
	}
	return res.Uids, res.HasMore, nil
}

func GetFollowCount(ctx context.Context, uid int64, fType int32) (int64, int64, error) {
	res, err := client.GetFollowCount(ctx, toCountRequest(uid, fType))
	if err != nil {
		global.ExcLog.Printf("ctx %v GetFollowCount uid %v ftype %v err %v", ctx, uid, fType, err)
		return 0, 0, err
	}
	return res.FollowCount, res.FollowerCount, nil
}

func GetFollowAll(ctx context.Context, uid int64) ([]int64, error) {
	stream, err := client.GetFollowAll(ctx, toFollowAllRequest(uid))
	if err != nil {
		global.ExcLog.Printf("ctx %v GetFollowAll uid %v err %v", ctx, uid, err)
		return nil, err
	}
	defer stream.Close()
	uids := make([]int64, 0)
	for {
		fls, err := stream.Recv()
		if err == nil {
			uids = append(uids, fls.Uids...)
		} else if err == io.EOF {
			break
		} else {
			global.ExcLog.Printf("ctx %v GetFollowAll uid %v err %v", ctx, uid, err)
			return nil, err
		}
	}
	return uids, nil
}

func GetFollowerAll(ctx context.Context, uid int64) ([]int64, error) {
	stream, err := client.GetFollowerAll(ctx, toFollowAllRequest(uid))
	if err != nil {
		global.ExcLog.Printf("ctx %v GetFollowerAll uid %v err %v", ctx, uid, err)
		return nil, err
	}
	defer stream.Close()
	uids := make([]int64, 0)
	for {
		fls, err := stream.Recv()
		if err == nil {
			uids = append(uids, fls.Uids...)
		} else if err == io.EOF {
			break
		} else {
			global.ExcLog.Printf("ctx %v GetFollowerAll uid %v err %v", ctx, uid, err)
			return nil, err
		}
	}
	return uids, nil
}
