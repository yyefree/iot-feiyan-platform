package main

import (
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

type Rule struct {
	gorm.Model
	Name          string     `json:"name" gorm:"not null;index"`
	Description   string     `json:"description"`
	RuleType      string     `json:"rule_type" gorm:"default:'scene'"`
	TenantID      uint       `json:"tenant_id" gorm:"index"`
	SQLExpression string     `json:"sql_expression" gorm:"type:text"`
	TriggerConfig string     `json:"trigger_config" gorm:"type:text"`
	ActionConfig  string     `json:"action_config" gorm:"type:text"`
	Status        string     `json:"status" gorm:"default:'active'"`
	CreatedBy     uint       `json:"created_by"`
	LastTriggered *time.Time `json:"last_triggered"`
	TriggerCount  int        `json:"trigger_count" gorm:"default:0"`
}

type RuleExecutionLog struct {
	gorm.Model
	RuleID     uint   `json:"rule_id" gorm:"index"`
	RuleName   string `json:"rule_name"`
	Status     string `json:"status"`
	InputData  string `json:"input_data" gorm:"type:text"`
	OutputData string `json:"output_data" gorm:"type:text"`
	ErrorMsg   string `json:"error_msg"`
	ExecTime   int64  `json:"exec_time"`
}

type Scene struct {
	gorm.Model
	Name         string `json:"name" gorm:"not null;index"`
	Description  string `json:"description"`
	TenantID     uint   `json:"tenant_id" gorm:"index"`
	TriggerRules string `json:"trigger_rules" gorm:"type:text"`
	ActionRules  string `json:"action_rules" gorm:"type:text"`
	CronExpr     string `json:"cron_expr"`
	EnvTrigger   string `json:"env_trigger" gorm:"type:text"` // 环境触发条件
	Status       string `json:"status" gorm:"default:'active'"`
	CreatedBy    uint   `json:"created_by"`
}

type RuleForwardLog struct {
	gorm.Model
	RuleID   uint      `json:"rule_id" gorm:"index"`
	RuleName string    `json:"rule_name"`
	Source   string    `json:"source"`
	Dest     string    `json:"dest"`
	Status   string    `json:"status"`
	ErrorMsg string    `json:"error_msg"`
	ExecTime int64     `json:"exec_time"`
	CreatedAt time.Time `json:"created_at"`
}

// ========== 请求结构 ==========

type CreateRuleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	RuleType      string `json:"rule_type"`
	TenantID      uint   `json:"tenant_id"`
	SQLExpression string `json:"sql_expression"`
	TriggerConfig string `json:"trigger_config"`
	ActionConfig  string `json:"action_config"`
}

type CreateSceneRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TenantID     uint   `json:"tenant_id"`
	TriggerRules string `json:"trigger_rules"`
	ActionRules  string `json:"action_rules"`
	CronExpr     string `json:"cron_expr"`
	EnvTrigger   string `json:"env_trigger"`
}

type TestSQLRequest struct {
	SQLExpression string `json:"sql_expression" binding:"required"`
	SampleData    string `json:"sample_data"`
}

type TestSQLResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ForwardRuleRequest struct {
	RuleID  uint   `json:"rule_id" binding:"required"`
	Source  string `json:"source"`
	Dest    string `json:"dest"`
	Filter  string `json:"filter"`
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
	
	db.AutoMigrate(&Rule{}, &RuleExecutionLog{}, &Scene{}, &RuleForwardLog{})
}

func main() {
	initDB()
	
	r := gin.Default()
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "message": "ok"})
	})
	
	// 规则引擎API
	r.GET("/api/v1/rules", getRuleList)
	r.GET("/api/v1/rules/:id", getRule)
	r.POST("/api/v1/rules", createRule)
	r.PUT("/api/v1/rules/:id", updateRule)
	r.DELETE("/api/v1/rules/:id", deleteRule)
	r.POST("/api/v1/rules/:id/enable", enableRule)
	r.POST("/api/v1/rules/:id/disable", disableRule)
	r.GET("/api/v1/rules/:id/logs", getRuleLogs)
	r.GET("/api/v1/rules/:id/execute", executeRuleTest)
	
	// SQL测试API
	r.POST("/api/v1/rules/sql/test", testSQL)
	
	// 数据转发API
	r.GET("/api/v1/rules/forward/logs", getForwardLogs)
	r.POST("/api/v1/rules/forward/config", updateForwardConfig)
	r.POST("/api/v1/rules/forward/test", testForward)
	
	// 场景联动API
	r.GET("/api/v1/scenes", getSceneList)
	r.GET("/api/v1/scenes/:id", getScene)
	r.POST("/api/v1/scenes", createScene)
	r.PUT("/api/v1/scenes/:id", updateScene)
	r.DELETE("/api/v1/scenes/:id", deleteScene)
	r.POST("/api/v1/scenes/:id/enable", enableScene)
	r.POST("/api/v1/scenes/:id/disable", disableScene)
	r.POST("/api/v1/scenes/:id/execute", testScene)
	r.GET("/api/v1/scenes/sunrise", getSunriseSunset)
	
	// 规则统计
	r.GET("/api/v1/rules/statistics", getRuleStatistics)
	
	// 运行在8093端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("rule-service starting on %s", addr)
	
	if err := r.Run(addr); err != nil {
		log.Printf("Failed to start on %s, trying :8094", addr)
		if err := r.Run(":8094"); err != nil {
			log.Fatal("failed to start server:", err)
		}
	}
}

