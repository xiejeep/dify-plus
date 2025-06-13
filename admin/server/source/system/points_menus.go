package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsMenus = initOrderMenu + 1

type initPointsMenus struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPointsMenus, &initPointsMenus{})
}

func (i *initPointsMenus) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initPointsMenus) TableCreated(ctx context.Context) bool {
	return false // always run
}

func (i initPointsMenus) InitializerName() string {
	return "sys_points_menus"
}

func (i *initPointsMenus) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 定义积分管理菜单
	pointsMenus := []sysModel.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "points",
			Name:      "pointsManagement",
			Hidden:    false,
			Component: "view/gaia/points/index.vue",
			Sort:      3,
			Meta: sysModel.Meta{
				Title: "积分管理",
				Icon:  "coin",
			},
		},
		{
			ParentId:  0, // 这里需要在创建后更新为父菜单的实际ID
			Path:      "users",
			Name:      "pointsUsers",
			Hidden:    false,
			Component: "view/gaia/points/users.vue",
			Sort:      1,
			Meta: sysModel.Meta{
				Title:     "用户积分管理",
				Icon:      "user",
				KeepAlive: true,
			},
		},
		{
			ParentId:  0, // 这里需要在创建后更新为父菜单的实际ID
			Path:      "records",
			Name:      "pointsRecords",
			Hidden:    false,
			Component: "view/gaia/points/records.vue",
			Sort:      2,
			Meta: sysModel.Meta{
				Title: "签到记录管理",
				Icon:  "document",
			},
		},
		{
			ParentId:  0, // 这里需要在创建后更新为父菜单的实际ID
			Path:      "transactions",
			Name:      "pointsTransactions",
			Hidden:    false,
			Component: "view/gaia/points/transactions.vue",
			Sort:      3,
			Meta: sysModel.Meta{
				Title: "积分流水管理",
				Icon:  "money",
			},
		},
		{
			ParentId:  0, // 这里需要在创建后更新为父菜单的实际ID
			Path:      "config",
			Name:      "pointsConfig",
			Hidden:    false,
			Component: "view/gaia/points/config.vue",
			Sort:      4,
			Meta: sysModel.Meta{
				Title:     "积分配置管理",
				Icon:      "setting",
				KeepAlive: true,
			},
		},
	}

	// 检查是否已存在积分菜单
	var existingMenu sysModel.SysBaseMenu
	err = db.Where("name = ?", "pointsManagement").First(&existingMenu).Error
	if err == nil {
		// 菜单已存在，跳过初始化
		return ctx, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx, errors.Wrap(err, "检查积分菜单是否存在失败")
	}

	// 先删除可能存在的旧菜单记录（处理主键冲突）
	db.Where("name IN ?", []string{"pointsManagement", "pointsUsers", "pointsRecords", "pointsTransactions", "pointsConfig"}).Delete(&sysModel.SysBaseMenu{})

	// 先创建父菜单
	parentMenu := pointsMenus[0]
	if err = db.Create(&parentMenu).Error; err != nil {
		return ctx, errors.Wrap(err, "创建积分管理父菜单失败")
	}

	// 更新子菜单的ParentId并创建
	for i := 1; i < len(pointsMenus); i++ {
		pointsMenus[i].ParentId = parentMenu.ID
		if err = db.Create(&pointsMenus[i]).Error; err != nil {
			return ctx, errors.Wrap(err, "创建积分管理子菜单失败")
		}
	}

	return ctx, nil
}

func (i *initPointsMenus) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	// 检查积分主菜单是否存在
	var menu sysModel.SysBaseMenu
	err := db.Where("name = ?", "pointsManagement").First(&menu).Error
	return err == nil
}