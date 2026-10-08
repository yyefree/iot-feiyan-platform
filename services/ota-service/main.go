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

type Firmware struct {
	gorm.Model
	ProductName  string    `json:"product_name" gorm:"not null;index"`
	FirmwareName string    `json:"firmware_name" gorm:"not null"`
	Version      string    `json:"version" gorm:"not null"`
	FileURL      string    `json:"file_url"`
	FileSize     int64     `json:"file_size"`
	Checksum     string    `json:"checksum"`
	ReleaseNotes string    `json:"release_notes" gorm:"type:text"`
	Status       string    `json:"status" gorm:"default:'draft'"` // draft/active/archived
	CreatedBy    uint      `json:"created_by"`
	TenantID     uint      `json:"tenant_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpgradeTask struct {
	gorm.Model
	FirmwareID     uint       `json:"firmware_id" gorm:"index"`
	TaskName       string     `json:"task_name" gorm:"not null"`
	TargetPercent  int        `json:"target_percent" gorm:"default:100"`
	Status         string     `json:"status" gorm:"default:'pending'"` // pending/running/completed/failed/rollback
	Description    string     `json:"description"`
	TenantID       uint       `json:"tenant_id"`
	StartedAt      *time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

type DeviceUpgradeLog struct {
	gorm.Model
	DeviceID       uint      `json:"device_id" gorm:"index"`
	UpgradeTaskID  uint      `json:"upgrade_task_id" gorm:"index"`
	FirmwareID     uint      `json:"firmware_id"`
	FromVersion    string    `json:"from_version"`
	ToVersion      string    `json:"to_version"`
	Status         string    `json:"status"` // pending/downloading/updated/failed/rollback
	Progress       int       `json:"progress"`
	ErrorMsg       string    `json:"error_msg" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at"`
}

// ========== 请求结构 ==========

type CreateFirmwareRequest struct {
	ProductName  string `json:"product_name" binding:"required"`
	FirmwareName string `json:"firmware_name" binding:"required"`
	Version      string `json:"version" binding:"required"`
	FileURL      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	Checksum     string `json:"checksum"`
	ReleaseNotes string `json:"release_notes"`
	TenantID     uint   `json:"tenant_id"`
}

type CreateUpgradeTaskRequest struct {
	FirmwareID    uint   `json:"firmware_id" binding:"required"`
	TaskName      string `json:"task_name" binding:"required"`
	TargetPercent int    `json:"target_percent"`
	Description   string `json:"description"`
	TenantID      uint   `json:"tenant_id"`
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

	db.AutoMigrate(&Firmware{}, &UpgradeTask{}, &DeviceUpgradeLog{})
}

func main() {
	initDB()

	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})

	// 固件管理API
	r.GET("/api/v1/ota/firmware", getFirmwareList)
	r.GET("/api/v1/ota/firmware/:id", getFirmware)
	r.POST("/api/v1/ota/firmware", createFirmware)
	r.PUT("/api/v1/ota/firmware/:id", updateFirmware)
	r.DELETE("/api/v1/ota/firmware/:id", deleteFirmware)
	r.POST("/api/v1/ota/firmware/:id/publish", publishFirmware)
	r.POST("/api/v1/ota/firmware/:id/archive", archiveFirmware)
	r.GET("/api/v1/ota/firmware/:id/devices", getFirmwareDevices)

	// 升级任务API
	r.GET("/api/v1/ota/tasks", getUpgradeTasks)
	r.POST("/api/v1/ota/tasks", createUpgradeTask)
	r.GET("/api/v1/ota/tasks/:id", getUpgradeTask)
	r.PUT("/api/v1/ota/tasks/:id", updateUpgradeTask)
	r.DELETE("/api/v1/ota/tasks/:id", deleteUpgradeTask)
	r.POST("/api/v1/ota/tasks/:id/start", startUpgradeTask)
	r.POST("/api/v1/ota/tasks/:id/rollback", rollbackUpgradeTask)
	r.GET("/api/v1/ota/tasks/:id/devices", getTaskDevices)

	// 设备升级日志API
	r.GET("/api/v1/ota/devices/:id/upgrade-logs", getDeviceUpgradeLogs)
	r.POST("/api/v1/ota/devices/:id/upgrade-status", updateDeviceUpgradeStatus)

	// 统计API
	r.GET("/api/v1/ota/statistics", getOTAStatistics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}
	log.Printf("ota-service starting on :%s", port)
	r.Run(":" + port)
}

