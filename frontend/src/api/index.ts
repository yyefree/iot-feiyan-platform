import request from '@/utils/request'

// ========== 设备相关API ==========
export function getDeviceList(params: Record<string, any>) {
  return request({ url: '/devices', method: 'get', params })
}

export function getDevice(id: number) {
  return request({ url: `/devices/${id}`, method: 'get' })
}

export function createDevice(data: Record<string, any>) {
  return request({ url: '/devices', method: 'post', data })
}

export function batchCreateDevices(data: Record<string, any>) {
  return request({ url: '/devices/batch', method: 'post', data })
}

export function updateDevice(id: number, data: Record<string, any>) {
  return request({ url: `/devices/${id}`, method: 'put', data })
}

export function deleteDevice(id: number) {
  return request({ url: `/devices/${id}`, method: 'delete' })
}

export function enableDevice(id: number) {
  return request({ url: `/devices/${id}/enable`, method: 'post' })
}

export function disableDevice(id: number) {
  return request({ url: `/devices/${id}/disable`, method: 'post' })
}

export function resetDevice(id: number) {
  return request({ url: `/devices/${id}/reset`, method: 'post' })
}

export function getDeviceCredentials(id: number) {
  return request({ url: `/devices/${id}/credentials`, method: 'get' })
}

export function getDeviceStatistics() {
  return request({ url: '/devices/statistics', method: 'get' })
}

export function getDeviceLogs(id: number, params?: Record<string, any>) {
  return request({ url: `/devices/${id}/logs`, method: 'get', params })
}

// ========== 产品相关API ==========
export function getProductList(params: Record<string, any>) {
  return request({ url: '/products', method: 'get', params })
}

export function getProduct(id: number) {
  return request({ url: `/products/${id}`, method: 'get' })
}

export function createProduct(data: Record<string, any>) {
  return request({ url: '/products', method: 'post', data })
}

export function updateProduct(id: number, data: Record<string, any>) {
  return request({ url: `/products/${id}`, method: 'put', data })
}

export function deleteProduct(id: number) {
  return request({ url: `/products/${id}`, method: 'delete' })
}

export function publishProduct(id: number) {
  return request({ url: `/products/${id}/publish`, method: 'post' })
}

export function unpublishProduct(id: number) {
  return request({ url: `/products/${id}/unpublish`, method: 'post' })
}

export function getProductStatistics() {
  return request({ url: '/products/statistics', method: 'get' })
}

// ========== 物模型API ==========
export function getPropertyList(productId?: number) {
  const url = productId ? `/products/${productId}/properties` : '/products/properties'
  return request({ url, method: 'get' })
}

export function addProperty(data: Record<string, any>) {
  return request({ url: '/properties', method: 'post', data })
}

export function updateProperty(id: number, data: Record<string, any>) {
  return request({ url: `/properties/${id}`, method: 'put', data })
}

export function deleteProperty(id: number) {
  return request({ url: `/properties/${id}`, method: 'delete' })
}

export function getEventList(productId?: number) {
  const url = productId ? `/products/${productId}/events` : '/products/events'
  return request({ url, method: 'get' })
}

export function addEvent(data: Record<string, any>) {
  return request({ url: '/events', method: 'post', data })
}

export function updateEvent(id: number, data: Record<string, any>) {
  return request({ url: `/events/${id}`, method: 'put', data })
}

export function deleteEvent(id: number) {
  return request({ url: `/events/${id}`, method: 'delete' })
}

export function getServiceList(productId?: number) {
  const url = productId ? `/products/${productId}/services` : '/products/services'
  return request({ url, method: 'get' })
}

export function addService(data: Record<string, any>) {
  return request({ url: '/services', method: 'post', data })
}

export function updateService(id: number, data: Record<string, any>) {
  return request({ url: `/services/${id}`, method: 'put', data })
}

export function deleteService(id: number) {
  return request({ url: `/services/${id}`, method: 'delete' })
}

export function getProductThingModel(productId: number) {
  return request({ url: `/products/${productId}/thing-model`, method: 'get' })
}

