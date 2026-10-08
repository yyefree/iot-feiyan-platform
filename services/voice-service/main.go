package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ========== 数据模型 ==========

type VoiceDevice struct {
	gorm.Model
	DeviceID     uint   `json:"device_id" gorm:"index"`
	Vendor       string `json:"vendor" gorm:"not null"` // 天猫精灵/Amazon Alexa/Google Home
	VendorDeviceID string `json:"vendor_device_id"` // 厂商设备ID
	VendorToken    string `json:"vendor_token" gorm:"type:text"`
	Status       string `json:"status" gorm:"default:'disconnected'"`
	LastSyncAt   *time.Time `json:"last_sync_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type VoiceCommandLog struct {
	gorm.Model
	DeviceID   uint   `json:"device_id" gorm:"index"`
	Vendor     string `json:"vendor"`
	Command    string `json:"command"` // 语音指令
	Result     string `json:"result"`
	ErrorMsg   string `json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`
}

// ========== 请求结构 ==========

type BindVoiceDeviceRequest struct {
	DeviceID         uint   `json:"device_id" binding:"required"`
	Vendor           string `json:"vendor" binding:"required"` // 天猫精灵/Amazon Alexa/Google Home
	VendorDeviceID   string `json:"vendor_device_id"`
	VendorToken      string `json:"vendor_token"`
}

type ExecuteVoiceCommandRequest struct {
	DeviceID uint   `json:"device_id" binding:"required"`
	Command  string `json:"command" binding:"required"`
	Vendor   string `json:"vendor"`
}

// ========== 服务初始化 ==========

var db *gorm.DB

func initDB() {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "iot-postgres"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "iot_admin"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "iot_admin_2024"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "iot_platform"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("WARNING: failed to connect database, will retry...: %v", err)
		// 尝试本地连接
		localDSN := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", port, user, password, dbname)
		db, err = gorm.Open(postgres.Open(localDSN), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database:", err)
		}
	}

	db.AutoMigrate(&VoiceDevice{}, &VoiceCommandLog{})
}

func main() {
	initDB()

	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})

	// 语音设备管理API
	r.GET("/api/v1/voice/devices", getVoiceDevices)
	r.POST("/api/v1/voice/devices", bindVoiceDevice)
	r.DELETE("/api/v1/voice/devices/:id", unbindVoiceDevice)
	r.PUT("/api/v1/voice/devices/:id/sync", syncVoiceDevice)
	r.GET("/api/v1/voice/devices/:id/status", getVoiceDeviceStatus)

	// 语音指令API
	r.POST("/api/v1/voice/commands", executeVoiceCommand)
	r.GET("/api/v1/voice/commands/logs", getCommandLogs)

	// 统计API
	r.GET("/api/v1/voice/statistics", getVoiceStatistics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	log.Printf("voice-service starting on :%s", port)
	r.Run(":" + port)
}

// ========== 语音设备管理处理函数 ==========

func getVoiceDevices(c *gin.Context) {
	var devices []VoiceDevice
	var total int64

	db.Model(&VoiceDevice{}).Count(&total)
	db.Order("created_at DESC").Limit(20).Find(&devices)

	c.JSON(200, gin.H{
		"code": 0,
		"data": devices,
		"pagination": gin.H{
			"page":  c.DefaultQuery("page", "1"),
			"size":  c.DefaultQuery("size", "20"),
			"total": total,
		},
	})
}

func bindVoiceDevice(c *gin.Context) {
	var req BindVoiceDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	device := VoiceDevice{
		DeviceID:       req.DeviceID,
		Vendor:         req.Vendor,
		VendorDeviceID: req.VendorDeviceID,
		VendorToken:    req.VendorToken,
		Status:         "connected",
	}

	if err := db.Create(&device).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": device})
}

func unbindVoiceDevice(c *gin.Context) {
	var device VoiceDevice
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "语音设备不存在"})
		return
	}

	db.Delete(&device)
	c.JSON(200, gin.H{"code": 0, "message": "解绑成功"})
}

func syncVoiceDevice(c *gin.Context) {
	var device VoiceDevice
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "语音设备不存在"})
		return
	}

	now := time.Now()
	device.LastSyncAt = &now
	db.Save(&device)
	c.JSON(200, gin.H{"code": 0, "message": "同步成功"})
}

func getVoiceDeviceStatus(c *gin.Context) {
	var device VoiceDevice
	if err := db.First(&device, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "语音设备不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": device})
}

// ========== 语音指令处理函数 ==========

func executeVoiceCommand(c *gin.Context) {
	var req ExecuteVoiceCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 模拟语音指令执行
	log := VoiceCommandLog{
		DeviceID: req.DeviceID,
		Vendor:   req.Vendor,
		Command:  req.Command,
		Result:   "success",
	}

	if err := db.Create(&log).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "message": "语音指令执行成功", "data": log})
}

func getCommandLogs(c *gin.Context) {
	var logs []VoiceCommandLog
	var total int64

	db.Model(&VoiceCommandLog{}).Count(&total)
	db.Order("created_at DESC").Limit(20).Find(&logs)

	c.JSON(200, gin.H{
		"code": 0,
		"data": logs,
		"pagination": gin.H{
			"page":  c.DefaultQuery("page", "1"),
			"size":  c.DefaultQuery("size", "20"),
			"total": total,
		},
	})
}

// ========== 语音统计 ==========

func getVoiceStatistics(c *gin.Context) {
	var connectedDevices int64
	var totalCommands int64
	var successCommands int64

	db.Model(&VoiceDevice{}).Where("status = 'connected'").Count(&connectedDevices)
	db.Model(&VoiceCommandLog{}).Count(&totalCommands)
	db.Model(&VoiceCommandLog{}).Where("result = 'success'").Count(&successCommands)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"connectedDevices": connectedDevices,
			"totalCommands":    totalCommands,
			"successCommands": successCommands,
		},
	})
}
