<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">设备分组</span>
          <el-button type="primary" @click="showDialog = true"><el-icon><Plus /></el-icon> 新建分组</el-button>
        </div>
      </template>
      
      <el-table :data="groups" stripe v-loading="loading">
        <el-table-column prop="name" label="分组名称" width="200" />
        <el-table-column prop="parentName" label="父分组" width="150">
          <template #default="{ row }">
            {{ row.parent_id ? parentNames[row.parent_id] || '未知' : '无' }}
          </template>
        </el-table-column>
        <el-table-column prop="deviceCount" label="设备数量" width="100" />
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleAddDevice(row)">添加设备</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 新建/编辑分组对话框 -->
    <el-dialog v-model="showDialog" :title="editMode ? '编辑分组' : '新建分组'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="分组名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="父分组">
          <el-cascader
            v-model="form.parentId"
            :options="groupOptions"
            :props="{ checkStrictly: true, value: 'id', label: 'name' }"
            placeholder="选择父分组（可选）"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 添加设备对话框 -->
    <el-dialog v-model="showDeviceDialog" title="添加设备" width="600px">
      <el-table :data="devices" stripe v-loading="deviceLoading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="deviceName" label="设备名称" />
        <el-table-column prop="productKey" label="产品Key" width="150" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.online ? 'success' : 'danger'" size="small">
              {{ row.online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showDeviceDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddDevices">确认添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getGroupList, createGroup, updateGroup, deleteGroup, addDeviceToGroup, getDeviceList } from '@/api'

const loading = ref(false)
const groups = ref<any[]>([])
const devices = ref<any[]>([])
const showDialog = ref(false)
const showDeviceDialog = ref(false)
const editMode = ref(false)
const editId = ref<number | null>(null)
const selectedDevices = ref<any[]>([])
const currentGroupId = ref<number | null>(null)
const deviceLoading = ref(false)

const form = ref({
  name: '',
  parentId: [] as number[]
})

const parentNames = ref<Record<number, string>>({})

const groupOptions = computed(() => {
  return groups.value.map(g => ({
    id: g.id,
    name: g.name,
    value: g.id,
    label: g.name
  }))
})

const loadGroups = async () => {
  loading.value = true
  try {
    const res: any = await getGroupList()
    groups.value = res.data || []
    // 构建父分组名称映射
    groups.value.forEach(g => {
      if (g.parent_id) {
        const parent = groups.value.find(p => p.id === g.parent_id)
        if (parent) {
          parentNames.value[g.parent_id] = parent.name
        }
      }
    })
  } finally {
    loading.value = false
  }
}

const showAdd = () => {
  editMode.value = false
  editId.value = null
  form.value = { name: '', parentId: [] }
  showDialog.value = true
}

const handleEdit = (row: any) => {
  editMode.value = true
  editId.value = row.id
  form.value = {
    name: row.name,
    parentId: row.parent_id ? [row.parent_id] : []
  }
  showDialog.value = true
}

const handleSave = async () => {
  try {
    const data = {
      ...form.value,
      parentId: form.value.parentId.length > 0 ? form.value.parentId[0] : null,
      tenantId: 1
    }
    if (editMode.value && editId.value) {
      await updateGroup(editId.value, data)
      ElMessage.success('分组更新成功')
    } else {
      await createGroup(data)
      ElMessage.success('分组创建成功')
    }
    showDialog.value = false
    loadGroups()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该分组？分组下的设备不会被删除。', '提示', { type: 'warning' })
    .then(async () => {
      await deleteGroup(id)
      ElMessage.success('删除成功')
      loadGroups()
    })
    .catch(() => {})
}

const handleAddDevice = (row: any) => {
  currentGroupId.value = row.id
  showDeviceDialog.value = true
  loadDevices()
}

const loadDevices = async () => {
  deviceLoading.value = true
  try {
    const res: any = await getDeviceList({ page: 1, size: 100 })
    devices.value = res.data || []
  } finally {
    deviceLoading.value = false
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedDevices.value = selection
}

const handleAddDevices = async () => {
  if (!currentGroupId.value || selectedDevices.value.length === 0) {
    ElMessage.warning('请选择要添加的设备')
    return
  }
  
  try {
    for (const device of selectedDevices.value) {
      await addDeviceToGroup(currentGroupId.value, device.id)
    }
    ElMessage.success(`成功添加 ${selectedDevices.value.length} 台设备`)
    showDeviceDialog.value = false
    selectedDevices.value = []
    loadGroups()
  } catch (e: any) {
    ElMessage.error(e.message || '添加失败')
  }
}

onMounted(loadGroups)
</script>

<style scoped lang="scss">
.page-container { }
</style>