export function updateProductThingModel(productId: number, data: Record<string, any>) {
  return request({ url: `/products/${productId}/thing-model`, method: 'put', data })
}

export function publishThingModel() {
  return request({ url: '/thing-model/publish', method: 'post' })
}

// ========== 品类API ==========
export function getCategoryList() {
  return request({ url: '/categories', method: 'get' })
}

export function createCategory(data: Record<string, any>) {
  return request({ url: '/categories', method: 'post', data })
}

export function updateCategory(id: number, data: Record<string, any>) {
  return request({ url: `/categories/${id}`, method: 'put', data })
}

export function deleteCategory(id: number) {
  return request({ url: `/categories/${id}`, method: 'delete' })
}

// ========== 规则引擎API ==========
export function getRuleList(params: Record<string, any>) {
  return request({ url: '/rules', method: 'get', params })
}

export function getRule(id: number) {
  return request({ url: `/rules/${id}`, method: 'get' })
}

export function createRule(data: Record<string, any>) {
  return request({ url: '/rules', method: 'post', data })
}

export function updateRule(id: number, data: Record<string, any>) {
  return request({ url: `/rules/${id}`, method: 'put', data })
}

export function deleteRule(id: number) {
  return request({ url: `/rules/${id}`, method: 'delete' })
}

export function enableRule(id: number) {
  return request({ url: `/rules/${id}/enable`, method: 'post' })
}

export function disableRule(id: number) {
  return request({ url: `/rules/${id}/disable`, method: 'post' })
}

export function getRuleLogs(id: number, params?: Record<string, any>) {
  return request({ url: `/rules/${id}/logs`, method: 'get', params })
}

export function testRule(id: number) {
  return request({ url: `/rules/${id}/execute`, method: 'get' })
}

export function getRuleStatistics() {
  return request({ url: '/rules/statistics', method: 'get' })
}

// ========== 场景联动API ==========
export function getSceneList(params: Record<string, any>) {
  return request({ url: '/scenes', method: 'get', params })
}

export function getScene(id: number) {
  return request({ url: `/scenes/${id}`, method: 'get' })
}

export function createScene(data: Record<string, any>) {
  return request({ url: '/scenes', method: 'post', data })
}

export function updateScene(id: number, data: Record<string, any>) {
  return request({ url: `/scenes/${id}`, method: 'put', data })
}

export function deleteScene(id: number) {
  return request({ url: `/scenes/${id}`, method: 'delete' })
}

export function enableScene(id: number) {
  return request({ url: `/scenes/${id}/enable`, method: 'post' })
}

export function disableScene(id: number) {
  return request({ url: `/scenes/${id}/disable`, method: 'post' })
}

export function testScene(id: number) {
  return request({ url: `/scenes/${id}/execute`, method: 'post' })
}

// ========== 数据服务API ==========
export function reportTelemetry(data: Record<string, any>) {
  return request({ url: '/data/telemetry', method: 'post', data })
}

export function batchReportTelemetry(data: Record<string, any>[]) {
  return request({ url: '/data/telemetry/batch', method: 'post', data })
}

export function getTelemetryData(params: Record<string, any>) {
  return request({ url: '/data/telemetry', method: 'get', params })
}

export function getDeviceTelemetry(deviceId: number, params?: Record<string, any>) {
  return request({ url: `/data/telemetry/device/${deviceId}`, method: 'get', params })
}

export function getLatestTelemetry(deviceId: number) {
  return request({ url: `/data/telemetry/device/${deviceId}/latest`, method: 'get' })
}

export function getDeviceStatistics(deviceId: number) {
  return request({ url: `/data/telemetry/device/${deviceId}/stats`, method: 'get' })
}

export function getLatestProperties(params?: Record<string, any>) {
  return request({ url: '/data/properties', method: 'get', params })
}

export function getDeviceLatestProperties(deviceId: number) {
  return request({ url: `/data/properties/device/${deviceId}`, method: 'get' })
}

export function getDataStatistics(params?: Record<string, any>) {
  return request({ url: '/data/statistics', method: 'get', params })
}

