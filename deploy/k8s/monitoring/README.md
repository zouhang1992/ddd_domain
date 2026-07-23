# DDD House 监控系统

本目录包含 Prometheus 监控和 Grafana Dashboard 的完整配置。

## 目录结构

```
deploy/k8s/
├── prometheus/
│   └── prometheus-service-monitor.yaml  # Prometheus ServiceMonitor 配置
├── grafana/
│   └── ddd-house-dashboard.yaml          # Grafana Dashboard 配置
├── loki/
│   └── grafana-datasource.yaml           # Loki 数据源配置
└── monitoring/
    └── README.md                          # 本文件
```

## 部署步骤

### 前置条件

1. 已部署 Prometheus Operator (kube-prometheus-stack)
2. 已部署 Grafana
3. 已部署 Loki (可选，用于日志查询)

### 1. 部署应用

```bash
# 部署 ddd-house 应用
kubectl apply -f ../app/ddd-house.yaml
```

### 2. 部署 Prometheus ServiceMonitor

```bash
# 创建 ServiceMonitor 让 Prometheus 抓取应用指标
kubectl apply -f ../prometheus/prometheus-service-monitor.yaml
```

### 3. 部署 Grafana Dashboard

```bash
# 部署 Grafana Dashboard 配置
kubectl apply -f ../grafana/ddd-house-dashboard.yaml
```

### 4. (可选) 部署 Loki 数据源

```bash
# 如果还没有部署 Loki 数据源
kubectl apply -f ../loki/grafana-datasource.yaml
```

## 监控指标说明

### HTTP 指标

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `http_requests_total` | Counter | HTTP 请求总数 |
| `http_request_duration_seconds` | Histogram | HTTP 请求延迟分布 |
| `http_requests_in_flight` | Gauge | 进行中的请求数 |
| `http_request_size_bytes` | Summary | HTTP 请求大小 |
| `http_response_size_bytes` | Summary | HTTP 响应大小 |

### 业务指标

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `commands_executed_total` | Counter | 命令执行总数 |
| `queries_executed_total` | Counter | 查询执行总数 |
| `events_published_total` | Counter | 领域事件发布总数 |

### 数据库指标

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `db_queries_total` | Counter | 数据库查询总数 |
| `db_query_duration_seconds` | Histogram | 数据库查询延迟分布 |

### Go 运行时指标

Prometheus client_golang 自动收集的指标：
- `go_goroutines` - Goroutine 数量
- `go_threads` - 线程数量
- `go_memory_heap_alloc_bytes` - 堆内存分配
- `go_gc_duration_seconds` - GC 延迟
- 等等...

## Grafana Dashboard 面板说明

### HTTP 服务概览 (第1行)

- **请求速率 (RPS)**: 每秒 HTTP 请求数
- **进行中的请求**: 当前正在处理的请求数
- **p95 延迟**: 95% 请求的延迟
- **5xx 错误速率**: 服务器错误的速率

### HTTP 服务详情 (第2行)

- **HTTP 请求速率（按路径）**: 各个 API 端点的请求速率
- **HTTP 请求延迟分布**: p50、p95、p99 延迟

### Go 运行时指标 (第3-4行)

- **内存使用**: 堆内存、系统内存分配
- **Goroutines & Threads**: 并发数
- **GC 频率**: 垃圾回收频率
- **GC 延迟分布**: GC 延迟的 p50、p95

### 业务指标 (第5-6行)

- **命令执行速率**: CQRS 命令执行速率
- **查询执行速率**: CQRS 查询执行速率
- **领域事件发布速率**: 领域事件发布速率
- **数据库查询速率**: 数据库操作速率

## 如何使用业务指标函数

在代码中使用以下函数记录业务指标：

```go
import "github.com/zouhang1992/ddd_domain/internal/infrastructure/middleware"

// 记录命令执行
middleware.RecordCommand("create_lease", true)  // success
middleware.RecordCommand("create_lease", false) // error

// 记录查询执行
middleware.RecordQuery("list_leases", true)

// 记录事件发布
middleware.RecordEvent("lease.created")

// 记录数据库查询
middleware.RecordDBQuery("SELECT", time.Since(start), true)
```

## 验证部署

### 1. 检查 ServiceMonitor

```bash
kubectl get servicemonitor -n monitoring
kubectl describe servicemonitor ddd-house-monitor -n monitoring
```

### 2. 检查 Prometheus Targets

在 Prometheus UI 中访问 `/targets`，确认 `ddd-house` 目标已出现。

### 3. 检查 Grafana Dashboard

1. 打开 Grafana UI
2. 点击 "Dashboards" → "Browse"
3. 搜索 "DDD House 应用监控"
4. 打开 Dashboard 查看监控数据

## 访问地址

根据你的 [CLAUDE.md](../../../../CLAUDE.md) 配置：

- **Grafana**: http://grafana.zouhang.com
- **Prometheus**: http://prometheus.zouhang.com
- **Alertmanager**: http://alertmanager.zouhang.com

## 故障排查

### 问题：Prometheus 没有抓取到指标

1. 检查 ServiceMonitor 是否创建成功
2. 检查 Service 的标签是否匹配 `app: ddd-house-backend`
3. 检查应用的 `/metrics` 端点是否正常
4. 查看 Prometheus 的 Targets 页面

### 问题：Grafana Dashboard 没有显示数据

1. 确认 Prometheus 数据源已配置
2. 检查时间范围是否正确
3. 查看 Prometheus 查询是否返回数据
4. 确认指标名称是否正确

### 问题：Dashboard 找不到

1. 检查 ConfigMap 是否正确部署在 `monitoring` 命名空间
2. 确认 ConfigMap 有标签 `grafana_dashboard: "1"`
3. 重启 Grafana Pod: `kubectl rollout restart deployment/grafana -n monitoring`
