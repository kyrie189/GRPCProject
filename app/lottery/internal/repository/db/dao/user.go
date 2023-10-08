package dao

import (
	"context"
	"gorm.io/gorm"
)

type LotteryDao struct {
	*gorm.DB
}

func NewLotteryDao(ctx context.Context) *LotteryDao {
	return &LotteryDao{NewDBClient(ctx)}
}


