<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">数据分析</span></template>
      <el-form :inline="true" style="margin-bottom: 16px;">
        <el-form-item label="设备"><el-select v-model="selectedDevice" placeholder="选择设备" style="width: 200px;">
          <el-option label="温度传感器-001" value="1" /><el-option label="智能灯泡-002" value="2" />
        </el-select></el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker v-model="dateRange" type="datetimerange" range-separator="至" start-placeholder="开始" end-placeholder="结束" />
        </el-form-item>
        <el-form-item><el-button type="primary" @click="loadChart">查询</el-button></el-form-item>
      </el-form>
      <div ref="chartRef" style="height: 400px;"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
const selectedDevice = ref('1')
const dateRange = ref<[string, string]>()
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const loadChart = () => {
  if (!chart) chart = echarts.init(chartRef.value!)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00','04:00','08:00','12:00','16:00','20:00','24:00'] },
    yAxis: { type: 'value', name: '温度(℃)' },
    series: [{ name: '温度', type: 'line', smooth: true, data: [22,21,23,26,25,23,22], areaStyle: {}, itemStyle: { color: '#667eea' } }]
  })
}

onMounted(loadChart)
onUnmounted(() => chart?.dispose())
</script>
