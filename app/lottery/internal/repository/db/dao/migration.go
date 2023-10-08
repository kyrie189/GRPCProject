package dao

import (
	"github.com/CocaineCong/grpc-todolist/app/lottery/internal/repository/db/model"
	"github.com/CocaineCong/grpc-todolist/pkg/util/logger"
	"os"
)

func migration() {
	//自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Lottery{},
		)
	if err != nil {
		logger.LogrusObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger.LogrusObj.Infoln("register table success")
}
