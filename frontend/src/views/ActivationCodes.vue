<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: 600;">激活码管理</span>
          <el-button type="primary" @click="showDialog = true"><el-icon><Plus /></el-icon> 生成激活码</el-button>
        </div>
      </template>
      
      <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
        激活码用于设备量产烧录，支持一型一密和一机一密两种模式。一型一密为同类产品生成相同激活码，一机一密为每台设备生成唯一激活码。
      </el-alert>
      
      <el-form :inline="true" :model="searchForm" style="margin-bottom: 16px;">
        <el-form-item label="产品">
          <el-select v-model="searchForm.productId" placeholder="选择产品" style="width: 200px;">
            <el-option v-for="p in products" :key="p.id" :label="p.product_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="未使用" value="unused" />
            <el-option label="已使用" value="used" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="tableData" stripe v-loading="loading">
        <el-table-column prop="code" label="激活码" width="200">
          <template #default="{ row }">
            <span style="font-family: monospace;">{{ row.code }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="产品" width="150" />
        <el-table-column prop="deviceName" label="设备名称" width="150">
          <template #default="{ row }">
            {{ row.device_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'used' ? 'success' : 'info'" size="small">
              {{ row.status === 'used' ? '已使用' : '未使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="usedBy" label="使用者" width="120">
          <template #default="{ row }">
            {{ row.used_by || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="usedAt" label="使用时间" width="170">
          <template #default="{ row }">
            {{ row.used_at || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleCopy(row.code)">复制</el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @change="loadData"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
    
    <!-- 生成激活码对话框 -->
    <el-dialog v-model="showDialog" title="生成激活码" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="选择产品">
          <el-select v-model="form.productId" placeholder="选择产品" style="width: 100%;">
            <el-option v-for="p in products" :key="p.id" :label="p.product_name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="生成模式">
          <el-radio-group v-model="form.mode">
            <el-radio label="one_model">一型一密</el-radio>
            <el-radio label="one_device">一机一密</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="生成数量" v-if="form.mode === 'one_model'">
          <el-input-number v-model="form.count" :min="1" :max="10000" />
        </el-form-item>
        <el-form-item label="设备前缀" v-if="form.mode === 'one_device'">
          <el-input v-model="form.prefix" placeholder="如: device" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleGenerate">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getActivationCodeList, generateActivationCodes, deleteActivationCode, getProductList } from '@/api'

const loading = ref(false)
const tableData = ref<any[]>([])
const products = ref<any[]>([])
const showDialog = ref(false)
const searchForm = ref({ productId: '', status: '' })
const form = ref({
  productId: '',
  mode: 'one_model',
  count: 10,
  prefix: 'device'
})
const pagination = ref({ page: 1, size: 20, total: 0 })

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getActivationCodeList({
      page: pagination.value.page,
      size: pagination.value.size,
      productId: searchForm.value.productId,
      status: searchForm.value.status
    })
    tableData.value = res.content || []
    pagination.value.total = res.totalElements || 0
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.value = { productId: '', status: '' }
  pagination.value.page = 1
  loadData()
}

const loadProducts = async () => {
  try {
    const res: any = await getProductList({ page: 1, size: 100 })
    products.value = res.content || []
  } catch (e) {
    console.error('Failed to load products', e)
  }
}

const handleGenerate = async () => {
  if (!form.value.productId) {
    ElMessage.warning('请选择产品')
    return
  }
  
  try {
    await generateActivationCodes({
      productId: form.value.productId,
      mode: form.value.mode,
      count: form.value.mode === 'one_model' ? form.value.count : form.value.count,
      prefix: form.value.prefix
    })
    ElMessage.success('激活码生成成功')
    showDialog.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '生成失败')
  }
}

const handleCopy = (code: string) => {
  navigator.clipboard.writeText(code).then(() => {
    ElMessage.success('已复制到剪贴板')
  })
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该激活码？', '提示', { type: 'warning' })
    .then(async () => {
      await deleteActivationCode(id)
      ElMessage.success('删除成功')
      loadData()
    })
    .catch(() => {})
}

onMounted(() => {
  loadData()
  loadProducts()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
