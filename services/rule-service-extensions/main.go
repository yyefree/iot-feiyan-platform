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

type SceneDeviceTrigger struct {
	gorm.Model
	SceneID      uint   `json:"scene_id" gorm:"index"`
	DeviceID     uint   `json:"device_id" gorm:"index"`
	PropertyID   uint   `json:"property_id"`
	Condition    string `json:"condition" gorm:"not null"` // eq/neq/lt/lte/gt/gte
	TargetValue  string `json:"target_value"`
	ActionConfig string `json:"action_config" gorm:"type:text"`
	Status       string `json:"status" gorm:"default:'active'"`
	CreatedAt    time.Time `json:"created_at"`
}

type DeviceLocation struct {
	gorm.Model
	DeviceID  uint    `json:"device_id" gorm:"index"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
	Region    string  `json:"region"`
	TenantID  uint    `json:"tenant_id" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========== 请求结构 ==========

type CreateDeviceTriggerRequest struct {
	SceneID      uint   `json:"scene_id" binding:"required"`
	DeviceID     uint   `json:"device_id" binding:"required"`
	PropertyID   uint   `json:"property_id"`
	Condition    string `json:"condition" binding:"required"`
	TargetValue  string `json:"target_value"`
	ActionConfig string `json:"action_config"`
}

type UpdateDeviceLocationRequest struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Location  string   `json:"location"`
	Region    string   `json:"region"`
}

// ========== 服务初始化 ==========

var db *gorm.DB

func initDB() {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "iot_admin")
	password := getEnv("POSTGRES_PASSWORD", "iot_admin_2024")
	dbname := getEnv("POSTGRES_DB", "iot_platform")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("WARNING: failed to connect database: %v", err)
		// 尝试docker容器名
		dsn = fmt.Sprintf("host=iot-postgres port=%s user=%s password=%s dbname=%s sslmode=disable", port, user, password, dbname)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database:", err)
		}
	}

	db.AutoMigrate(&SceneDeviceTrigger{}, &DeviceLocation{})
}

func main() {
	initDB()

	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})

	// 场景设备触发API
	r.GET("/api/v1/scenes/:id/triggers", getSceneTriggers)
	r.POST("/api/v1/scenes/:id/triggers", createSceneTrigger)
	r.PUT("/api/v1/scenes/triggers/:id", updateSceneTrigger)
	r.DELETE("/api/v1/scenes/triggers/:id", deleteSceneTrigger)
	r.POST("/api/v1/scenes/triggers/:id/enable", enableSceneTrigger)
	r.POST("/api/v1/scenes/triggers/:id/disable", disableSceneTrigger)
	r.GET("/api/v1/scenes/triggers/statistics", getTriggerStatistics)

	// 设备地理位置API
	r.GET("/api/v1/devices/:id/location", getDeviceLocation)
	r.PUT("/api/v1/devices/:id/location", updateDeviceLocation)
	r.GET("/api/v1/devices/locations", getDeviceLocations)
	r.GET("/api/v1/devices/locations/geography", getGeographyDistribution)

	// 运行在8094端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8094"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("rule-service extensions starting on %s", addr)
	
	if err := r.Run(addr); err != nil {
		log.Printf("Failed to start on %s, trying :8095", addr)
		if err := r.Run(":8095"); err != nil {
			log.Fatal("failed to start server:", err)
		}
	}
}

// ========== 场景设备触发处理函数 ==========

func getSceneTriggers(c *gin.Context) {
	var triggers []SceneDeviceTrigger
	var total int64

	db.Model(&SceneDeviceTrigger{}).Where("scene_id = ?", c.Param("id")).Count(&total)
	db.Where("scene_id = ?", c.Param("id")).Order("created_at DESC").Limit(20).Find(&triggers)

	c.JSON(200, gin.H{
		"code": 0,
		"data": triggers,
		"pagination": gin.H{
			"page":  c.DefaultQuery("page", "1"),
			"size":  c.DefaultQuery("size", "20"),
			"total": total,
		},
	})
}

func createSceneTrigger(c *gin.Context) {
	var req CreateDeviceTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	trigger := SceneDeviceTrigger{
		SceneID:      req.SceneID,
		DeviceID:     req.DeviceID,
		PropertyID:   req.PropertyID,
		Condition:    req.Condition,
		TargetValue:  req.TargetValue,
		ActionConfig: req.ActionConfig,
		Status:       "active",
	}

	if err := db.Create(&trigger).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": trigger})
}

