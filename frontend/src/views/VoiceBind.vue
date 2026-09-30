<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">语音平台绑定</span>
          <el-button type="primary" @click="showBindDialog = true"><el-icon><Link /></el-icon> 绑定语音平台</el-button>
        </div>
      </template>
      
      <el-table :data="bindings" stripe v-loading="loading">
        <el-table-column prop="productName" label="产品名称" width="180" />
        <el-table-column prop="platform" label="语音平台" width="150">
          <template #default="{ row }">
            <el-tag :type="getPlatformColor(row.platform)" size="small">
              {{ getPlatformLabel(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="skillId" label="技能ID" width="150" />
        <el-table-column prop="bindStatus" label="绑定状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.bindStatus)" size="small">
              {{ getStatusLabel(row.bindStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="绑定时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleUnbind(row)" v-if="row.bindStatus === 'bound'">解绑</el-button>
            <el-button link type="primary" @click="handleRefresh(row)" v-if="row.bindStatus === 'pending'">刷新</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 绑定对话框 -->
    <el-dialog v-model="showBindDialog" title="绑定语音平台" width="500px">
      <el-form :model="bindForm" label-width="100px">
        <el-form-item label="选择产品">
          <el-select v-model="bindForm.productId" placeholder="选择产品" style="width: 100%;">
            <el-option v-for="p in products" :key="p.id" :label="p.productName" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="语音平台">
          <el-select v-model="bindForm.platform" placeholder="选择语音平台" style="width: 100%;">
            <el-option label="天猫精灵" value="tmall_genie" />
            <el-option label="Amazon Alexa" value="alexa" />
            <el-option label="Google Home" value="google_home" />
            <el-option label="小爱同学" value="xiaomi" />
          </el-select>
        </el-form-item>
        <el-form-item label="技能ID" v-if="bindForm.platform">
          <el-input v-model="bindForm.skillId" placeholder="请输入技能ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBindDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBind">确认绑定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getVoiceBindList, bindVoicePlatform, unbindVoicePlatform, deleteVoiceBind, getProductList } from '@/api'

const loading = ref(false)
const bindings = ref<any[]>([])
const products = ref<any[]>([])
const showBindDialog = ref(false)

const bindForm = ref({
  productId: '',
  platform: '',
  skillId: ''
})

const loadBindings = async () => {
  loading.value = true
  try {
    const res: any = await getVoiceBindList()
    bindings.value = res.content || []
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

const handleBind = async () => {
  try {
    await bindVoicePlatform({ ...bindForm.value, tenantId: 1 })
    ElMessage.success('绑定成功')
    showBindDialog.value = false
    bindForm.value = { productId: '', platform: '', skillId: '' }
    loadBindings()
  } catch (e: any) {
    ElMessage.error(e.message || '绑定失败')
  }
}

const handleUnbind = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定解绑该语音平台？', '提示', { type: 'warning' })
    await unbindVoicePlatform(row.id)
    ElMessage.success('解绑成功')
    loadBindings()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e.message || '解绑失败')
  }
}

const handleRefresh = async (row: any) => {
  ElMessage.info('刷新绑定状态...')
  loadBindings()
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该绑定记录？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteVoiceBind(id)
      ElMessage.success('删除成功')
      loadBindings()
    })
    .catch(() => {})
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    tmall_genie: '天猫精灵',
    alexa: 'Amazon Alexa',
    google_home: 'Google Home',
    xiaomi: '小爱同学'
  }
  return labels[platform] || platform
}

const getPlatformColor = (platform: string) => {
  const colors: Record<string, string> = {
    tmall_genie: 'danger',
    alexa: 'warning',
    google_home: 'success',
    xiaomi: 'primary'
  }
  return colors[platform] || ''
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待激活',
    bound: '已绑定',
    failed: '绑定失败'
  }
  return labels[status] || status
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'warning',
    bound: 'success',
    failed: 'danger'
  }
  return colors[status] || ''
}

onMounted(() => {
  loadBindings()
  loadProducts()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
