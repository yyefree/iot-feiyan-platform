<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">语音映射</span></template>
      
      <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
        语音映射用于将设备的属性或服务映射到语音助手的指令，实现语音控制设备。
      </el-alert>
      
      <el-table :data="mappings" stripe v-loading="loading">
        <el-table-column prop="voiceName" label="语音指令" width="150" />
        <el-table-column prop="identify" label="属性/服务" width="150" />
        <el-table-column prop="productName" label="产品" width="150" />
        <el-table-column prop="platform" label="语音平台" width="120">
          <template #default="{ row }">
            <el-tag :type="getPlatformColor(row.platform)" size="small">{{ getPlatformLabel(row.platform) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="params" label="参数">
          <template #default="{ row }">
            {{ JSON.stringify(row.params) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 编辑映射对话框 -->
    <el-dialog v-model="showDialog" title="编辑语音映射" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="语音指令"><el-input v-model="form.voiceName" placeholder="如: 打开灯" /></el-form-item>
        <el-form-item label="属性/服务"><el-input v-model="form.identify" placeholder="如: switch" /></el-form-item>
        <el-form-item label="语音平台">
          <el-select v-model="form.platform" style="width: 100%;">
            <el-option label="天猫精灵" value="tmall_genie" />
            <el-option label="Amazon Alexa" value="alexa" />
            <el-option label="Google Home" value="google_home" />
          </el-select>
        </el-form-item>
        <el-form-item label="参数">
          <el-input v-model="form.params" type="textarea" :rows="3" placeholder='JSON格式: {"value": true}' />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getVoiceMappingList, createVoiceMapping, updateVoiceMapping, deleteVoiceMapping } from '@/api'

const loading = ref(false)
const mappings = ref<any[]>([])
const showDialog = ref(false)
const form = ref({
  voiceName: '',
  identify: '',
  platform: 'tmall_genie',
  params: '{}'
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getVoiceMappingList()
    mappings.value = res.data || []
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  form.value = {
    voiceName: row.voice_name,
    identify: row.identify,
    platform: row.platform,
    params: typeof row.params === 'string' ? row.params : JSON.stringify(row.params)
  }
  showDialog.value = true
}

const handleSave = async () => {
  try {
    const data = { ...form.value, params: JSON.parse(form.value.params) }
    if (form.value.id) {
      await updateVoiceMapping(form.value.id, data)
    } else {
      await createVoiceMapping(data)
    }
    ElMessage.success('保存成功')
    showDialog.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该语音映射？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteVoiceMapping(id)
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    tmall_genie: '天猫精灵',
    alexa: 'Alexa',
    google_home: 'Google Home'
  }
  return labels[platform] || platform
}

const getPlatformColor = (platform: string) => {
  const colors: Record<string, string> = {
    tmall_genie: 'danger',
    alexa: 'warning',
    google_home: 'success'
  }
  return colors[platform] || ''
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.page-container { }
</style>
