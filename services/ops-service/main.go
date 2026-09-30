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

type OperationLog struct {
	gorm.Model
	TenantID     uint      `json:"tenant_id" gorm:"index"`
	UserID       uint      `json:"user_id"`
	Operation    string    `json:"operation" gorm:"not null"`
	Module       string    `json:"module"`
	TargetID     uint      `json:"target_id"`
	TargetType   string    `json:"target_type"`
	RequestData  string    `json:"request_data" gorm:"type:text"`
	ResponseData string    `json:"response_data" gorm:"type:text"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type MessagePush struct {
	gorm.Model
	TenantID uint      `json:"tenant_id" gorm:"index"`
	UserID   uint      `json:"user_id"`
	DeviceID uint      `json:"device_id"`
	MsgType  string    `json:"msg_type" gorm:"not null"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Params   string    `json:"params" gorm:"type:text"`
	Status   string    `json:"status" gorm:"default:'pending'"`
	SentAt   *time.Time `json:"sent_at"`
}

type DeviceStat struct {
	gorm.Model
	TenantID     uint      `json:"tenant_id"`
	DeviceID     uint      `json:"device_id"`
	Identify     string    `json:"identify"`
	StatPeriod   string    `json:"stat_period"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Sum          float64   `json:"sum"`
	Avg          float64   `json:"avg"`
	Min          float64   `json:"min"`
	Max          float64   `json:"max"`
	Count        int64     `json:"count"`
	ComputedAt   time.Time `json:"computed_at"`
}

type Alert struct {
	gorm.Model
	TenantID   uint      `json:"tenant_id" gorm:"index"`
	DeviceID   uint      `json:"device_id"`
	AlertType  string    `json:"alert_type"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Severity   string    `json:"severity"` // low/medium/high/critical
	Status     string    `json:"status"`   // pending/acknowledged/resolved
	HandledBy  uint      `json:"handled_by"`
	HandledAt  *time.Time `json:"handled_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type DashboardConfig struct {
	gorm.Model
	TenantID uint   `json:"tenant_id" gorm:"index"`
	Name     string `json:"name"`
	Config   string `json:"config" gorm:"type:text"`
	CreatedBy uint  `json:"created_by"`
}

// ========== 请求结构 ==========

type CreateOperationLogRequest struct {
	TenantID     uint   `json:"tenant_id" binding:"required"`
	UserID       uint   `json:"user_id"`
	Operation    string `json:"operation" binding:"required"`
	Module       string `json:"module"`
	TargetID     uint   `json:"target_id"`
	TargetType   string `json:"target_type"`
	RequestData  string `json:"request_data"`
	ResponseData string `json:"response_data"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
}

type SendPushRequest struct {
	TenantID uint   `json:"tenant_id" binding:"required"`
	UserID   uint   `json:"user_id"`
	DeviceID uint   `json:"device_id"`
	MsgType  string `json:"msg_type" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Params   string `json:"params"`
}

type CreateAlertRequest struct {
	TenantID uint   `json:"tenant_id" binding:"required"`
	DeviceID uint   `json:"device_id"`
	AlertType string `json:"alert_type" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Severity string `json:"severity"`
}

type HandleAlertRequest struct {
	HandledBy uint   `json:"handled_by"`
	Status    string `json:"status"`
	Note      string `json:"note"`
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
	
	db.AutoMigrate(&OperationLog{}, &MessagePush{}, &DeviceStat{}, &Alert{}, &DashboardConfig{})
}

func main() {
	initDB()
	
	r := gin.Default()
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})
	
	// 运营中心API
	r.GET("/api/v1/ops/dashboard", getDashboard)
	r.GET("/api/v1/ops/devices", getDeviceStats)
	r.GET("/api/v1/ops/users", getUserStats)
	r.GET("/api/v1/ops/alerts", getAlerts)
	r.PUT("/api/v1/ops/alerts/:id/handle", handleAlert)
	r.GET("/api/v1/ops/logs", getOperationLogs)
	r.GET("/api/v1/ops/dashboard/config", getDashboardConfig)
	r.PUT("/api/v1/ops/dashboard/config", updateDashboardConfig)
	
	// 消息推送API
	r.POST("/api/v1/push/app", sendAppPush)
	r.POST("/api/v1/push/wechat", sendWechatPush)
	r.GET("/api/v1/push/templates", getPushTemplates)
	r.POST("/api/v1/push/config", updatePushConfig)
	
	log.Printf("ops-service starting on :8085")
	r.Run(":8085")
}

// ========== 运营中心处理函数 ==========

