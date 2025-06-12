package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/gofrs/uuid/v5"
	"time"
)

// CheckinRequest 签到请求
type CheckinRequest struct {
	AccountId uuid.UUID `json:"accountId" form:"accountId" binding:"required"`
}

// PointsExchangeRequest 积分兑换请求
type PointsExchangeRequest struct {
	AccountId    uuid.UUID `json:"accountId" form:"accountId" binding:"required"`
	ExchangeType string    `json:"exchangeType" form:"exchangeType" binding:"required"` // quota
	PointsCost   float64   `json:"pointsCost" form:"pointsCost" binding:"required,gt=0"`
	QuotaAmount  *float64  `json:"quotaAmount" form:"quotaAmount"`
	Description  string    `json:"description" form:"description"`
}

// GetCheckinRecordsRequest 获取签到记录请求
type GetCheckinRecordsRequest struct {
	request.PageInfo
	AccountId   *uuid.UUID `json:"accountId" form:"accountId"`
	StartDate   *time.Time `json:"startDate" form:"startDate"`
	EndDate     *time.Time `json:"endDate" form:"endDate"`
	IsBonus     *bool      `json:"isBonus" form:"isBonus"`
}

// GetPointsTransactionRequest 获取积分流水请求
type GetPointsTransactionRequest struct {
	request.PageInfo
	AccountId       *uuid.UUID `json:"accountId" form:"accountId"`
	TransactionType *string    `json:"transactionType" form:"transactionType"`
	StartDate       *time.Time `json:"startDate" form:"startDate"`
	EndDate         *time.Time `json:"endDate" form:"endDate"`
}

// GetPointsExchangeRequest 获取积分兑换记录请求
type GetPointsExchangeRequest struct {
	request.PageInfo
	AccountId    *uuid.UUID `json:"accountId" form:"accountId"`
	ExchangeType *string    `json:"exchangeType" form:"exchangeType"`
	Status       *string    `json:"status" form:"status"`
	StartDate    *time.Time `json:"startDate" form:"startDate"`
	EndDate      *time.Time `json:"endDate" form:"endDate"`
}

// GetUserPointsRequest 获取用户积分请求
type GetUserPointsRequest struct {
	request.PageInfo
	AccountId *uuid.UUID `json:"accountId" form:"accountId"`
	MinPoints *float64   `json:"minPoints" form:"minPoints"`
	MaxPoints *float64   `json:"maxPoints" form:"maxPoints"`
}

// UpdatePointsConfigRequest 更新积分配置请求
type UpdatePointsConfigRequest struct {
	ConfigKey   string  `json:"configKey" form:"configKey" binding:"required"`
	ConfigValue float64 `json:"configValue" form:"configValue" binding:"required"`
	Description string  `json:"description" form:"description"`
}

// ManualAdjustPointsRequest 手动调整积分请求
type ManualAdjustPointsRequest struct {
	AccountId    uuid.UUID `json:"accountId" form:"accountId" binding:"required"`
	PointsChange float64   `json:"pointsChange" form:"pointsChange" binding:"required"`
	Description  string    `json:"description" form:"description" binding:"required"`
} 