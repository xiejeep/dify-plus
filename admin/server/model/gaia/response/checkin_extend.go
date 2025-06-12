package response

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// CheckinResponse 签到响应
type CheckinResponse struct {
	Success         bool    `json:"success"`
	Message         string  `json:"message"`
	PointsEarned    float64 `json:"pointsEarned"`
	ConsecutiveDays int     `json:"consecutiveDays"`
	IsBonus         bool    `json:"isBonus"`
	TotalPoints     float64 `json:"totalPoints"`
	AvailablePoints float64 `json:"availablePoints"`
}

// UserPointsResponse 用户积分信息响应
type UserPointsResponse struct {
	Id              uuid.UUID `json:"id"`
	AccountId       uuid.UUID `json:"accountId"`
	AccountName     string    `json:"accountName,omitempty"`
	TotalPoints     float64   `json:"totalPoints"`
	AvailablePoints float64   `json:"availablePoints"`
	UsedPoints      float64   `json:"usedPoints"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CheckinRecordResponse 签到记录响应
type CheckinRecordResponse struct {
	Id              uuid.UUID `json:"id"`
	AccountId       uuid.UUID `json:"accountId"`
	AccountName     string    `json:"accountName,omitempty"`
	CheckinDate     time.Time `json:"checkinDate"`
	PointsEarned    float64   `json:"pointsEarned"`
	ConsecutiveDays int       `json:"consecutiveDays"`
	IsBonus         bool      `json:"isBonus"`
	CreatedAt       time.Time `json:"createdAt"`
}

// PointsTransactionResponse 积分流水响应
type PointsTransactionResponse struct {
	Id              uuid.UUID  `json:"id"`
	AccountId       uuid.UUID  `json:"accountId"`
	AccountName     string     `json:"accountName,omitempty"`
	TransactionType string     `json:"transactionType"`
	PointsChange    float64    `json:"pointsChange"`
	PointsBefore    float64    `json:"pointsBefore"`
	PointsAfter     float64    `json:"pointsAfter"`
	Description     string     `json:"description"`
	RelatedId       *uuid.UUID `json:"relatedId"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// PointsExchangeResponse 积分兑换记录响应
type PointsExchangeResponse struct {
	Id           uuid.UUID `json:"id"`
	AccountId    uuid.UUID `json:"accountId"`
	AccountName  string    `json:"accountName,omitempty"`
	ExchangeType string    `json:"exchangeType"`
	PointsCost   float64   `json:"pointsCost"`
	QuotaAmount  *float64  `json:"quotaAmount"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PointsConfigResponse 积分配置响应
type PointsConfigResponse struct {
	Id          uuid.UUID `json:"id"`
	ConfigKey   string    `json:"configKey"`
	ConfigValue float64   `json:"configValue"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CheckinStatusResponse 签到状态响应
type CheckinStatusResponse struct {
	HasCheckedIn     bool    `json:"hasCheckedIn"`     // 今日是否已签到
	ConsecutiveDays  int     `json:"consecutiveDays"`  // 连续签到天数
	NextBonusDay     int     `json:"nextBonusDay"`     // 下次奖励签到天数
	TotalPoints      float64 `json:"totalPoints"`      // 总积分
	AvailablePoints  float64 `json:"availablePoints"`  // 可用积分
	LastCheckinDate  *string `json:"lastCheckinDate"`  // 最后签到日期
}

// PointsStatisticsResponse 积分统计响应
type PointsStatisticsResponse struct {
	TotalUsers           int64   `json:"totalUsers"`           // 总用户数
	TotalPoints          float64 `json:"totalPoints"`          // 总积分
	TotalUsedPoints      float64 `json:"totalUsedPoints"`      // 总已使用积分
	TotalAvailablePoints float64 `json:"totalAvailablePoints"` // 总可用积分
	TodayCheckins        int64   `json:"todayCheckins"`        // 今日签到数
	TodayExchanges       int64   `json:"todayExchanges"`       // 今日兑换数
	TodayPointsEarned    float64 `json:"todayPointsEarned"`    // 今日获得积分
	TodayPointsUsed      float64 `json:"todayPointsUsed"`      // 今日使用积分
} 