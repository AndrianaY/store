package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	customErrors "github.com/AndrianaY/store/customErrors"
	"github.com/AndrianaY/store/models"
	"github.com/AndrianaY/store/util/router"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	notEmptyErrMessage = "Should not be empty"
)

func encoder(w http.ResponseWriter) *json.Encoder {
	encoder := json.NewEncoder(w)
	return encoder
}

func getServerOptions(errorEncoderFn httptransport.ErrorEncoder) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		httptransport.ServerBefore(makeRequestContext),
		httptransport.ServerErrorEncoder(errorEncoderFn),
	}
}

func MakeHandler(s Service) http.Handler {
	r := router.New(createRoutes(s))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.ServeHTTP(w, req)
	})
}

func encodeErrors(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	enc := encoder(w)

	if customErrors.IsNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)

		enc.Encode(errorResponse{
			Message: err.Error(),
		})
	} else if customErrors.IsBadRequest(err) {
		w.WriteHeader(http.StatusBadRequest)

		enc.Encode(errorResponse{
			Message: err.Error(),
		})
	} else {
		w.WriteHeader(http.StatusInternalServerError)

		enc.Encode(errorResponse{
			Message: err.Error(),
		})
	}
}

func makeRequestContext(_ context.Context, r *http.Request) context.Context {
	return r.Context()
}

func encodeEmptyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})

	return nil
}

func decodeGoodsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return r, nil

}

func encodeGoodsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(goodsResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return encoder(w).Encode(res.Goods)
}

func decodeDeleteGoodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteGoodRequest

	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["ID"])
	if err != nil {
		return nil, customErrors.ErrGoodNotFound
	}

	request.ID = ID
	return request, nil
}

func decodeCreateGoodRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request createGoodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.Name == "" {
		return nil, customErrors.ErrWrongBodyRequest
	}

	return request, nil
}

func encodeCreateGoodResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(models.Good)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return encoder(w).Encode(res)
}

func decodeEditGoodRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request editGoodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["ID"])
	if err != nil {
		return nil, customErrors.ErrIncorrectGoodID
	}
	request.ID = ID

	if request.Name != nil {
		if *request.Name == "" {
			return nil, customErrors.ErrIncorrectArguments
		}
	}
	return request, nil
}

func encodeEditGoodResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(models.Good)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return encoder(w).Encode(res)
}

func decodeUploadFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["ID"])
	if err != nil {
		return nil, customErrors.ErrIncorrectGoodID
	}

	files, err := getFiles(r, "files")
	if err != nil {
		errors.New("files")
	} else if len(files) == 0 {
		errors.New("Should upload at least one file")
	}

	return uploadFilesRequest{
		ID:    ID,
		Files: files,
	}, nil
}

func getFiles(r *http.Request, key string) ([]models.File, error) {
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, err
		}
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, customErrors.ErrNoFilesUploaded
	}

	fileHeaders := r.MultipartForm.File[key]
	if len(fileHeaders) == 0 {
		return nil, customErrors.ErrNoFilesUploaded
	}

	files := make([]models.File, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		var err error
		files, err = getFile(fileHeader, files)
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func getFile(fileHeader *multipart.FileHeader, files []models.File) ([]models.File, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	return append(files, models.File{
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header["Content-Type"][0],
		Content:     buf.Bytes(),
	}), nil
}

func parseFileContentFromRequest(r *http.Request) ([]byte, string, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), fileHeader.Filename, nil
}
