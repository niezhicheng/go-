package handler

import (
	"awesomeProject8/goods_srv/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"context"
)

func (g GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	panic("implement me")
}

func (g GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	panic("implement me")
}

func (g GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (g GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
