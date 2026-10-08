package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ========== 数据模型 ==========

type Product struct {
	gorm.Model
	ProductName  string `json:"product_name" gorm:"not null;index"`
	ProductKey   string `json:"product_key" gorm:"uniqueIndex;not null"`
	TenantID     uint   `json:"tenant_id" gorm:"index"`
	ProjectID    uint   `json:"project_id"`
	CategoryID   uint   `json:"category_id"`
	ProductType  string `json:"product_type" gorm:"default:'direct_device'"`
	CommType     string `json:"comm_type" gorm:"default:'wifi'"`
	AuthType     string `json:"auth_type" gorm:"default:'one_device_one_key'"`
	Status       string `json:"status" gorm:"default:'draft'"`
	Icon         string `json:"icon"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	ThingModel   string `json:"thing_model" gorm:"type:text"`
}

type ThingModelDefine struct {
	gorm.Model
	ProductID    uint   `json:"product_id" gorm:"index"`
	DefineType   string `json:"define_type" gorm:"index"`
	Identifier   string `json:"identifier" gorm:"uniqueIndex:idx_product_identify;not null"`
	Name         string `json:"name" gorm:"not null"`
	Desc         string `json:"desc"`
	AccessMode   string `json:"access_mode"`
	Type         string `json:"type" gorm:"not null"`
	Unit         string `json:"unit"`
	UnitSymbol   string `json:"unit_symbol"`
	MinValue     *float64 `json:"min_value"`
	MaxValue     *float64 `json:"max_value"`
	Step         *float64 `json:"step"`
	EnumValues   string   `json:"enum_values" gorm:"type:text"`
	InputParams  string   `json:"input_params" gorm:"type:text"`
	OutputParams string   `json:"output_params" gorm:"type:text"`
	Level        string   `json:"level"`
}

type TSLVersion struct {
	gorm.Model
	ProductID uint      `json:"product_id" gorm:"index"`
	Version   string    `json:"version" gorm:"not null"`
	TSLData   string    `json:"tsl_data" gorm:"type:text"`
	ChangeLog string    `json:"change_log"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ========== 请求结构 ==========

type CreateProductRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	TenantID    uint   `json:"tenant_id"`
	ProjectID   uint   `json:"project_id"`
	CategoryID  uint   `json:"category_id"`
	CommType    string `json:"comm_type"`
	AuthType    string `json:"auth_type"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ThingModel  string `json:"thing_model"`
}

type ImportTSLRequest struct {
	TSLData string `json:"tsl_data" binding:"required"`
}

type ExportTSLResponse struct {
	Schema     string      `json:"schema"`
	Profile    ProfileInfo `json:"profile"`
	Properties []Property  `json:"properties"`
	Services   []Service   `json:"services"`
	Events     []Event     `json:"events"`
}

type ProfileInfo struct {
	ProductKey string `json:"productKey"`
}

type Property struct {
	Identifier   string      `json:"identifier"`
	Name         string      `json:"name"`
	Required     bool        `json:"required"`
	Type         string      `json:"type"`
	AccessMode   string      `json:"accessMode"`
	Unit         string      `json:"unit"`
	UnitSymbol   string      `json:"unitSymbol"`
	DefaultValue interface{} `json:"defaultValue"`
	MinValue     *float64    `json:"minValue"`
	MaxValue     *float64    `json:"maxValue"`
	Step         *float64    `json:"step"`
	Enum         []EnumValue `json:"enum"`
}

type EnumValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Service struct {
	Identifier   string     `json:"identifier"`
	Name         string     `json:"name"`
	Method       string     `json:"method"`
	InputParams  []ParamDef `json:"inputParams"`
	OutputParams []ParamDef `json:"outputParams"`
	AsyncResult  bool       `json:"asyncResult"`
}

type Event struct {
	Identifier   string     `json:"identifier"`
	Name         string     `json:"name"`
	Level        string     `json:"level"`
	OutputParams []ParamDef `json:"outputParams"`
}

type ParamDef struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Required   bool   `json:"required"`
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

	db.AutoMigrate(&Product{}, &ThingModelDefine{}, &TSLVersion{})
}

func main() {
	initDB()

	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})

	// 产品相关API
	r.GET("/api/v1/products", getProductList)
	r.GET("/api/v1/products/:id", getProduct)
	r.POST("/api/v1/products", createProduct)
	r.PUT("/api/v1/products/:id", updateProduct)
	r.DELETE("/api/v1/products/:id", deleteProduct)
	r.POST("/api/v1/products/:id/publish", publishProduct)
	r.POST("/api/v1/products/:id/tsl/import", importTSL)
	r.GET("/api/v1/products/:id/tsl/export", exportTSL)
	r.GET("/api/v1/products/:id/tsl/history", getTSLHistory)
	r.GET("/api/v1/products/statistics", getProductStatistics)

	// 物模型定义API
	r.GET("/api/v1/products/:id/properties", getProductProperties)
	r.POST("/api/v1/products/:id/properties", addProperty)
	r.PUT("/api/v1/products/:id/properties/:pid", updateProperty)
	r.DELETE("/api/v1/products/:id/properties/:pid", deleteProperty)

	r.GET("/api/v1/products/:id/services", getProductServices)
	r.POST("/api/v1/products/:id/services", addService)
	r.PUT("/api/v1/products/:id/services/:sid", updateService)
	r.DELETE("/api/v1/products/:id/services/:sid", deleteService)

	r.GET("/api/v1/products/:id/events", getProductEvents)
	r.POST("/api/v1/products/:id/events", addEvent)
	r.PUT("/api/v1/products/:id/events/:eid", updateEvent)
	r.DELETE("/api/v1/products/:id/events/:eid", deleteEvent)

	// 运行在8082端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("product-service starting on %s", addr)
	
	if err := r.Run(addr); err != nil {
		log.Printf("Failed to start on %s, trying :8095", addr)
		if err := r.Run(":8095"); err != nil {
			log.Fatal("failed to start server:", err)
		}
	}
}

// ========== 产品相关处理函数 ==========

func getProductList(c *gin.Context) {
	var products []Product
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	offset := (page - 1) * pageSize

	db.Model(&Product{}).Count(&total)
	db.Offset(offset).Limit(pageSize).Find(&products)

	c.JSON(200, gin.H{
		"code": 0,
		"data": products,
		"pagination": gin.H{
			"page":       page,
			"size":       pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func getProduct(c *gin.Context) {
	var product Product
	if err := db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": product})
}

func createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	product := Product{
		ProductName: req.ProductName,
		ProductKey:  generateProductKey(),
		TenantID:    req.TenantID,
		ProjectID:   req.ProjectID,
		CategoryID:  req.CategoryID,
		CommType:    req.CommType,
		AuthType:    req.AuthType,
		Status:      "draft",
		Icon:        req.Icon,
		Description: req.Description,
		ThingModel:  req.ThingModel,
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": product})
}

func updateProduct(c *gin.Context) {
	var product Product
	if err := db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	product.ProductName = req.ProductName
	product.CommType = req.CommType
	product.AuthType = req.AuthType
	product.Description = req.Description
	product.Icon = req.Icon
	product.ThingModel = req.ThingModel

	db.Save(&product)
	c.JSON(200, gin.H{"code": 0, "data": product})
}

func deleteProduct(c *gin.Context) {
	var product Product
	if err := db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	db.Delete(&product)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func publishProduct(c *gin.Context) {
	var product Product
	if err := db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	product.Status = "published"
	db.Save(&product)

	c.JSON(200, gin.H{"code": 0, "message": "产品已发布"})
}

// ========== TSL管理API ==========

func importTSL(c *gin.Context) {
	productID := c.Param("id")
	
	var product Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	var req ImportTSLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 解析并验证TSL数据
	var tslData map[string]interface{}
	if err := json.Unmarshal([]byte(req.TSLData), &tslData); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "TSL数据格式错误: " + err.Error()})
		return
	}

	// 保存TSL数据
	product.ThingModel = req.TSLData
	db.Save(&product)

	// 解析并保存物模型定义
	saveTSLDefinitions(productID, tslData)

	// 保存版本历史
	version := TSLVersion{
		ProductID: uintStrToUint(productID),
		Version:   formatTimestamp(),
		TSLData:   req.TSLData,
		ChangeLog: "导入TSL",
	}
	db.Create(&version)

	c.JSON(200, gin.H{"code": 0, "message": "TSL导入成功"})
}

func saveTSLDefinitions(productID string, tslData map[string]interface{}) {
	// 清除旧的定义
	db.Where("product_id = ?", uintStrToUint(productID)).Delete(&ThingModelDefine{})

	// 保存属性
	if properties, ok := tslData["properties"].([]interface{}); ok {
		for _, prop := range properties {
			if propMap, ok := prop.(map[string]interface{}); ok {
				define := ThingModelDefine{
					ProductID:  uintStrToUint(productID),
					DefineType: "property",
					Identifier: getStringValue(propMap, "identifier"),
					Name:       getStringValue(propMap, "name"),
					Desc:       getStringValue(propMap, "desc"),
					AccessMode: getStringValue(propMap, "accessMode"),
					Type:       getStringValue(propMap, "type"),
					Unit:       getStringValue(propMap, "unit"),
					UnitSymbol: getStringValue(propMap, "unitSymbol"),
				}
				if v, ok := propMap["minValue"].(float64); ok {
					define.MinValue = &v
				}
				if v, ok := propMap["maxValue"].(float64); ok {
					define.MaxValue = &v
				}
				if v, ok := propMap["step"].(float64); ok {
					define.Step = &v
				}
				db.Create(&define)
			}
		}
	}

	// 保存服务
	if services, ok := tslData["services"].([]interface{}); ok {
		for _, svc := range services {
			if svcMap, ok := svc.(map[string]interface{}); ok {
				define := ThingModelDefine{
					ProductID:  uintStrToUint(productID),
					DefineType: "service",
					Identifier: getStringValue(svcMap, "identifier"),
					Name:       getStringValue(svcMap, "name"),
					Desc:       getStringValue(svcMap, "desc"),
					InputParams: getStringValue(svcMap, "inputParams"),
					OutputParams: getStringValue(svcMap, "outputParams"),
				}
				db.Create(&define)
			}
		}
	}

	// 保存事件
	if events, ok := tslData["events"].([]interface{}); ok {
		for _, evt := range events {
			if evtMap, ok := evt.(map[string]interface{}); ok {
				define := ThingModelDefine{
					ProductID:  uintStrToUint(productID),
					DefineType: "event",
					Identifier: getStringValue(evtMap, "identifier"),
					Name:       getStringValue(evtMap, "name"),
					Desc:       getStringValue(evtMap, "desc"),
					Level:      getStringValue(evtMap, "level"),
					OutputParams: getStringValue(evtMap, "outputParams"),
				}
				db.Create(&define)
			}
		}
	}
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func exportTSL(c *gin.Context) {
	productID := c.Param("id")
	
	var product Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	// 获取物模型定义
	var properties []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'property'", productID).Find(&properties)

	var services []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'service'", productID).Find(&services)

	var events []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'event'", productID).Find(&events)

	// 构建AlLink格式
	resp := ExportTSLResponse{
		Schema:  "https://iot.aliyun.com/tsl/v1",
		Profile: ProfileInfo{ProductKey: product.ProductKey},
	}

	for _, p := range properties {
		prop := Property{
			Identifier: p.Identifier,
			Name:       p.Name,
			Type:       p.Type,
			AccessMode: p.AccessMode,
			Unit:       p.Unit,
		}
		if p.MinValue != nil {
			prop.MinValue = p.MinValue
		}
		if p.MaxValue != nil {
			prop.MaxValue = p.MaxValue
		}
		if p.Step != nil {
			prop.Step = p.Step
		}
		resp.Properties = append(resp.Properties, prop)
	}

	for _, s := range services {
		service := Service{
			Identifier: s.Identifier,
			Name:       s.Name,
			Method:     "sync",
		}
		resp.Services = append(resp.Services, service)
	}

	for _, e := range events {
		event := Event{
			Identifier: e.Identifier,
			Name:       e.Name,
			Level:      e.Level,
		}
		resp.Events = append(resp.Events, event)
	}

	c.JSON(200, gin.H{"code": 0, "data": resp})
}

func getTSLHistory(c *gin.Context) {
	productID := c.Param("id")
	
	var versions []TSLVersion
	db.Where("product_id = ?", productID).Order("created_at DESC").Limit(10).Find(&versions)

	c.JSON(200, gin.H{"code": 0, "data": versions})
}

// ========== 物模型定义API ==========

func getProductProperties(c *gin.Context) {
	var properties []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'property'", c.Param("id")).Find(&properties)
	c.JSON(200, gin.H{"code": 0, "data": properties})
}

func addProperty(c *gin.Context) {
	var req ThingModelDefine
	req.DefineType = "property"
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	req.ProductID = uintStrToUint(c.Param("id"))

	if err := db.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": req})
}

func updateProperty(c *gin.Context) {
	var property ThingModelDefine
	if err := db.First(&property, c.Param("pid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "属性不存在"})
		return
	}

	var req ThingModelDefine
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	property.Name = req.Name
	property.Desc = req.Desc
	property.AccessMode = req.AccessMode
	property.Type = req.Type
	property.Unit = req.Unit
	property.MinValue = req.MinValue
	property.MaxValue = req.MaxValue
	property.Step = req.Step
	property.EnumValues = req.EnumValues

	db.Save(&property)
	c.JSON(200, gin.H{"code": 0, "data": property})
}

func deleteProperty(c *gin.Context) {
	var property ThingModelDefine
	if err := db.First(&property, c.Param("pid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "属性不存在"})
		return
	}

	db.Delete(&property)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func getProductServices(c *gin.Context) {
	var services []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'service'", c.Param("id")).Find(&services)
	c.JSON(200, gin.H{"code": 0, "data": services})
}

func addService(c *gin.Context) {
	var req ThingModelDefine
	req.DefineType = "service"
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	req.ProductID = uintStrToUint(c.Param("id"))

	if err := db.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": req})
}

func updateService(c *gin.Context) {
	var service ThingModelDefine
	if err := db.First(&service, c.Param("sid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "服务不存在"})
		return
	}

	var req ThingModelDefine
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	service.Name = req.Name
	service.Desc = req.Desc
	service.InputParams = req.InputParams
	service.OutputParams = req.OutputParams

	db.Save(&service)
	c.JSON(200, gin.H{"code": 0, "data": service})
}

func deleteService(c *gin.Context) {
	var service ThingModelDefine
	if err := db.First(&service, c.Param("sid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "服务不存在"})
		return
	}

	db.Delete(&service)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func getProductEvents(c *gin.Context) {
	var events []ThingModelDefine
	db.Where("product_id = ? AND define_type = 'event'", c.Param("id")).Find(&events)
	c.JSON(200, gin.H{"code": 0, "data": events})
}

func addEvent(c *gin.Context) {
	var req ThingModelDefine
	req.DefineType = "event"
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	req.ProductID = uintStrToUint(c.Param("id"))

	if err := db.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 0, "data": req})
}

func updateEvent(c *gin.Context) {
	var event ThingModelDefine
	if err := db.First(&event, c.Param("eid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "事件不存在"})
		return
	}

	var req ThingModelDefine
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	event.Name = req.Name
	event.Desc = req.Desc
	event.Level = req.Level
	event.OutputParams = req.OutputParams

	db.Save(&event)
	c.JSON(200, gin.H{"code": 0, "data": event})
}

func deleteEvent(c *gin.Context) {
	var event ThingModelDefine
	if err := db.First(&event, c.Param("eid")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "事件不存在"})
		return
	}

	db.Delete(&event)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

// ========== 产品统计 ==========

func getProductStatistics(c *gin.Context) {
	var publishedCount, draftCount, totalCount int64

	db.Raw("SELECT COUNT(*) FROM products WHERE status = 'published'").Scan(&publishedCount)
	db.Raw("SELECT COUNT(*) FROM products WHERE status = 'draft'").Scan(&draftCount)
	db.Raw("SELECT COUNT(*) FROM products").Scan(&totalCount)

	stats := gin.H{
		"publishedCount": publishedCount,
		"draftCount":     draftCount,
		"totalCount":     totalCount,
	}

	c.JSON(200, gin.H{"code": 0, "data": stats})
}

// ========== 工具函数 ==========

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func generateProductKey() string {
	return fmt.Sprintf("p%s", formatTimestamp())
}

func formatTimestamp() string {
	return fmt.Sprintf("%010d", time.Now().UnixNano()/1e6)
}

func uintStrToUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint(n)
}
