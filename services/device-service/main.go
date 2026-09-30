package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ========== 数据模型 ==========

type Device struct {
	gorm.Model
	DeviceName      string     `json:"device_name" gorm:"uniqueIndex;not null"`
	DeviceKey       string     `json:"device_key" gorm:"uniqueIndex"`
	DeviceSecret    string     `json:"-" gorm:"not null"`
	ProductID       uint       `json:"product_id" gorm:"index"`
	ProductKey      string     `json:"product_key" gorm:"index"`
	TenantID        uint       `json:"tenant_id" gorm:"index"`
	Status          string     `json:"status" gorm:"default:'offline'"`
	Online          bool       `json:"online"`
	LastOnlineAt    *time.Time `json:"last_online_at"`
	LastOfflineAt   *time.Time `json:"last_offline_at"`
	IPAddress       string     `json:"ip_address"`
	FirmwareVersion string     `json:"firmware_version"`
	Region          string     `json:"region"`
	GroupID         uint       `json:"group_id"`
	Tags            string     `json:"tags" gorm:"type:text"`
	Description     string     `json:"description"`
	IsVirtual       bool       `json:"is_virtual" gorm:"default:false"` // 虚拟设备标识
	VirtualInterval int        `json:"virtual_interval" gorm:"default:60"` // 虚拟设备上报间隔（秒）
}

