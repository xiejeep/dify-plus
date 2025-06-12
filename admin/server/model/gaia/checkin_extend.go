package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// UserPointsExtend 用户积分账户表
type UserPointsExtend struct {
	Id              uuid.UUID  `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;type:uuid;default:gen_random_uuid();"`
	AccountId       uuid.UUID  `json:"accountId" form:"accountId" gorm:"uniqueIndex;column:account_id;comment:关联用户账户;"`
	TotalPoints     float64    `json:"totalPoints" form:"totalPoints" gorm:"column:total_points;comment:总积分;default:0"`
	AvailablePoints float64    `json:"availablePoints" form:"availablePoints" gorm:"column:available_points;comment:可用积分;default:0"`
	UsedPoints      float64    `json:"usedPoints" form:"usedPoints" gorm:"column:used_points;comment:已使用积分;default:0"`
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;size:6;"`
	UpdatedAt       *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;size:6;"`
}

func (UserPointsExtend) TableName() string {
	return "user_points_extend"
}

// CheckinRecordExtend 用户签到记录表
type CheckinRecordExtend struct {
	Id              uuid.UUID  `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;type:uuid;default:gen_random_uuid();"`
	AccountId       uuid.UUID  `json:"accountId" form:"accountId" gorm:"column:account_id;comment:关联用户账户;index"`
	CheckinDate     time.Time  `json:"checkinDate" form:"checkinDate" gorm:"column:checkin_date;comment:签到日期;type:date;index;uniqueIndex:uk_account_checkin_date,composite:account_id"`
	PointsEarned    float64    `json:"pointsEarned" form:"pointsEarned" gorm:"column:points_earned;comment:签到获得积分;"`
	ConsecutiveDays int        `json:"consecutiveDays" form:"consecutiveDays" gorm:"column:consecutive_days;comment:连续签到天数;default:1"`
	IsBonus         bool       `json:"isBonus" form:"isBonus" gorm:"column:is_bonus;comment:是否为奖励签到;default:false"`
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;size:6;"`
}

func (CheckinRecordExtend) TableName() string {
	return "checkin_record_extend"
}

// PointsTransactionExtend 积分流水记录表
type PointsTransactionExtend struct {
	Id              uuid.UUID  `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;type:uuid;default:gen_random_uuid();"`
	AccountId       uuid.UUID  `json:"accountId" form:"accountId" gorm:"column:account_id;comment:关联用户账户;index"`
	TransactionType string     `json:"transactionType" form:"transactionType" gorm:"column:transaction_type;comment:交易类型;size:50;index"`
	PointsChange    float64    `json:"pointsChange" form:"pointsChange" gorm:"column:points_change;comment:积分变化;"`
	PointsBefore    float64    `json:"pointsBefore" form:"pointsBefore" gorm:"column:points_before;comment:交易前积分;"`
	PointsAfter     float64    `json:"pointsAfter" form:"pointsAfter" gorm:"column:points_after;comment:交易后积分;"`
	Description     string     `json:"description" form:"description" gorm:"column:description;comment:交易描述;size:200"`
	RelatedId       *uuid.UUID `json:"relatedId" form:"relatedId" gorm:"column:related_id;comment:关联ID;"`
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;size:6;index"`
}

func (PointsTransactionExtend) TableName() string {
	return "points_transaction_extend"
}

// PointsExchangeExtend 积分兑换记录表
type PointsExchangeExtend struct {
	Id           uuid.UUID  `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;type:uuid;default:gen_random_uuid();"`
	AccountId    uuid.UUID  `json:"accountId" form:"accountId" gorm:"column:account_id;comment:关联用户账户;index"`
	ExchangeType string     `json:"exchangeType" form:"exchangeType" gorm:"column:exchange_type;comment:兑换类型;size:50"`
	PointsCost   float64    `json:"pointsCost" form:"pointsCost" gorm:"column:points_cost;comment:消耗积分;"`
	QuotaAmount  *float64   `json:"quotaAmount" form:"quotaAmount" gorm:"column:quota_amount;comment:兑换的额度金额;"`
	Status       string     `json:"status" form:"status" gorm:"column:status;comment:兑换状态;size:20;default:completed;index"`
	Description  string     `json:"description" form:"description" gorm:"column:description;comment:兑换描述;size:200"`
	CreatedAt    *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;size:6;"`
	UpdatedAt    *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;size:6;"`
}

func (PointsExchangeExtend) TableName() string {
	return "points_exchange_extend"
}

// PointsConfigExtend 积分配置表
type PointsConfigExtend struct {
	Id          uuid.UUID  `json:"id" form:"id" gorm:"primarykey;column:id;comment:主键;type:uuid;default:gen_random_uuid();"`
	ConfigKey   string     `json:"configKey" form:"configKey" gorm:"column:config_key;comment:配置键;size:100;uniqueIndex"`
	ConfigValue float64    `json:"configValue" form:"configValue" gorm:"column:config_value;comment:配置值;"`
	Description string     `json:"description" form:"description" gorm:"column:description;comment:配置描述;size:200"`
	CreatedAt   *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间;size:6;"`
	UpdatedAt   *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;comment:更新时间;size:6;"`
}

func (PointsConfigExtend) TableName() string {
	return "points_config_extend"
} 