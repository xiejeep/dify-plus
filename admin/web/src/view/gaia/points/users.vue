<template>
  <div class="gva-card">
    <div class="gva-card-box">
      <div class="gva-card-header">
        <span class="gva-card-title">用户积分管理</span>
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
            <el-input
              v-model="searchInfo.minPoints"
              type="number"
              placeholder="最小积分"
              clearable
            />
          </el-col>
          <el-col :span="6">
            <el-input
              v-model="searchInfo.maxPoints"
              type="number"
              placeholder="最大积分"
              clearable
            />
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

      <!-- 表格区域 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="accountId" label="用户ID" width="280" />
        <el-table-column prop="accountName" label="用户名称" min-width="120" />
        <el-table-column prop="totalPoints" label="总积分" width="120">
          <template #default="scope">
            <el-tag type="info">{{ formatNumber(scope.row.totalPoints) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="availablePoints" label="可用积分" width="120">
          <template #default="scope">
            <el-tag type="success">{{ formatNumber(scope.row.availablePoints) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="usedPoints" label="已用积分" width="120">
          <template #default="scope">
            <el-tag type="warning">{{ formatNumber(scope.row.usedPoints) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="viewUserDetails(scope.row)"
            >
              详情
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="adjustPoints(scope.row)"
            >
              调整积分
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

    <!-- 调整积分对话框 -->
    <el-dialog
      v-model="adjustDialogVisible"
      title="调整用户积分"
      width="500px"
      @close="closeAdjustDialog"
    >
      <el-form
        ref="adjustFormRef"
        :model="adjustForm"
        :rules="adjustRules"
        label-width="100px"
      >
        <el-form-item label="用户ID">
          <el-input v-model="adjustForm.accountId" disabled />
        </el-form-item>
        <el-form-item label="当前积分">
          <el-input v-model="currentUserPoints" disabled />
        </el-form-item>
        <el-form-item label="调整积分" prop="pointsChange">
          <el-input
            v-model.number="adjustForm.pointsChange"
            type="number"
            placeholder="正数增加，负数减少"
          />
        </el-form-item>
        <el-form-item label="调整原因" prop="description">
          <el-input
            v-model="adjustForm.description"
            type="textarea"
            placeholder="请输入调整原因"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="adjustDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmAdjustPoints">确认调整</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { getUserPoints, manualAdjustPoints } from '@/api/gaia/checkin'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref([])
const page = ref(1)
const pageSize = ref(25)
const total = ref(0)
const multipleSelection = ref([])

const searchInfo = reactive({
  accountId: '',
  minPoints: '',
  maxPoints: ''
})

const adjustDialogVisible = ref(false)
const adjustFormRef = ref()
const adjustForm = reactive({
  accountId: '',
  pointsChange: 0,
  description: ''
})

const currentUserPoints = computed(() => {
  const user = tableData.value.find(u => u.accountId === adjustForm.accountId)
  return user ? formatNumber(user.availablePoints) : '0'
})

const adjustRules = {
  pointsChange: [
    { required: true, message: '请输入调整积分', trigger: 'blur' },
    { type: 'number', message: '积分必须是数字', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入调整原因', trigger: 'blur' },
    { min: 5, max: 200, message: '调整原因长度在 5 到 200 个字符', trigger: 'blur' }
  ]
}

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getTableData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo
    }
    
    // 清理空字符串参数
    Object.keys(params).forEach(key => {
      if (params[key] === '') {
        delete params[key]
      }
    })

    const res = await getUserPoints(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
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
  searchInfo.minPoints = ''
  searchInfo.maxPoints = ''
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

const viewUserDetails = (row) => {
  // 跳转到用户详情页面，显示该用户的积分详细信息
  ElMessage.info('功能开发中...')
}

const adjustPoints = (row) => {
  adjustForm.accountId = row.accountId
  adjustForm.pointsChange = 0
  adjustForm.description = ''
  adjustDialogVisible.value = true
}

const closeAdjustDialog = () => {
  adjustDialogVisible.value = false
  adjustFormRef.value?.resetFields()
}

const confirmAdjustPoints = async () => {
  try {
    await adjustFormRef.value.validate()
    
    await ElMessageBox.confirm(
      `确认${adjustForm.pointsChange > 0 ? '增加' : '减少'}用户积分 ${Math.abs(adjustForm.pointsChange)} 分吗？`,
      '确认调整',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await manualAdjustPoints(adjustForm)
    if (res.code === 0) {
      ElMessage.success('积分调整成功')
      adjustDialogVisible.value = false
      getTableData()
    } else {
      ElMessage.error(res.msg || '积分调整失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('积分调整失败:', error)
      ElMessage.error('积分调整失败')
    }
  }
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
</style> 