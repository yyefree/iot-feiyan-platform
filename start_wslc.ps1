# IoT飞燕平台 - WSLC 启动脚本
# 基于 WSL2 + Docker Desktop 的容器化部署方案

Write-Host "=========================================="
Write-Host "IoT飞燕平台 - WSLC 启动脚本"
Write-Host "=========================================="
Write-Host ""

# 检查 WSL2 环境
if (-not (wsl --list -q 2>$null)) {
    Write-Host "警告: 未检测到 WSL2 环境"
    Write-Host "请先安装 WSL2: wsl --install"
}

# 检查 Docker Desktop
$dockerRunning = docker info 2>$null
if (-not $dockerRunning) {
    Write-Host "错误: Docker Desktop 未运行"
    Write-Host "请启动 Docker Desktop"
    exit 1
}

Write-Host "✓ WSL2 环境检查通过"
Write-Host "✓ Docker Desktop 运行中"
Write-Host ""

# 停止现有容器
Write-Host "停止现有容器..."
docker compose -f docker-compose.wslc.yml down 2>$null
Write-Host ""

# 启动所有服务
Write-Host "启动所有服务..."
docker compose -f docker-compose.wslc.yml up -d --build 2>&1
Write-Host ""

# 等待服务就绪
Write-Host "等待服务就绪..."
Start-Sleep -Seconds 15

# 显示容器状态
Write-Host ""
Write-Host "=========================================="
Write-Host "容器状态"
Write-Host "=========================================="
docker ps --filter "name=iot-" --format "table {{.Names}}`t{{.Status}}`t{{.Ports}}"

Write-Host ""
Write-Host "=========================================="
Write-Host "服务访问地址"
Write-Host "=========================================="
Write-Host "前端界面:   http://localhost:3000"
Write-Host "API网关:    http://localhost:8080/health"
Write-Host "设备服务:   http://localhost:8081/health"
Write-Host "产品服务:   http://localhost:8082/health"
Write-Host "规则服务:   http://localhost:8083/health"
Write-Host "数据服务:   http://localhost:8084/health"
Write-Host "运营服务:   http://localhost:8085/health"
Write-Host "EMQX Dashboard: http://localhost:18083"
Write-Host "=========================================="