type DeviceShadow struct {
	gorm.Model
	DeviceID  uint   `json:"device_id" gorm:"uniqueIndex;index"`
	Reported  string `json:"reported" gorm:"type:jsonb"`
	Desired   string `json:"desired" gorm:"type:jsonb"`
	Version   int    `json:"version" gorm:"default:1"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeviceGroup struct {
	gorm.Model
	Name          string `json:"name" gorm:"not null"`
	ParentID      uint   `json:"parent_id"`
	TenantID      uint   `json:"tenant_id" gorm:"index"`
	DeviceCount   int    `json:"device_count" gorm:"default:0"`
}

type ActivationCode struct {
	gorm.Model
	ProductID uint   `json:"product_id" gorm:"index"`
	Code      string `json:"code" gorm:"uniqueIndex;not null"`
	Status    string `json:"status" gorm:"default:'unused'"`
	UsedAt    *time.Time `json:"used_at"`
	UsedBy    uint   `json:"used_by"`
}

// ========== 请求结构 ==========

type CreateDeviceRequest struct {
	DeviceName    string `json:"device_name" binding:"required"`
	ProductID     uint   `json:"product_id" binding:"required"`
	TenantID      uint   `json:"tenant_id"`
	Description   string `json:"description"`
	IsVirtual     bool   `json:"is_virtual"`
	VirtualInterval int  `json:"virtual_interval"`
}

type BatchCreateRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	TenantID  uint   `json:"tenant_id"`
	Count     int    `json:"count" binding:"required,min=1,max=1000"`
	Prefix    string `json:"prefix"`
}

type CreateGroupRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint   `json:"parent_id"`
}

// ========== 服务初始化 ==========

var db *gorm.DB

func initDB() {
	host := getEnv("POSTGRES_HOST", "iot-postgres")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "iot_admin")
	password := getEnv("POSTGRES_PASSWORD", "iot_admin_2024")
	dbname := getEnv("POSTGRES_DB", "iot_platform")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	
	db.AutoMigrate(&Device{}, &DeviceShadow{}, &DeviceGroup{}, &ActivationCode{})
}

func main() {
	initDB()
	
	r := gin.Default()
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})
	
	// 设备管理API
	r.GET("/api/v1/devices", getDeviceList)
	r.GET("/api/v1/devices/:id", getDevice)
	r.POST("/api/v1/devices", createDevice)
	r.POST("/api/v1/devices/batch", batchCreateDevices)
	r.PUT("/api/v1/devices/:id", updateDevice)
	r.DELETE("/api/v1/devices/:id", deleteDevice)
	r.POST("/api/v1/devices/:id/enable", enableDevice)
	r.POST("/api/v1/devices/:id/disable", disableDevice)
	r.POST("/api/v1/devices/:id/reset", resetDevice)
	
	// 虚拟设备API（新增）
	r.POST("/api/v1/devices/virtual/create", createVirtualDevice)
	r.POST("/api/v1/devices/virtual/:id/report", reportVirtualProperty)
	r.POST("/api/v1/devices/virtual/:id/event", triggerVirtualEvent)
	r.GET("/api/v1/devices/virtual/list", getVirtualDeviceList)
	
	// 设备认证API
	r.POST("/api/v1/devices/auth", deviceAuth)
	r.GET("/api/v1/devices/:id/credentials", getDeviceCredentials)
	
	// 设备影子API
	r.GET("/api/v1/devices/:id/shadow", getDeviceShadow)
	r.PUT("/api/v1/devices/:id/shadow", updateDeviceShadow)
	
	// 设备统计
	r.GET("/api/v1/devices/statistics", getDeviceStatistics)
	r.GET("/api/v1/devices/:id/logs", getDeviceLogs)
	
	// 设备分组API
	r.GET("/api/v1/device-groups", getDeviceGroups)
	r.POST("/api/v1/device-groups", createDeviceGroup)
	r.PUT("/api/v1/device-groups/:id", updateDeviceGroup)
	r.DELETE("/api/v1/device-groups/:id", deleteDeviceGroup)
	r.POST("/api/v1/device-groups/:id/devices", addDeviceToGroup)
	r.DELETE("/api/v1/device-groups/:id/devices/:device_id", removeDeviceFromGroup)
	
	// 激活码API
	r.GET("/api/v1/activation-codes", getActivationCodes)
	r.POST("/api/v1/activation-codes/generate", generateActivationCodes)
	r.DELETE("/api/v1/activation-codes/:id", deleteActivationCode)
	
	log.Printf("device-service starting on :8081")
	r.Run(":8081")
}

// ========== 设备管理处理函数 ==========

func getDeviceList(c *gin.Context) {
	var devices []Device
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	offset := (page - 1) * pageSize
	
	db.Model(&Device{}).Count(&total)
	db.Offset(offset).Limit(pageSize).Find(&devices)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": devices,
		"pagination": gin.H{
			"page": page,
			"size": pageSize,
			"total": total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func getDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": device})
}

func createDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var productKey string
	db.Raw("SELECT product_key FROM products WHERE id = ?", req.ProductID).Scan(&productKey)
	if productKey == "" {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	
	device := Device{
		DeviceName:    req.DeviceName,
		DeviceKey:     generateDeviceKey(),
		DeviceSecret:  generateDeviceSecret(),
		ProductID:     req.ProductID,
		ProductKey:    productKey,
		TenantID:      req.TenantID,
		Status:        "offline",
		Online:        false,
		Description:   req.Description,
		IsVirtual:     req.IsVirtual,
		VirtualInterval: req.VirtualInterval,
	}
	
	if err := db.Create(&device).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	shadow := DeviceShadow{
		DeviceID: device.ID,
		Version:  1,
	}
	db.Create(&shadow)
	
	c.JSON(201, gin.H{"code": 0, "data": device})
}

func batchCreateDevices(c *gin.Context) {
	var req BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var productKey string
	db.Raw("SELECT product_key FROM products WHERE id = ?", req.ProductID).Scan(&productKey)
	if productKey == "" {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	
	var devices []Device
	prefix := req.Prefix
	if prefix == "" {
		prefix = "device"
	}
	
	for i := 0; i < req.Count; i++ {
		device := Device{
			DeviceName: fmt.Sprintf("%s_%04d", prefix, i+1),
			DeviceKey:  generateDeviceKey(),
			DeviceSecret: generateDeviceSecret(),
			ProductID:  req.ProductID,
			ProductKey: productKey,
			TenantID:   req.TenantID,
			Status:     "offline",
			Online:     false,
		}
		devices = append(devices, device)
	}
	
	if err := db.Create(&devices).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "message": fmt.Sprintf("成功创建 %d 台设备", len(devices)), "data": devices})
}

func updateDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	if req.DeviceName != "" {
		device.DeviceName = req.DeviceName
	}
	if req.Description != "" {
		device.Description = req.Description
	}
	
	db.Save(&device)
	c.JSON(200, gin.H{"code": 0, "data": device})
}

func deleteDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	db.Delete(&device)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func enableDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	device.Status = "online"
	db.Save(&device)
	c.JSON(200, gin.H{"code": 0, "message": "设备已启用"})
}

func disableDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	device.Status = "offline"
	db.Save(&device)
	c.JSON(200, gin.H{"code": 0, "message": "设备已禁用"})
}

func resetDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	device.DeviceSecret = generateDeviceSecret()
	db.Save(&device)
	c.JSON(200, gin.H{"code": 0, "message": "设备已重置", "data": device})
}

// ========== 虚拟设备API ==========

func createVirtualDevice(c *gin.Context) {
	var req CreateDeviceRequest
	req.IsVirtual = true
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var productKey string
	db.Raw("SELECT product_key FROM products WHERE id = ?", req.ProductID).Scan(&productKey)
	if productKey == "" {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	
	device := Device{
		DeviceName:      req.DeviceName,
		DeviceKey:       generateDeviceKey(),
		DeviceSecret:    generateDeviceSecret(),
		ProductID:       req.ProductID,
		ProductKey:      productKey,
		TenantID:        req.TenantID,
		Status:          "offline",
		Online:          true, // 虚拟设备默认在线
		Description:     req.Description,
		IsVirtual:       true,
		VirtualInterval: 60,
	}
	
	if err := db.Create(&device).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "data": device})
}

func reportVirtualProperty(c *gin.Context) {
	deviceID := c.Param("id")
	
	var device Device
	if err := db.First(&device, deviceID).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	if !device.IsVirtual {
		c.JSON(400, gin.H{"code": 400, "message": "不是虚拟设备"})
		return
	}
	
	// 模拟属性上报
	var req struct {
		Identify string  `json:"identify"`
		Value    float64 `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	// 更新影子
	shadow := DeviceShadow{
		DeviceID: device.ID,
		Reported: fmt.Sprintf(`{"%s":%f}`, req.Identify, req.Value),
		Version:  1,
	}
	db.Where("device_id = ?", deviceID).FirstOrCreate(&shadow)
	db.Model(&shadow).Updates(map[string]interface{}{
		"reported": shadow.Reported,
		"version":  shadow.Version + 1,
		"updated_at": time.Now(),
	})
	
	c.JSON(200, gin.H{"code": 0, "message": "属性上报成功"})
}

