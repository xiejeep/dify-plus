package system

import (
	"context"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"gorm.io/gorm"
)

const initOrderEmailVerification = initOrderUser + 1

type initEmailVerification struct{}

// auto run
func init() {
	system.RegisterInit(initOrderEmailVerification, &initEmailVerification{})
}

func (i initEmailVerification) InitializerName() string {
	return sysModel.SysEmailVerification{}.TableName()
}

func (i *initEmailVerification) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysEmailVerification{})
}

func (i *initEmailVerification) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysEmailVerification{})
}

func (i *initEmailVerification) InitializeData(ctx context.Context) (context.Context, error) {
	// 邮箱验证码表不需要初始数据，只需要创建表结构
	return ctx, nil
}

func (i *initEmailVerification) DataInserted(ctx context.Context) bool {
	// 邮箱验证码表不需要初始数据，总是返回true
	return true
} 