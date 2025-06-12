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
func init() {
	system.RegisterInit(initOrderPointsMenus, &initPointsMenus{})
}

func (i initPointsMenus) InitializerName() string {
	return "points_menus"
}

func (i *initPointsMenus) MigrateTable(ctx context.Context) (context.Context, error) {
	// 菜单表已经在initMenu中创建，这里不需要额外迁移
	return ctx, nil
}

func (i *initPointsMenus) TableCreated(ctx context.Context) bool {
	// 菜单表已经存在，这里总是返回true
	return true
}

func (i *initPointsMenus) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 积分管理菜单配置
	pointsMenus := []sysModel.SysBaseMenu{
		{
			GVA_MODEL: global.GVA_MODEL{ID: 41}, 
			MenuLevel: 0, 
			Hidden: false, 
			ParentId: 0, 
			Path: "points", 
			Name: "pointsManagement", 
			Component: "view/gaia/points/index.vue", 
			Sort: 3, 
			Meta: sysModel.Meta{
				Title: "积分管理", 
				Icon: "coin",
			},
		},
		{
			GVA_MODEL: global.GVA_MODEL{ID: 42}, 
			MenuLevel: 0, 
			Hidden: false, 
			ParentId: 41, 
			Path: "users", 
			Name: "pointsUsers", 
			Component: "view/gaia/points/users.vue", 
			Sort: 1, 
			Meta: sysModel.Meta{
				Title: "用户积分管理", 
				Icon: "user",
				KeepAlive: true,
			},
		},
		{
			GVA_MODEL: global.GVA_MODEL{ID: 43}, 
			MenuLevel: 0, 
			Hidden: false, 
			ParentId: 41, 
			Path: "records", 
			Name: "pointsRecords", 
			Component: "view/gaia/points/records.vue", 
			Sort: 2, 
			Meta: sysModel.Meta{
				Title: "签到记录管理", 
				Icon: "document",
			},
		},
		{
			GVA_MODEL: global.GVA_MODEL{ID: 44}, 
			MenuLevel: 0, 
			Hidden: false, 
			ParentId: 41, 
			Path: "transactions", 
			Name: "pointsTransactions", 
			Component: "view/gaia/points/transactions.vue", 
			Sort: 3, 
			Meta: sysModel.Meta{
				Title: "积分流水管理", 
				Icon: "money",
			},
		},
		{
			GVA_MODEL: global.GVA_MODEL{ID: 45}, 
			MenuLevel: 0, 
			Hidden: false, 
			ParentId: 41, 
			Path: "config", 
			Name: "pointsConfig", 
			Component: "view/gaia/points/config.vue", 
			Sort: 4, 
			Meta: sysModel.Meta{
				Title: "积分配置管理", 
				Icon: "setting",
				KeepAlive: true,
			},
		},
	}

	// 创建积分管理菜单
	if err = db.Create(&pointsMenus).Error; err != nil {
		return ctx, errors.Wrap(err, "积分管理菜单数据初始化失败!")
	}

	next = context.WithValue(ctx, i.InitializerName(), pointsMenus)
	return next, nil
}

func (i *initPointsMenus) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	// 检查积分管理主菜单是否存在
	if errors.Is(db.Where("name = ?", "pointsManagement").First(&sysModel.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
} 