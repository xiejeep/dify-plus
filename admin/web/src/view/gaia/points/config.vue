<template>
  <div class="gva-card">
    <div class="gva-card-box">
      <div class="gva-card-header">
        <span class="gva-card-title">积分配置管理</span>
      </div>
      
      <!-- 配置列表 -->
      <el-table
        v-loading="loading"
        :data="configData"
        style="width: 100%"
      >
        <el-table-column prop="configKey" label="配置键" width="200" />
        <el-table-column prop="description" label="配置说明" min-width="250" />
        <el-table-column prop="configValue" label="配置值" width="150">
          <template #default="scope">
            <el-tag type="primary">{{ formatConfigValue(scope.row) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="editConfig(scope.row)"
            >
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 配置说明 -->
      <div class="config-help">
        <h3>配置说明</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="daily_checkin_points">
            <span>每日签到获得的基础积分数量</span>
          </el-descriptions-item>
          <el-descriptions-item label="consecutive_bonus_days">
            <span>连续签到多少天后可获得奖励积分</span>
          </el-descriptions-item>
          <el-descriptions-item label="consecutive_bonus_points">
            <span>连续签到奖励的积分数量</span>
          </el-descriptions-item>
          <el-descriptions-item label="points_to_quota_rate">
            <span>积分兑换额度的比例（多少积分=1美元）</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- 编辑配置对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑积分配置"
      width="500px"
      @close="closeEditDialog"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="120px"
      >
        <el-form-item label="配置键">
          <el-input v-model="editForm.configKey" disabled />
        </el-form-item>
        <el-form-item label="配置说明">
          <el-input v-model="editForm.description" disabled />
        </el-form-item>
        <el-form-item label="配置值" prop="configValue">
          <el-input-number
            v-model="editForm.configValue"
            :precision="1"
            :step="0.1"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="editForm.configKey === 'daily_checkin_points'" label="">
          <div class="config-tip">
            <el-icon><InfoFilled /></el-icon>
            建议设置为 5-20 积分，过低影响用户积极性，过高可能造成积分通胀
          </div>
        </el-form-item>
        <el-form-item v-if="editForm.configKey === 'consecutive_bonus_days'" label="">
          <div class="config-tip">
            <el-icon><InfoFilled /></el-icon>
            建议设置为 3-10 天，过短奖励效果不明显，过长用户难以坚持
          </div>
        </el-form-item>
        <el-form-item v-if="editForm.configKey === 'consecutive_bonus_points'" label="">
          <div class="config-tip">
            <el-icon><InfoFilled /></el-icon>
            建议设置为日签积分的 3-5 倍，确保奖励有足够吸引力
          </div>
        </el-form-item>
        <el-form-item v-if="editForm.configKey === 'points_to_quota_rate'" label="">
          <div class="config-tip">
            <el-icon><InfoFilled /></el-icon>
            建议设置为 50-200，数值越大兑换成本越高
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmEditConfig">确认修改</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getPointsConfig, updatePointsConfig } from '@/api/gaia/checkin'
import { ElMessage, ElMessageBox } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'

const loading = ref(false)
const configData = ref([])

const editDialogVisible = ref(false)
const editFormRef = ref()
const editForm = reactive({
  configKey: '',
  configValue: 0,
  description: ''
})

const editRules = {
  configValue: [
    { required: true, message: '请输入配置值', trigger: 'blur' },
    { type: 'number', min: 0, message: '配置值必须大于等于0', trigger: 'blur' }
  ]
}

const formatConfigValue = (row) => {
  switch (row.configKey) {
    case 'daily_checkin_points':
      return `${row.configValue} 积分`
    case 'consecutive_bonus_days':
      return `${row.configValue} 天`
    case 'consecutive_bonus_points':
      return `${row.configValue} 积分`
    case 'points_to_quota_rate':
      return `${row.configValue} 积分 = 1美元`
    default:
      return row.configValue
  }
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getConfigData = async () => {
  loading.value = true
  try {
    const res = await getPointsConfig()
    if (res.code === 0) {
      configData.value = res.data || []
    } else {
      ElMessage.error(res.msg || '获取配置失败')
    }
  } catch (error) {
    console.error('获取配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const editConfig = (row) => {
  editForm.configKey = row.configKey
  editForm.configValue = row.configValue
  editForm.description = row.description
  editDialogVisible.value = true
}

const closeEditDialog = () => {
  editDialogVisible.value = false
  editFormRef.value?.resetFields()
}

const confirmEditConfig = async () => {
  try {
    await editFormRef.value.validate()
    
    // 根据配置类型进行额外验证
    let warningMessage = ''
    if (editForm.configKey === 'daily_checkin_points' && (editForm.configValue < 1 || editForm.configValue > 50)) {
      warningMessage = '每日签到积分建议设置在 1-50 之间，'
    } else if (editForm.configKey === 'consecutive_bonus_days' && (editForm.configValue < 2 || editForm.configValue > 30)) {
      warningMessage = '连续签到天数建议设置在 2-30 之间，'
    } else if (editForm.configKey === 'consecutive_bonus_points' && (editForm.configValue < 5 || editForm.configValue > 500)) {
      warningMessage = '连续签到奖励建议设置在 5-500 之间，'
    } else if (editForm.configKey === 'points_to_quota_rate' && (editForm.configValue < 10 || editForm.configValue > 1000)) {
      warningMessage = '兑换比例建议设置在 10-1000 之间，'
    }

    await ElMessageBox.confirm(
      `${warningMessage}确认将【${editForm.description}】修改为 ${editForm.configValue} 吗？`,
      '确认修改配置',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await updatePointsConfig(editForm)
    if (res.code === 0) {
      ElMessage.success('配置修改成功')
      editDialogVisible.value = false
      getConfigData()
    } else {
      ElMessage.error(res.msg || '配置修改失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('配置修改失败:', error)
      ElMessage.error('配置修改失败')
    }
  }
}

onMounted(() => {
  getConfigData()
})
</script>

<style scoped>
.config-help {
  margin-top: 30px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.config-help h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.config-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
  font-size: 14px;
  margin-top: 5px;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.el-descriptions {
  margin-top: 10px;
}
</style> 