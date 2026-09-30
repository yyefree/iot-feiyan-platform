import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { title: '工作台' }
  },
  {
    path: '/projects',
    name: 'Projects',
    component: () => import('../views/Projects.vue'),
    meta: { title: '项目管理' }
  },
  {
    path: '/products',
    name: 'Products',
    component: () => import('../views/Products.vue'),
    meta: { title: '产品管理' }
  },
  {
    path: '/thing-model',
    name: 'ThingModel',
    component: () => import('../views/ThingModel.vue'),
    meta: { title: '物模型' }
  },
  {
    path: '/devices',
    name: 'Devices',
    component: () => import('../views/Devices.vue'),
    meta: { title: '设备管理' }
  },
  {
    path: '/device-groups',
    name: 'DeviceGroups',
    component: () => import('../views/DeviceGroups.vue'),
    meta: { title: '设备分组' }
  },
  {
    path: '/activation-codes',
    name: 'ActivationCodes',
    component: () => import('../views/ActivationCodes.vue'),
    meta: { title: '激活码管理' }
  },
  {
    path: '/scenes',
    name: 'Scenes',
    component: () => import('../views/Scenes.vue'),
    meta: { title: '场景联动' }
  },
  {
    path: '/rules',
    name: 'Rules',
    component: () => import('../views/Rules.vue'),
    meta: { title: '规则引擎' }
  },
  {
    path: '/data-analysis',
    name: 'DataAnalysis',
    component: () => import('../views/DataAnalysis.vue'),
    meta: { title: '数据分析' }
  },
  {
    path: '/telemetry',
    name: 'Telemetry',
    component: () => import('../views/Telemetry.vue'),
    meta: { title: '遥测数据' }
  },
  {
    path: '/logs',
    name: 'Logs',
    component: () => import('../views/Logs.vue'),
    meta: { title: '日志查询' }
  },
  {
    path: '/ota-firmware',
    name: 'OtaFirmware',
    component: () => import('../views/OtaFirmware.vue'),
    meta: { title: '固件管理' }
  },
  {
    path: '/ota-tasks',
    name: 'OtaTasks',
    component: () => import('../views/OtaTasks.vue'),
    meta: { title: '升级任务' }
  },
  {
    path: '/voice-bind',
    name: 'VoiceBind',
    component: () => import('../views/VoiceBind.vue'),
    meta: { title: '语音平台绑定' }
  },
  {
    path: '/voice-mapping',
    name: 'VoiceMapping',
    component: () => import('../views/VoiceMapping.vue'),
    meta: { title: '语音映射' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue'),
    meta: { title: '系统设置' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || ''} - IoT飞燕平台`
  next()
})

export default router
