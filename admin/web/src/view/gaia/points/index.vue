<template>
  <div class="gva-card">
    <div class="gva-card-box">
      <div class="gva-card-header">
        <span class="gva-card-title">积分管理</span>
      </div>
      
      <!-- 统计卡片区域 -->
      <div class="stats-container">
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-content">
              <div class="stat-title">总用户数</div>
              <div class="stat-number">{{ statistics.totalUsers }}</div>
            </div>
            <div class="stat-icon">
              <el-icon size="32"><User /></el-icon>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-content">
              <div class="stat-title">总积分</div>
              <div class="stat-number">{{ formatNumber(statistics.totalPoints) }}</div>
            </div>
            <div class="stat-icon">
              <el-icon size="32"><Coin /></el-icon>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-content">
              <div class="stat-title">可用积分</div>
              <div class="stat-number">{{ formatNumber(statistics.totalAvailablePoints) }}</div>
            </div>
            <div class="stat-icon">
              <el-icon size="32"><Money /></el-icon>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-content">
              <div class="stat-title">今日签到</div>
              <div class="stat-number">{{ statistics.todayCheckins }}</div>
            </div>
            <div class="stat-icon">
              <el-icon size="32"><Calendar /></el-icon>
            </div>
          </div>
        </div>
      </div>

      <!-- 今日数据统计 -->
      <div class="today-stats">
        <h3>今日数据</h3>
        <div class="today-grid">
          <div class="today-item">
            <span>今日签到数：</span>
            <strong>{{ statistics.todayCheckins }}</strong>
          </div>
          <div class="today-item">
            <span>今日兑换数：</span>
            <strong>{{ statistics.todayExchanges }}</strong>
          </div>
          <div class="today-item">
            <span>今日获得积分：</span>
            <strong>{{ formatNumber(statistics.todayPointsEarned) }}</strong>
          </div>
          <div class="today-item">
            <span>今日使用积分：</span>
            <strong>{{ formatNumber(statistics.todayPointsUsed) }}</strong>
          </div>
        </div>
      </div>

      <!-- 快捷操作区域 -->
      <div class="quick-actions">
        <h3>快捷操作</h3>
        <div class="action-buttons">
          <el-button type="danger" @click="handleTestCheckin" :loading="checkinLoading">
            <el-icon><Calendar /></el-icon>
            测试签到
          </el-button>
          <el-button type="primary" @click="$router.push('/gaia/points/users')">
            <el-icon><User /></el-icon>
            用户积分管理
          </el-button>
          <el-button type="success" @click="$router.push('/gaia/points/records')">
            <el-icon><Document /></el-icon>
            签到记录管理
          </el-button>
          <el-button type="warning" @click="$router.push('/gaia/points/transactions')">
            <el-icon><Money /></el-icon>
            积分流水管理
          </el-button>
          <el-button type="info" @click="$router.push('/gaia/points/config')">
            <el-icon><Setting /></el-icon>
            积分配置管理
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPointsStatistics, checkin } from '@/api/gaia/checkin'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User, Coin, Money, Calendar, Document, Setting } from '@element-plus/icons-vue'

const statistics = ref({
  totalUsers: 0,
  totalPoints: 0,
  totalUsedPoints: 0,
  totalAvailablePoints: 0,
  todayCheckins: 0,
  todayExchanges: 0,
  todayPointsEarned: 0,
  todayPointsUsed: 0
})

const checkinLoading = ref(false)

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const loadStatistics = async () => {
  try {
    const res = await getPointsStatistics()
    if (res.code === 0) {
      statistics.value = res.data
    } else {
      ElMessage.error(res.msg || '获取统计数据失败')
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  }
}

const handleTestCheckin = async () => {
  try {
    const result = await ElMessageBox.prompt(
      '请输入要测试签到的用户ID (UUID格式)：',
      '测试签到功能',
      {
        confirmButtonText: '确定签到',
        cancelButtonText: '取消',
        inputPlaceholder: '例如：59883ad3-81f9-4650-bef1-5f1ff0856419',
        inputValidator: (value) => {
          if (!value) {
            return '请输入用户ID'
          }
          // 简单的UUID格式验证
          const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
          if (!uuidRegex.test(value)) {
            return 'UUID格式不正确'
          }
          return true
        }
      }
    )

    checkinLoading.value = true
    
    const checkinResult = await checkin({
      accountId: result.value
    })

    if (checkinResult.code === 0) {
      ElMessage.success(`签到成功！获得积分：${checkinResult.data.pointsEarned}`)
      // 刷新统计数据
      await loadStatistics()
    } else {
      ElMessage.error(checkinResult.msg || '签到失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('签到测试失败:', error)
      ElMessage.error('签到测试失败：' + (error.response?.data?.msg || error.message))
    }
  } finally {
    checkinLoading.value = false
  }
}

onMounted(() => {
  loadStatistics()
})
</script>

<style scoped>
.stats-container {
  margin: 20px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card:nth-child(2) {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-card:nth-child(3) {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-card:nth-child(4) {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-content {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  opacity: 0.8;
  margin-bottom: 8px;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
}

.stat-icon {
  opacity: 0.6;
}

.today-stats {
  margin: 30px 0;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.today-stats h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.today-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.today-item {
  padding: 10px;
  background: white;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.quick-actions {
  margin-top: 30px;
}

.quick-actions h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.action-buttons {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  display: flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .today-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style> 