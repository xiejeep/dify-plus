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
func init() {
	system.RegisterInit(initOrderPointsAuthoritiesMenus, &initPointsAuthoritiesMenus{})
}

func (i initPointsAuthoritiesMenus) InitializerName() string {
	return "points_authorities_menus"
}

func (i *initPointsAuthoritiesMenus) MigrateTable(ctx context.Context) (context.Context, error) {
	// 权限菜单关联表已经存在，不需要额外迁移
	return ctx, nil
}

func (i *initPointsAuthoritiesMenus) TableCreated(ctx context.Context) bool {
	// 关联表已经存在，这里总是返回true
	return true
}

func (i *initPointsAuthoritiesMenus) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 获取权限数据
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建积分管理 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}

	// 获取积分管理菜单数据
	pointsMenus, ok := ctx.Value(initPointsMenus{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建积分管理 [菜单-权限] 关联失败, 未找到积分菜单表初始化数据")
	}

	next = ctx

	// 为超级管理员角色 (888) 分配所有积分管理菜单权限
	var authority888 sysModel.SysAuthority
	if err = db.Where("authority_id = ?", 888).First(&authority888).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return next, err
		}
		// 如果角色888不存在，尝试使用authorities数组中的第一个（通常是888）
		if len(authorities) > 0 {
			authority888 = authorities[0]
		} else {
			return next, errors.New("未找到超级管理员角色")
		}
	}
	
	// 为888角色添加所有积分管理菜单
	if err = db.Model(&authority888).Association("SysBaseMenus").Append(pointsMenus); err != nil {
		return next, errors.Wrap(err, "为超级管理员角色分配积分管理菜单失败")
	}

	// 为管理员角色 (9528) 分配所有积分管理菜单权限
	var authority9528 sysModel.SysAuthority
	if err = db.Where("authority_id = ?", 9528).First(&authority9528).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return next, err
		}
		// 如果角色9528不存在，尝试使用authorities数组中的对应角色
		for _, auth := range authorities {
			if auth.AuthorityId == 9528 {
				authority9528 = auth
				break
			}
		}
		if authority9528.AuthorityId == 0 {
			// 如果9528角色不存在，跳过该步骤但不报错
			return next, nil
		}
	}
	
	// 为9528角色添加所有积分管理菜单
	if authority9528.AuthorityId != 0 {
		if err = db.Model(&authority9528).Association("SysBaseMenus").Append(pointsMenus); err != nil {
			return next, errors.Wrap(err, "为管理员角色分配积分管理菜单失败")
		}
	}

	// 为其他角色分配基础积分查看权限（可选）
	// 这里可以根据需要为其他角色分配特定的积分管理菜单

	return next, nil
}

func (i *initPointsAuthoritiesMenus) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	
	// 检查超级管理员是否已有积分管理菜单权限
	var authority sysModel.SysAuthority
	if err := db.Model(&authority).Where("authority_id = ?", 888).
		Preload("SysBaseMenus", "name = ?", "pointsManagement").Find(&authority).Error; err != nil {
		return false
	}
	
	return len(authority.SysBaseMenus) > 0
} 