func triggerVirtualEvent(c *gin.Context) {
	deviceID := c.Param("id")
	
	var device Device
	if err := db.First(&device, deviceID).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	if !device.IsVirtual {
		c.JSON(400, gin.H{"code": 400, "message": "不是虚拟设备"})
		return
	}
	
	var req struct {
		Identify string `json:"identify"`
		Level    string `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"code": 0, "message": fmt.Sprintf("事件触发成功: %s", req.Identify)})
}

func getVirtualDeviceList(c *gin.Context) {
	var devices []Device
	db.Where("is_virtual = true").Find(&devices)
	
	c.JSON(200, gin.H{"code": 0, "data": devices})
}

// ========== 设备认证API ==========

func deviceAuth(c *gin.Context) {
	var authReq struct {
		ProductKey string `json:"productKey" binding:"required"`
		DeviceName string `json:"deviceName" binding:"required"`
		Timestamp  int64  `json:"timestamp"`
		SignMethod string `json:"signMethod" default:"hmacsha1"`
		Sign       string `json:"sign" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var device Device
	if err := db.Where("product_key = ? AND device_name = ?", authReq.ProductKey, authReq.DeviceName).First(&device).Error; err != nil {
		c.JSON(401, gin.H{"code": 401, "message": "认证失败"})
		return
	}
	
	now := time.Now()
	device.Online = true
	device.LastOnlineAt = &now
	device.Status = "online"
	device.IPAddress = c.ClientIP()
	db.Save(&device)
	
	c.JSON(200, gin.H{
		"code": 0,
		"message": "认证成功",
		"data": gin.H{
			"deviceName": device.DeviceName,
			"productKey": device.ProductKey,
		},
	})
}

func getDeviceCredentials(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"productKey":   device.ProductKey,
			"deviceName":   device.DeviceName,
			"deviceSecret": device.DeviceSecret,
			"authType":     "one_device_one_key",
		},
	})
}

// ========== 设备影子API ==========

func getDeviceShadow(c *gin.Context) {
	var shadow DeviceShadow
	if err := db.Where("device_id = ?", c.Param("id")).First(&shadow).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "影子不存在"})
		return
	}
	
	c.JSON(200, gin.H{"code": 0, "data": shadow})
}

