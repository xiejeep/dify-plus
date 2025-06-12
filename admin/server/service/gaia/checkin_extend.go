package gaia

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	gaiaRes "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/response"
	"github.com/gofrs/uuid/v5"
)

type CheckinService struct{}

// GetPointsConfig 获取积分配置
func (cs *CheckinService) GetPointsConfig(configKey string) (float64, error) {
	var config gaia.PointsConfigExtend
	err := global.GVA_DB.Where("config_key = ?", configKey).First(&config).Error
	if err != nil {
		// 返回默认配置
		switch configKey {
		case "daily_checkin_points":
			return 10.0, nil // 每日签到积分
		case "consecutive_bonus_days":
			return 7.0, nil // 连续签到奖励天数
		case "consecutive_bonus_points":
			return 50.0, nil // 连续签到奖励积分
		case "points_to_quota_rate":
			return 100.0, nil // 积分兑换额度比例(100积分=1美元)
		default:
			return 0.0, errors.New("未知的配置项")
		}
	}
	return config.ConfigValue, nil
}

// InitUserPoints 初始化用户积分账户
func (cs *CheckinService) InitUserPoints(accountId uuid.UUID) (*gaia.UserPointsExtend, error) {
	var userPoints gaia.UserPointsExtend
	err := global.GVA_DB.Where("account_id = ?", accountId).First(&userPoints).Error
	if err != nil {
			// 创建新的积分账户
	now := time.Now()
	newId, _ := uuid.NewV4()
	userPoints = gaia.UserPointsExtend{
		Id:              newId,
		AccountId:       accountId,
		TotalPoints:     0.0,
		AvailablePoints: 0.0,
		UsedPoints:      0.0,
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}
		err = global.GVA_DB.Create(&userPoints).Error
		if err != nil {
			return nil, err
		}
	}
	return &userPoints, nil
}

// CheckTodayCheckin 检查今日是否已签到
func (cs *CheckinService) CheckTodayCheckin(accountId uuid.UUID) (bool, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := global.GVA_DB.Model(&gaia.CheckinRecordExtend{}).
		Where("account_id = ? AND checkin_date = ?", accountId, today).
		Count(&count).Error
	return count > 0, err
}

// GetConsecutiveCheckinDays 获取连续签到天数
func (cs *CheckinService) GetConsecutiveCheckinDays(accountId uuid.UUID) (int, error) {
	var records []gaia.CheckinRecordExtend
	err := global.GVA_DB.Where("account_id = ?", accountId).
		Order("checkin_date DESC").
		Limit(30). // 最多查询30天
		Find(&records).Error
	if err != nil {
		return 0, err
	}

	if len(records) == 0 {
		return 0, nil
	}

	consecutiveDays := 0
	today := time.Now()
	
	for i, record := range records {
		expectedDate := today.AddDate(0, 0, -i)
		recordDate := record.CheckinDate
		
		// 比较日期（只比较年月日）
		if recordDate.Year() == expectedDate.Year() && 
		   recordDate.Month() == expectedDate.Month() && 
		   recordDate.Day() == expectedDate.Day() {
			consecutiveDays++
		} else {
			break
		}
	}

	return consecutiveDays, nil
}

