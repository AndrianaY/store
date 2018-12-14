package main

import (
	"context"

	"github.com/AndrianaY/store/models"
	"github.com/go-kit/kit/endpoint"
)

type errorResponse struct {
	Message string `json:"message"`
}

type goodsResponse struct {
	Goods []models.Good
}

func makeGetGoodsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var goods []models.Good
		_, err := svc.Goods(ctx, &goods)

		return goodsResponse{
			Goods: goods,
		}, err
	}
}

type createGoodRequest struct {
	Name  string
	Price int
}

func makeCreateGoodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createGoodRequest)
		good, err := svc.CreateGood(ctx, req.Name, req.Price)
		if err != nil {
			return nil, err
		}

		return *good, nil
	}
}

type editGoodRequest struct {
	ID    int
	Name  *string
	Price *int
}

func makeEditGoodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(editGoodRequest)

		good, err := svc.EditGood(ctx, req.ID, req.Name, req.Price)
		if err != nil {
			return nil, err
		}

		return *good, nil
	}
}

type deleteGoodRequest struct {
	ID int
}

func makeDeleteGoodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteGoodRequest)

		err := svc.DeleteGood(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type uploadFilesRequest struct {
	ID    int
	Files []models.File
}

func makeUploadFilesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadFilesRequest)

		err := svc.UploadFiles(ctx, req.ID, req.Files)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
