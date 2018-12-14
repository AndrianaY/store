package mysqldb

import (
	"github.com/AndrianaY/store/models"

	"github.com/jinzhu/gorm"
)

type GoodsRepository interface {
	FindFirstByName(name string, out *models.Good) error
	GetGoods(out *[]models.Good) (uint, error)
}

const nameOrUWILikeCondition = "`name` LIKE ?"

type GoodsTable struct {
	MysqlDB *gorm.DB
}

func (t *GoodsTable) FindFirstByName(name string, outGood *models.Good) error {
	return t.MysqlDB.First(outGood, "name = ?", name).Error
}

func (t *GoodsTable) GetGoods(out *[]models.Good) (totalCount uint, err error) {
	t.MysqlDB.Find(out)
	err = t.MysqlDB.Model(&models.Good{}).Count(&totalCount).Error
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}
