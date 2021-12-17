// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/social/social.proto

package social_service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for SocialServer server

type SocialServerService interface {
	Follow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*EmptyResponse, error)
	Unfollow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*EmptyResponse, error)
	GetFollow(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	GetFollower(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	GetFollowCount(ctx context.Context, in *CountRequest, opts ...client.CallOption) (*CountResponse, error)
	GetFollowAll(ctx context.Context, in *FollowAllRequest, opts ...client.CallOption) (SocialServer_GetFollowAllService, error)
	GetFollowerAll(ctx context.Context, in *FollowAllRequest, opts ...client.CallOption) (SocialServer_GetFollowerAllService, error)
}

type socialServerService struct {
	c    client.Client
	name string
}

func NewSocialServerService(name string, c client.Client) SocialServerService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "social"
	}
	return &socialServerService{
		c:    c,
		name: name,
	}
}

func (c *socialServerService) Follow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*EmptyResponse, error) {
	req := c.c.NewRequest(c.name, "SocialServer.Follow", in)
	out := new(EmptyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServerService) Unfollow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*EmptyResponse, error) {
	req := c.c.NewRequest(c.name, "SocialServer.Unfollow", in)
	out := new(EmptyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServerService) GetFollow(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "SocialServer.GetFollow", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServerService) GetFollower(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "SocialServer.GetFollower", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServerService) GetFollowCount(ctx context.Context, in *CountRequest, opts ...client.CallOption) (*CountResponse, error) {
	req := c.c.NewRequest(c.name, "SocialServer.GetFollowCount", in)
	out := new(CountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialServerService) GetFollowAll(ctx context.Context, in *FollowAllRequest, opts ...client.CallOption) (SocialServer_GetFollowAllService, error) {
	req := c.c.NewRequest(c.name, "SocialServer.GetFollowAll", &FollowAllRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &socialServerServiceGetFollowAll{stream}, nil
}

type SocialServer_GetFollowAllService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*FollowAllResponse, error)
}

type socialServerServiceGetFollowAll struct {
	stream client.Stream
}

func (x *socialServerServiceGetFollowAll) Close() error {
	return x.stream.Close()
}

func (x *socialServerServiceGetFollowAll) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *socialServerServiceGetFollowAll) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *socialServerServiceGetFollowAll) Recv() (*FollowAllResponse, error) {
	m := new(FollowAllResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *socialServerService) GetFollowerAll(ctx context.Context, in *FollowAllRequest, opts ...client.CallOption) (SocialServer_GetFollowerAllService, error) {
	req := c.c.NewRequest(c.name, "SocialServer.GetFollowerAll", &FollowAllRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &socialServerServiceGetFollowerAll{stream}, nil
}

type SocialServer_GetFollowerAllService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*FollowAllResponse, error)
}

type socialServerServiceGetFollowerAll struct {
	stream client.Stream
}

func (x *socialServerServiceGetFollowerAll) Close() error {
	return x.stream.Close()
}

func (x *socialServerServiceGetFollowerAll) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *socialServerServiceGetFollowerAll) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *socialServerServiceGetFollowerAll) Recv() (*FollowAllResponse, error) {
	m := new(FollowAllResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SocialServer server

type SocialServerHandler interface {
	Follow(context.Context, *FollowRequest, *EmptyResponse) error
	Unfollow(context.Context, *FollowRequest, *EmptyResponse) error
	GetFollow(context.Context, *ListRequest, *ListResponse) error
	GetFollower(context.Context, *ListRequest, *ListResponse) error
	GetFollowCount(context.Context, *CountRequest, *CountResponse) error
	GetFollowAll(context.Context, *FollowAllRequest, SocialServer_GetFollowAllStream) error
	GetFollowerAll(context.Context, *FollowAllRequest, SocialServer_GetFollowerAllStream) error
}

func RegisterSocialServerHandler(s server.Server, hdlr SocialServerHandler, opts ...server.HandlerOption) error {
	type socialServer interface {
		Follow(ctx context.Context, in *FollowRequest, out *EmptyResponse) error
		Unfollow(ctx context.Context, in *FollowRequest, out *EmptyResponse) error
		GetFollow(ctx context.Context, in *ListRequest, out *ListResponse) error
		GetFollower(ctx context.Context, in *ListRequest, out *ListResponse) error
		GetFollowCount(ctx context.Context, in *CountRequest, out *CountResponse) error
		GetFollowAll(ctx context.Context, stream server.Stream) error
		GetFollowerAll(ctx context.Context, stream server.Stream) error
	}
	type SocialServer struct {
		socialServer
	}
	h := &socialServerHandler{hdlr}
	return s.Handle(s.NewHandler(&SocialServer{h}, opts...))
}

type socialServerHandler struct {
	SocialServerHandler
}

func (h *socialServerHandler) Follow(ctx context.Context, in *FollowRequest, out *EmptyResponse) error {
	return h.SocialServerHandler.Follow(ctx, in, out)
}

func (h *socialServerHandler) Unfollow(ctx context.Context, in *FollowRequest, out *EmptyResponse) error {
	return h.SocialServerHandler.Unfollow(ctx, in, out)
}

func (h *socialServerHandler) GetFollow(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.SocialServerHandler.GetFollow(ctx, in, out)
}

func (h *socialServerHandler) GetFollower(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.SocialServerHandler.GetFollower(ctx, in, out)
}

func (h *socialServerHandler) GetFollowCount(ctx context.Context, in *CountRequest, out *CountResponse) error {
	return h.SocialServerHandler.GetFollowCount(ctx, in, out)
}

func (h *socialServerHandler) GetFollowAll(ctx context.Context, stream server.Stream) error {
	m := new(FollowAllRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.SocialServerHandler.GetFollowAll(ctx, m, &socialServerGetFollowAllStream{stream})
}

type SocialServer_GetFollowAllStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*FollowAllResponse) error
}

type socialServerGetFollowAllStream struct {
	stream server.Stream
}

func (x *socialServerGetFollowAllStream) Close() error {
	return x.stream.Close()
}

func (x *socialServerGetFollowAllStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *socialServerGetFollowAllStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *socialServerGetFollowAllStream) Send(m *FollowAllResponse) error {
	return x.stream.Send(m)
}

func (h *socialServerHandler) GetFollowerAll(ctx context.Context, stream server.Stream) error {
	m := new(FollowAllRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.SocialServerHandler.GetFollowerAll(ctx, m, &socialServerGetFollowerAllStream{stream})
}

type SocialServer_GetFollowerAllStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*FollowAllResponse) error
}

type socialServerGetFollowerAllStream struct {
	stream server.Stream
}

func (x *socialServerGetFollowerAllStream) Close() error {
	return x.stream.Close()
}

func (x *socialServerGetFollowerAllStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *socialServerGetFollowerAllStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *socialServerGetFollowerAllStream) Send(m *FollowAllResponse) error {
	return x.stream.Send(m)
}