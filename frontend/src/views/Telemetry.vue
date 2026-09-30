<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">遥测数据</span></template>
      
      <el-form :inline="true" :model="searchForm" style="margin-bottom: 16px;">
        <el-form-item label="设备">
          <el-select v-model="searchForm.deviceId" placeholder="选择设备" style="width: 200px;" @change="loadData">
            <el-option v-for="d in devices" :key="d.id" :label="d.device_name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="属性">
          <el-select v-model="searchForm.identify" placeholder="选择属性" style="width: 150px;" @change="loadChart">
            <el-option v-for="p in properties" :key="p.identifier" :label="p.name" :value="p.identifier" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始"
            end-placeholder="结束"
            style="width: 240px;"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="exportData">导出数据</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 数据表格 -->
      <el-table :data="telemetryData" stripe v-loading="loading" max-height="400">
        <el-table-column prop="timestamp" label="时间" width="180" />
        <el-table-column prop="identify" label="属性" width="120" />
        <el-table-column prop="valueFloat" label="数值" width="100">
          <template #default="{ row }">{{ row.value_float?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="valueStr" label="字符串值" />
      </el-table>
      
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @change="loadData"
        style="margin-top: 16px; justify-content: flex-end;"
      />
    </el-card>
    
    <!-- 图表展示 -->
    <el-card shadow="hover" style="margin-top: 16px;">
      <template #header>
        <span style="font-weight: 600;">数据趋势</span>
      </template>
      <div ref="chartRef" style="height: 400px;"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { getTelemetryData, getDeviceList, getPropertyList } from '@/api'

const loading = ref(false)
const devices = ref<any[]>([])
const properties = ref<any[]>([])
const telemetryData = ref<any[]>([])
const searchForm = ref({
  deviceId: '',
  identify: '',
  dateRange: null as any
})
const pagination = ref({ page: 1, size: 20, total: 0 })
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const loadData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.value.page,
      size: pagination.value.size
    }
    if (searchForm.value.deviceId) params.deviceId = searchForm.value.deviceId
    if (searchForm.value.identify) params.identify = searchForm.value.identify
    if (searchForm.value.dateRange) {
      params.startTime = searchForm.value.dateRange[0]
      params.endTime = searchForm.value.dateRange[1]
    }
    
    const res: any = await getTelemetryData(params)
    telemetryData.value = res.content || []
    pagination.value.total = res.totalElements || 0
    
    if (telemetryData.value.length > 0) {
      loadChart()
    }
  } catch (e: any) {
    ElMessage.error(e.message || '查询失败')
  } finally {
    loading.value = false
  }
}

const loadChart = () => {
  if (!chartRef.value) return
  
  if (chart) {
    chart.dispose()
  }
  
  chart = echarts.init(chartRef.value)
  
  const data = telemetryData.value.map(item => ({
    name: item.timestamp,
    value: item.value_float
  }))
  
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.name), rotate: 45 },
    yAxis: { type: 'value', name: '数值' },
    series: [{
      name: searchForm.value.identify || '数据',
      type: 'line',
      smooth: true,
      data: data.map(d => d.value),
      areaStyle: {},
      itemStyle: { color: '#667eea' }
    }]
  })
}

const exportData = () => {
  // TODO: 实现导出功能
  ElMessage.info('导出功能开发中...')
}

const loadDevices = async () => {
  try {
    const res: any = await getDeviceList({ page: 1, size: 100 })
    devices.value = res.content || []
  } catch (e) {
    console.error('Failed to load devices', e)
  }
}

const loadProperties = async () => {
  try {
    const res: any = await getPropertyList()
    properties.value = res.data || []
  } catch (e) {
    console.error('Failed to load properties', e)
  }
}

onMounted(() => {
  loadDevices()
  loadProperties()
})

onUnmounted(() => {
  chart?.dispose()
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
