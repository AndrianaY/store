package main

import (
	"context"
	"errors"

	"github.com/AndrianaY/store/customErrors"

	"github.com/AndrianaY/store/models"
	"github.com/AndrianaY/store/mysqldb"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Goods(ctx context.Context, outGoods *[]models.Good) (goods *[]models.Good, err error)

	CreateGood(ctx context.Context, name string, price int) (*models.Good, error)

	EditGood(ctx context.Context, ID int, name *string, price *int) (*models.Good, error)

	DeleteGood(ctx context.Context, ID int) error

	UploadFiles(ctx context.Context, price int, files []models.File) error
}

type Storage interface {
	Put(ctx context.Context, id int, files []models.File) (*models.Good, error)
	Upload(ctx context.Context, goodID int, files []models.File) error
	Get(ctx context.Context, goodID int) ([]byte, error)
}

type service struct {
	storage Storage
	log     log.Logger
	DB      mysqldb.DB
}

func MakeService(bucket Storage, db mysqldb.DB, log log.Logger) Service {
	return &service{
		storage: bucket,
		log:     log,
		DB:      db,
	}
}

func (s *service) CreateGood(ctx context.Context, name string, price int) (*models.Good, error) {

	var duplicatedGood models.Good

	if err := s.DB.Goods.FindFirstByName(name, &duplicatedGood); err == nil {
		return nil, customErrors.ErrGoodWithNameExists
	} else if err != gorm.ErrRecordNotFound {
		s.log.Log("Unable to find a good by name, goodName = %v, err = %v", name, err)
		return nil, customErrors.ErrGoodNotFound
	}

	var good models.Good
	good = models.Good{
		Name:  name,
		Price: price,
	}

	if err := s.DB.Common.Create(&good); err != nil {
		s.log.Log("Unable to create a good = %+v, err = %v", good, err)
		return nil, customErrors.ErrUnableCreateGood
	}

	return &models.Good{
		ID:    good.ID,
		Name:  good.Name,
		Price: good.Price,
	}, nil
}

func (s *service) UploadFiles(ctx context.Context, id int, files []models.File) error {
	err := s.storage.Upload(ctx, id, files)
	if err != nil {
		if err == customErrors.ErrGoodNotFound {
			s.log.Log("Service.UploadFiles: Good not found in the DB, error = %v", err)
			return err
		}

		s.log.Log("Service.UploadFIles: Unable upload a file to the storage, error = %v", err)
		return err
	}
	return nil
}

func (s *service) Goods(ctx context.Context, outGoods *[]models.Good) (goods *[]models.Good, err error) {
	var savedGoods []models.Good
	if _, err := s.DB.Goods.GetGoods(&savedGoods); err != nil {
		s.log.Log("Service.GetGoods: Unable to find goods in the DB, error = %v", err)
		return &savedGoods, errors.New("Unable to get goods")
	}

	*outGoods = make([]models.Good, len(savedGoods))
	for i, good := range savedGoods {
		(*outGoods)[i] = models.Good{
			ID:    good.ID,
			Name:  good.Name,
			Price: good.Price,
		}
	}

	return outGoods, nil
}

func (s *service) EditGood(ctx context.Context, ID int, name *string, price *int) (*models.Good, error) {
	//todo: add verification if name already exists
	var good models.Good

	if err := s.DB.Common.FindByID(ID, &good); err != nil {
		return nil, customErrors.ErrGoodNotFound
	}
	good.Name = *name
	good.Price = *price
	if err := s.DB.Common.Update(&good); err != nil {
		s.log.Log("Service.EditGood: Unable to update a good in the DB, id = %v, name = %v, error = %v", ID, name, err)
		return nil, customErrors.ErrUnableUpdateGood
	}

	return &models.Good{
		ID:    good.ID,
		Name:  good.Name,
		Price: good.Price,
	}, nil
}

func (s *service) DeleteGood(ctx context.Context, ID int) error {
	var good models.Good
	if err := s.DB.Common.FindByID(ID, &good); err != nil {
		msg := "Service.DeleteGood: Good not found in the DB, id = %v, error = %v"

		if err == gorm.ErrRecordNotFound {
			s.log.Log(msg, ID, err)
			return customErrors.ErrGoodNotFound
		}

		s.log.Log(msg, ID, err)
		return customErrors.ErrUnableToDeleteGood
	}

	if err := s.DB.Common.Delete(&good); err != nil {
		s.log.Log("Service.DeleteGood: Unable to delete a good from the DB, ID = %v, error = %v", ID, err)
		return customErrors.ErrUnableToDeleteGood
	}

	return nil
}
