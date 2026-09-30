# IoT飞燕平台 - Rust 服务一键启动脚本
# 需要先启动 Docker Desktop

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "IoT飞燕平台 - Rust 服务启动" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Docker Desktop 是否运行
Write-Host "检查 Docker Desktop..." -ForegroundColor Yellow
$dockerRunning = docker info 2>$null
if (-not $dockerRunning) {
    Write-Host "Docker Desktop 未运行，正在启动..." -ForegroundColor Red
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    Write-Host "请等待 Docker Desktop 启动完成（约30秒）..." -ForegroundColor Yellow
    Start-Sleep -Seconds 35
    
    # 再次检查
    $dockerRunning = docker info 2>$null
    if (-not $dockerRunning) {
        Write-Host "错误: Docker Desktop 启动失败，请手动启动后重试" -ForegroundColor Red
        exit 1
    }
}
Write-Host "✓ Docker Desktop 运行中" -ForegroundColor Green
Write-Host ""

# 停止现有容器
Write-Host "停止现有容器..." -ForegroundColor Yellow
docker compose -f docker-compose.wslc.yml down 2>$null
Write-Host ""

# 启动所有服务
Write-Host "启动所有 Rust 服务..." -ForegroundColor Yellow
docker compose -f docker-compose.wslc.yml up -d --build 2>&1
Write-Host ""

# 等待服务就绪
Write-Host "等待服务就绪..." -ForegroundColor Yellow
Start-Sleep -Seconds 25

# 显示容器状态
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "容器状态" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
docker ps --filter "name=iot-" --format "table {{.Names}}`t{{.Status}}`t{{.Ports}}"
Write-Host ""

# 健康检查
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "服务健康检查" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
$services = @(
    @("网关", "8080"),
    @("设备", "8081"),
    @("产品", "8082"),
    @("规则", "8083"),
    @("数据", "8084"),
    @("运营", "8085")
)

foreach ($svc in $services) {
    $name = $svc[0]
    $port = $svc[1]
    try {
        $health = curl -s --max-time 5 "http://localhost:$port/health" 2>$null
        if ($health -match '"code":0') {
            Write-Host "✓ $name 服务正常 (:$port)" -ForegroundColor Green
        } else {
            Write-Host "✗ $name 服务异常 (:$port)" -ForegroundColor Red
        }
    } catch {
        Write-Host "✗ $name 服务未启动 (:$port)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "服务访问地址" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "前端界面:        http://localhost:3000" -ForegroundColor White
Write-Host "API网关:         http://localhost:8080/health" -ForegroundColor White
Write-Host "设备服务:        http://localhost:8081/health" -ForegroundColor White
Write-Host "产品服务:        http://localhost:8082/health" -ForegroundColor White
Write-Host "规则服务:        http://localhost:8083/health" -ForegroundColor White
Write-Host "数据服务:        http://localhost:8084/health" -ForegroundColor White
Write-Host "运营服务:        http://localhost:8085/health" -ForegroundColor White
Write-Host "EMQX Dashboard:  http://localhost:18083" -ForegroundColor White
Write-Host "==========================================" -ForegroundColor Cyan
