package handler

import (
	"awesomeProject8/goods_srv/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"context"
)


func (g GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	panic("implement me")
}

func (g GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	panic("implement me")
}

func (g GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (g GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

