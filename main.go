package main

import (
	"net/http"
	"os"

	"github.com/AndrianaY/store/bucket"
	"github.com/AndrianaY/store/config"
	"github.com/AndrianaY/store/mysqldb"
	"github.com/go-kit/kit/log"
)

func main() {

	config.InitConfig()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db, err := mysqldb.NewDatabase(logger)
	if err != nil {
		panic(err)
	}
	bucket := bucket.MakeStorage(config.Keys.BucketID, logger, db.Goods)

	service := MakeService(bucket, db, logger)
	http.Handle("/", MakeHandler(service))

	port := config.Keys.AppPort
	logger.Log("Web server is running on :" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Log("Unable to start the server, %v", err)
	}
}
