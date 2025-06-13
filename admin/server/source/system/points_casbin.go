package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsCasbin = initOrderPointsAuthoritiesMenus + 1

type initPointsCasbin struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPointsCasbin, &initPointsCasbin{})
}

func (i *initPointsCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initPointsCasbin) TableCreated(ctx context.Context) bool {
	return false // always run
}

func (i initPointsCasbin) InitializerName() string {
	return "sys_points_casbin"
}

func (i *initPointsCasbin) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 定义积分管理权限规则
	pointsCasbinRules := [][]string{
		// 超级管理员权限 (888)
		{"888", "/gaia/checkin/checkin", "POST"},
		{"888", "/gaia/checkin/getStatus", "GET"},
		{"888", "/gaia/checkin/getUserPointsByAccountId/:accountId", "GET"},
		{"888", "/gaia/checkin/exchangePoints", "POST"},
		{"888", "/gaia/checkin/getUserPoints", "GET"},
		{"888", "/gaia/checkin/getCheckinRecords", "GET"},
		{"888", "/gaia/checkin/getPointsTransaction", "GET"},
		{"888", "/gaia/checkin/getPointsExchange", "GET"},
		{"888", "/gaia/checkin/getPointsConfig", "GET"},
		{"888", "/gaia/checkin/updatePointsConfig", "POST"},
		{"888", "/gaia/checkin/manualAdjustPoints", "POST"},
		{"888", "/gaia/checkin/getPointsStatistics", "GET"},

		// 管理员权限 (9528)
		{"9528", "/gaia/checkin/checkin", "POST"},
		{"9528", "/gaia/checkin/getStatus", "GET"},
		{"9528", "/gaia/checkin/getUserPointsByAccountId/:accountId", "GET"},
		{"9528", "/gaia/checkin/exchangePoints", "POST"},
		{"9528", "/gaia/checkin/getUserPoints", "GET"},
		{"9528", "/gaia/checkin/getCheckinRecords", "GET"},
		{"9528", "/gaia/checkin/getPointsTransaction", "GET"},
		{"9528", "/gaia/checkin/getPointsExchange", "GET"},
		{"9528", "/gaia/checkin/getPointsConfig", "GET"},
		{"9528", "/gaia/checkin/updatePointsConfig", "POST"},
		{"9528", "/gaia/checkin/manualAdjustPoints", "POST"},
		{"9528", "/gaia/checkin/getPointsStatistics", "GET"},
	}

	// 删除现有的积分相关casbin规则
	for _, rule := range pointsCasbinRules {
		err = db.Where("v0 = ? AND v1 = ? AND v2 = ?", rule[0], rule[1], rule[2]).Delete(&CasbinRule{}).Error
		if err != nil {
			return ctx, errors.Wrapf(err, "删除现有casbin规则失败: %v", rule)
		}
	}

	// 添加新的casbin规则
	for _, rule := range pointsCasbinRules {
		casbinRule := CasbinRule{
			Ptype: "p",
			V0:    rule[0],
			V1:    rule[1],
			V2:    rule[2],
		}
		err = db.Create(&casbinRule).Error
		if err != nil {
			return ctx, errors.Wrapf(err, "创建casbin规则失败: %v", rule)
		}
	}

	return ctx, nil
}

func (i *initPointsCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	// 检查积分casbin规则是否存在
	var count int64
	err := db.Model(&CasbinRule{}).Where("v1 LIKE ?", "/gaia/checkin%").Count(&count).Error
	return err == nil && count > 0
}

// CasbinRule casbin规则模型
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Ptype string `gorm:"size:100" json:"ptype"`
	V0    string `gorm:"size:100" json:"v0"`
	V1    string `gorm:"size:100" json:"v1"`
	V2    string `gorm:"size:100" json:"v2"`
	V3    string `gorm:"size:100" json:"v3"`
	V4    string `gorm:"size:100" json:"v4"`
	V5    string `gorm:"size:100" json:"v5"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
} 