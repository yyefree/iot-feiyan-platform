<template>
  <div class="page-container">
    <el-card shadow="hover">
      <template #header><span style="font-weight: 600;">系统设置</span></template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基础配置" name="basic">
          <el-form :model="basicForm" label-width="120px" style="max-width: 600px;">
            <el-form-item label="平台名称">
              <el-input v-model="basicForm.platformName" />
            </el-form-item>
            <el-form-item label="平台Logo">
              <el-upload action="/api/v1/settings/upload-logo" :on-success="handleLogoSuccess" :show-file-list="false">
                <el-avatar :size="80" :src="basicForm.logoUrl" />
                <template #trigger><el-button type="primary">上传图片</el-button></template>
              </el-upload>
            </el-form-item>
            <el-form-item label="平台描述">
              <el-input v-model="basicForm.description" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveBasic">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="设备配置" name="device">
          <el-form :model="deviceForm" label-width="120px" style="max-width: 600px;">
            <el-form-item label="默认认证方式">
              <el-select v-model="deviceForm.defaultAuthType">
                <el-option label="一机一密" value="one_device_one_key" />
                <el-option label="一型一密" value="one_model_one_key" />
                <el-option label="ID²认证" value="id2" />
              </el-select>
            </el-form-item>
            <el-form-item label="默认通信协议">
              <el-select v-model="deviceForm.defaultCommType">
                <el-option label="MQTT" value="mqtt" />
                <el-option label="CoAP" value="coap" />
                <el-option label="HTTP" value="http" />
              </el-select>
            </el-form-item>
            <el-form-item label="设备命名规则">
              <el-input v-model="deviceForm.deviceNamePattern" placeholder="如: device_{productId}_{index}" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveDevice">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="消息配置" name="message">
          <el-form :model="messageForm" label-width="120px" style="max-width: 600px;">
            <el-form-item label="消息推送渠道">
              <el-checkbox-group v-model="messageForm.pushChannels">
                <el-checkbox label="app">App推送</el-checkbox>
                <el-checkbox label="wechat">微信公众号</el-checkbox>
                <el-checkbox label="sms">短信</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="告警阈值">
              <el-input-number v-model="messageForm.alertThreshold" :min="1" :max="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveMessage">保存配置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('basic')

const basicForm = ref({
  platformName: 'IoT物联网平台',
  logoUrl: '',
  description: '面向消费级智能设备的物联网云平台'
})

const deviceForm = ref({
  defaultAuthType: 'one_device_one_key',
  defaultCommType: 'mqtt',
  deviceNamePattern: 'device_{productId}_{index}'
})

const messageForm = ref({
  pushChannels: ['app', 'wechat'],
  alertThreshold: 10
})

const handleSaveBasic = () => {
  ElMessage.success('基础配置保存成功')
}

const handleSaveDevice = () => {
  ElMessage.success('设备配置保存成功')
}

const handleSaveMessage = () => {
  ElMessage.success('消息配置保存成功')
}

const handleLogoSuccess = (response: any) => {
  basicForm.value.logoUrl = response.data.url
}

onMounted(() => {
  // 加载当前配置
})
</script>

<style scoped lang="scss">
.page-container { }
</style>
