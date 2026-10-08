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
	IsVirtual       bool       `json:"is_virtual" gorm:"default:false"`
	VirtualInterval int        `json:"virtual_interval" gorm:"default:60"`
}

type DeviceLog struct {
	gorm.Model
	DeviceID   uint      `json:"device_id" gorm:"index"`
	TenantID   uint      `json:"tenant_id" gorm:"index"`
	Level      string    `json:"level"` // info/warn/error
	Message    string    `json:"message" gorm:"type:text"`
	Module     string    `json:"module"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
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
	
	db.AutoMigrate(&Device{}, &DeviceLog{})
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
	
	// 虚拟设备API
	r.POST("/api/v1/devices/virtual/create", createVirtualDevice)
	r.GET("/api/v1/devices/virtual/list", getVirtualDeviceList)
	
	// 设备认证API
	r.POST("/api/v1/devices/auth", deviceAuth)
	r.GET("/api/v1/devices/:id/credentials", getDeviceCredentials)
	
	// 设备统计
	r.GET("/api/v1/devices/statistics", getDeviceStatistics)
	
	// 设备日志API
	r.GET("/api/v1/devices/:id/logs", getDeviceLogs)
	r.POST("/api/v1/devices/:id/logs", createDeviceLog)
	
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
			"page":       page,
			"size":       pageSize,
			"total":      total,
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
	
	// 记录操作日志
	logDeviceOperation(device.ID, device.TenantID, "info", "设备创建成功", "device", device.IPAddress, "")
	
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
	
	now := time.Now()
	device.Status = "online"
	device.Online = true
	device.LastOnlineAt = &now
	db.Save(&device)
	
	logDeviceOperation(device.ID, device.TenantID, "info", "设备已启用", "device", "", "")
	
	c.JSON(200, gin.H{"code": 0, "message": "设备已启用"})
}

func disableDevice(c *gin.Context) {
	var device Device
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	
	now := time.Now()
	device.Status = "offline"
	device.Online = false
	device.LastOfflineAt = &now
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
		Status:          "online",
		Online:          true,
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
	
	logDeviceOperation(device.ID, device.TenantID, "info", "设备认证成功", "auth", device.IPAddress, c.GetHeader("User-Agent"))
	
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

// ========== 设备日志API ==========

func getDeviceLogs(c *gin.Context) {
	deviceID := c.Param("id")
	
	var logs []DeviceLog
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&DeviceLog{}).Where("device_id = ?", deviceID).Count(&total)
	db.Where("device_id = ?", deviceID).Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": logs,
		"pagination": gin.H{
			"page":  page,
			"size":  pageSize,
			"total": total,
		},
	})
}

func createDeviceLog(c *gin.Context) {
	var req struct {
		DeviceID  uint   `json:"device_id" binding:"required"`
		TenantID  uint   `json:"tenant_id"`
		Level     string `json:"level" binding:"required"`
		Message   string `json:"message" binding:"required"`
		Module    string `json:"module"`
		IPAddress string `json:"ip_address"`
		UserAgent string `json:"user_agent"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	log := DeviceLog{
		DeviceID:   req.DeviceID,
		TenantID:   req.TenantID,
		Level:      req.Level,
		Message:    req.Message,
		Module:     req.Module,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		CreatedAt:  time.Now(),
	}
	
	if err := db.Create(&log).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "data": log})
}

// ========== 工具函数 ==========

func logDeviceOperation(deviceID uint, tenantID uint, level, message, module, ip, ua string) {
	entry := DeviceLog{
		DeviceID:  deviceID,
		TenantID:  tenantID,
		Level:     level,
		Message:   message,
		Module:    module,
		IPAddress: ip,
		UserAgent: ua,
		CreatedAt: time.Now(),
	}
	db.Create(&entry)
}

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
