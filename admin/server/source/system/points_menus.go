package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPointsMenus = initOrderMenu + 1

type initPointsMenus struct{}

// auto run
// 启用积分菜单自动初始化
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
			GVA_MODEL: global.GVA_MODEL{ID: 41}, // 明确指定ID，避免与基础菜单冲突
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
			GVA_MODEL: global.GVA_MODEL{ID: 42}, // 明确指定ID
			ParentId:  41, // 直接使用父菜单ID
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
			GVA_MODEL: global.GVA_MODEL{ID: 43}, // 明确指定ID
			ParentId:  41, // 直接使用父菜单ID
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
			GVA_MODEL: global.GVA_MODEL{ID: 44}, // 明确指定ID
			ParentId:  41, // 直接使用父菜单ID
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
			GVA_MODEL: global.GVA_MODEL{ID: 45}, // 明确指定ID
			ParentId:  41, // 直接使用父菜单ID
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

	// 批量创建所有积分菜单（由于已指定明确的ID和ParentId，可以直接批量创建）
	if err = db.Create(&pointsMenus).Error; err != nil {
		return ctx, errors.Wrap(err, "创建积分管理菜单失败")
	}

	return ctx, nil
}

func (i *initPointsMenus) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	// 检查所有积分菜单是否都存在
	var count int64
	db.Model(&sysModel.SysBaseMenu{}).Where("name IN ?", 
		[]string{"pointsManagement", "pointsUsers", "pointsRecords", "pointsTransactions", "pointsConfig"}).Count(&count)
	return count == 5 // 确保所有5个菜单都存在
}