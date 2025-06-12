package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type CheckinApi struct{}

// Checkin
// @Tags CheckinExtend
// @Summary 用户签到
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gaiaReq.CheckinRequest true "签到请求"
// @Success 200 {object} response.Response{data=gaiaRes.CheckinResponse,msg=string} "签到成功"
// @Router /gaia/checkin/checkin [post]
func (ca *CheckinApi) Checkin(c *gin.Context) {
	var req gaiaReq.CheckinRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := CheckinService.Checkin(req)
	if err != nil {
		global.GVA_LOG.Error("签到失败!", zap.Error(err))
		response.FailWithMessage("签到失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "签到成功", c)
}

// GetCheckinStatus
// @Tags CheckinExtend
// @Summary 获取签到状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param accountId query string true "用户账户ID"
// @Success 200 {object} response.Response{data=gaiaRes.CheckinStatusResponse,msg=string} "获取成功"
// @Router /gaia/checkin/getStatus [get]
func (ca *CheckinApi) GetCheckinStatus(c *gin.Context) {
	accountIdStr := c.Query("accountId")
	if accountIdStr == "" {
		response.FailWithMessage("用户账户ID不能为空", c)
		return
	}

	accountId, err := uuid.FromString(accountIdStr)
	if err != nil {
		response.FailWithMessage("用户账户ID格式不正确", c)
		return
	}

	result, err := CheckinService.GetCheckinStatus(accountId)
	if err != nil {
		global.GVA_LOG.Error("获取签到状态失败!", zap.Error(err))
		response.FailWithMessage("获取签到状态失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// ExchangePoints
// @Tags CheckinExtend
// @Summary 积分兑换
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gaiaReq.PointsExchangeRequest true "积分兑换请求"
// @Success 200 {object} response.Response{data=gaiaRes.PointsExchangeResponse,msg=string} "兑换成功"
// @Router /gaia/checkin/exchangePoints [post]
func (ca *CheckinApi) ExchangePoints(c *gin.Context) {
	var req gaiaReq.PointsExchangeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := CheckinService.ExchangePoints(req)
	if err != nil {
		global.GVA_LOG.Error("积分兑换失败!", zap.Error(err))
		response.FailWithMessage("积分兑换失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "兑换成功", c)
}

// GetUserPoints
// @Tags CheckinExtend
// @Summary 获取用户积分列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.GetUserPointsRequest true "分页获取用户积分列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/checkin/getUserPoints [get]
func (ca *CheckinApi) GetUserPoints(c *gin.Context) {
	var pageInfo gaiaReq.GetUserPointsRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := CheckinService.GetUserPoints(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetCheckinRecords
// @Tags CheckinExtend
// @Summary 获取签到记录列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.GetCheckinRecordsRequest true "分页获取签到记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/checkin/getCheckinRecords [get]
func (ca *CheckinApi) GetCheckinRecords(c *gin.Context) {
	var pageInfo gaiaReq.GetCheckinRecordsRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := CheckinService.GetCheckinRecords(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetPointsTransaction
// @Tags CheckinExtend
// @Summary 获取积分流水记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.GetPointsTransactionRequest true "分页获取积分流水记录"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/checkin/getPointsTransaction [get]
func (ca *CheckinApi) GetPointsTransaction(c *gin.Context) {
	var pageInfo gaiaReq.GetPointsTransactionRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := CheckinService.GetPointsTransaction(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetPointsExchange
// @Tags CheckinExtend
// @Summary 获取积分兑换记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaiaReq.GetPointsExchangeRequest true "分页获取积分兑换记录"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /gaia/checkin/getPointsExchange [get]
func (ca *CheckinApi) GetPointsExchange(c *gin.Context) {
	var pageInfo gaiaReq.GetPointsExchangeRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := CheckinService.GetPointsExchange(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// UpdatePointsConfig
// @Tags CheckinExtend
// @Summary 更新积分配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gaiaReq.UpdatePointsConfigRequest true "更新积分配置请求"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /gaia/checkin/updatePointsConfig [post]
func (ca *CheckinApi) UpdatePointsConfig(c *gin.Context) {
	var req gaiaReq.UpdatePointsConfigRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = CheckinService.UpdatePointsConfig(req)
	if err != nil {
		global.GVA_LOG.Error("更新积分配置失败!", zap.Error(err))
		response.FailWithMessage("更新积分配置失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetPointsConfig
// @Tags CheckinExtend
// @Summary 获取积分配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]gaiaRes.PointsConfigResponse,msg=string} "获取成功"
// @Router /gaia/checkin/getPointsConfig [get]
func (ca *CheckinApi) GetPointsConfig(c *gin.Context) {
	list, err := CheckinService.GetAllPointsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取积分配置失败!", zap.Error(err))
		response.FailWithMessage("获取积分配置失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// ManualAdjustPoints
// @Tags CheckinExtend
// @Summary 手动调整积分
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body gaiaReq.ManualAdjustPointsRequest true "手动调整积分请求"
// @Success 200 {object} response.Response{msg=string} "调整成功"
// @Router /gaia/checkin/manualAdjustPoints [post]
func (ca *CheckinApi) ManualAdjustPoints(c *gin.Context) {
	var req gaiaReq.ManualAdjustPointsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = CheckinService.ManualAdjustPoints(req)
	if err != nil {
		global.GVA_LOG.Error("手动调整积分失败!", zap.Error(err))
		response.FailWithMessage("手动调整积分失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("调整成功", c)
}

// GetPointsStatistics
// @Tags CheckinExtend
// @Summary 获取积分统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=gaiaRes.PointsStatisticsResponse,msg=string} "获取成功"
// @Router /gaia/checkin/getPointsStatistics [get]
func (ca *CheckinApi) GetPointsStatistics(c *gin.Context) {
	result, err := CheckinService.GetPointsStatistics()
	if err != nil {
		global.GVA_LOG.Error("获取积分统计失败!", zap.Error(err))
		response.FailWithMessage("获取积分统计失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// GetUserPointsByAccountId
// @Tags CheckinExtend
// @Summary 根据用户ID获取积分信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param accountId path string true "用户账户ID"
// @Success 200 {object} response.Response{data=gaiaRes.UserPointsResponse,msg=string} "获取成功"
// @Router /gaia/checkin/getUserPointsByAccountId/{accountId} [get]
func (ca *CheckinApi) GetUserPointsByAccountId(c *gin.Context) {
	accountIdStr := c.Param("accountId")
	if accountIdStr == "" {
		response.FailWithMessage("用户账户ID不能为空", c)
		return
	}

	accountId, err := uuid.FromString(accountIdStr)
	if err != nil {
		response.FailWithMessage("用户账户ID格式不正确", c)
		return
	}

	result, err := CheckinService.GetUserPointsByAccountId(accountId)
	if err != nil {
		global.GVA_LOG.Error("获取用户积分失败!", zap.Error(err))
		response.FailWithMessage("获取用户积分失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
} 