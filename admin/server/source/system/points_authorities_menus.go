package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsAuthoritiesMenus = initOrderPointsMenus + 1

type initPointsAuthoritiesMenus struct{}

// auto run
// 禁用积分权限菜单自动初始化 - 改为手动配置
// func init() {
// 	system.RegisterInit(initOrderPointsAuthoritiesMenus, &initPointsAuthoritiesMenus{})
// }

func (i *initPointsAuthoritiesMenus) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initPointsAuthoritiesMenus) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initPointsAuthoritiesMenus) InitializerName() string {
	return "sys_points_authorities_menus"
}

func (i *initPointsAuthoritiesMenus) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 获取积分管理相关菜单
	var pointsMenus []sysModel.SysBaseMenu
	err = db.Where("name IN ?", []string{"pointsManagement", "pointsUsers", "pointsRecords", "pointsTransactions", "pointsConfig"}).Find(&pointsMenus).Error
	if err != nil {
		return ctx, errors.Wrap(err, "获取积分菜单失败")
	}

	if len(pointsMenus) == 0 {
		// 积分菜单不存在，跳过权限分配
		return ctx, nil
	}

	// 获取权限角色
	var authorities []sysModel.SysAuthority
	err = db.Where("authority_id IN ?", []uint{888, 9528}).Find(&authorities).Error
	if err != nil {
		return ctx, errors.Wrap(err, "获取权限角色失败")
	}

	// 为超级管理员(888)和管理员(9528)分配积分管理菜单权限
	for _, authority := range authorities {
		// 只清除该角色现有的积分菜单权限，不影响其他菜单
		var existingPointsMenus []sysModel.SysBaseMenu
		err = db.Model(&authority).Association("SysBaseMenus").Find(&existingPointsMenus, "name IN ?", []string{"pointsManagement", "pointsUsers", "pointsRecords", "pointsTransactions", "pointsConfig"})
		if err == nil && len(existingPointsMenus) > 0 {
			err = db.Model(&authority).Association("SysBaseMenus").Delete(existingPointsMenus)
			if err != nil {
				return ctx, errors.Wrapf(err, "清除角色%d的积分菜单权限失败", authority.AuthorityId)
			}
		}

		// 添加积分菜单权限（使用Append而不是Replace，保留现有其他菜单权限）
		err = db.Model(&authority).Association("SysBaseMenus").Append(pointsMenus)
		if err != nil {
			return ctx, errors.Wrapf(err, "为角色%d分配积分菜单权限失败", authority.AuthorityId)
		}
	}

	return ctx, nil
}

func (i *initPointsAuthoritiesMenus) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	// 检查超级管理员是否拥有积分菜单权限
	var authority sysModel.SysAuthority
	err := db.Where("authority_id = ?", 888).First(&authority).Error
	if err != nil {
		return false
	}

	var menus []sysModel.SysBaseMenu
	err = db.Model(&authority).Association("SysBaseMenus").Find(&menus)
	if err != nil {
		return false
	}

	// 检查是否包含积分管理菜单
	for _, menu := range menus {
		if menu.Name == "pointsManagement" {
			return true
		}
	}

	return false
}