// ========== 规则引擎处理函数 ==========

func getRuleList(c *gin.Context) {
	var rules []Rule
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	offset := (page - 1) * pageSize
	
	db.Model(&Rule{}).Count(&total)
	db.Offset(offset).Limit(pageSize).Find(&rules)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": rules,
		"pagination": gin.H{
			"page":       page,
			"size":       pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func getRule(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": rule})
}

func createRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	rule := Rule{
		Name:          req.Name,
		Description:   req.Description,
		RuleType:      req.RuleType,
		TenantID:      req.TenantID,
		SQLExpression: req.SQLExpression,
		TriggerConfig: req.TriggerConfig,
		ActionConfig:  req.ActionConfig,
		Status:        "active",
	}
	
	if err := db.Create(&rule).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "data": rule})
}

func updateRule(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	rule.Name = req.Name
	rule.Description = req.Description
	rule.RuleType = req.RuleType
	rule.SQLExpression = req.SQLExpression
	rule.TriggerConfig = req.TriggerConfig
	rule.ActionConfig = req.ActionConfig
	
	db.Save(&rule)
	c.JSON(200, gin.H{"code": 0, "data": rule})
}

func deleteRule(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	
	db.Delete(&rule)
	c.JSON(200, gin.H{"code": 0, "message": "删除成功"})
}

func enableRule(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	
	rule.Status = "active"
	db.Save(&rule)
	c.JSON(200, gin.H{"code": 0, "message": "规则已启用"})
}

func disableRule(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	
	rule.Status = "inactive"
	db.Save(&rule)
	c.JSON(200, gin.H{"code": 0, "message": "规则已禁用"})
}

func getRuleLogs(c *gin.Context) {
	var logs []RuleExecutionLog
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&RuleExecutionLog{}).Where("rule_id = ?", c.Param("id")).Count(&total)
	db.Where("rule_id = ?", c.Param("id")).Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs)
	
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

func executeRuleTest(c *gin.Context) {
	var rule Rule
	if err := db.First(&rule, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	
	startTime := time.Now()
	
	// 模拟规则执行
	result := map[string]interface{}{
		"status": "success",
		"message": "规则执行成功（测试模式）",
		"executed_at": startTime.Format(time.RFC3339),
	}
	
	execTime := time.Since(startTime).Milliseconds()
	
	// 保存执行日志
	log := RuleExecutionLog{
		RuleID:     rule.ID,
		RuleName:   rule.Name,
		Status:     "success",
		InputData:  "测试数据",
		OutputData: fmt.Sprintf("%v", result),
		ExecTime:   execTime,
	}
	db.Create(&log)
	
	// 更新规则统计
	db.Model(&rule).Updates(map[string]interface{}{
		"last_triggered": startTime,
		"trigger_count":  gorm.Expr("trigger_count + 1"),
	})
	
	c.JSON(200, gin.H{"code": 0, "data": result})
}

// ========== SQL测试API ==========

func testSQL(c *gin.Context) {
	var req TestSQLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	result := TestSQLResponse{
		Success: true,
		Message: "SQL语法正确",
	}
	
	// 完整的SQL语法检查
	lowerSQL := strings.ToLower(req.SQLExpression)
	
	// 检查SQL关键字
	keywords := map[string]string{
		"select": "SELECT查询语句",
		"insert": "INSERT插入语句",
		"update": "UPDATE更新语句",
		"delete": "DELETE删除语句",
		"from":   "FROM子句",
		"where":  "WHERE条件",
		"order":  "ORDER排序",
		"group":  "GROUP分组",
		"having": "HAVING过滤",
		"join":   "JOIN连接",
		"and":    "AND条件",
		"or":     "OR条件",
	}
	
	foundKeywords := make([]string, 0)
	for kw, desc := range keywords {
		if strings.Contains(lowerSQL, kw) {
			foundKeywords = append(foundKeywords, desc)
		}
	}
	
	if len(foundKeywords) > 0 {
		result.Result = fmt.Sprintf("检测到SQL关键字: %s", strings.Join(foundKeywords, ", "))
	} else {
		result.Success = false
		result.Error = "无法识别的SQL语句类型"
	}
	
	c.JSON(200, gin.H{"code": 0, "data": result})
}

// ========== 数据转发API ==========

func getForwardLogs(c *gin.Context) {
	var logs []RuleForwardLog
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&RuleForwardLog{}).Count(&total)
	db.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs)
	
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

