<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">产品列表</span>
          <el-button type="primary" @click="showDialog = true"><el-icon><Plus /></el-icon> 新建产品</el-button>
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="productName" label="产品名称" />
        <el-table-column prop="productKey" label="Product Key" width="160" />
        <el-table-column prop="productType" label="设备类型" width="120">
          <template #default="{ row }">{{ row.productType === 'direct_device' ? '直连设备' : '网关子设备' }}</template>
        </el-table-column>
        <el-table-column prop="commType" label="通信方式" width="100" />
        <el-table-column prop="authType" label="认证方式" width="140" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">{{ row.status === 'online' ? '已上线' : '草稿' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push('/thing-model')">物模型</el-button>
            <el-button link type="primary" @click="handlePublish(row.id)">发布</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.size"
        :total="pagination.total" layout="total, prev, pager, next" @change="loadData"
        style="margin-top: 16px; justify-content: flex-end;" />
    </el-card>
    <el-dialog v-model="showDialog" title="新建产品" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="产品名称"><el-input v-model="form.productName" /></el-form-item>
        <el-form-item label="通信方式">
          <el-select v-model="form.commType"><el-option label="WiFi" value="wifi" /><el-option label="BLE" value="ble" /><el-option label="Zigbee" value="zigbee" /></el-select>
        </el-form-item>
        <el-form-item label="认证方式">
          <el-select v-model="form.authType"><el-option label="一设备一密" value="one_device_one_key" /><el-option label="一型一密" value="one_model_one_key" /></el-select>
        </el-form-item>
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
import { getProductList, createProduct, publishProduct, deleteProduct } from '@/api'

const loading = ref(false)
const tableData = ref<any[]>([])
const showDialog = ref(false)
const form = ref({ productName: '', commType: 'wifi', authType: 'one_device_one_key', description: '' })
const pagination = ref({ page: 1, size: 20, total: 0 })

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getProductList({ page: pagination.value.page, size: pagination.value.size })
    tableData.value = res.content || []
    pagination.value.total = res.totalElements || 0
  } finally { loading.value = false }
}

const handleCreate = async () => {
  await createProduct({ ...form.value, tenantId: 1, projectId: 1 })
  ElMessage.success('产品创建成功')
  showDialog.value = false
  loadData()
}

const handlePublish = async (id: number) => {
  await publishProduct(id)
  ElMessage.success('产品已发布')
  loadData()
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该产品？', '提示', { type: 'warning' })
    .then(async () => { await deleteProduct(id); ElMessage.success('删除成功'); loadData() })
    .catch(() => {})
}

onMounted(loadData)
</script>
