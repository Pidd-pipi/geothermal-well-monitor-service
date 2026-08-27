# geothermal-well-monitor-service__004

基于 Go 实现的地热井监测 Web 项目，一款后端服务，完成井口工况、测点读数、异常告警、巡检任务、保养计划与运维记录管理。
## 构建镜像

请从**仓库根目录**执行；`benzhi.Dockerfile`、`build_benzhi_docker.sh`、`BENZHI_README.md` 均固定在该目录：

```bash
./build_benzhi_docker.sh <image-name> [linux/amd64|linux/arm64]
```

## 标准命令

```bash
go build ./...     # 编译
go run .   # 启动
go test ./...      # 测试（如有）
```

## 环境

- 基础镜像: golang:1.23
- Go 模块目录: `.`
- 依赖已在镜像构建阶段预下载，容器内离线可用。
- 容器内工作目录: `/app`
