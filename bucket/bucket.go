// Package bucket provides functionality for working with google cloud bucket.
package bucket

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"

	"github.com/AndrianaY/store/models"
)

// Storage is struct for storing id of Google cloud storage
type Storage struct {
	bucketID string
}

// MakeStorage returns Google cloud storage struct by provided id
func MakeStorage(bucketID string) *Storage {
	return &Storage{
		bucketID: bucketID,
	}
}

func getFileStoragePath(goodID int, filename string) string {
	return fmt.Sprintf("%v/allpics/%v", goodID, filename)
}

// Put creates filder for good by its id and puts files in storage.
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

// Get gets files available in storage by given id
func (s *Storage) Get(ctx context.Context, goodID int) ([]byte, error) {

	return nil, nil
}

// Delete deletes folder for good by given id
func (s *Storage) Delete(ctx context.Context, goodID int) error {

	return nil
}
