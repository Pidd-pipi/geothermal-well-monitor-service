# Geothermal Well Monitor Service

标准库地热井监测 HTTP 服务，覆盖井口工况、测点读数、异常告警、巡检任务、保养计划、运维记录（OPS）、历史回填与值班登录。`GET /health` 健康检查，根路径为内置监测页面，默认端口 8080。

## 目录结构

```text
backend/                 # Go module、源码、内嵌前端页面
  main.go / http.go      # 入口与路由
  store.go / service.go  # 井口工况（wells）
  reading*.go            # 测点读数采集、查询与统计
  alert*.go              # 告警规则、评估与确认
  task*.go               # 巡检任务工作流（queued→assigned→in_progress→completed）
  plan*.go               # 保养计划与到期排程
  ops_*.go               # 运维记录、状态机、审计、分页查询
  backfill*.go           # 历史读数批量回填（后台 worker）
  operator.go / session.go / auth_http.go  # 值班登录与会话
  runtime.go             # 服务生命周期、中间件、优雅关闭
  web/                   # 内嵌监测页面（index.html / app.js）
database/                # 持久化说明与未来表结构
runtime_smoke.json       # 真实启动健康检查声明
```

## 运行与测试

```bash
cd backend
go build ./...
go test ./...
PORT=8080 go run .
```

健康检查 `GET /health`；页面 `GET /`、`GET /app.js`。

## API 一览

- `GET /api/wells` 井列表；`POST /api/wells/status` 更新井状态（producing/inspection/isolated）
- `POST /api/readings` 上报读数；`GET /api/readings/recent?well_id=&n=` 最近读数；`GET /api/readings/stats?well_id=&hours=` 统计
- `GET /api/alerts` 告警列表；`POST /api/alerts/{id}/ack` 确认告警
- `POST /api/tasks` 建巡检任务；`POST /api/tasks/{id}/assign`、`/start`、`/complete` 流转；`GET /api/tasks`
- `GET /api/plans` 保养计划；`POST /api/plans`；`POST /api/plans/{id}/done` 完成并重排；`GET /api/plans?due=true` 到期列表
- `GET /api/ops/records`、`POST /api/ops/records`、`POST /api/ops/records/{id}/transition`、`GET /api/ops/records/{id}/audit`、`GET /api/ops/snapshot`
- `POST /api/backfill` 历史读数回填；`GET /api/backfill/{id}` 回填状态
- `POST /api/operators/login` 登录（默认 liang/well2026、zhao/ops2026）；`GET /api/operators/me` 需 `X-Session-Token`

## 环境变量

| 变量 | 默认 | 说明 |
|---|---|---|
| PORT | 8080 | 监听端口 |
| SHUTDOWN_TIMEOUT_SECONDS | 10 | 优雅关闭超时 |
| REQUEST_TIMEOUT_SECONDS | 5 | 运维操作请求超时 |
