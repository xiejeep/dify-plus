package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type CheckinRouter struct{}

// InitCheckinRouter 初始化签到路由
func (cr *CheckinRouter) InitCheckinRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	checkinRouterWithoutRecord := Router.Group("gaia/checkin")
	checkinApi := v1.ApiGroupApp.GaiaApiGroup.CheckinApi
	{
		// 用户签到相关
		checkinRouterWithoutRecord.POST("checkin", checkinApi.Checkin)                              // 用户签到
		checkinRouterWithoutRecord.GET("getStatus", checkinApi.GetCheckinStatus)                   // 获取签到状态
		checkinRouterWithoutRecord.GET("getUserPointsByAccountId/:accountId", checkinApi.GetUserPointsByAccountId) // 根据账户ID获取积分信息
		
		// 积分兑换相关
		checkinRouterWithoutRecord.POST("exchangePoints", checkinApi.ExchangePoints)               // 积分兑换
		
		// 查询相关
		checkinRouterWithoutRecord.GET("getUserPoints", checkinApi.GetUserPoints)                  // 获取用户积分列表
		checkinRouterWithoutRecord.GET("getCheckinRecords", checkinApi.GetCheckinRecords)          // 获取签到记录
		checkinRouterWithoutRecord.GET("getPointsTransaction", checkinApi.GetPointsTransaction)    // 获取积分流水
		checkinRouterWithoutRecord.GET("getPointsExchange", checkinApi.GetPointsExchange)          // 获取积分兑换记录
		
		// 配置管理相关
		checkinRouterWithoutRecord.GET("getPointsConfig", checkinApi.GetPointsConfig)              // 获取积分配置
		checkinRouterWithoutRecord.POST("updatePointsConfig", checkinApi.UpdatePointsConfig)       // 更新积分配置
		
		// 管理员操作
		checkinRouterWithoutRecord.POST("manualAdjustPoints", checkinApi.ManualAdjustPoints)       // 手动调整积分
		checkinRouterWithoutRecord.GET("getPointsStatistics", checkinApi.GetPointsStatistics)      // 获取积分统计
	}
} 