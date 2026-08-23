# 森林火情值守服务

标准库实现的森林火情值守 HTTP 服务，默认监听 `8080`，可由 `PORT` 环境变量调整。接口为 `GET /healthz`、`GET /api/v1/alerts`、`POST /api/v1/alerts/{id}/status`，根路径提供静态值守页面。

状态值包括 `new`、`monitoring`、`contained`、`resolved`。

## Verification

在 `backend/` 目录执行 `gofmt -w .`、`go build ./...`、`go test ./...`，结果全部通过；HTTP 定向测试通过。

真实服务使用端口 `18103` 启动后验证：`GET /healthz` 返回 `{"service":"forest-fire-watch","status":"ok"}`；`GET /api/v1/alerts` 返回 2 条记录；`POST /api/v1/alerts/fa-301/status` 将状态更新为 `contained` 并返回 200；非法状态返回 400；未知火情返回 404；`/` 返回 797 字节，`/app.js` 返回 836 字节。验证后服务已关闭。

## Engineering Notes

森林火情流程代码按领域模型、校验、状态转换、并发安全存储、审计事件和 HTTP 生命周期分层。请求会保留请求标识并经过恢复与超时保护；状态写入使用版本校验，错误通过可识别的领域错误返回。

除现有接口回归测试外，项目还保留可复用的分页、过滤、策略、工作流和运行健康能力，便于后续扩展而不把业务规则堆积到处理器中。

## Enterprise Layout

```text
.
├── backend/       # Go module, source, tests, and embedded web assets
├── database/      # persistence documentation and future schemas
├── output/        # verification records
├── prompt.txt
└── runtime_smoke.json
```

Run `cd backend && go test ./...`, `cd backend && go build ./...`, or `cd backend && PORT=8080 go run .`. The health check is `GET /healthz`; the main API endpoints are `GET /api/v1/alerts` and `POST /api/v1/alerts/{id}/status`.
