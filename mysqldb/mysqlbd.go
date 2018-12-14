// Package mysqldb Provide functionality to work with mysql
package mysqldb

import (
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/AndrianaY/store/config"

	_ "database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"

	"github.com/jinzhu/gorm"
)

type Common interface {
	Create(record interface{}) error

	FindByID(id interface{}, out interface{}) error

	Update(record interface{}) error

	Delete(record interface{}) error
}

type common struct {
	MysqlDB *gorm.DB
}

type DB struct {
	Goods   GoodsRepository
	MysqlDB *gorm.DB
	Common  Common
	log     log.Logger
}

func getConnectionString() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=UTC",
		config.Keys.DbUser,
		config.Keys.DbPassword,
		config.Keys.DbServer,
		config.Keys.DbPort,
		config.Keys.DbSchema,
	)
}

func NewDatabase(log log.Logger) (DB, error) {
	db, err := gorm.Open("mysql", getConnectionString())
	if err != nil {
		return DB{}, err
	}

	return DB{
		Goods: &GoodsTable{
			MysqlDB: db,
		},
		MysqlDB: db,
		Common: &common{
			MysqlDB: db,
		},
		log: log,
	}, nil
}

func (db *common) Create(record interface{}) error {
	if !db.MysqlDB.NewRecord(record) {
		return errors.New("Record already exists")
	}

	return db.MysqlDB.Create(record).Error
}

func (db *common) FindByID(id interface{}, out interface{}) error {
	return db.MysqlDB.First(out, "id = ?", id).Error
}

func (db *common) Update(record interface{}) error {
	return db.MysqlDB.Save(record).Error
}

func (db *common) Delete(record interface{}) error {
	return db.MysqlDB.Delete(record).Error
}
