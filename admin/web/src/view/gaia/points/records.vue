<template>
  <div class="gva-card">
    <div class="gva-card-box">
      <div class="gva-card-header">
        <span class="gva-card-title">签到记录管理</span>
      </div>
      
      <!-- 搜索区域 -->
      <div class="search-container">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-input
              v-model="searchInfo.accountId"
              placeholder="请输入用户ID"
              clearable
            />
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
          <el-col :span="6">
            <el-select
              v-model="searchInfo.isConsecutive"
              placeholder="是否连续签到"
              clearable
              style="width: 100%"
            >
              <el-option label="连续签到" :value="true" />
              <el-option label="普通签到" :value="false" />
            </el-select>
          </el-col>
          <el-col :span="6">
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
            <div class="stat-item">
              <div class="stat-value">{{ statistics.totalRecords }}</div>
              <div class="stat-label">总签到次数</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-value">{{ statistics.todayRecords }}</div>
              <div class="stat-label">今日签到</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-value">{{ statistics.consecutiveRecords }}</div>
              <div class="stat-label">连续签到奖励</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-value">{{ formatNumber(statistics.totalPoints) }}</div>
              <div class="stat-label">总获得积分</div>
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
        <el-table-column prop="checkinDate" label="签到日期" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.checkinDate, 'YYYY-MM-DD') }}
          </template>
        </el-table-column>
        <el-table-column prop="points" label="获得积分" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.isConsecutive ? 'warning' : 'success'">
              +{{ scope.row.points }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isConsecutive" label="签到类型" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.isConsecutive ? 'warning' : 'info'">
              {{ scope.row.isConsecutive ? '连续奖励' : '普通签到' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="consecutiveDays" label="连续天数" width="120">
          <template #default="scope">
            <span v-if="scope.row.consecutiveDays > 0">
              {{ scope.row.consecutiveDays }} 天
            </span>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="签到时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120">
          <template #default="scope">
            <span v-if="scope.row.remark" class="text-gray-600">
              {{ scope.row.remark }}
            </span>
            <span v-else class="text-gray-400">-</span>
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
      title="签到记录详情"
      width="500px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户ID">
          {{ currentRecord.accountId }}
        </el-descriptions-item>
        <el-descriptions-item label="签到日期">
          {{ formatDate(currentRecord.checkinDate, 'YYYY-MM-DD') }}
        </el-descriptions-item>
        <el-descriptions-item label="签到时间">
          {{ formatDate(currentRecord.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="获得积分">
          <el-tag :type="currentRecord.isConsecutive ? 'warning' : 'success'">
            +{{ currentRecord.points }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="签到类型">
          <el-tag :type="currentRecord.isConsecutive ? 'warning' : 'info'">
            {{ currentRecord.isConsecutive ? '连续奖励' : '普通签到' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.consecutiveDays > 0" label="连续天数">
          {{ currentRecord.consecutiveDays }} 天
        </el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.remark" label="备注">
          {{ currentRecord.remark }}
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
import { getCheckinRecords } from '@/api/gaia/checkin'
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
  dateRange: [],
  isConsecutive: null
})

const statistics = reactive({
  totalRecords: 0,
  todayRecords: 0,
  consecutiveRecords: 0,
  totalPoints: 0
})

const detailDialogVisible = ref(false)
const currentRecord = ref({})

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const formatDate = (dateStr, format = 'YYYY-MM-DD HH:mm:ss') => {
  return dayjs(dateStr).format(format)
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

    const res = await getCheckinRecords(params)
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
  searchInfo.dateRange = []
  searchInfo.isConsecutive = null
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.stat-item:nth-child(1) .stat-item {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-item:nth-child(2) .stat-item {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-item:nth-child(3) .stat-item {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-item:nth-child(4) .stat-item {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
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

.text-gray-400 {
  color: #9ca3af;
}

.text-gray-600 {
  color: #6b7280;
}

@media (max-width: 768px) {
  .stats-container .el-row .el-col {
    margin-bottom: 15px;
  }
}
</style> 