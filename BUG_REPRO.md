# BUG_REPRO: 保养计划到期判定口径错位

## Bug 是什么
- 逾期计划完成后不推进 `NextDueAt`（`MarkDone` 带 `p.NextDueAt.After(now)` 守卫），做完仍显示到期。
- `Due` 用 `Status != PlanArchived` 判定，暂停（paused）计划也被算入到期。
- HTTP `?due=true` 用 `LastDoneAt` 谓词自行过滤，与 service 的 `NextDueAt` 口径不一致，刚完成的计划又出现在到期列表。

## 如何触发
1. 创建 30 天间隔的保养计划，40 天后 `POST /api/plans/{id}/done` 完成 → `GET /api/plans?due=true` 仍列出该计划。
2. 将计划置为 paused 并让其到期 → `Due()` 仍返回该计划。
3. 刚创建/刚完成的计划经 `GET /api/plans?due=true` 与 service `Due()` 结果不一致。

## 错误信息
- 逾期计划完成后 `NextDueAt` 仍在过去，`Due()` 持续返回它。
- 暂停计划出现在 `Due()` 结果中。
- HTTP 到期列表与 service 到期统计数字对不上（两套谓词：`LastDoneAt` vs `NextDueAt`）。
