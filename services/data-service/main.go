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

type Telemetry struct {
	gorm.Model
	DeviceID   uint      `json:"device_id" gorm:"index"`
	TenantID   uint      `json:"tenant_id" gorm:"index"`
	ProductID  uint      `json:"product_id" gorm:"index"`
	Identify   string    `json:"identify" gorm:"index"`
	ValueType  string    `json:"value_type"`
	ValueFloat float64   `json:"value_float"`
	ValueInt   int64     `json:"value_int"`
	ValueStr   string    `json:"value_string"`
	ValueBool  bool      `json:"value_bool"`
	Timestamp  time.Time `json:"timestamp" gorm:"index"`
}

type DeviceLatestProperty struct {
	gorm.Model
	DeviceID   uint      `json:"device_id" gorm:"uniqueIndex;index"`
	Identify   string    `json:"identify" gorm:"uniqueIndex;index"`
	ValueFloat float64   `json:"value_float"`
	ValueInt   int64     `json:"value_int"`
	ValueStr   string    `json:"value_string"`
	ValueBool  bool      `json:"value_bool"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DataStatistics struct {
	gorm.Model
	TenantID   uint      `json:"tenant_id"`
	DeviceID   uint      `json:"device_id"`
	Identify   string    `json:"identify"`
	StatPeriod string    `json:"stat_period"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Sum        float64   `json:"sum"`
	Avg        float64   `json:"avg"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
	Count      int64     `json:"count"`
	ComputedAt time.Time `json:"computed_at"`
}

type ChartDataPoint struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// ========== 请求结构 ==========

type ReportTelemetryRequest struct {
	DeviceID   uint    `json:"device_id" binding:"required"`
	Identify   string  `json:"identify" binding:"required"`
	ValueFloat float64 `json:"value_float"`
	ValueInt   int64   `json:"value_int"`
	ValueStr   string  `json:"value_string"`
	ValueBool  bool    `json:"value_bool"`
	Timestamp  string  `json:"timestamp"`
}

type QueryTelemetryRequest struct {
	DeviceID  uint   `form:"device_id" binding:"required"`
	Identify  string `form:"identify"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Limit     int    `form:"limit" default:"1000"`
}

