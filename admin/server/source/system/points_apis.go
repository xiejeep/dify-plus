package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsApis = initOrderApi + 1

type initPointsApis struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPointsApis, &initPointsApis{})
}

func (i *initPointsApis) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initPointsApis) TableCreated(ctx context.Context) bool {
	return false // always run
}

func (i initPointsApis) InitializerName() string {
	return "sys_points_apis"
}

func (i *initPointsApis) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 定义积分管理API
	pointsApis := []sysModel.SysApi{
		{
			Path:        "/gaia/checkin/checkin",
			Description: "用户签到",
			ApiGroup:    "积分管理",
			Method:      "POST",
		},
		{
			Path:        "/gaia/checkin/getStatus",
			Description: "获取签到状态",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/getUserPointsByAccountId/:accountId",
			Description: "根据账户ID获取积分信息",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/exchangePoints",
			Description: "积分兑换",
			ApiGroup:    "积分管理",
			Method:      "POST",
		},
		{
			Path:        "/gaia/checkin/getUserPoints",
			Description: "获取用户积分列表",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/getCheckinRecords",
			Description: "获取签到记录",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/getPointsTransaction",
			Description: "获取积分流水",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/getPointsExchange",
			Description: "获取积分兑换记录",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/getPointsConfig",
			Description: "获取积分配置",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
		{
			Path:        "/gaia/checkin/updatePointsConfig",
			Description: "更新积分配置",
			ApiGroup:    "积分管理",
			Method:      "POST",
		},
		{
			Path:        "/gaia/checkin/manualAdjustPoints",
			Description: "手动调整积分",
			ApiGroup:    "积分管理",
			Method:      "POST",
		},
		{
			Path:        "/gaia/checkin/getPointsStatistics",
			Description: "获取积分统计",
			ApiGroup:    "积分管理",
			Method:      "GET",
		},
	}

	// 检查是否已存在积分API
	var existingApi sysModel.SysApi
	err = db.Where("path = ? AND method = ?", "/gaia/checkin/checkin", "POST").First(&existingApi).Error
	if err == nil {
		// API已存在，跳过初始化
		return ctx, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx, errors.Wrap(err, "检查积分API是否存在失败")
	}

	// 批量创建积分API
	for _, api := range pointsApis {
		// 检查每个API是否已存在
		var existing sysModel.SysApi
		err = db.Where("path = ? AND method = ?", api.Path, api.Method).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// API不存在，创建它
			if err = db.Create(&api).Error; err != nil {
				return ctx, errors.Wrapf(err, "创建积分API失败: %s %s", api.Method, api.Path)
			}
		} else if err != nil {
			return ctx, errors.Wrapf(err, "检查积分API是否存在失败: %s %s", api.Method, api.Path)
		}
		// 如果API已存在，跳过创建
	}

	return ctx, nil
}

func (i *initPointsApis) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	// 检查积分API是否存在
	var api sysModel.SysApi
	err := db.Where("path = ? AND method = ?", "/gaia/checkin/checkin", "POST").First(&api).Error
	return err == nil
} 