// Package bucket provides functionality for working with google cloud bucket.
package bucket

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/go-kit/kit/log"

	"github.com/AndrianaY/store/models"
	db "github.com/AndrianaY/store/mysqldb"
)

type Storage struct {
	bucketID string
	log      log.Logger
	goods    db.GoodsRepository
}

func MakeStorage(bucketID string, log log.Logger, goods db.GoodsRepository) *Storage {
	return &Storage{
		bucketID: bucketID,
		log:      log,
		goods:    goods,
	}
}

func getFileStoragePath(goodID int, filename string) string {
	// extension := path.Ext(filename)
	// filenameWithoutExtension := filename[0 : len(filename)-len(extension)]

	return fmt.Sprintf("%v/allpics/%v", goodID, filename)
}

func (s *Storage) Put(ctx context.Context, goodID int, files []models.File) (*models.Good, error) {

	return nil, nil
}

// Upload puts files in storage.
func (s *Storage) Upload(ctx context.Context, goodID int, files []models.File) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketID)

	for _, f := range files {
		wc := bucket.Object(getFileStoragePath(goodID, f.Name)).NewWriter(ctx)
		wc.ContentType = f.ContentType

		if _, err := wc.Write(f.Content); err != nil {
			return err
		}

		if err := wc.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, goodID int) ([]byte, error) {

	return nil, nil
}

func (s *Storage) Delete(ctx context.Context, goodID int) error {

	return nil
}