func getDashboard(c *gin.Context) {
	// 获取dashboard数据
	var totalDevices, onlineDevices, offlineDevices int64
	var totalDataPoints, last24hDataPoints int64
	var totalAlerts, pendingAlerts int64
	
	db.Raw("SELECT COUNT(*) FROM devices").Scan(&totalDevices)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = true").Scan(&onlineDevices)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = false").Scan(&offlineDevices)
	db.Raw("SELECT COUNT(*) FROM telemetry").Scan(&totalDataPoints)
	db.Raw("SELECT COUNT(*) FROM telemetry WHERE timestamp >= NOW() - INTERVAL '24 hours'").Scan(&last24hDataPoints)
	db.Raw("SELECT COUNT(*) FROM alerts").Scan(&totalAlerts)
	db.Raw("SELECT COUNT(*) FROM alerts WHERE status = 'pending'").Scan(&pendingAlerts)
	
	data := gin.H{
		"deviceStats": gin.H{
			"total":      totalDevices,
			"online":     onlineDevices,
			"offline":    offlineDevices,
		},
		"dataStats": gin.H{
			"total":          totalDataPoints,
			"last24h":        last24hDataPoints,
		},
		"alertStats": gin.H{
			"total":    totalAlerts,
			"pending":  pendingAlerts,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	c.JSON(200, gin.H{"code": 0, "data": data})
}

func getDeviceStats(c *gin.Context) {
	var stats []struct {
		TenantID     uint   `json:"tenant_id"`
		TotalDevices int64  `json:"total_devices"`
		OnlineDevices int64 `json:"online_devices"`
		OfflineDevices int64 `json:"offline_devices"`
	}
	
	db.Raw(`
		SELECT 
			tenant_id,
			COUNT(*) as total_devices,
			COUNT(CASE WHEN online = true THEN 1 END) as online_devices,
			COUNT(CASE WHEN online = false THEN 1 END) as offline_devices
		FROM devices
		GROUP BY tenant_id
	`).Scan(&stats)
	
	c.JSON(200, gin.H{"code": 0, "data": stats})
}

func getUserStats(c *gin.Context) {
	var stats []struct {
		TenantID    uint   `json:"tenant_id"`
		TotalUsers  int64  `json:"total_users"`
		ActiveUsers int64  `json:"active_users"`
	}
	
	db.Raw(`
		SELECT 
			tenant_id,
			COUNT(*) as total_users,
			COUNT(CASE WHEN last_login_at >= NOW() - INTERVAL '7 days' THEN 1 END) as active_users
		FROM "user"
		GROUP BY tenant_id
	`).Scan(&stats)
	
	c.JSON(200, gin.H{"code": 0, "data": stats})
}

func getAlerts(c *gin.Context) {
	var alerts []Alert
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&Alert{}).Count(&total)
	db.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&alerts)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": alerts,
		"pagination": gin.H{
			"page":  page,
			"size":  pageSize,
			"total": total,
		},
	})
}

func handleAlert(c *gin.Context) {
	var alert Alert
	if err := db.First(&alert, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "告警不存在"})
		return
	}
	
	var req HandleAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	alert.Status = req.Status
	alert.HandledBy = req.HandledBy
	now := time.Now()
	alert.HandledAt = &now
	db.Save(&alert)
	
	c.JSON(200, gin.H{"code": 0, "message": "告警已处理"})
}

func getOperationLogs(c *gin.Context) {
	var logs []OperationLog
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	query := db.Model(&OperationLog{})
	
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		if tid, err := strconv.ParseUint(tenantID, 10, 64); err == nil {
			query = query.Where("tenant_id = ?", tid)
		}
	}
	
	if userID := c.Query("user_id"); userID != "" {
		if uid, err := strconv.ParseUint(userID, 10, 64); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}
	
	if module := c.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	
	if startTime := c.Query("start_time"); startTime != "" {
		if start, err := time.Parse(time.RFC3339, startTime); err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}
	
	if endTime := c.Query("end_time"); endTime != "" {
		if end, err := time.Parse(time.RFC3339, endTime); err == nil {
			query = query.Where("created_at <= ?", end)
		}
	}
	
	query.Count(&total)
	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs)
	
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

func getDashboardConfig(c *gin.Context) {
	var config DashboardConfig
	result := db.Where("tenant_id = ?", 1).First(&config)
	
	if result.Error == gorm.ErrRecordNotFound {
		config = DashboardConfig{
			TenantID: 1,
			Name:     "默认大屏",
			Config:   "{}",
		}
		db.Create(&config)
	}
	
	c.JSON(200, gin.H{"code": 0, "data": config})
}

func updateDashboardConfig(c *gin.Context) {
	var config DashboardConfig
	result := db.Where("tenant_id = ?", 1).First(&config)
	
	if result.Error == gorm.ErrRecordNotFound {
		config = DashboardConfig{
			TenantID: 1,
			Name:     "默认大屏",
			Config:   "{}",
		}
	}
	
	var req struct {
		Name     string `json:"name"`
		Config   string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	config.Name = req.Name
	config.Config = req.Config
	db.Save(&config)
	
	c.JSON(200, gin.H{"code": 0, "message": "配置更新成功"})
}

// ========== 消息推送API ==========

func sendAppPush(c *gin.Context) {
	var req SendPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	push := MessagePush{
		TenantID: req.TenantID,
		UserID:   req.UserID,
		DeviceID: req.DeviceID,
		MsgType:  "app",
		Title:    req.Title,
		Content:  req.Content,
		Params:   req.Params,
		Status:   "sent",
	}
	
	if err := db.Create(&push).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"code": 0, "message": "推送发送成功"})
}

func sendWechatPush(c *gin.Context) {
	var req SendPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	push := MessagePush{
		TenantID: req.TenantID,
		UserID:   req.UserID,
		DeviceID: req.DeviceID,
		MsgType:  "wechat",
		Title:    req.Title,
		Content:  req.Content,
		Params:   req.Params,
		Status:   "sent",
	}
	
	if err := db.Create(&push).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"code": 0, "message": "微信推送发送成功"})
}

func getPushTemplates(c *gin.Context) {
	templates := []map[string]string{
		{"id": "device_alert", "name": "设备告警", "content": "设备{{device_name}}发生告警：{{alert_content}}"},
		{"id": "device_online", "name": "设备上线", "content": "设备{{device_name}}已上线"},
		{"id": "device_offline", "name": "设备离线", "content": "设备{{device_name}}已离线"},
		{"id": "ota_complete", "name": "升级完成", "content": "设备{{device_name}}固件升级完成"},
	}
	
	c.JSON(200, gin.H{"code": 0, "data": templates})
}

func updatePushConfig(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "推送配置更新成功"})
}

// ========== 工具函数 ==========

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
