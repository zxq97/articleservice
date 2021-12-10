package social

import (
	"articleservice/rpc/social/pb"
)

func toFollowItem(uid, targetID int64, fType int32) *social_service.FollowItem {
	return &social_service.FollowItem{
		Uid:        uid,
		TargetId:   targetID,
		FollowType: fType,
	}
}

func toFollowRequest(uid, targetID int64, fType int32) *social_service.FollowRequest {
	return &social_service.FollowRequest{
		FollowItem: toFollowItem(uid, targetID, fType),
	}
}

func toListRequest(uid, lastID, offset int64, fType int32) *social_service.ListRequest {
	return &social_service.ListRequest{
		Uid:        uid,
		LastId:     lastID,
		Offset:     offset,
		FollowType: fType,
	}
}

func toCountRequest(uid int64, fType int32) *social_service.CountRequest {
	return &social_service.CountRequest{
		Uid:        uid,
		FollowType: fType,
	}
}

func toFollowAllRequest(uid int64) *social_service.FollowAllRequest {
	return &social_service.FollowAllRequest{
		Uid: uid,
	}
}