func updateDeviceShadow(c *gin.Context) {
	var shadow DeviceShadow
	if err := db.Where("device_id = ?", c.Param("id")).First(&shadow).Error; err != nil {
		shadow = DeviceShadow{
			DeviceID: uintStrToUint(c.Param("id")),
			Version:  1,
		}
	}
	
	var req struct {
		Reported string `json:"reported"`
		Desired  string `json:"desired"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	if req.Reported != "" {
		shadow.Reported = req.Reported
	}
	if req.Desired != "" {
		shadow.Desired = req.Desired
	}
	shadow.Version++
	shadow.UpdatedAt = time.Now()
	
	db.Save(&shadow)
	c.JSON(200, gin.H{"code": 0, "data": shadow})
}

// ========== 设备统计API ==========

func getDeviceStatistics(c *gin.Context) {
	var totalCount, onlineCount, offlineCount int64
	
	db.Raw("SELECT COUNT(*) FROM devices").Scan(&totalCount)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = true").Scan(&onlineCount)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = false").Scan(&offlineCount)
	
	stats := gin.H{
		"totalCount": totalCount,
		"onlineCount": onlineCount,
		"offlineCount": offlineCount,
	}
	
	c.JSON(200, gin.H{"code": 0, "data": stats})
}

func getDeviceLogs(c *gin.Context) {
	logs := []map[string]string{
		{"time": time.Now().Add(-1 * time.Hour).Format(time.RFC3339), "level": "INFO", "msg": "设备上线"},
		{"time": time.Now().Add(-30 * time.Minute).Format(time.RFC3339), "level": "INFO", "msg": "属性上报: temperature=25.5"},
		{"time": time.Now().Add(-15 * time.Minute).Format(time.RFC3339), "level": "WARN", "msg": "信号强度较弱"},
		{"time": time.Now().Add(-5 * time.Minute).Format(time.RFC3339), "level": "INFO", "msg": "固件版本检查"},
	}
	
	c.JSON(200, gin.H{"code": 0, "data": logs})
}

// ========== 设备分组API ==========

func getDeviceGroups(c *gin.Context) {
	var groups []DeviceGroup
	db.Where("tenant_id = 1").Order("id ASC").Find(&groups)
	c.JSON(200, gin.H{"code": 0, "data": groups})
}

func createDeviceGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	group := DeviceGroup{
		Name:     req.Name,
		ParentID: req.ParentID,
		TenantID: 1,
	}
	
	if err := db.Create(&group).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "data": group})
}

func updateDeviceGroup(c *gin.Context) {
	var group DeviceGroup
	if err := db.First(&group, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "分组不存在"})
		return
	}
	
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	group.Name = req.Name
	group.ParentID = req.ParentID
	db.Save(&group)
	
	c.JSON(200, gin.H{"code": 0, "data": group})
}

func deleteDeviceGroup(c *gin.Context) {
	var group DeviceGroup
	if err := db.First(&group, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "分组不存在"})
		return
	}
	
	db.Delete(&group)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func addDeviceToGroup(c *gin.Context) {
	groupID := uintStrToUint(c.Param("id"))
	deviceID := uintStrToUint(c.Param("device_id"))
	
	db.Model(&Device{}).Where("id = ?", deviceID).Update("group_id", groupID)
	
	c.JSON(200, gin.H{"code": 0, "message": "添加成功"})
}

func removeDeviceFromGroup(c *gin.Context) {
	deviceID := uintStrToUint(c.Param("device_id"))
	
	db.Model(&Device{}).Where("id = ?", deviceID).Update("group_id", 0)
	
	c.JSON(200, gin.H{"code": 0, "message": "移除成功"})
}

// ========== 激活码API ==========

func getActivationCodes(c *gin.Context) {
	var codes []ActivationCode
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&ActivationCode{}).Count(&total)
	db.Offset((page-1)*pageSize).Limit(pageSize).Find(&codes)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": codes,
		"pagination": gin.H{
			"page": page,
			"size": pageSize,
			"total": total,
		},
	})
}

func generateActivationCodes(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Count     int  `json:"count" binding:"required,min=1,max=10000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var productKey string
	db.Raw("SELECT product_key FROM products WHERE id = ?", req.ProductID).Scan(&productKey)
	if productKey == "" {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	
	var codes []ActivationCode
	for i := 0; i < req.Count; i++ {
		code := ActivationCode{
			ProductID: req.ProductID,
			Code:      generateCode(),
			Status:    "unused",
		}
		codes = append(codes, code)
	}
	
	if err := db.Create(&codes).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "message": fmt.Sprintf("成功生成 %d 个激活码", len(codes))})
}

func deleteActivationCode(c *gin.Context) {
	var code ActivationCode
	if err := db.First(&code, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "激活码不存在"})
		return
	}
	
	db.Delete(&code)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

// ========== 工具函数 ==========

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func generateDeviceKey() string {
	return fmt.Sprintf("%08x%04x%04x%04x%012x", 
		os.Getpid(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano())
}

func generateDeviceSecret() string {
	return fmt.Sprintf("%064x", time.Now().UnixNano())
}

func generateCode() string {
	return fmt.Sprintf("%08x%04x%04x%04x%012x", 
		os.Getpid(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano(), 
		time.Now().UnixNano())
}

func uintStrToUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint(n)
}
