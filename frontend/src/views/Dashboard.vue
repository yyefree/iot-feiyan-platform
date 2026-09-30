<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <el-icon size="28"><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.totalCount || 0 }}</div>
            <div class="stat-label">设备总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);">
            <el-icon size="28"><Connected /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ statistics.onlineCount || 0 }}</div>
            <div class="stat-label">在线设备</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);">
            <el-icon size="28"><Box /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ productCount }}</div>
            <div class="stat-label">产品数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <el-icon size="28"><Operation /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ ruleCount }}</div>
            <div class="stat-label">规则数量</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <span style="font-weight: 600;">设备在线趋势</span>
          </template>
          <div ref="chartRef" style="height: 320px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <span style="font-weight: 600;">设备状态分布</span>
          </template>
          <div ref="pieRef" style="height: 320px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <span style="font-weight: 600;">最近设备动态</span>
            <el-button type="primary" size="small" style="float: right;" @click="$router.push('/devices')">
              查看全部
            </el-button>
          </template>
          <el-table :data="recentDevices" style="width: 100%">
            <el-table-column prop="deviceName" label="设备名称" />
            <el-table-column prop="productKey" label="产品Key" width="180" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ row.status === 'online' ? '在线' : '离线' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="firmwareVersion" label="固件版本" width="120" />
            <el-table-column prop="lastOnlineAt" label="最后上线时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getDeviceStatistics, getDeviceList } from '@/api'

const statistics = ref({ onlineCount: 0, totalCount: 0, offlineCount: 0 })
const productCount = ref(0)
const ruleCount = ref(0)
const recentDevices = ref<any[]>([])
const chartRef = ref<HTMLElement>()
const pieRef = ref<HTMLElement>()
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

onMounted(async () => {
  const stats = await getDeviceStatistics()
  statistics.value = stats

  const devices = await getDeviceList({ page: 1, size: 8, status: 'online' })
  recentDevices.value = devices.content?.slice(0, 8) || []

  // 设备在线趋势图
  lineChart = echarts.init(chartRef.value!)
  lineChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['在线', '离线'] },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [
      { name: '在线', type: 'line', smooth: true, data: [12, 15, 28, 45, 38, 22, 18], areaStyle: {}, itemStyle: { color: '#667eea' } },
      { name: '离线', type: 'line', smooth: true, data: [5, 3, 2, 1, 3, 6, 4], areaStyle: {}, itemStyle: { color: '#f5576c' } }
    ]
  })

  // 设备状态饼图
  pieChart = echarts.init(pieRef.value!)
  pieChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      data: [
        { value: stats.onlineCount || 0, name: '在线', itemStyle: { color: '#38ef7d' } },
        { value: stats.offlineCount || 0, name: '离线', itemStyle: { color: '#f45c43' } }
      ]
    }]
  })

  window.addEventListener('resize', () => {
    lineChart?.resize()
    pieChart?.resize()
  })
})

onUnmounted(() => {
  lineChart?.dispose()
  pieChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard { padding: 0; }
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  .stat-icon {
    width: 56px; height: 56px; border-radius: 12px;
    display: flex; align-items: center; justify-content: center;
    color: #fff; flex-shrink: 0;
  }
  .stat-info { flex: 1; }
  .stat-value { font-size: 28px; font-weight: 700; color: #1a1a1a; }
  .stat-label { font-size: 13px; color: #888; margin-top: 4px; }
}
</style>
