import service from '@/utils/request'

// 用户签到
export const checkin = (data) => {
  return service({
    url: '/gaia/checkin/checkin',
    method: 'post',
    data
  })
}

// 获取签到状态
export const getCheckinStatus = (params) => {
  return service({
    url: '/gaia/checkin/getStatus',
    method: 'get',
    params
  })
}

// 根据账户ID获取积分信息
export const getUserPointsByAccountId = (accountId) => {
  return service({
    url: `/gaia/checkin/getUserPointsByAccountId/${accountId}`,
    method: 'get'
  })
}

// 积分兑换
export const exchangePoints = (data) => {
  return service({
    url: '/gaia/checkin/exchangePoints',
    method: 'post',
    data
  })
}

// 获取用户积分列表
export const getUserPoints = (params) => {
  return service({
    url: '/gaia/checkin/getUserPoints',
    method: 'get',
    params
  })
}

// 获取签到记录
export const getCheckinRecords = (params) => {
  return service({
    url: '/gaia/checkin/getCheckinRecords',
    method: 'get',
    params
  })
}

// 获取积分流水
export const getPointsTransaction = (params) => {
  return service({
    url: '/gaia/checkin/getPointsTransaction',
    method: 'get',
    params
  })
}

// 获取积分兑换记录
export const getPointsExchange = (params) => {
  return service({
    url: '/gaia/checkin/getPointsExchange',
    method: 'get',
    params
  })
}

// 获取积分配置
export const getPointsConfig = () => {
  return service({
    url: '/gaia/checkin/getPointsConfig',
    method: 'get'
  })
}

// 更新积分配置
export const updatePointsConfig = (data) => {
  return service({
    url: '/gaia/checkin/updatePointsConfig',
    method: 'post',
    data
  })
}

// 手动调整积分
export const manualAdjustPoints = (data) => {
  return service({
    url: '/gaia/checkin/manualAdjustPoints',
    method: 'post',
    data
  })
}

// 获取积分统计
export const getPointsStatistics = () => {
  return service({
    url: '/gaia/checkin/getPointsStatistics',
    method: 'get'
  })
} 