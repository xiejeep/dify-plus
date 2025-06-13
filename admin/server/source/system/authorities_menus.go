package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	// 为超级管理员(888)分配所有菜单权限
	// 超级管理员应该拥有系统中所有菜单的访问权限
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(menus); err != nil {
		return next, err
	}

	// 为开发者角色(8881)分配基础菜单权限
	// 包括：仪表盘、关于我们、个人信息
	menu8881 := menus[:2]  // 仪表盘、关于我们
	menu8881 = append(menu8881, menus[9]) // 个人信息 (ID: 10, index: 9)
	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu8881); err != nil {
		return next, err
	}

	// 为管理员角色(9528)分配扩展菜单权限
	// 包括：基础菜单 + 部分管理功能，但不包括超级管理员专用功能
	var menu9528 []sysModel.SysBaseMenu
	
	// 添加基础功能菜单 (仪表盘、关于我们等)
	menu9528 = append(menu9528, menus[:3]...)  // ID: 1-3 (仪表盘、关于我们、超级管理员)
	
	// 添加超级管理员下的子菜单 (角色管理、菜单管理等)
	menu9528 = append(menu9528, menus[3:10]...) // ID: 4-10 (角色管理到个人信息)
	
	// 添加示例文件及其子菜单
	menu9528 = append(menu9528, menus[10:14]...) // ID: 11-14 (示例文件相关)
	
	// 添加系统工具及部分子菜单 (不包括所有自动化代码功能)
	menu9528 = append(menu9528, menus[14:18]...) // ID: 15-18 (系统工具部分)
	menu9528 = append(menu9528, menus[20:22]...) // ID: 21-22 (模板配置、官方网站)
	
	// 添加服务器状态
	menu9528 = append(menu9528, menus[22]) // ID: 23 (服务器状态)
	
	// 添加插件系统及子菜单
	menu9528 = append(menu9528, menus[23:30]...) // ID: 24-30 (插件系统相关)
	
	// 添加导出模板和公告管理
	menu9528 = append(menu9528, menus[28:31]...) // ID: 29-31 (导出模板、公告管理、参数管理)
	
	// 添加二开功能：费用报表、用户列表、额度管理
	menu9528 = append(menu9528, menus[31:34]...) // ID: 32-34 (费用报表、用户列表、额度管理)
	
	// 添加测试管理相关
	menu9528 = append(menu9528, menus[34:37]...) // ID: 35-37 (测试管理相关)
	
	// 添加系统集成功能
	menu9528 = append(menu9528, menus[37:]...) // ID: 38-40 (系统集成相关)
	
	if err = db.Model(&authorities[2]).Association("SysBaseMenus").Replace(menu9528); err != nil {
		return next, err
	}
	
	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 888).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
