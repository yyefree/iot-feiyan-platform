<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">固件管理</span>
          <el-button type="primary" @click="showDialog = true"><el-icon><Upload /></el-icon> 上传固件</el-button>
        </div>
      </template>
      
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
            <el-button link type="primary" @click="handleDownload(row)" v-if="row.status === 'published'">下载</el-button>
            <el-button link type="primary" @click="handlePublish(row)" v-if="row.status === 'draft'">发布</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 上传固件对话框 -->
    <el-dialog v-model="showDialog" title="上传固件" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="固件名称"><el-input v-model="form.firmwareName" /></el-form-item>
        <el-form-item label="版本号"><el-input v-model="form.version" placeholder="如: 1.0.0" /></el-form-item>
        <el-form-item label="选择产品">
          <el-select v-model="form.productId" placeholder="选择产品" style="width: 100%;">
            <el-option v-for="p in products" :key="p.id" :label="p.product_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="固件文件">
          <el-upload
            ref="uploadRef"
            action="/api/v1/ota/firmware/upload"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :limit="1"
            accept=".bin,.zip"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持 .bin, .zip 文件，最大 50MB</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="发布说明"><el-input v-model="form.releaseNotes" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFirmwareList, createFirmware, deleteFirmware, getProductList } from '@/api'

const loading = ref(false)
const firmwares = ref<any[]>([])
const products = ref<any[]>([])
const showDialog = ref(false)
const uploadRef = ref()

const form = ref({
  firmwareName: '',
  version: '',
  productId: '',
  fileUrl: '',
  releaseNotes: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getFirmwareList()
    firmwares.value = res.content || []
  } finally {
    loading.value = false
  }
}

const loadProducts = async () => {
  try {
    const res: any = await getProductList({ page: 1, size: 100 })
    products.value = res.content || []
  } catch (e) {
    console.error('Failed to load products', e)
  }
}

const handleCreate = async () => {
  if (!form.value.firmwareName || !form.value.version || !form.value.productId) {
    ElMessage.warning('请填写完整信息')
    return
  }
  
  try {
    await createFirmware({ ...form.value, tenantId: 1 })
    ElMessage.success('固件上传成功')
    showDialog.value = false
    form.value = { firmwareName: '', version: '', productId: '', fileUrl: '', releaseNotes: '' }
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '上传失败')
  }
}

const handlePublish = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定发布该固件？发布后将通知所有关联设备进行升级。', '提示', { type: 'info' })
    ElMessage.success('固件已发布')
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该固件？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteFirmware(id)
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

const handleDownload = (row: any) => {
  // TODO: 实现下载功能
  ElMessage.info('下载功能开发中...')
}

const handleUploadSuccess = (response: any) => {
  form.value.fileUrl = response.data?.url || ''
  ElMessage.success('文件上传成功')
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

onMounted(() => {
  loadData()
  loadProducts()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
