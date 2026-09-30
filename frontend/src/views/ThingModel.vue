<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">物模型定义</span>
          <el-button type="primary" @click="handlePublish">
            <el-icon><Upload /></el-icon> 发布物模型
          </el-button>
        </div>
      </template>
      
      <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
        物模型定义了产品的数据格式，包含属性（Property）、服务（Service）、事件（Event）三种类型。
        属性用于描述设备状态，服务用于设备控制，事件用于设备上报。
      </el-alert>
      
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <!-- 属性定义 -->
        <el-tab-pane label="属性" name="property">
          <div style="margin-bottom: 12px;">
            <el-button type="primary" size="small" @click="showAddProperty">+ 新增属性</el-button>
            <el-button size="small" @click="showImportProperty">导入属性</el-button>
          </div>
          <el-table :data="properties" stripe size="small" v-loading="loading">
            <el-table-column prop="identifier" label="标识符" width="150" />
            <el-table-column prop="name" label="名称" width="120" />
            <el-table-column prop="type" label="类型" width="80">
              <template #default="{ row }">
                <el-tag size="small" :type="getTypeColor(row.type)">
                  {{ row.type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="accessMode" label="访问权限" width="100">
              <template #default="{ row }">
                {{ row.accessMode === 'read_write' ? '读写' : row.accessMode === 'read' ? '只读' : '只写' }}
              </template>
            </el-table-column>
            <el-table-column prop="unit" label="单位" width="80" />
            <el-table-column prop="desc" label="描述" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="editProperty(row)">编辑</el-button>
                <el-button link type="danger" size="small" @click="deleteProperty(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <!-- 事件定义 -->
        <el-tab-pane label="事件" name="event">
          <div style="margin-bottom: 12px;">
            <el-button type="primary" size="small" @click="showAddEvent">+ 新增事件</el-button>
          </div>
          <el-table :data="events" stripe size="small" v-loading="loading">
            <el-table-column prop="identifier" label="标识符" width="150" />
            <el-table-column prop="name" label="名称" width="120" />
            <el-table-column prop="level" label="级别" width="80">
              <template #default="{ row }">
                <el-tag :type="getEventLevelColor(row.level)" size="small">
                  {{ row.level === 'info' ? '信息' : row.level === 'warn' ? '警告' : '错误' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="desc" label="描述" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="editEvent(row)">编辑</el-button>
                <el-button link type="danger" size="small" @click="deleteEvent(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <!-- 服务定义 -->
        <el-tab-pane label="服务" name="service">
          <div style="margin-bottom: 12px;">
            <el-button type="primary" size="small" @click="showAddService">+ 新增服务</el-button>
          </div>
          <el-table :data="services" stripe size="small" v-loading="loading">
            <el-table-column prop="identifier" label="标识符" width="150" />
            <el-table-column prop="name" label="名称" width="120" />
            <el-table-column prop="desc" label="描述" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="editService(row)">编辑</el-button>
                <el-button link type="danger" size="small" @click="deleteService(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <!-- 属性编辑对话框 -->
    <el-dialog v-model="propertyDialog.visible" :title="propertyDialog.isEdit ? '编辑属性' : '新增属性'" width="600px">
      <el-form :model="propertyDialog.form" label-width="100px">
        <el-form-item label="标识符"><el-input v-model="propertyDialog.form.identifier" placeholder="如: switch, temperature" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="propertyDialog.form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="propertyDialog.form.type">
            <el-option label="字符串" value="string" />
            <el-option label="整数" value="int" />
            <el-option label="浮点数" value="float" />
            <el-option label="布尔" value="bool" />
            <el-option label="枚举" value="enums" />
          </el-select>
        </el-form-item>
        <el-form-item label="访问权限">
          <el-select v-model="propertyDialog.form.accessMode">
            <el-option label="读写" value="read_write" />
            <el-option label="只读" value="read" />
            <el-option label="只写" value="write" />
          </el-select>
        </el-form-item>
        <el-form-item label="单位"><el-input v-model="propertyDialog.form.unit" placeholder="如: ℃, %"</el-input></el-form-item>
        <el-form-item label="描述"><el-input v-model="propertyDialog.form.desc" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="默认值"><el-input v-model="propertyDialog.form.defaultValue" /></el-form-item>
        <el-form-item label="最小值" v-if="propertyDialog.form.type === 'int' || propertyDialog.form.type === 'float'">
          <el-input-number v-model="propertyDialog.form.minValue" :precision="2" />
        </el-form-item>
        <el-form-item label="最大值" v-if="propertyDialog.form.type === 'int' || propertyDialog.form.type === 'float'">
          <el-input-number v-model="propertyDialog.form.maxValue" :precision="2" />
        </el-form-item>
        <el-form-item label="步进值" v-if="propertyDialog.form.type === 'int' || propertyDialog.form.type === 'float'">
          <el-input-number v-model="propertyDialog.form.step" :precision="2" />
        </el-form-item>
        <el-form-item label="枚举值" v-if="propertyDialog.form.type === 'enums'">
          <el-input v-model="propertyDialog.form.enumValues" placeholder="JSON格式: [{label:'开',value:0},{label:'关',value:1}]" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="propertyDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveProperty">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 事件编辑对话框 -->
    <el-dialog v-model="eventDialog.visible" :title="eventDialog.isEdit ? '编辑事件' : '新增事件'" width="600px">
      <el-form :model="eventDialog.form" label-width="100px">
        <el-form-item label="标识符"><el-input v-model="eventDialog.form.identifier" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="eventDialog.form.name" /></el-form-item>
        <el-form-item label="级别">
          <el-select v-model="eventDialog.form.level">
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warn" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="eventDialog.form.desc" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="eventDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveEvent">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 服务编辑对话框 -->
    <el-dialog v-model="serviceDialog.visible" :title="serviceDialog.isEdit ? '编辑服务' : '新增服务'" width="600px">
      <el-form :model="serviceDialog.form" label-width="100px">
        <el-form-item label="标识符"><el-input v-model="serviceDialog.form.identifier" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="serviceDialog.form.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="serviceDialog.form.desc" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="serviceDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveService">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPropertyList, addProperty, updateProperty, deleteProperty, getEventList, addEvent, updateEvent, deleteEvent, getServiceList, addService, updateService, deleteService, publishThingModel } from '@/api'

const activeTab = ref('property')
const loading = ref(false)
const properties = ref<any[]>([])
const events = ref<any[]>([])
const services = ref<any[]>([])

// 属性对话框
const propertyDialog = ref({
  visible: false,
  isEdit: false,
  form: {
    identifier: '',
    name: '',
    type: 'float',
    accessMode: 'read_write',
    unit: '',
    desc: '',
    defaultValue: '',
    minValue: 0,
    maxValue: 100,
    step: 0.1,
    enumValues: ''
  }
})

// 事件对话框
const eventDialog = ref({
  visible: false,
  isEdit: false,
  form: {
    identifier: '',
    name: '',
    level: 'info',
    desc: ''
  }
})

// 服务对话框
const serviceDialog = ref({
  visible: false,
  isEdit: false,
  form: {
    identifier: '',
    name: '',
    desc: ''
  }
})

const handleTabClick = (tab: any) => {
  const name = tab.props.name
  if (name === 'property') loadProperties()
  if (name === 'event') loadEvents()
  if (name === 'service') loadServices()
}

const loadProperties = async () => {
  loading.value = true
  try {
    const res: any = await getPropertyList()
    properties.value = res.data || []
  } finally {
    loading.value = false
  }
}

const loadEvents = async () => {
  loading.value = true
  try {
    const res: any = await getEventList()
    events.value = res.data || []
  } finally {
    loading.value = false
  }
}

const loadServices = async () => {
  loading.value = true
  try {
    const res: any = await getServiceList()
    services.value = res.data || []
  } finally {
    loading.value = false
  }
}

const showAddProperty = () => {
  propertyDialog.value = { visible: true, isEdit: false, form: { identifier: '', name: '', type: 'float', accessMode: 'read_write', unit: '', desc: '', defaultValue: '', minValue: 0, maxValue: 100, step: 0.1, enumValues: '' } }
}

const editProperty = (row: any) => {
  propertyDialog.value = { visible: true, isEdit: true, form: { ...row } }
}

const saveProperty = async () => {
  try {
    if (propertyDialog.value.isEdit) {
      await updateProperty(propertyDialog.value.form.id, propertyDialog.value.form)
      ElMessage.success('属性更新成功')
    } else {
      await addProperty(propertyDialog.value.form)
      ElMessage.success('属性添加成功')
    }
    propertyDialog.value.visible = false
    loadProperties()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const deleteProperty = (id: number) => {
  ElMessageBox.confirm('确定删除该属性？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteProperty(id)
      ElMessage.success('删除成功')
      loadProperties()
    })
    .catch(() => {})
}

const showAddEvent = () => {
  eventDialog.value = { visible: true, isEdit: false, form: { identifier: '', name: '', level: 'info', desc: '' } }
}

const editEvent = (row: any) => {
  eventDialog.value = { visible: true, isEdit: true, form: { ...row } }
}

const saveEvent = async () => {
  try {
    if (eventDialog.value.isEdit) {
      await updateEvent(eventDialog.value.form.id, eventDialog.value.form)
      ElMessage.success('事件更新成功')
    } else {
      await addEvent(eventDialog.value.form)
      ElMessage.success('事件添加成功')
    }
    eventDialog.value.visible = false
    loadEvents()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const deleteEvent = (id: number) => {
  ElMessageBox.confirm('确定删除该事件？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteEvent(id)
      ElMessage.success('删除成功')
      loadEvents()
    })
    .catch(() => {})
}

const showAddService = () => {
  serviceDialog.value = { visible: true, isEdit: false, form: { identifier: '', name: '', desc: '' } }
}

const editService = (row: any) => {
  serviceDialog.value = { visible: true, isEdit: true, form: { ...row } }
}

const saveService = async () => {
  try {
    if (serviceDialog.value.isEdit) {
      await updateService(serviceDialog.value.form.id, serviceDialog.value.form)
      ElMessage.success('服务更新成功')
    } else {
      await addService(serviceDialog.value.form)
      ElMessage.success('服务添加成功')
    }
    serviceDialog.value.visible = false
    loadServices()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const deleteService = (id: number) => {
  ElMessageBox.confirm('确定删除该服务？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteService(id)
      ElMessage.success('删除成功')
      loadServices()
    })
    .catch(() => {})
}

const showImportProperty = () => {
  ElMessage.info('导入功能开发中...')
}

const handlePublish = async () => {
  try {
    await publishThingModel()
    ElMessage.success('物模型发布成功')
  } catch (e: any) {
    ElMessage.error(e.message || '发布失败')
  }
}

const getTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    string: '',
    int: 'success',
    float: 'warning',
    bool: 'info',
    enums: 'danger'
  }
  return colors[type] || ''
}

const getEventLevelColor = (level: string) => {
  const colors: Record<string, string> = {
    info: '',
    warn: 'warning',
    error: 'danger'
  }
  return colors[level] || ''
}

onMounted(() => {
  loadProperties()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
