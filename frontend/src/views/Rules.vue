<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">规则引擎</span></template>
      <el-button type="primary" style="margin-bottom: 16px;" @click="showDialog = true"><el-icon><Plus /></el-icon> 新建规则</el-button>
      <el-table :data="rules" stripe v-loading="loading">
        <el-table-column prop="name" label="规则名称" />
        <el-table-column prop="ruleType" label="规则类型" width="120">
          <template #default="{ row }">{{ row.ruleType === 'scene' ? '场景联动' : '数据流转' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status === 'active' ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button link type="primary" @click="toggleRule(row)">
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog v-model="showDialog" title="新建规则" width="520px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="规则名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="规则类型">
          <el-select v-model="form.ruleType"><el-option label="场景联动" value="scene" /><el-option label="数据流转" value="stream" /></el-select>
        </el-form-item>
        <el-form-item label="触发条件"><el-input v-model="form.triggerConfig" placeholder='{"topic":"device/+/post"}' /></el-form-item>
        <el-form-item label="执行动作"><el-input v-model="form.actionConfig" placeholder='{"action":"publish","topic":"scene/execute"}' /></el-form-item>
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
import { getRuleList, createRule, deleteRule, enableRule, disableRule } from '@/api'

const loading = ref(false)
const rules = ref<any[]>([])
const showDialog = ref(false)
const form = ref({ name: '', ruleType: 'scene', triggerConfig: '', actionConfig: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getRuleList({ page: 1, size: 50 })
    rules.value = res.data || []
  } finally { loading.value = false }
}

const handleCreate = async () => {
  await createRule({ ...form.value, tenantId: 1 })
  ElMessage.success('规则创建成功')
  showDialog.value = false
  loadData()
}

const toggleRule = async (row: any) => {
  const fn = row.status === 'active' ? disableRule : enableRule
  await fn(row.id)
  loadData()
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该规则？', '提示', { type: 'warning' })
    .then(async () => { await deleteRule(id); ElMessage.success('删除成功'); loadData() })
    .catch(() => {})
}

onMounted(loadData)
</script>