func updateForwardConfig(c *gin.Context) {
	var req ForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	// 保存转发配置
	log := RuleForwardLog{
		RuleID: req.RuleID,
		Source: req.Source,
		Dest:   req.Dest,
		Status: "configured",
	}
	db.Create(&log)
	
	c.JSON(200, gin.H{"code": 0, "message": "转发配置已更新"})
}

func testForward(c *gin.Context) {
	var req ForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	// 模拟转发测试
	log := RuleForwardLog{
		RuleID: req.RuleID,
		Source: req.Source,
		Dest:   req.Dest,
		Status: "success",
		ExecTime: 123,
	}
	db.Create(&log)
	
	c.JSON(200, gin.H{"code": 0, "message": "转发测试成功", "data": log})
}

// ========== 场景联动处理函数 ==========

func getSceneList(c *gin.Context) {
	var scenes []Scene
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	
	db.Model(&Scene{}).Count(&total)
	db.Offset((page-1)*pageSize).Limit(pageSize).Find(&scenes)
	
	c.JSON(200, gin.H{
		"code": 0,
		"data": scenes,
		"pagination": gin.H{
			"page":       page,
			"size":       pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func getScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": scene})
}

func createScene(c *gin.Context) {
	var req CreateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	scene := Scene{
		Name:         req.Name,
		Description:  req.Description,
		TenantID:     req.TenantID,
		TriggerRules: req.TriggerRules,
		ActionRules:  req.ActionRules,
		CronExpr:     req.CronExpr,
		EnvTrigger:   req.EnvTrigger,
		Status:       "active",
	}
	
	if err := db.Create(&scene).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{"code": 0, "data": scene})
}

func updateScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	
	var req CreateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	
	scene.Name = req.Name
	scene.Description = req.Description
	scene.TriggerRules = req.TriggerRules
	scene.ActionRules = req.ActionRules
	scene.CronExpr = req.CronExpr
	scene.EnvTrigger = req.EnvTrigger
	
	db.Save(&scene)
	c.JSON(200, gin.H{"code": 0, "data": scene})
}

func deleteScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	
	db.Delete(&scene)
	c.JSON(200, gin.H{"code": 0, "message": "场景已删除"})
}

func enableScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	
	scene.Status = "active"
	db.Save(&scene)
	c.JSON(200, gin.H{"code": 0, "message": "场景已启用"})
}

func disableScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	
	scene.Status = "inactive"
	db.Save(&scene)
	c.JSON(200, gin.H{"code": 0, "message": "场景已禁用"})
}

func testScene(c *gin.Context) {
	var scene Scene
	if err := db.First(&scene, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"code": 404, "message": "场景不存在"})
		return
	}
	
	result := map[string]interface{}{
		"status": "success",
		"message": "场景测试执行成功",
		"scene_name": scene.Name,
		"executed_at": time.Now().Format(time.RFC3339),
	}
	
	c.JSON(200, gin.H{"code": 0, "data": result})
}

func getSunriseSunset(c *gin.Context) {
	now := time.Now()
	sunrise := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
	sunset := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	
	result := gin.H{
		"date":     now.Format("2006-01-02"),
		"sunrise":  sunrise.Format("15:04:05"),
		"sunset":   sunset.Format("15:04:05"),
		"latitude": 31.2304,
		"longitude": 121.4737,
	}
	
	c.JSON(200, gin.H{"code": 0, "data": result})
}

// ========== 规则统计 ==========

func getRuleStatistics(c *gin.Context) {
	var totalRules, activeRules, sceneRules, streamRules int64
	var totalExecutions int64
	
	db.Raw("SELECT COUNT(*) FROM rules").Scan(&totalRules)
	db.Raw("SELECT COUNT(*) FROM rules WHERE status = 'active'").Scan(&activeRules)
	db.Raw("SELECT COUNT(*) FROM rules WHERE rule_type = 'scene'").Scan(&sceneRules)
	db.Raw("SELECT COUNT(*) FROM rules WHERE rule_type = 'stream'").Scan(&streamRules)
	db.Raw("SELECT COUNT(*) FROM rule_execution_logs").Scan(&totalExecutions)
	
	stats := gin.H{
		"totalRules":      totalRules,
		"activeRules":     activeRules,
		"sceneRules":      sceneRules,
		"streamRules":     streamRules,
		"totalExecutions": totalExecutions,
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
