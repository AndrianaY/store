package main

import (
	"net/http"

	"github.com/AndrianaY/store/util/router"

	httptransport "github.com/go-kit/kit/transport/http"
)

func createRoutes(s Service) []router.Route {
	return []router.Route{
		router.Route{
			"Create good",
			http.MethodPost,
			"/goods",
			false,
			httptransport.NewServer(
				makeCreateGoodEndpoint(s),
				decodeCreateGoodRequest,
				encodeCreateGoodResponse,
				getServerOptions(encodeErrors)...,
			),
		},
		router.Route{
			"Get goods",
			http.MethodGet,
			"/goods",
			true,
			httptransport.NewServer(
				makeGetGoodsEndpoint(s),
				decodeGoodsRequest,
				encodeGoodsResponse,
				getServerOptions(encodeErrors)...,
			),
		},
		router.Route{
			"Update good",
			http.MethodPatch,
			"/goods/{ID}",
			true,
			httptransport.NewServer(
				makeEditGoodEndpoint(s),
				decodeEditGoodRequest,
				encodeEditGoodResponse,
				getServerOptions(encodeErrors)...,
			),
		},
		router.Route{
			"Delete good",
			http.MethodDelete,
			"/goods/{ID}",
			true,
			httptransport.NewServer(
				makeDeleteGoodEndpoint(s),
				decodeDeleteGoodRequest,
				encodeEmptyResponse,
				getServerOptions(encodeErrors)...,
			),
		},
		router.Route{
			"Upload good photo",
			http.MethodPut,
			"/goods/{ID}/upload",
			false,
			httptransport.NewServer(
				makeUploadFilesEndpoint(s),
				decodeUploadFileRequest,
				encodeEmptyResponse,
				getServerOptions(encodeErrors)...,
			),
		},
	}
}