// ========== 固件管理处理函数 ==========

func getFirmwareList(c *gin.Context) {
	var firmware []Firmware
	var total int64

	db.Model(&Firmware{}).Count(&total)
	db.Order("created_at DESC").Limit(20).Find(&firmware)

	c.JSON(200, gin.H{
		"code": 0,
		"data": firmware,
		"pagination": gin.H{
			"page":  c.DefaultQuery("page", "1"),
			"size":  c.DefaultQuery("size", "20"),
			"total": total,
		},
	})
}

func getFirmware(c *gin.Context) {
	var firmware Firmware
	if err := db.First(&firmware, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "固件不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": firmware})
}

func createFirmware(c *gin.Context) {
	var req CreateFirmwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	firmware := Firmware{
		ProductName:  req.ProductName,
		FirmwareName: req.FirmwareName,
		Version:      req.Version,
		FileURL:      req.FileURL,
		FileSize:     req.FileSize,
		Checksum:     req.Checksum,
		ReleaseNotes: req.ReleaseNotes,
		Status:       "draft",
		TenantID:     req.TenantID,
	}

	if err := db.Create(&firmware).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": firmware})
}

func updateFirmware(c *gin.Context) {
	var firmware Firmware
	if err := db.First(&firmware, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "固件不存在"})
		return
	}

	var req CreateFirmwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	firmware.FirmwareName = req.FirmwareName
	firmware.Version = req.Version
	firmware.FileURL = req.FileURL
	firmware.FileSize = req.FileSize
	firmware.Checksum = req.Checksum
	firmware.ReleaseNotes = req.ReleaseNotes

	db.Save(&firmware)
	c.JSON(200, gin.H{"code": 0, "data": firmware})
}

func deleteFirmware(c *gin.Context) {
	var firmware Firmware
	if err := db.First(&firmware, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "固件不存在"})
		return
	}

	db.Delete(&firmware)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func publishFirmware(c *gin.Context) {
	var firmware Firmware
	if err := db.First(&firmware, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "固件不存在"})
		return
	}

	firmware.Status = "active"
	db.Save(&firmware)
	c.JSON(200, gin.H{"code": 0, "message": "固件已发布"})
}

func archiveFirmware(c *gin.Context) {
	var firmware Firmware
	if err := db.First(&firmware, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "固件不存在"})
		return
	}

	firmware.Status = "archived"
	db.Save(&firmware)
	c.JSON(200, gin.H{"code": 0, "message": "固件已归档"})
}

func getFirmwareDevices(c *gin.Context) {
	var devices []struct {
		ID           uint   `json:"id"`
		DeviceName   string `json:"device_name"`
		FirmwareVer  string `json:"firmware_version"`
		NeedUpgrade  bool   `json:"need_upgrade"`
		LastUpgrade  string `json:"last_upgrade"`
	}

	c.JSON(200, gin.H{"code": 0, "data": devices})
}

// ========== 升级任务处理函数 ==========

func getUpgradeTasks(c *gin.Context) {
	var tasks []UpgradeTask
	var total int64

	db.Model(&UpgradeTask{}).Count(&total)
	db.Order("created_at DESC").Limit(20).Find(&tasks)

	c.JSON(200, gin.H{
		"code": 0,
		"data": tasks,
		"pagination": gin.H{
			"page":  c.DefaultQuery("page", "1"),
			"size":  c.DefaultQuery("size", "20"),
			"total": total,
		},
	})
}

func createUpgradeTask(c *gin.Context) {
	var req CreateUpgradeTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	task := UpgradeTask{
		FirmwareID:     req.FirmwareID,
		TaskName:       req.TaskName,
		TargetPercent:  req.TargetPercent,
		Description:    req.Description,
		Status:         "pending",
		TenantID:       req.TenantID,
	}

	if err := db.Create(&task).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": task})
}

