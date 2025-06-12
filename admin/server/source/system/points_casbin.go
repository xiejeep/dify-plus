package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsCasbin = initOrderPointsManagement + 1

type initPointsCasbin struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPointsCasbin, &initPointsCasbin{})
}

func (i initPointsCasbin) InitializerName() string {
	return "points_management_casbin"
}

func (i *initPointsCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	// 使用现有的casbin_rule表，不需要新建表
	return ctx, nil
}

func (i *initPointsCasbin) TableCreated(ctx context.Context) bool {
	// 总是返回true，因为我们使用现有的casbin_rule表
	return true
}

func (i *initPointsCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 为超级管理员角色(888)添加积分管理API权限
	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/checkin", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getStatus", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getUserPointsByAccountId/:accountId", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/exchangePoints", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getUserPoints", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getCheckinRecords", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getPointsTransaction", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getPointsExchange", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getPointsConfig", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/updatePointsConfig", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/manualAdjustPoints", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/getPointsStatistics", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/gaia/checkin/deleteCheckinExtend", V2: "DELETE"},
		
		// 为管理员角色(9528)添加积分管理API权限
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/checkin", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getStatus", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getUserPointsByAccountId/:accountId", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/exchangePoints", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getUserPoints", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getCheckinRecords", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getPointsTransaction", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getPointsExchange", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getPointsConfig", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/updatePointsConfig", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/manualAdjustPoints", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/getPointsStatistics", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/gaia/checkin/deleteCheckinExtend", V2: "DELETE"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, "积分管理Casbin权限初始化失败!")
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initPointsCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	
	// 检查是否已经存在积分管理相关的权限规则
	var count int64
	db.Model(&adapter.CasbinRule{}).Where("v1 LIKE ?", "/gaia/checkin/%").Count(&count)
	return count > 0
} 