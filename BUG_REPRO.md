# BUG_REPRO: 请求取消不传播、goroutine 泄漏与优雅关闭不排空

## Bug 是什么
- 请求取消不会向下游传播：`opsContext` 用 `context.Background()` 重建父级、`opsDelay` 忽略 ctx 且不 Stop 定时器，回填任务上下文与请求脱钩。
- 服务关闭不等待在途请求：`serverLifecycle.inflight` 只 Add/Done 从不 Wait，压测中关闭服务后 goroutine 持续上涨、进程不干净退出。
- 关闭超时/请求超时配置读取口径错误：`SHUTDOWN_TIMEOUT_SECONDS` 被钳制、`REQUEST_TIMEOUT_SECONDS` 读错变量导致配置不生效。

## 如何触发
1. 压测期间频繁提交/取消回填任务，goroutine 从个位数涨到数十万，停止后不回收。
2. 请求取消后，下游 `OpsService.Transition` / 回填 job 仍在继续处理。
3. 压测中 Ctrl-C 关闭服务，在飞 handler 未被排空，进程退出不干净。

## 错误信息
- `opsContext(parent, ...)` 对已取消 parent 派生的 ctx 仍为 live（deadline 未继承）。
- `opsDelay` 忽略 ctx.Done 且定时器不回收。
- `serveHTTPListener` 关闭后 `lifecycle.Active()` 不为 0（在飞 handler 未等待）。
- `SHUTDOWN_TIMEOUT_SECONDS=5` 实际生效为 10；`REQUEST_TIMEOUT_SECONDS=2` 不生效。
