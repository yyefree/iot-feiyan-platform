<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">日志查询</span>
          <el-button @click="loadData"><el-icon><Refresh /></el-icon> 刷新</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" style="margin-bottom: 16px;">
        <el-form-item label="设备">
          <el-select v-model="searchForm.deviceId" placeholder="选择设备" clearable style="width: 200px;">
            <el-option v-for="d in devices" :key="d.id" :label="d.device_name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日志级别">
          <el-select v-model="searchForm.level" clearable>
            <el-option label="INFO" value="INFO" />
            <el-option label="WARN" value="WARN" />
            <el-option label="ERROR" value="ERROR" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="logs" stripe v-loading="loading" size="small">
        <el-table-column prop="timestamp" label="时间" width="180" />
        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelColor(row.level)" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deviceName" label="设备" width="150" />
        <el-table-column prop="message" label="消息" />
        <el-table-column prop="module" label="模块" width="100" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @change="loadData"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
    
    <!-- 日志详情对话框 -->
    <el-dialog v-model="showDetailDialog" title="日志详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="时间">{{ currentLog?.timestamp }}</el-descriptions-item>
        <el-descriptions-item label="级别">
          <el-tag :type="getLevelColor(currentLog?.level)" size="small">{{ currentLog?.level }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="设备">{{ currentLog?.device_name }}</el-descriptions-item>
        <el-descriptions-item label="模块">{{ currentLog?.module }}</el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <div style="white-space: pre-wrap; font-family: monospace; font-size: 12px; background: #f5f5f5; padding: 12px; border-radius: 4px;">
        {{ currentLog?.message }}
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getLogList, getDeviceList } from '@/api'

const loading = ref(false)
const logs = ref<any[]>([])
const devices = ref<any[]>([])
const showDetailDialog = ref(false)
const currentLog = ref<any>(null)
const searchForm = ref({
  deviceId: '',
  level: '',
  dateRange: null as any
})
const pagination = ref({ page: 1, size: 20, total: 0 })

const loadData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.value.page,
      size: pagination.value.size
    }
    if (searchForm.value.deviceId) params.deviceId = searchForm.value.deviceId
    if (searchForm.value.level) params.level = searchForm.value.level
    if (searchForm.value.dateRange) {
      params.startTime = searchForm.value.dateRange[0]
      params.endTime = searchForm.value.dateRange[1]
    }
    
    const res: any = await getLogList(params)
    logs.value = res.data || []
    pagination.value.total = res.pagination?.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '查询失败')
  } finally {
    loading.value = false
  }
}

const loadDevices = async () => {
  try {
    const res: any = await getDeviceList({ page: 1, size: 100 })
    devices.value = res.data || []
  } catch (e) {
    console.error('Failed to load devices', e)
  }
}

const handleViewDetail = (row: any) => {
  currentLog.value = row
  showDetailDialog.value = true
}

const getLevelColor = (level: string) => {
  const colors: Record<string, string> = {
    INFO: '',
    WARN: 'warning',
    ERROR: 'danger'
  }
  return colors[level] || ''
}

onMounted(() => {
  loadData()
  loadDevices()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
