<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">设备列表</span>
          <el-button type="primary" @click="showDialog = true">
            <el-icon><Plus /></el-icon> 新增设备
          </el-button>
        </div>
      </template>
      <el-form :inline="true" :model="searchForm" style="margin-bottom: 16px;">
        <el-form-item label="设备名称"><el-input v-model="searchForm.deviceName" placeholder="请输入" clearable /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" clearable placeholder="全部">
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="device_name" label="设备名称" />
        <el-table-column prop="product_id" label="产品ID" width="100" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.online ? 'success' : 'danger'" size="small">
              {{ row.online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="firmware_version" label="固件版本" width="120" />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="last_online_at" label="最后在线" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleDetail(row)">详情</el-button>
            <el-button link type="primary" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @change="loadData"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>

    <!-- 新增设备对话框 -->
    <el-dialog v-model="showDialog" title="新增设备" width="480px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="设备名称"><el-input v-model="form.deviceName" /></el-form-item>
        <el-form-item label="产品ID"><el-input-number v-model="form.productId" :min="1" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDeviceList, createDevice, deleteDevice } from '@/api'

const loading = ref(false)
const tableData = ref<any[]>([])
const showDialog = ref(false)
const searchForm = ref({ deviceName: '', status: '' })
const form = ref({ deviceName: '', productId: 1, description: '' })
const pagination = ref({ page: 1, size: 20, total: 0 })

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getDeviceList({
      page: pagination.value.page,
      size: pagination.value.size,
      status: searchForm.value.status
    })
    console.log('Device list response:', res)
    // Handle both formats: { data: [...], pagination: {...} } and { content: [...], totalElements: number }
    if (res.data && Array.isArray(res.data)) {
      tableData.value = res.data
      pagination.value.total = res.pagination?.total || res.data.length
    } else if (res.content) {
      tableData.value = res.content
      pagination.value.total = res.totalElements || 0
    } else {
      tableData.value = []
      pagination.value.total = 0
    }
  } catch (error: any) {
    console.error('Failed to load devices:', error)
    ElMessage.error('加载设备列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.value = { deviceName: '', status: '' }
  pagination.value.page = 1
  loadData()
}

const handleCreate = async () => {
  await createDevice({ ...form.value, tenantId: 1 })
  ElMessage.success('设备创建成功')
  showDialog.value = false
  loadData()
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该设备？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteDevice(id)
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

const handleDetail = (row: any) => {
  ElMessage.info(`设备: ${row.device_name || row.deviceName}`)
}

onMounted(loadData)
</script>

<style scoped lang="scss">
.page-container { }
</style>
