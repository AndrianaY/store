package main

import (
	"errors"
	"testing"

	"github.com/AndrianaY/store/customErrors"
	"github.com/AndrianaY/store/mocks"
	"github.com/AndrianaY/store/models"
	"github.com/AndrianaY/store/mysqldb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetGoods(t *testing.T) {
	goodsRepository := &mocks.GoodsRepository{}

	var savedGoods = []models.Good{
		models.Good{
			ID:    1,
			Name:  "name",
			Price: 100,
		},
	}
	goodsRepository.On("GetGoods", mock.AnythingOfType("*[]models.Good")).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]models.Good)
		(*arg) = savedGoods
	}).Return(uint(1), nil).Once()

	s := MakeService(
		nil,
		mysqldb.DB{
			Goods: goodsRepository,
		},
		nil,
	)
	var actualGoods []models.Good
	var expectedGoods = []models.Good{
		models.Good{
			ID:    1,
			Name:  "name1",
			Price: 100,
		},
	}
	_, err := s.Goods(mocks.Context{}, &actualGoods)
	assert.Nil(t, err)
	assert.Equal(t, expectedGoods, actualGoods)
	// if &expectedGoods != &actualGoods {
	// t.Errorf("service.Goods(context, []models.Good) = %v, want %v", actualGoods, expectedGoods)
	// }
	assert.NotEmpty(t, &actualGoods)
	goodsRepository.AssertExpectations(t)
}

func Test_GetGoods_Error(t *testing.T) {
	goodsRepository := &mocks.GoodsRepository{}
	goodsRepository.On("GetGoods", mock.Anything).Return(uint(0), errors.New("")).Once()
	log := &mocks.Logger{}
	log.On(
		"Log",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*errors.errorString"),
	).Return(nil)

	s := MakeService(
		nil,
		mysqldb.DB{
			Goods: goodsRepository,
		},
		log,
	)
	var actualGoods []models.Good

	_, err := s.Goods(mocks.Context{}, &actualGoods)

	assert.NotNil(t, err)
	assert.Empty(t, &actualGoods)
	goodsRepository.AssertExpectations(t)
	assert.EqualError(t, err, "Unable to get goods")
	log.AssertExpectations(t)
}

func Test_CreateGood(t *testing.T) {
	goodsRepository := &mocks.GoodsRepository{}
	common := &mocks.Common{}

	goodsRepository.On("FindFirstByName", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Good")).
		Return(gorm.ErrRecordNotFound).Once()

	common.On("Create", mock.AnythingOfType("*models.Good")).Return(nil)

	s := MakeService(
		nil,
		mysqldb.DB{
			Goods:  goodsRepository,
			Common: common,
		},
		nil,
	)
	var expectedGood = models.Good{
		ID:    0,
		Name:  "name",
		Price: 100,
	}
	actualGood, err := s.CreateGood(mocks.Context{}, "name", 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedGood, *actualGood)
	assert.NotEmpty(t, &actualGood)
	goodsRepository.AssertExpectations(t)
}

func Test_CreateGood_Existed_Name(t *testing.T) {
	goodsRepository := &mocks.GoodsRepository{}
	common := &mocks.Common{}

	goodsRepository.On("FindFirstByName", mock.AnythingOfType("string"), mock.AnythingOfType("*models.Good")).
		Return(customErrors.ErrGoodWithNameExists).Once()

	log := &mocks.Logger{}
	log.On(
		"Log",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*errors.errorString"),
	).Return(nil)

	s := MakeService(
		nil,
		mysqldb.DB{
			Goods:  goodsRepository,
			Common: common,
		},
		log,
	)

	actualGood, err := s.CreateGood(mocks.Context{}, "name", 100)

	assert.NotNil(t, err)
	assert.Nil(t, actualGood)
	goodsRepository.AssertExpectations(t)
	log.AssertExpectations(t)
}
