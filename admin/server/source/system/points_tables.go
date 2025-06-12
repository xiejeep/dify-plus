package system

import (
	"context"
	"time"
	gaiaModel "github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type initPointsTables struct{}

const initOrderPointsTables = initOrderUser + 1

// auto run
func init() {
	system.RegisterInit(initOrderPointsTables, &initPointsTables{})
}

func (i initPointsTables) InitializerName() string {
	return "points_management_tables"
}

func (i *initPointsTables) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	
	// 自动创建积分管理相关表
	return ctx, db.AutoMigrate(
		&gaiaModel.UserPointsExtend{},
		&gaiaModel.CheckinRecordExtend{},
		&gaiaModel.PointsTransactionExtend{},
		&gaiaModel.PointsExchangeExtend{},
		&gaiaModel.PointsConfigExtend{},
	)
}

func (i *initPointsTables) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	
	// 检查所有积分管理表是否已创建
	return db.Migrator().HasTable(&gaiaModel.UserPointsExtend{}) &&
		db.Migrator().HasTable(&gaiaModel.CheckinRecordExtend{}) &&
		db.Migrator().HasTable(&gaiaModel.PointsTransactionExtend{}) &&
		db.Migrator().HasTable(&gaiaModel.PointsExchangeExtend{}) &&
		db.Migrator().HasTable(&gaiaModel.PointsConfigExtend{})
}

func (i *initPointsTables) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 初始化默认积分配置
	now := time.Now()
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	id4, _ := uuid.NewV4()
	
	configs := []gaiaModel.PointsConfigExtend{
		{
			Id:          id1,
			ConfigKey:   "daily_checkin_points",
			ConfigValue: 10.0,
			Description: "每日签到基础积分",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
		{
			Id:          id2,
			ConfigKey:   "consecutive_bonus_days",
			ConfigValue: 7.0,
			Description: "连续签到奖励间隔天数",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
		{
			Id:          id3,
			ConfigKey:   "consecutive_bonus_points",
			ConfigValue: 50.0,
			Description: "连续签到奖励积分",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
		{
			Id:          id4,
			ConfigKey:   "points_to_quota_rate",
			ConfigValue: 100.0,
			Description: "积分兑换额度比例(100积分=1美元)",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		},
	}

	// 使用 ON CONFLICT 处理重复键（对于支持的数据库）
	for _, config := range configs {
		var existingConfig gaiaModel.PointsConfigExtend
		err := db.Where("config_key = ?", config.ConfigKey).First(&existingConfig).Error
		if err == gorm.ErrRecordNotFound {
			// 配置不存在，创建新配置
			if err := db.Create(&config).Error; err != nil {
				return ctx, err
			}
		}
	}

	return ctx, nil
}

func (i *initPointsTables) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	
	// 检查是否已有积分配置数据
	var count int64
	db.Model(&gaiaModel.PointsConfigExtend{}).Count(&count)
	return count > 0
} 