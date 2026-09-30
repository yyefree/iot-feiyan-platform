<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">OTA升级管理</span>
          <el-button type="primary" @click="showFirmwareDialog = true"><el-icon><Upload /></el-icon> 上传固件</el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="固件管理" name="firmware">
          <el-table :data="firmwares" stripe v-loading="loading">
            <el-table-column prop="firmwareName" label="固件名称" width="180" />
            <el-table-column prop="version" label="版本号" width="120" />
            <el-table-column prop="productName" label="产品" width="150" />
            <el-table-column prop="fileSize" label="文件大小" width="100">
              <template #default="{ row }">{{ formatFileSize(row.fileSize) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusColor(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="releaseNotes" label="发布说明" />
            <el-table-column prop="createdAt" label="创建时间" width="170" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="handlePublish(row)" v-if="row.status === 'draft'">发布</el-button>
                <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="升级任务" name="tasks">
          <el-table :data="tasks" stripe v-loading="loading">
            <el-table-column prop="taskName" label="任务名称" width="180" />
            <el-table-column prop="firmwareName" label="固件版本" width="120" />
            <el-table-column prop="targetCount" label="目标设备数" width="100" />
            <el-table-column prop="strategy" label="升级策略" width="120">
              <template #default="{ row }">{{ row.strategy === 'full' ? '全量升级' : '灰度升级' }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getTaskStatusColor(row.status)" size="small">{{ getTaskStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="progress" label="进度" width="200">
              <template #default="{ row }">
                <el-progress :percentage="row.progress?.success / (row.progress?.success + row.progress?.failed + row.progress?.pending) * 100 || 0" :status="row.status === 'completed' ? 'success' : ''" />
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="170" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="handleViewDetail(row)">详情</el-button>
                <el-button link type="danger" @click="handleDeleteTask(row.id)" v-if="row.status === 'pending'">取消</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <!-- 上传固件对话框 -->
    <el-dialog v-model="showFirmwareDialog" title="上传固件" width="600px">
      <el-form :model="firmwareForm" label-width="100px">
        <el-form-item label="固件名称"><el-input v-model="firmwareForm.firmwareName" /></el-form-item>
        <el-form-item label="版本号"><el-input v-model="firmwareForm.version" placeholder="如: 1.0.0" /></el-form-item>
        <el-form-item label="选择产品">
          <el-select v-model="firmwareForm.productId" placeholder="选择产品" style="width: 100%;">
            <el-option label="智能灯泡" value="1" />
            <el-option label="温湿度传感器" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="固件文件">
          <el-upload action="/api/v1/ota/firmware/upload" :on-success="handleUploadSuccess" :on-error="handleUploadError">
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持 .bin, .zip 文件，最大 50MB</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="发布说明"><el-input v-model="firmwareForm.releaseNotes" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFirmwareDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateFirmware">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 创建升级任务对话框 -->
    <el-dialog v-model="showTaskDialog" title="创建升级任务" width="600px">
      <el-form :model="taskForm" label-width="100px">
        <el-form-item label="任务名称"><el-input v-model="taskForm.taskName" /></el-form-item>
        <el-form-item label="选择固件">
          <el-select v-model="taskForm.firmwareId" placeholder="选择固件" style="width: 100%;">
            <el-option v-for="f in firmwares.filter(f => f.status === 'published')" :key="f.id" :label="`${f.firmwareName} (${f.version})`" :value="f.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="升级策略">
          <el-radio-group v-model="taskForm.strategy">
            <el-radio label="full">全量升级</el-radio>
            <el-radio label="grayscale">灰度升级</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="目标设备" v-if="taskForm.strategy === 'grayscale'">
          <el-input-number v-model="taskForm.grayscalePercent" :min="1" :max="100" />%
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showTaskDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTask">创建任务</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFirmwareList, createFirmware, deleteFirmware, getTaskList, createTask, deleteTask } from '@/api'

const activeTab = ref('firmware')
const loading = ref(false)
const firmwares = ref<any[]>([])
const tasks = ref<any[]>([])
const showFirmwareDialog = ref(false)
const showTaskDialog = ref(false)

const firmwareForm = ref({
  firmwareName: '',
  version: '',
  productId: '',
  releaseNotes: ''
})

const taskForm = ref({
  taskName: '',
  firmwareId: '',
  strategy: 'full',
  grayscalePercent: 10
})

const filteredFirmwares = computed(() => firmwares.value.filter(f => f.status === 'published'))

const loadFirmwares = async () => {
  loading.value = true
  try {
    const res: any = await getFirmwareList()
    firmwares.value = res.content || []
  } finally {
    loading.value = false
  }
}

const loadTasks = async () => {
  loading.value = true
  try {
    const res: any = await getTaskList()
    tasks.value = res.content || []
  } finally {
    loading.value = false
  }
}

const handleCreateFirmware = async () => {
  try {
    await createFirmware({ ...firmwareForm.value, tenantId: 1 })
    ElMessage.success('固件上传成功')
    showFirmwareDialog.value = false
    loadFirmwares()
  } catch (e: any) {
    ElMessage.error(e.message || '上传失败')
  }
}

const handlePublish = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定发布该固件？', '提示', { type: 'info' })
    // 发布操作
    ElMessage.success('固件已发布')
    loadFirmwares()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该固件？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteFirmware(id)
      ElMessage.success('删除成功')
      loadFirmwares()
    })
    .catch(() => {})
}

const handleCreateTask = async () => {
  try {
    await createTask({ ...taskForm.value, tenantId: 1 })
    ElMessage.success('升级任务创建成功')
    showTaskDialog.value = false
    loadTasks()
  } catch (e: any) {
    ElMessage.error(e.message || '创建失败')
  }
}

const handleDeleteTask = (id: number) => {
  ElMessageBox.confirm('确定取消该升级任务？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteTask(id)
      ElMessage.success('任务已取消')
      loadTasks()
    })
    .catch(() => {})
}

const handleViewDetail = (row: any) => {
  ElMessage.info('任务详情功能开发中...')
}

const handleUploadSuccess = (response: any) => {
  firmwareForm.value.fileUrl = response.data.url
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

const formatFileSize = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    archived: '已归档'
  }
  return labels[status] || status
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    draft: 'info',
    published: 'success',
    archived: ''
  }
  return colors[status] || ''
}

const getTaskStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待执行',
    running: '执行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return labels[status] || status
}

const getTaskStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'info',
    running: 'primary',
    completed: 'success',
    failed: 'danger',
    cancelled: ''
  }
  return colors[status] || ''
}

onMounted(() => {
  loadFirmwares()
  loadTasks()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