func startUpgradeTask(c *gin.Context) {
	var task UpgradeTask
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "任务不存在"})
		return
	}

	now := time.Now()
	task.Status = "running"
	task.StartedAt = &now
	db.Save(&task)

	// 模拟升级进度
	time.Sleep(1 * time.Second)

	task.Status = "completed"
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	db.Save(&task)

	c.JSON(200, gin.H{"code": 0, "message": "升级任务已开始", "data": task})
}

func rollbackUpgradeTask(c *gin.Context) {
	var task UpgradeTask
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "任务不存在"})
		return
	}

	task.Status = "rollback"
	db.Save(&task)
	c.JSON(200, gin.H{"code": 0, "message": "升级回滚成功"})
}

func getUpgradeTask(c *gin.Context) {
	var task UpgradeTask
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "任务不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": task})
}

func updateUpgradeTask(c *gin.Context) {
	var task UpgradeTask
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "任务不存在"})
		return
	}

	var req CreateUpgradeTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	task.TaskName = req.TaskName
	task.TargetPercent = req.TargetPercent
	task.Description = req.Description

	db.Save(&task)
	c.JSON(200, gin.H{"code": 0, "data": task})
}

func deleteUpgradeTask(c *gin.Context) {
	var task UpgradeTask
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "任务不存在"})
		return
	}

	db.Delete(&task)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func getTaskDevices(c *gin.Context) {
	var devices []struct {
		ID           uint   `json:"id"`
		DeviceName   string `json:"device_name"`
		CurrentVer   string `json:"current_version"`
		TargetVer    string `json:"target_version"`
		Status       string `json:"status"`
		Progress     int    `json:"progress"`
	}

	c.JSON(200, gin.H{"code": 0, "data": devices})
}

// ========== 设备升级日志 ==========

func getDeviceUpgradeLogs(c *gin.Context) {
	var logs []DeviceUpgradeLog
	var total int64

	db.Model(&DeviceUpgradeLog{}).Where("device_id = ?", c.Param("id")).Count(&total)
	db.Where("device_id = ?", c.Param("id")).Order("created_at DESC").Limit(20).Find(&logs)

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

func updateDeviceUpgradeStatus(c *gin.Context) {
	var req struct {
		DeviceID      uint   `json:"device_id"`
		UpgradeTaskID uint   `json:"upgrade_task_id"`
		FirmwareID    uint   `json:"firmware_id"`
		FromVersion   string `json:"from_version"`
		ToVersion     string `json:"to_version"`
		Status        string `json:"status"`
		Progress      int    `json:"progress"`
		ErrorMsg      string `json:"error_msg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	log := DeviceUpgradeLog{
		DeviceID:      req.DeviceID,
		UpgradeTaskID: req.UpgradeTaskID,
		FirmwareID:    req.FirmwareID,
		FromVersion:   req.FromVersion,
		ToVersion:     req.ToVersion,
		Status:        req.Status,
		Progress:      req.Progress,
		ErrorMsg:      req.ErrorMsg,
	}

	db.Create(&log)
	c.JSON(200, gin.H{"code": 0, "message": "状态更新成功"})
}

// ========== OTA统计 ==========

func getOTAStatistics(c *gin.Context) {
	var firmwareCount, activeCount, draftCount int64
	var taskCount, completedCount, failedCount int64

	db.Model(&Firmware{}).Count(&firmwareCount)
	db.Model(&Firmware{}).Where("status = 'active'").Count(&activeCount)
	db.Model(&Firmware{}).Where("status = 'draft'").Count(&draftCount)
	db.Model(&UpgradeTask{}).Count(&taskCount)
	db.Model(&UpgradeTask{}).Where("status = 'completed'").Count(&completedCount)
	db.Model(&UpgradeTask{}).Where("status = 'failed'").Count(&failedCount)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"firmwareCount":    firmwareCount,
			"activeCount":      activeCount,
			"draftCount":       draftCount,
			"taskCount":        taskCount,
			"completedCount":   completedCount,
			"failedCount":      failedCount,
		},
	})
}
