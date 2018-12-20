package main

import (
	"testing"

	"github.com/AndrianaY/store/mocks"
	"github.com/AndrianaY/store/models"
	"github.com/AndrianaY/store/mysqldb"
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
	}).Return(savedGoods, nil).Once()

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
			Name:  "name",
			Price: 100,
		},
	}

	_, err := s.Goods(mocks.Context{}, &actualGoods)

	assert.Nil(t, err)
	assert.NotEmpty(t, &actualGoods)
	assert.Equal(t, expectedGoods, actualGoods)
	goodsRepository.AssertExpectations(t)
}
