package system

import (
	"context"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type initPointsManagement struct{}

const initOrderPointsManagement = initOrderPointsTables + 1

// auto run
func init() {
	system.RegisterInit(initOrderPointsManagement, &initPointsManagement{})
}

func (i initPointsManagement) InitializerName() string {
	return "points_management_apis"
}

func (i *initPointsManagement) MigrateTable(ctx context.Context) (context.Context, error) {
	// 积分管理API使用现有的sys_apis表，不需要新建表
	return ctx, nil
}

func (i *initPointsManagement) TableCreated(ctx context.Context) bool {
	// 总是返回true，因为我们使用现有的sys_apis表
	return true
}

func (i *initPointsManagement) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 积分管理API列表
	entities := []sysModel.SysApi{
		{ApiGroup: "积分管理", Method: "POST", Path: "/gaia/checkin/checkin", Description: "用户签到"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getStatus", Description: "获取签到状态"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getUserPointsByAccountId/:accountId", Description: "根据账户ID获取积分信息"},
		{ApiGroup: "积分管理", Method: "POST", Path: "/gaia/checkin/exchangePoints", Description: "积分兑换"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getUserPoints", Description: "获取用户积分列表"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getCheckinRecords", Description: "获取签到记录"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getPointsTransaction", Description: "获取积分流水"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getPointsExchange", Description: "获取积分兑换记录"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getPointsConfig", Description: "获取积分配置"},
		{ApiGroup: "积分管理", Method: "POST", Path: "/gaia/checkin/updatePointsConfig", Description: "更新积分配置"},
		{ApiGroup: "积分管理", Method: "POST", Path: "/gaia/checkin/manualAdjustPoints", Description: "手动调整积分"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/gaia/checkin/getPointsStatistics", Description: "获取积分统计"},
		{ApiGroup: "积分管理", Method: "DELETE", Path: "/gaia/checkin/deleteCheckinExtend", Description: "删除积分管理"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, "积分管理API初始化失败!")
	}

	// 获取插入的API记录
	var insertedApis []sysModel.SysApi
	if err := db.Where("api_group = ?", "积分管理").Find(&insertedApis).Error; err != nil {
		return ctx, errors.Wrap(err, "获取积分管理API失败!")
	}

	next := context.WithValue(ctx, i.InitializerName(), insertedApis)
	return next, nil
}

func (i *initPointsManagement) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	
	// 检查是否已经存在积分管理API
	var count int64
	db.Model(&sysModel.SysApi{}).Where("api_group = ?", "积分管理").Count(&count)
	return count > 0
} 