// Checkin 执行签到
func (cs *CheckinService) Checkin(req gaiaReq.CheckinRequest) (*gaiaRes.CheckinResponse, error) {
	// 检查今日是否已签到
	hasCheckedIn, err := cs.CheckTodayCheckin(req.AccountId)
	if err != nil {
		return nil, err
	}
	if hasCheckedIn {
		return &gaiaRes.CheckinResponse{
			Success: false,
			Message: "今日已签到",
		}, nil
	}

	// 初始化用户积分账户
	userPoints, err := cs.InitUserPoints(req.AccountId)
	if err != nil {
		return nil, err
	}

	// 获取连续签到天数
	consecutiveDays, err := cs.GetConsecutiveCheckinDays(req.AccountId)
	if err != nil {
		return nil, err
	}
	consecutiveDays++ // 包含今天

	// 获取积分配置
	dailyPoints, err := cs.GetPointsConfig("daily_checkin_points")
	if err != nil {
		return nil, err
	}

	bonusDays, err := cs.GetPointsConfig("consecutive_bonus_days")
	if err != nil {
		return nil, err
	}

	bonusPoints, err := cs.GetPointsConfig("consecutive_bonus_points")
	if err != nil {
		return nil, err
	}

	// 计算获得积分
	pointsEarned := dailyPoints
	isBonus := false
	
	// 检查是否获得连续签到奖励
	if consecutiveDays%int(bonusDays) == 0 {
		pointsEarned += bonusPoints
		isBonus = true
	}

	// 开始数据库事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建签到记录
	now := time.Now()
	checkinId, _ := uuid.NewV4()
	checkinRecord := gaia.CheckinRecordExtend{
		Id:              checkinId,
		AccountId:       req.AccountId,
		CheckinDate:     now,
		PointsEarned:    pointsEarned,
		ConsecutiveDays: consecutiveDays,
		IsBonus:         isBonus,
		CreatedAt:       &now,
	}
	err = tx.Create(&checkinRecord).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户积分
	pointsBefore := userPoints.AvailablePoints
	err = tx.Model(&gaia.UserPointsExtend{}).
		Where("account_id = ?", req.AccountId).
		Updates(map[string]interface{}{
			"total_points":     userPoints.TotalPoints + pointsEarned,
			"available_points": userPoints.AvailablePoints + pointsEarned,
			"updated_at":       now,
		}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建积分流水记录
	pointsAfter := pointsBefore + pointsEarned
	transactionId, _ := uuid.NewV4()
	transaction := gaia.PointsTransactionExtend{
		Id:              transactionId,
		AccountId:       req.AccountId,
		TransactionType: "checkin",
		PointsChange:    pointsEarned,
		PointsBefore:    pointsBefore,
		PointsAfter:     pointsAfter,
		Description:     fmt.Sprintf("每日签到获得积分，连续签到%d天", consecutiveDays),
		RelatedId:       &checkinRecord.Id,
		CreatedAt:       &now,
	}
	err = tx.Create(&transaction).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// 更新userPoints以返回最新数据
	userPoints.TotalPoints += pointsEarned
	userPoints.AvailablePoints += pointsEarned

	return &gaiaRes.CheckinResponse{
		Success:         true,
		Message:         "签到成功",
		PointsEarned:    pointsEarned,
		ConsecutiveDays: consecutiveDays,
		IsBonus:         isBonus,
		TotalPoints:     userPoints.TotalPoints,
		AvailablePoints: userPoints.AvailablePoints,
	}, nil
}

// ExchangePoints 积分兑换
func (cs *CheckinService) ExchangePoints(req gaiaReq.PointsExchangeRequest) (*gaiaRes.PointsExchangeResponse, error) {
	// 检查用户积分账户
	userPoints, err := cs.InitUserPoints(req.AccountId)
	if err != nil {
		return nil, err
	}

	// 检查积分是否足够
	if userPoints.AvailablePoints < req.PointsCost {
		return nil, errors.New("积分不足")
	}

	// 计算兑换额度
	var quotaAmount float64
	if req.ExchangeType == "quota" {
		rate, err := cs.GetPointsConfig("points_to_quota_rate")
		if err != nil {
			return nil, err
		}
		quotaAmount = req.PointsCost / rate
		req.QuotaAmount = &quotaAmount
	}

	// 开始数据库事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建兑换记录
	now := time.Now()
	exchangeId, _ := uuid.NewV4()
	exchangeRecord := gaia.PointsExchangeExtend{
		Id:           exchangeId,
		AccountId:    req.AccountId,
		ExchangeType: req.ExchangeType,
		PointsCost:   req.PointsCost,
		QuotaAmount:  req.QuotaAmount,
		Status:       "completed",
		Description:  req.Description,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	err = tx.Create(&exchangeRecord).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户积分
	pointsBefore := userPoints.AvailablePoints
	err = tx.Model(&gaia.UserPointsExtend{}).
		Where("account_id = ?", req.AccountId).
		Updates(map[string]interface{}{
			"available_points": userPoints.AvailablePoints - req.PointsCost,
			"used_points":      userPoints.UsedPoints + req.PointsCost,
			"updated_at":       now,
		}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建积分流水记录
	pointsAfter := pointsBefore - req.PointsCost
	transactionId2, _ := uuid.NewV4()
	transaction := gaia.PointsTransactionExtend{
		Id:              transactionId2,
		AccountId:       req.AccountId,
		TransactionType: "exchange",
		PointsChange:    -req.PointsCost,
		PointsBefore:    pointsBefore,
		PointsAfter:     pointsAfter,
		Description:     fmt.Sprintf("兑换%s，消耗%.2f积分", req.ExchangeType, req.PointsCost),
		RelatedId:       &exchangeRecord.Id,
		CreatedAt:       &now,
	}
	err = tx.Create(&transaction).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果是额度兑换，更新用户额度
	if req.ExchangeType == "quota" && req.QuotaAmount != nil {
		var accountMoney gaia.AccountMoneyExtend
		err = tx.Where("account_id = ?", req.AccountId).First(&accountMoney).Error
		if err != nil {
			// 如果不存在，创建新记录
			accountMoney = gaia.AccountMoneyExtend{
				AccountId:  req.AccountId,
				TotalQuota: *req.QuotaAmount,
				UsedQuota:  0.0,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			}
			err = tx.Create(&accountMoney).Error
		} else {
			// 更新现有记录
			err = tx.Model(&accountMoney).Updates(map[string]interface{}{
				"total_quota": accountMoney.TotalQuota + *req.QuotaAmount,
				"updated_at":  now,
			}).Error
		}
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &gaiaRes.PointsExchangeResponse{
		Id:           exchangeRecord.Id,
		AccountId:    exchangeRecord.AccountId,
		ExchangeType: exchangeRecord.ExchangeType,
		PointsCost:   exchangeRecord.PointsCost,
		QuotaAmount:  exchangeRecord.QuotaAmount,
		Status:       exchangeRecord.Status,
		Description:  exchangeRecord.Description,
		CreatedAt:    *exchangeRecord.CreatedAt,
		UpdatedAt:    *exchangeRecord.UpdatedAt,
	}, nil
}

// GetCheckinStatus 获取签到状态
func (cs *CheckinService) GetCheckinStatus(accountId uuid.UUID) (*gaiaRes.CheckinStatusResponse, error) {
	// 检查今日是否已签到
	hasCheckedIn, err := cs.CheckTodayCheckin(accountId)
	if err != nil {
		return nil, err
	}

	// 获取连续签到天数
	consecutiveDays, err := cs.GetConsecutiveCheckinDays(accountId)
	if err != nil {
		return nil, err
	}

	// 获取用户积分
	userPoints, err := cs.InitUserPoints(accountId)
	if err != nil {
		return nil, err
	}

	// 获取最后签到日期
	var lastRecord gaia.CheckinRecordExtend
	var lastCheckinDate *string
	err = global.GVA_DB.Where("account_id = ?", accountId).
		Order("checkin_date DESC").
		First(&lastRecord).Error
	if err == nil {
		dateStr := lastRecord.CheckinDate.Format("2006-01-02")
		lastCheckinDate = &dateStr
	}

	// 计算下次奖励签到天数
	bonusDays, _ := cs.GetPointsConfig("consecutive_bonus_days")
	nextBonusDay := int(bonusDays) - (consecutiveDays % int(bonusDays))
	if nextBonusDay == int(bonusDays) {
		nextBonusDay = 0
	}

	return &gaiaRes.CheckinStatusResponse{
		HasCheckedIn:    hasCheckedIn,
		ConsecutiveDays: consecutiveDays,
		NextBonusDay:    nextBonusDay,
		TotalPoints:     userPoints.TotalPoints,
		AvailablePoints: userPoints.AvailablePoints,
		LastCheckinDate: lastCheckinDate,
	}, nil
}

// GetUserPoints 获取用户积分列表
func (cs *CheckinService) GetUserPoints(req gaiaReq.GetUserPointsRequest) ([]gaiaRes.UserPointsResponse, int64, error) {
	var userPoints []gaia.UserPointsExtend
	var total int64

	// 构建查询条件
	db := global.GVA_DB.Model(&gaia.UserPointsExtend{})
	
	if req.AccountId != nil {
		db = db.Where("account_id = ?", *req.AccountId)
	}
	if req.MinPoints != nil {
		db = db.Where("available_points >= ?", *req.MinPoints)
	}
	if req.MaxPoints != nil {
		db = db.Where("available_points <= ?", *req.MaxPoints)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&userPoints).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var result []gaiaRes.UserPointsResponse
	for _, up := range userPoints {
		response := gaiaRes.UserPointsResponse{
			Id:              up.Id,
			AccountId:       up.AccountId,
			TotalPoints:     up.TotalPoints,
			AvailablePoints: up.AvailablePoints,
			UsedPoints:      up.UsedPoints,
			CreatedAt:       *up.CreatedAt,
			UpdatedAt:       *up.UpdatedAt,
		}
		
		// 获取账户名称
		var account gaia.Account
		err = global.GVA_DB.Where("id = ?", up.AccountId).First(&account).Error
		if err == nil {
			response.AccountName = account.Name
		}
		
		result = append(result, response)
	}

	return result, total, nil
}

// UpdatePointsConfig 更新积分配置
func (cs *CheckinService) UpdatePointsConfig(req gaiaReq.UpdatePointsConfigRequest) error {
	now := time.Now()
	
	var config gaia.PointsConfigExtend
	err := global.GVA_DB.Where("config_key = ?", req.ConfigKey).First(&config).Error
	if err != nil {
			// 创建新配置
	configId, _ := uuid.NewV4()
	config = gaia.PointsConfigExtend{
		Id:          configId,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
		return global.GVA_DB.Create(&config).Error
	} else {
		// 更新现有配置
		return global.GVA_DB.Model(&config).Updates(map[string]interface{}{
			"config_value": req.ConfigValue,
			"description":  req.Description,
			"updated_at":   now,
		}).Error
	}
}

// ManualAdjustPoints 手动调整积分
func (cs *CheckinService) ManualAdjustPoints(req gaiaReq.ManualAdjustPointsRequest) error {
	// 检查用户积分账户
	userPoints, err := cs.InitUserPoints(req.AccountId)
	if err != nil {
		return err
	}

	// 检查积分是否足够（如果是扣减）
	if req.PointsChange < 0 && userPoints.AvailablePoints < -req.PointsChange {
		return errors.New("可用积分不足")
	}

	// 开始数据库事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户积分
	pointsBefore := userPoints.AvailablePoints
	now := time.Now()
	
	updates := map[string]interface{}{
		"updated_at": now,
	}
	
	if req.PointsChange > 0 {
		// 增加积分
		updates["total_points"] = userPoints.TotalPoints + req.PointsChange
		updates["available_points"] = userPoints.AvailablePoints + req.PointsChange
	} else {
		// 减少积分
		updates["available_points"] = userPoints.AvailablePoints + req.PointsChange
		updates["used_points"] = userPoints.UsedPoints - req.PointsChange
	}

	err = tx.Model(&gaia.UserPointsExtend{}).
		Where("account_id = ?", req.AccountId).
		Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 创建积分流水记录
	pointsAfter := pointsBefore + req.PointsChange
	transactionId3, _ := uuid.NewV4()
	transaction := gaia.PointsTransactionExtend{
		Id:              transactionId3,
		AccountId:       req.AccountId,
		TransactionType: "manual",
		PointsChange:    req.PointsChange,
		PointsBefore:    pointsBefore,
		PointsAfter:     pointsAfter,
		Description:     req.Description,
		CreatedAt:       &now,
	}
	err = tx.Create(&transaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// GetCheckinRecords 获取签到记录列表
func (cs *CheckinService) GetCheckinRecords(req gaiaReq.GetCheckinRecordsRequest) ([]gaiaRes.CheckinRecordResponse, int64, error) {
	var records []gaia.CheckinRecordExtend
	var total int64

	// 构建查询条件
	db := global.GVA_DB.Model(&gaia.CheckinRecordExtend{})
	
	if req.AccountId != nil {
		db = db.Where("account_id = ?", *req.AccountId)
	}
	if req.StartDate != nil {
		db = db.Where("checkin_date >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("checkin_date <= ?", *req.EndDate)
	}
	if req.IsBonus != nil {
		db = db.Where("is_bonus = ?", *req.IsBonus)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("checkin_date DESC").Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var result []gaiaRes.CheckinRecordResponse
	for _, record := range records {
		response := gaiaRes.CheckinRecordResponse{
			Id:              record.Id,
			AccountId:       record.AccountId,
			CheckinDate:     record.CheckinDate,
			PointsEarned:    record.PointsEarned,
			ConsecutiveDays: record.ConsecutiveDays,
			IsBonus:         record.IsBonus,
			CreatedAt:       *record.CreatedAt,
		}
		
		// 获取账户名称
		var account gaia.Account
		err = global.GVA_DB.Where("id = ?", record.AccountId).First(&account).Error
		if err == nil {
			response.AccountName = account.Name
		}
		
		result = append(result, response)
	}

	return result, total, nil
}

// GetPointsTransaction 获取积分流水记录
func (cs *CheckinService) GetPointsTransaction(req gaiaReq.GetPointsTransactionRequest) ([]gaiaRes.PointsTransactionResponse, int64, error) {
	var transactions []gaia.PointsTransactionExtend
	var total int64

	// 构建查询条件
	db := global.GVA_DB.Model(&gaia.PointsTransactionExtend{})
	
	if req.AccountId != nil {
		db = db.Where("account_id = ?", *req.AccountId)
	}
	if req.TransactionType != nil {
		db = db.Where("transaction_type = ?", *req.TransactionType)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", *req.EndDate)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var result []gaiaRes.PointsTransactionResponse
	for _, transaction := range transactions {
		response := gaiaRes.PointsTransactionResponse{
			Id:              transaction.Id,
			AccountId:       transaction.AccountId,
			TransactionType: transaction.TransactionType,
			PointsChange:    transaction.PointsChange,
			PointsBefore:    transaction.PointsBefore,
			PointsAfter:     transaction.PointsAfter,
			Description:     transaction.Description,
			RelatedId:       transaction.RelatedId,
			CreatedAt:       *transaction.CreatedAt,
		}
		
		// 获取账户名称
		var account gaia.Account
		err = global.GVA_DB.Where("id = ?", transaction.AccountId).First(&account).Error
		if err == nil {
			response.AccountName = account.Name
		}
		
		result = append(result, response)
	}

	return result, total, nil
}

// GetPointsExchange 获取积分兑换记录
func (cs *CheckinService) GetPointsExchange(req gaiaReq.GetPointsExchangeRequest) ([]gaiaRes.PointsExchangeResponse, int64, error) {
	var exchanges []gaia.PointsExchangeExtend
	var total int64

	// 构建查询条件
	db := global.GVA_DB.Model(&gaia.PointsExchangeExtend{})
	
	if req.AccountId != nil {
		db = db.Where("account_id = ?", *req.AccountId)
	}
	if req.ExchangeType != nil {
		db = db.Where("exchange_type = ?", *req.ExchangeType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", *req.EndDate)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&exchanges).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var result []gaiaRes.PointsExchangeResponse
	for _, exchange := range exchanges {
		response := gaiaRes.PointsExchangeResponse{
			Id:           exchange.Id,
			AccountId:    exchange.AccountId,
			ExchangeType: exchange.ExchangeType,
			PointsCost:   exchange.PointsCost,
			QuotaAmount:  exchange.QuotaAmount,
			Status:       exchange.Status,
			Description:  exchange.Description,
			CreatedAt:    *exchange.CreatedAt,
			UpdatedAt:    *exchange.UpdatedAt,
		}
		
		// 获取账户名称
		var account gaia.Account
		err = global.GVA_DB.Where("id = ?", exchange.AccountId).First(&account).Error
		if err == nil {
			response.AccountName = account.Name
		}
		
		result = append(result, response)
	}

	return result, total, nil
}

// GetAllPointsConfig 获取所有积分配置
func (cs *CheckinService) GetAllPointsConfig() ([]gaiaRes.PointsConfigResponse, error) {
	var configs []gaia.PointsConfigExtend
	err := global.GVA_DB.Find(&configs).Error
	if err != nil {
		return nil, err
	}

	var result []gaiaRes.PointsConfigResponse
	for _, config := range configs {
		result = append(result, gaiaRes.PointsConfigResponse{
			Id:          config.Id,
			ConfigKey:   config.ConfigKey,
			ConfigValue: config.ConfigValue,
			Description: config.Description,
			CreatedAt:   *config.CreatedAt,
			UpdatedAt:   *config.UpdatedAt,
		})
	}

	return result, nil
}

// GetUserPointsByAccountId 根据用户ID获取积分信息
func (cs *CheckinService) GetUserPointsByAccountId(accountId uuid.UUID) (*gaiaRes.UserPointsResponse, error) {
	userPoints, err := cs.InitUserPoints(accountId)
	if err != nil {
		return nil, err
	}

	result := &gaiaRes.UserPointsResponse{
		Id:              userPoints.Id,
		AccountId:       userPoints.AccountId,
		TotalPoints:     userPoints.TotalPoints,
		AvailablePoints: userPoints.AvailablePoints,
		UsedPoints:      userPoints.UsedPoints,
		CreatedAt:       *userPoints.CreatedAt,
		UpdatedAt:       *userPoints.UpdatedAt,
	}

	// 获取账户名称
	var account gaia.Account
	err = global.GVA_DB.Where("id = ?", accountId).First(&account).Error
	if err == nil {
		result.AccountName = account.Name
	}

	return result, nil
}

// GetPointsStatistics 获取积分统计
func (cs *CheckinService) GetPointsStatistics() (*gaiaRes.PointsStatisticsResponse, error) {
	today := time.Now().Format("2006-01-02")
	
	// 总用户数
	var totalUsers int64
	err := global.GVA_DB.Model(&gaia.UserPointsExtend{}).Count(&totalUsers).Error
	if err != nil {
		return nil, err
	}

	// 总积分统计
	var totalStats struct {
		TotalPoints     *float64
		TotalUsedPoints *float64
		TotalAvailable  *float64
	}
	err = global.GVA_DB.Model(&gaia.UserPointsExtend{}).
		Select("SUM(total_points) as total_points, SUM(used_points) as total_used_points, SUM(available_points) as total_available").
		Scan(&totalStats).Error
	if err != nil {
		return nil, err
	}

	// 处理NULL值，默认为0
	var totalPoints, totalUsedPoints, totalAvailable float64
	if totalStats.TotalPoints != nil {
		totalPoints = *totalStats.TotalPoints
	}
	if totalStats.TotalUsedPoints != nil {
		totalUsedPoints = *totalStats.TotalUsedPoints
	}
	if totalStats.TotalAvailable != nil {
		totalAvailable = *totalStats.TotalAvailable
	}

	// 今日签到数
	var todayCheckins int64
	err = global.GVA_DB.Model(&gaia.CheckinRecordExtend{}).
		Where("checkin_date = ?", today).
		Count(&todayCheckins).Error
	if err != nil {
		return nil, err
	}

	// 今日兑换数
	var todayExchanges int64
	err = global.GVA_DB.Model(&gaia.PointsExchangeExtend{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayExchanges).Error
	if err != nil {
		return nil, err
	}

	// 今日获得积分
	var todayPointsEarnedPtr *float64
	err = global.GVA_DB.Model(&gaia.CheckinRecordExtend{}).
		Where("checkin_date = ?", today).
		Select("SUM(points_earned)").
		Scan(&todayPointsEarnedPtr).Error
	if err != nil {
		return nil, err
	}

	var todayPointsEarned float64
	if todayPointsEarnedPtr != nil {
		todayPointsEarned = *todayPointsEarnedPtr
	}

	// 今日使用积分
	var todayPointsUsedPtr *float64
	err = global.GVA_DB.Model(&gaia.PointsExchangeExtend{}).
		Where("DATE(created_at) = ?", today).
		Select("SUM(points_cost)").
		Scan(&todayPointsUsedPtr).Error
	if err != nil {
		return nil, err
	}

	var todayPointsUsed float64
	if todayPointsUsedPtr != nil {
		todayPointsUsed = *todayPointsUsedPtr
	}

	return &gaiaRes.PointsStatisticsResponse{
		TotalUsers:           totalUsers,
		TotalPoints:          totalPoints,
		TotalUsedPoints:      totalUsedPoints,
		TotalAvailablePoints: totalAvailable,
		TodayCheckins:        todayCheckins,
		TodayExchanges:       todayExchanges,
		TodayPointsEarned:    todayPointsEarned,
		TodayPointsUsed:      todayPointsUsed,
	}, nil
} 