type ChartRequest struct {
	DeviceID  uint   `form:"device_id" binding:"required"`
	Identify  string `form:"identify"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Interval  string `form:"interval" default:"hour"` // minute/hour/day
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
	
	db.AutoMigrate(&Telemetry{}, &DeviceLatestProperty{}, &DataStatistics{})
}

func main() {
	initDB()
	
	r := gin.Default()
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})
	
	// 遥测数据API
	r.POST("/api/v1/data/telemetry", reportTelemetry)
	r.POST("/api/v1/data/telemetry/batch", batchReportTelemetry)
	r.GET("/api/v1/data/telemetry", queryTelemetry)
	r.GET("/api/v1/data/telemetry/device/:id", getDeviceTelemetry)
	r.GET("/api/v1/data/telemetry/device/:id/latest", getLatestTelemetry)
	r.GET("/api/v1/data/telemetry/device/:id/stats", getDeviceStatistics)
	
	// 设备最新属性
	r.GET("/api/v1/data/properties", getLatestProperties)
	r.GET("/api/v1/data/properties/device/:id", getDeviceLatestProperties)
	
	// 数据统计
	r.GET("/api/v1/data/statistics", getDataStatistics)
	r.GET("/api/v1/data/statistics/overview", getStatisticsOverview)
	
	// 图表数据API
	r.GET("/api/v1/data/charts/line", getLineChart)
	r.GET("/api/v1/data/charts/bar", getBarChart)
	r.GET("/api/v1/data/charts/pie", getPieChart)
	r.GET("/api/v1/data/charts/radar", getRadarChart)
	
	// 数据导出API
	r.POST("/api/v1/data/export", exportData)
	
	log.Printf("data-service starting on :8084")
	r.Run(":8084")
}

// ========== 遥测数据处理函数 ==========

func reportTelemetry(c *gin.Context) {
	var req ReportTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	timestamp := time.Now()
	if req.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, req.Timestamp); err == nil {
			timestamp = t
		}
	}
	
	// 确定值类型
	valueType := "float"
	if req.ValueBool {
		valueType = "bool"
	} else if req.ValueInt != 0 {
		valueType = "int"
	} else if req.ValueStr != "" {
		valueType = "string"
	}
	
	telemetry := Telemetry{
		DeviceID:   req.DeviceID,
		Identify:   req.Identify,
		ValueType:  valueType,
		ValueFloat: req.ValueFloat,
		ValueInt:   req.ValueInt,
		ValueStr:   req.ValueStr,
		ValueBool:  req.ValueBool,
		Timestamp:  timestamp,
	}
	
	if err := db.Create(&telemetry).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	// 更新最新属性
	updateLatestProperty(&telemetry)
	
	c.JSON(201, gin.H{"code": 0, "message": "上报成功", "data": telemetry})
}

func batchReportTelemetry(c *gin.Context) {
	var reqs []ReportTelemetryRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var telemetries []Telemetry
	for _, req := range reqs {
		timestamp := time.Now()
		if req.Timestamp != "" {
			if t, err := time.Parse(time.RFC3339, req.Timestamp); err == nil {
				timestamp = t
			}
		}
		
		valueType := "float"
		if req.ValueBool {
			valueType = "bool"
		} else if req.ValueInt != 0 {
			valueType = "int"
		} else if req.ValueStr != "" {
			valueType = "string"
		}
		
		telemetry := Telemetry{
			DeviceID:   req.DeviceID,
			Identify:   req.Identify,
			ValueType:  valueType,
			ValueFloat: req.ValueFloat,
			ValueInt:   req.ValueInt,
			ValueStr:   req.ValueStr,
			ValueBool:  req.ValueBool,
			Timestamp:  timestamp,
		}
		telemetries = append(telemetries, telemetry)
		updateLatestProperty(&telemetry)
	}
	
	if err := db.Create(&telemetries).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "message": fmt.Sprintf("成功上报 %d 条数据", len(telemetries))})
}

func queryTelemetry(c *gin.Context) {
	var req QueryTelemetryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var telemetries []Telemetry
	query := db.Model(&Telemetry{}).Where("device_id = ?", req.DeviceID)
	
	if req.Identify != "" {
		query = query.Where("identify = ?", req.Identify)
	}
	
	if req.StartTime != "" {
		if start, err := time.Parse(time.RFC3339, req.StartTime); err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	
	if req.EndTime != "" {
		if end, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}
	
	query.Order("timestamp ASC").Limit(req.Limit).Find(&telemetries)
	
	c.JSON(200, gin.H{"code": 0, "data": telemetries})
}

func getDeviceTelemetry(c *gin.Context) {
	deviceID := c.Param("id")
	
	var telemetries []Telemetry
	db.Where("device_id = ?", deviceID).Order("timestamp DESC").Limit(100).Find(&telemetries)
	
	c.JSON(200, gin.H{"code": 0, "data": telemetries})
}

func getLatestTelemetry(c *gin.Context) {
	deviceID := c.Param("id")
	
	var properties []DeviceLatestProperty
	db.Where("device_id = ?", deviceID).Find(&properties)
	
	c.JSON(200, gin.H{"code": 0, "data": properties})
}

func getDeviceStatistics(c *gin.Context) {
	deviceID := c.Param("id")
	
	var stats []DataStatistics
	db.Where("device_id = ?", deviceID).Order("start_time DESC").Limit(10).Find(&stats)
	
	c.JSON(200, gin.H{"code": 0, "data": stats})
}

// ========== 最新属性API ==========

func getLatestProperties(c *gin.Context) {
	var properties []DeviceLatestProperty
	db.Find(&properties)
	
	c.JSON(200, gin.H{"code": 0, "data": properties})
}

func getDeviceLatestProperties(c *gin.Context) {
	deviceID := c.Param("id")
	
	var properties []DeviceLatestProperty
	db.Where("device_id = ?", deviceID).Find(&properties)
	
	c.JSON(200, gin.H{"code": 0, "data": properties})
}

// ========== 数据统计API ==========

func getDataStatistics(c *gin.Context) {
	var stats []DataStatistics
	db.Order("start_time DESC").Limit(100).Find(&stats)
	
	c.JSON(200, gin.H{"code": 0, "data": stats})
}

func getStatisticsOverview(c *gin.Context) {
	var totalDevices, onlineDevices, offlineDevices int64
	var totalDataPoints, last24hDataPoints, totalIdentifiers int64
	
	db.Raw("SELECT COUNT(*) FROM devices").Scan(&totalDevices)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = true").Scan(&onlineDevices)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = false").Scan(&offlineDevices)
	db.Raw("SELECT COUNT(*) FROM telemetry").Scan(&totalDataPoints)
	db.Raw("SELECT COUNT(*) FROM telemetry WHERE timestamp >= NOW() - INTERVAL '24 hours'").Scan(&last24hDataPoints)
	db.Raw("SELECT COUNT(DISTINCT identify) FROM telemetry").Scan(&totalIdentifiers)
	
	overview := gin.H{
		"totalDevices":      totalDevices,
		"onlineDevices":     onlineDevices,
		"offlineDevices":    offlineDevices,
		"totalDataPoints":   totalDataPoints,
		"last24hDataPoints": last24hDataPoints,
		"totalIdentifiers":  totalIdentifiers,
	}
	
	c.JSON(200, gin.H{"code": 0, "data": overview})
}

// ========== 图表数据API ==========

func getLineChart(c *gin.Context) {
	var req ChartRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var telemetries []Telemetry
	query := db.Where("device_id = ?", req.DeviceID)
	
	if req.Identify != "" {
		query = query.Where("identify = ?", req.Identify)
	}
	
	if req.StartTime != "" {
		if start, err := time.Parse(time.RFC3339, req.StartTime); err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	
	if req.EndTime != "" {
		if end, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}
	
	query.Order("timestamp ASC").Limit(1000).Find(&telemetries)
	
	points := make([]ChartDataPoint, 0, len(telemetries))
	for _, t := range telemetries {
		points = append(points, ChartDataPoint{
			Time:  t.Timestamp.Format("2006-01-02 15:04:05"),
			Value: t.ValueFloat,
		})
	}
	
	c.JSON(200, gin.H{"code": 0, "data": points})
}

func getBarChart(c *gin.Context) {
	var req ChartRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var stats []DataStatistics
	query := db.Where("device_id = ?", req.DeviceID)
	
	if req.Identify != "" {
		query = query.Where("identify = ?", req.Identify)
	}
	
	if req.StartTime != "" {
		if start, err := time.Parse(time.RFC3339, req.StartTime); err == nil {
			query = query.Where("start_time >= ?", start)
		}
	}
	
	if req.EndTime != "" {
		if end, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			query = query.Where("end_time <= ?", end)
		}
	}
	
	query.Order("start_time ASC").Limit(24).Find(&stats)
	
	points := make([]ChartDataPoint, 0, len(stats))
	for _, s := range stats {
		points = append(points, ChartDataPoint{
			Time:  s.StartTime.Format("2006-01-02 15:04"),
			Value: s.Avg,
		})
	}
	
	c.JSON(200, gin.H{"code": 0, "data": points})
}

func getPieChart(c *gin.Context) {
	// 设备状态分布
	type PiePoint struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	
	var onlineCount, offlineCount int64
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = true").Scan(&onlineCount)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = false").Scan(&offlineCount)
	
	points := []PiePoint{
		{"在线", float64(onlineCount)},
		{"离线", float64(offlineCount)},
	}
	
	c.JSON(200, gin.H{"code": 0, "data": points})
}

func getRadarChart(c *gin.Context) {
	// 设备性能指标
	type RadarPoint struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	
	var totalDevices, onlineDevices int64
	db.Raw("SELECT COUNT(*) FROM devices").Scan(&totalDevices)
	db.Raw("SELECT COUNT(*) FROM devices WHERE online = true").Scan(&onlineDevices)
	
	onlineRate := float64(0)
	if totalDevices > 0 {
		onlineRate = float64(onlineDevices) / float64(totalDevices) * 100
	}
	
	points := []RadarPoint{
		{"在线率", onlineRate},
		{"数据上报率", 85.0},
		{"指令响应率", 90.0},
		{"故障率", 5.0},
		{"用户满意度", 92.0},
	}
	
	c.JSON(200, gin.H{"code": 0, "data": points})
}

// ========== 数据导出API ==========

func exportData(c *gin.Context) {
	var req QueryTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	var telemetries []Telemetry
	query := db.Where("device_id = ?", req.DeviceID)
	
	if req.StartTime != "" {
		if start, err := time.Parse(time.RFC3339, req.StartTime); err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	
	if req.EndTime != "" {
		if end, err := time.Parse(time.RFC3339, req.EndTime); err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}
	
	query.Order("timestamp ASC").Limit(10000).Find(&telemetries)
	
	// 生成CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=telemetry_export.csv")
	
	c.String(200, "time,device_id,identify,value_type,value_float,value_int,value_string\n")
	for _, t := range telemetries {
		c.String(200, fmt.Sprintf("%s,%d,%s,%s,%f,%d,%s\n", 
			t.Timestamp.Format(time.RFC3339),
			t.DeviceID,
			t.Identify,
			t.ValueType,
			t.ValueFloat,
			t.ValueInt,
			t.ValueStr))
	}
}

// ========== 辅助函数 ==========

func updateLatestProperty(telemetry *Telemetry) {
	var prop DeviceLatestProperty
	result := db.Where("device_id = ? AND identify = ?", telemetry.DeviceID, telemetry.Identify).First(&prop)
	
	prop.DeviceID = telemetry.DeviceID
	prop.Identify = telemetry.Identify
	prop.ValueFloat = telemetry.ValueFloat
	prop.ValueInt = telemetry.ValueInt
	prop.ValueStr = telemetry.ValueStr
	prop.ValueBool = telemetry.ValueBool
	prop.UpdatedAt = time.Now()
	
	if result.Error == gorm.ErrRecordNotFound {
		db.Create(&prop)
	} else {
		db.Save(&prop)
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func roundTo2Decimals(num float64) float64 {
	return float64(int(num*100)) / 100
}
