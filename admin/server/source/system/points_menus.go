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
			ID:        41,
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
			ID:        42,
			ParentId:  41,
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
			ID:        43,
			ParentId:  41,
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
			ID:        44,
			ParentId:  41,
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
			ID:        45,
			ParentId:  41,
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

	// 批量创建积分菜单
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

	// 检查积分主菜单是否存在
	var menu sysModel.SysBaseMenu
	err := db.Where("name = ?", "pointsManagement").First(&menu).Error
	return err == nil
} 