func updateSceneTrigger(c *gin.Context) {
	var trigger SceneDeviceTrigger
	if err := db.First(&trigger, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "触发规则不存在"})
		return
	}

	var req CreateDeviceTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	trigger.DeviceID = req.DeviceID
	trigger.PropertyID = req.PropertyID
	trigger.Condition = req.Condition
	trigger.TargetValue = req.TargetValue
	trigger.ActionConfig = req.ActionConfig

	db.Save(&trigger)
	c.JSON(200, gin.H{"code": 0, "data": trigger})
}

func deleteSceneTrigger(c *gin.Context) {
	var trigger SceneDeviceTrigger
	if err := db.First(&trigger, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "触发规则不存在"})
		return
	}

	db.Delete(&trigger)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func enableSceneTrigger(c *gin.Context) {
	var trigger SceneDeviceTrigger
	if err := db.First(&trigger, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "触发规则不存在"})
		return
	}

	trigger.Status = "active"
	db.Save(&trigger)
	c.JSON(200, gin.H{"code": 0, "message": "已启用"})
}

func disableSceneTrigger(c *gin.Context) {
	var trigger SceneDeviceTrigger
	if err := db.First(&trigger, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "触发规则不存在"})
		return
	}

	trigger.Status = "inactive"
	db.Save(&trigger)
	c.JSON(200, gin.H{"code": 0, "message": "已禁用"})
}

func getTriggerStatistics(c *gin.Context) {
	var activeTriggers, inactiveTriggers, totalTriggers int64

	db.Model(&SceneDeviceTrigger{}).Where("status = 'active'").Count(&activeTriggers)
	db.Model(&SceneDeviceTrigger{}).Where("status = 'inactive'").Count(&inactiveTriggers)
	db.Model(&SceneDeviceTrigger{}).Count(&totalTriggers)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"activeTriggers":   activeTriggers,
			"inactiveTriggers": inactiveTriggers,
			"totalTriggers":    totalTriggers,
		},
	})
}

// ========== 设备地理位置处理函数 ==========

func getDeviceLocation(c *gin.Context) {
	var location DeviceLocation
	if err := db.Where("device_id = ?", c.Param("id")).First(&location).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "位置信息不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": location})
}

func updateDeviceLocation(c *gin.Context) {
	var req UpdateDeviceLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var location DeviceLocation
	if err := db.Where("device_id = ?", c.Param("id")).First(&location).Error; err != nil {
		// 创建新位置
		location = DeviceLocation{
			DeviceID:  uintStrToUint(c.Param("id")),
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Location:  req.Location,
			Region:    req.Region,
		}
		if err := db.Create(&location).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "message": err.Error()})
			return
		}
	} else {
		location.Latitude = req.Latitude
		location.Longitude = req.Longitude
		location.Location = req.Location
		location.Region = req.Region
		db.Save(&location)
	}

	c.JSON(200, gin.H{"code": 0, "data": location})
}

func getDeviceLocations(c *gin.Context) {
	var locations []DeviceLocation
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	offset := (page - 1) * pageSize

	db.Model(&DeviceLocation{}).Count(&total)
	db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&locations)

	c.JSON(200, gin.H{
		"code": 0,
		"data": locations,
		"pagination": gin.H{
			"page":  page,
			"size":  pageSize,
			"total": total,
		},
	})
}

func getGeographyDistribution(c *gin.Context) {
	var locations []struct {
		Region string  `json:"region"`
		Count  int64   `json:"count"`
		Lat    float64 `json:"latitude"`
		Lng    float64 `json:"longitude"`
	}

	// 获取每个地区的统计和示例坐标
	db.Raw(`
		SELECT 
			region, 
			COUNT(*) as count,
			AVG(latitude) as latitude,
			AVG(longitude) as longitude
		FROM device_locations 
		WHERE region IS NOT NULL AND region != '' 
		GROUP BY region 
		ORDER BY count DESC
	`).Scan(&locations)

	c.JSON(200, gin.H{"code": 0, "data": locations})
}

// ========== 工具函数 ==========

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func uintStrToUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint(n)
}
