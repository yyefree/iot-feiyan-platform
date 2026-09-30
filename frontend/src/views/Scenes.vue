<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">场景联动</span>
          <el-button type="primary" @click="showDialog = true"><el-icon><Plus /></el-icon> 新建场景</el-button>
        </div>
      </template>
      
      <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
        场景联动支持多种触发方式：属性触发、事件触发、定时触发。
        当满足触发条件时，自动执行动作列表中的任务。
      </el-alert>
      
      <el-table :data="scenes" stripe v-loading="loading">
        <el-table-column prop="name" label="场景名称" width="200" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="triggerType" label="触发类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTriggerTypeColor(row.triggerType)" size="small">
              {{ getTriggerTypeLabel(row.triggerType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="primary" @click="handleTest(row)">测试</el-button>
            <el-button link type="primary" @click="toggleScene(row)">
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 场景编辑对话框 -->
    <el-dialog v-model="showDialog" :title="editMode ? '编辑场景' : '新建场景'" width="700px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="场景名称"><el-input v-model="form.name" placeholder="请输入场景名称" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入场景描述" /></el-form-item>
        
        <el-form-item label="触发条件">
          <el-select v-model="form.triggerType" placeholder="选择触发类型" style="width: 200px;">
            <el-option label="属性触发" value="property" />
            <el-option label="事件触发" value="event" />
            <el-option label="定时触发" value="timer" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="触发配置" v-if="form.triggerType">
          <el-input v-model="form.triggerConfig" type="textarea" :rows="3" placeholder='{"device_id":1,"property":"temperature","operator":">","value":30}' />
        </el-form-item>
        
        <el-form-item label="执行动作">
          <div style="width: 100%;">
            <el-button type="primary" size="small" @click="addAction" style="margin-bottom: 8px;">+ 添执行动作</el-button>
            <div v-for="(action, index) in form.actions" :key="index" style="display: flex; align-items: center; margin-bottom: 8px;">
              <el-input v-model="action.deviceId" placeholder="设备ID" style="width: 100px;" />
              <el-input v-model="action.identify" placeholder="属性/服务" style="width: 150px;" />
              <el-input v-model="action.value" placeholder="值" style="width: 100px;" />
              <el-button type="danger" size="small" @click="removeAction(index)">删除</el-button>
            </div>
          </div>
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
import { getSceneList, createScene, updateScene, deleteScene, enableScene, disableScene } from '@/api'

const loading = ref(false)
const scenes = ref<any[]>([])
const showDialog = ref(false)
const editMode = ref(false)
const editId = ref<number | null>(null)

const form = ref({
  name: '',
  description: '',
  triggerType: '',
  triggerConfig: '',
  actions: [] as any[]
})

const loadScenes = async () => {
  loading.value = true
  try {
    const res: any = await getSceneList({ page: 1, size: 50 })
    scenes.value = res.content || []
  } finally {
    loading.value = false
  }
}

const showAdd = () => {
  editMode.value = false
  editId.value = null
  form.value = { name: '', description: '', triggerType: '', triggerConfig: '', actions: [] }
  showDialog.value = true
}

const handleEdit = (row: any) => {
  editMode.value = true
  editId.value = row.id
  form.value = {
    name: row.name,
    description: row.description,
    triggerType: row.triggerType,
    triggerConfig: row.triggerConfig,
    actions: row.actions || []
  }
  showDialog.value = true
}

const handleTest = async (row: any) => {
  ElMessage.info('场景测试功能开发中...')
}

const toggleScene = async (row: any) => {
  try {
    if (row.status === 'active') {
      await disableScene(row.id)
      ElMessage.success('场景已禁用')
    } else {
      await enableScene(row.id)
      ElMessage.success('场景已启用')
    }
    loadScenes()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该场景？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteScene(id)
      ElMessage.success('删除成功')
      loadScenes()
    })
    .catch(() => {})
}

const handleSave = async () => {
  try {
    const data = { ...form.value, tenantId: 1 }
    if (editMode.value && editId.value) {
      await updateScene(editId.value, data)
      ElMessage.success('场景更新成功')
    } else {
      await createScene(data)
      ElMessage.success('场景创建成功')
    }
    showDialog.value = false
    loadScenes()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

const addAction = () => {
  form.value.actions.push({ deviceId: '', identify: '', value: '' })
}

const removeAction = (index: number) => {
  form.value.actions.splice(index, 1)
}

const getTriggerTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    property: '属性触发',
    event: '事件触发',
    timer: '定时触发'
  }
  return labels[type] || type
}

const getTriggerTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    property: '',
    event: 'warning',
    timer: 'info'
  }
  return colors[type] || ''
}

onMounted(loadScenes)
</script>

<style scoped lang="scss">
.page-container { }
</style>