export function getDeviceDataStatistics(deviceId: number) {
  return request({ url: `/data/statistics/device/${deviceId}`, method: 'get' })
}

export function getStatisticsOverview() {
  return request({ url: '/data/statistics/overview', method: 'get' })
}

// ========== 设备分组API ==========
export function getGroupList() {
  return request({ url: '/device-groups', method: 'get' })
}

export function createGroup(data: Record<string, any>) {
  return request({ url: '/device-groups', method: 'post', data })
}

export function updateGroup(id: number, data: Record<string, any>) {
  return request({ url: `/device-groups/${id}`, method: 'put', data })
}

export function deleteGroup(id: number) {
  return request({ url: `/device-groups/${id}`, method: 'delete' })
}

export function addDeviceToGroup(groupId: number, deviceId: number) {
  return request({ url: `/device-groups/${groupId}/devices`, method: 'post', data: { deviceId } })
}

export function removeDeviceFromGroup(groupId: number, deviceId: number) {
  return request({ url: `/device-groups/${groupId}/devices/${deviceId}`, method: 'delete' })
}

// ========== 激活码API ==========
export function getActivationCodeList(params: Record<string, any>) {
  return request({ url: '/activation-codes', method: 'get', params })
}

export function generateActivationCodes(data: Record<string, any>) {
  return request({ url: '/activation-codes/generate', method: 'post', data })
}

export function deleteActivationCode(id: number) {
  return request({ url: `/activation-codes/${id}`, method: 'delete' })
}

// ========== OTA升级API ==========
export function getFirmwareList(params: Record<string, any>) {
  return request({ url: '/ota/firmwares', method: 'get', params })
}

export function createFirmware(data: Record<string, any>) {
  return request({ url: '/ota/firmwares', method: 'post', data })
}

export function deleteFirmware(id: number) {
  return request({ url: `/ota/firmwares/${id}`, method: 'delete' })
}

export function publishFirmware(id: number) {
  return request({ url: `/ota/firmwares/${id}/publish`, method: 'post' })
}

export function getTaskList(params: Record<string, any>) {
  return request({ url: '/ota/tasks', method: 'get', params })
}

export function createTask(data: Record<string, any>) {
  return request({ url: '/ota/tasks', method: 'post', data })
}

export function deleteTask(id: number) {
  return request({ url: `/ota/tasks/${id}`, method: 'delete' })
}

// ========== 语音平台绑定API ==========
export function getVoiceBindList(params: Record<string, any>) {
  return request({ url: '/voice-bindings', method: 'get', params })
}

export function bindVoicePlatform(data: Record<string, any>) {
  return request({ url: '/voice-bindings/bind', method: 'post', data })
}

export function unbindVoicePlatform(id: number) {
  return request({ url: `/voice-bindings/${id}/unbind`, method: 'post' })
}

export function deleteVoiceBind(id: number) {
  return request({ url: `/voice-bindings/${id}`, method: 'delete' })
}

// ========== 语音映射API ==========
export function getVoiceMappingList(params: Record<string, any>) {
  return request({ url: '/voice-mappings', method: 'get', params })
}

export function createVoiceMapping(data: Record<string, any>) {
  return request({ url: '/voice-mappings', method: 'post', data })
}

export function updateVoiceMapping(id: number, data: Record<string, any>) {
  return request({ url: `/voice-mappings/${id}`, method: 'put', data })
}

export function deleteVoiceMapping(id: number) {
  return request({ url: `/voice-mappings/${id}`, method: 'delete' })
}

// ========== 项目API ==========
export function getProjectList() {
  return request({ url: '/projects', method: 'get' })
}

export function createProject(data: Record<string, any>) {
  return request({ url: '/projects', method: 'post', data })
}

export function updateProject(id: number, data: Record<string, any>) {
  return request({ url: `/projects/${id}`, method: 'put', data })
}

export function deleteProject(id: number) {
  return request({ url: `/projects/${id}`, method: 'delete' })
}

// ========== 日志API ==========
export function getLogList(params: Record<string, any>) {
  return request({ url: '/logs', method: 'get', params })
}
