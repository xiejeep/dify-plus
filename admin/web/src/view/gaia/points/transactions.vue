<template>
  <div class="gva-card">
    <div class="gva-card-box">
      <div class="gva-card-header">
        <span class="gva-card-title">积分流水管理</span>
      </div>
      
      <!-- 搜索区域 -->
      <div class="search-container">
        <el-row :gutter="20">
          <el-col :span="5">
            <el-input
              v-model="searchInfo.accountId"
              placeholder="请输入用户ID"
              clearable
            />
          </el-col>
          <el-col :span="4">
            <el-select
              v-model="searchInfo.type"
              placeholder="流水类型"
              clearable
              style="width: 100%"
            >
              <el-option label="签到获得" value="checkin" />
              <el-option label="连续奖励" value="consecutive_bonus" />
              <el-option label="兑换消费" value="exchange" />
              <el-option label="管理员调整" value="admin_adjust" />
              <el-option label="系统补偿" value="system_compensation" />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-select
              v-model="searchInfo.changeType"
              placeholder="变动类型"
              clearable
              style="width: 100%"
            >
              <el-option label="增加" value="increase" />
              <el-option label="减少" value="decrease" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-date-picker
              v-model="searchInfo.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-col>
          <el-col :span="5">
            <el-button type="primary" @click="getTableData">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetSearch">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 统计信息 -->
      <div class="stats-container">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="stat-item increase">
              <div class="stat-value">{{ formatNumber(statistics.totalIncrease) }}</div>
              <div class="stat-label">总增加积分</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item decrease">
              <div class="stat-value">{{ formatNumber(statistics.totalDecrease) }}</div>
              <div class="stat-label">总减少积分</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item today">
              <div class="stat-value">{{ statistics.todayTransactions }}</div>
              <div class="stat-label">今日流水数</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item net">
              <div class="stat-value">{{ formatNumber(statistics.netChange) }}</div>
              <div class="stat-label">净变动积分</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 表格区域 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="accountId" label="用户ID" width="280" />
        <el-table-column prop="type" label="流水类型" width="120">
          <template #default="scope">
            <el-tag :type="getTypeColor(scope.row.type)">
              {{ getTypeLabel(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="changeType" label="变动类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.changeType === 'increase' ? 'success' : 'danger'">
              {{ scope.row.changeType === 'increase' ? '增加' : '减少' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="points" label="积分变动" width="120">
          <template #default="scope">
            <span :class="scope.row.changeType === 'increase' ? 'text-green-600' : 'text-red-600'">
              {{ scope.row.changeType === 'increase' ? '+' : '-' }}{{ scope.row.points }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="beforePoints" label="变动前积分" width="120">
          <template #default="scope">
            <span class="text-gray-600">{{ formatNumber(scope.row.beforePoints) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="afterPoints" label="变动后积分" width="120">
          <template #default="scope">
            <span class="text-gray-600">{{ formatNumber(scope.row.afterPoints) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="150">
          <template #default="scope">
            <span class="text-gray-700">{{ scope.row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="流水时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="viewDetails(scope.row)"
            >
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 25, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="积分流水详情"
      width="600px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户ID" :span="2">
          {{ currentRecord.accountId }}
        </el-descriptions-item>
        <el-descriptions-item label="流水类型">
          <el-tag :type="getTypeColor(currentRecord.type)">
            {{ getTypeLabel(currentRecord.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="变动类型">
          <el-tag :type="currentRecord.changeType === 'increase' ? 'success' : 'danger'">
            {{ currentRecord.changeType === 'increase' ? '增加' : '减少' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="积分变动">
          <span :class="currentRecord.changeType === 'increase' ? 'text-green-600' : 'text-red-600'">
            {{ currentRecord.changeType === 'increase' ? '+' : '-' }}{{ currentRecord.points }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="流水时间">
          {{ formatDate(currentRecord.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="变动前积分">
          {{ formatNumber(currentRecord.beforePoints) }}
        </el-descriptions-item>
        <el-descriptions-item label="变动后积分">
          {{ formatNumber(currentRecord.afterPoints) }}
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ currentRecord.description || '-' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getPointsTransaction } from '@/api/gaia/checkin'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(25)
const total = ref(0)
const multipleSelection = ref([])

const searchInfo = reactive({
  accountId: '',
  type: '',
  changeType: '',
  dateRange: []
})

const statistics = reactive({
  totalIncrease: 0,
  totalDecrease: 0,
  todayTransactions: 0,
  netChange: 0
})

const detailDialogVisible = ref(false)
const currentRecord = ref({})

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const formatDate = (dateStr) => {
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const getTypeLabel = (type) => {
  const typeMap = {
    'checkin': '签到获得',
    'consecutive_bonus': '连续奖励',
    'exchange': '兑换消费',
    'admin_adjust': '管理员调整',
    'system_compensation': '系统补偿'
  }
  return typeMap[type] || type
}

const getTypeColor = (type) => {
  const colorMap = {
    'checkin': 'success',
    'consecutive_bonus': 'warning',
    'exchange': 'danger',
    'admin_adjust': 'info',
    'system_compensation': 'primary'
  }
  return colorMap[type] || 'info'
}

const getTableData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo
    }
    
    // 处理日期范围
    if (searchInfo.dateRange && searchInfo.dateRange.length === 2) {
      params.startDate = searchInfo.dateRange[0]
      params.endDate = searchInfo.dateRange[1]
      delete params.dateRange
    }
    
    // 清理空值参数
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })

    const res = await getPointsTransaction(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
      
      // 更新统计信息
      if (res.data.statistics) {
        Object.assign(statistics, res.data.statistics)
      }
    } else {
      ElMessage.error(res.msg || '获取数据失败')
    }
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchInfo.accountId = ''
  searchInfo.type = ''
  searchInfo.changeType = ''
  searchInfo.dateRange = []
  page.value = 1
  getTableData()
}

const handleSelectionChange = (selection) => {
  multipleSelection.value = selection
}

const handleSizeChange = (size) => {
  pageSize.value = size
  page.value = 1
  getTableData()
}

const handleCurrentChange = (currentPage) => {
  page.value = currentPage
  getTableData()
}

const viewDetails = (row) => {
  currentRecord.value = row
  detailDialogVisible.value = true
}

onMounted(() => {
  getTableData()
})
</script>

<style scoped>
.search-container {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.stats-container {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.stat-item {
  text-align: center;
  padding: 15px;
  border-radius: 8px;
  color: white;
}

.stat-item.increase {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-item.decrease {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-item.today {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-item.net {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.text-green-600 {
  color: #059669;
  font-weight: bold;
}

.text-red-600 {
  color: #dc2626;
  font-weight: bold;
}

.text-gray-600 {
  color: #6b7280;
}

.text-gray-700 {
  color: #374151;
}

@media (max-width: 768px) {
  .stats-container .el-row .el-col {
    margin-bottom: 15px;
  }
}
</style> 