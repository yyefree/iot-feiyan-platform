<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">项目列表</span></template>
      
      <el-button type="primary" style="margin-bottom: 16px;" @click="showDialog = true"><el-icon><Plus /></el-icon> 新建项目</el-button>
      
      <el-table :data="projects" stripe v-loading="loading">
        <el-table-column prop="projectName" label="项目名称" width="200" />
        <el-table-column prop="projectKey" label="项目Key" width="180" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEnter(row)">进入</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 新建/编辑项目对话框 -->
    <el-dialog v-model="showDialog" :title="editMode ? '编辑项目' : '新建项目'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="项目名称"><el-input v-model="form.projectName" /></el-form-item>
        <el-form-item label="项目Key">
          <el-input v-model="form.projectKey" placeholder="自动生成" :disabled="editMode" />
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
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
import { useRouter } from 'vue-router'
import { getProjectList, createProject, updateProject, deleteProject } from '@/api'

const router = useRouter()
const loading = ref(false)
const projects = ref<any[]>([])
const showDialog = ref(false)
const editMode = ref(false)
const editId = ref<number | null>(null)

const form = ref({
  projectName: '',
  projectKey: '',
  description: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getProjectList()
    projects.value = res.data || []
  } finally {
    loading.value = false
  }
}

const showAdd = () => {
  editMode.value = false
  editId.value = null
  form.value = { projectName: '', projectKey: '', description: '' }
  showDialog.value = true
}

const handleEdit = (row: any) => {
  editMode.value = true
  editId.value = row.id
  form.value = {
    projectName: row.project_name,
    projectKey: row.project_key,
    description: row.description
  }
  showDialog.value = true
}

const handleSave = async () => {
  try {
    const data = { ...form.value, tenantId: 1 }
    if (editMode.value && editId.value) {
      await updateProject(editId.value, data)
      ElMessage.success('项目更新成功')
    } else {
      await createProject(data)
      ElMessage.success('项目创建成功')
    }
    showDialog.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

const handleEnter = (row: any) => {
  ElMessage.info(`进入项目: ${row.project_name}`)
  // TODO: 跳转到项目详情页面
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该项目？项目下的产品将保留。', '提示', { type: 'warning' })
    .then(async () => {
      await deleteProject(id)
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.page-container { }
</style>
