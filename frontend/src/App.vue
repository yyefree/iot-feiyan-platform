<template>
  <el-container style="height: 100vh;">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar" style="background: #001529;">
      <div class="logo">
        <span v-if="!isCollapse">IoT 飞燕平台</span>
        <span v-else>IoT</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :router="true"
        background-color="#001529"
        text-color="#ffffffb3"
        active-text-color="#ffffff"
        style="border: none;"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>工作台</template>
        </el-menu-item>
        <el-sub-menu index="project">
          <template #title>
            <el-icon><Folder /></el-icon>
            <span>项目管理</span>
          </template>
          <el-menu-item index="/projects">项目列表</el-menu-item>
          <el-menu-item index="/categories">品类管理</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="product">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>产品管理</span>
          </template>
          <el-menu-item index="/products">产品列表</el-menu-item>
          <el-menu-item index="/thing-model">物模型</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="device">
          <template #title>
            <el-icon><Monitor /></el-icon>
            <span>设备管理</span>
          </template>
          <el-menu-item index="/devices">设备列表</el-menu-item>
          <el-menu-item index="/device-groups">设备分组</el-menu-item>
          <el-menu-item index="/activation-codes">激活码管理</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="scene">
          <template #title>
            <el-icon><Operation /></el-icon>
            <span>场景联动</span>
          </template>
          <el-menu-item index="/scenes">场景列表</el-menu-item>
          <el-menu-item index="/rules">规则引擎</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="data">
          <template #title>
            <el-icon><TrendCharts /></el-icon>
            <span>数据管理</span>
          </template>
          <el-menu-item index="/data-analysis">数据分析</el-menu-item>
          <el-menu-item index="/telemetry">遥测数据</el-menu-item>
          <el-menu-item index="/logs">日志查询</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="ota">
          <template #title>
            <el-icon><Upload /></el-icon>
            <span>OTA升级</span>
          </template>
          <el-menu-item index="/ota-firmware">固件管理</el-menu-item>
          <el-menu-item index="/ota-tasks">升级任务</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="voice">
          <template #title>
            <el-icon><Microphone /></el-icon>
            <span>语音控制</span>
          </template>
          <el-menu-item index="/voice-bind">语音平台绑定</el-menu-item>
          <el-menu-item index="/voice-mapping">语音映射</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- 顶栏 -->
      <el-header class="header" style="background: #fff; box-shadow: 0 1px 4px rgba(0,0,0,.08);">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="isCollapse = !isCollapse">
            <component :is="isCollapse ? 'Expand' : 'Fold'" />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-for="item in breadcrumbs" :key="item">{{ item }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-badge :value="3" class="header-icon">
            <el-icon size="18"><Bell /></el-icon>
          </el-badge>
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" icon="User" />
              <span class="user-name">管理员</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item>账号设置</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isCollapse = ref(false)
const activeMenu = computed(() => route.path)
const breadcrumbs = computed(() => {
  const map: Record<string, string> = {
    '/devices': '设备管理',
    '/products': '产品管理',
    '/rules': '规则引擎',
    '/data-analysis': '数据分析',
  }
  return [map[route.path] || route.path.replace('/', '')]
})
</script>

<style lang="scss">
.sidebar {
  transition: width 0.3s;
  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 18px;
    font-weight: 600;
    border-bottom: 1px solid rgba(255,255,255,0.1);
  }
  .el-menu {
    border-right: none;
  }
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  .header-left { display: flex; align-items: center; gap: 16px; }
  .header-right { display: flex; align-items: center; gap: 20px; }
  .collapse-btn { font-size: 20px; cursor: pointer; }
  .user-info { display: flex; align-items: center; gap: 8px; cursor: pointer; }
  .user-name { font-size: 14px; color: #333; }
}
.main-content { background: #f5f7fa; padding: 20px; }
.fade-enter-active, .fade-exit-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-exit-to { opacity: 0; }
</style>
