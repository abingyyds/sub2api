# 代理系统（Agent/Affiliate System）实现计划

## 需求概述

基于现有的邀请码（referral）系统，包装为"代理"系统：
- 用户可以申请成为代理
- 代理生成专属推广链接（复用 invite_code）
- 代理拥有独立后台，查看下级用户的所有资金变动（充值、消费、套餐购买）
- 管理员可以审核/管理代理

## 核心设计思路

**不引入新角色**，在现有 `user` 角色基础上增加 `is_agent` 标记。代理本质上是一个有特殊权限的普通用户，通过 referrals 表关联下级用户，通过查询 payment_orders 和 usage_logs 展示下级用户的资金变动。

---

## Phase 1：数据库迁移

### 新建 migration `070_add_agent_system.sql`

```sql
-- 用户表增加代理字段
ALTER TABLE users ADD COLUMN is_agent BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN agent_status VARCHAR(20) NOT NULL DEFAULT '';
-- agent_status: '' (非代理), 'pending' (待审核), 'approved' (已批准), 'rejected' (已拒绝)
ALTER TABLE users ADD COLUMN agent_commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0;
-- 代理佣金比例 (0~1)，如 0.1 表示 10% 返佣
ALTER TABLE users ADD COLUMN agent_note TEXT NOT NULL DEFAULT '';
-- 代理申请备注
ALTER TABLE users ADD COLUMN agent_approved_at TIMESTAMPTZ;

-- 代理佣金记录表
CREATE TABLE agent_commissions (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT NOT NULL REFERENCES users(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    order_id BIGINT REFERENCES payment_orders(id),
    source_type VARCHAR(20) NOT NULL,
    -- source_type: 'payment' (充值/购买套餐), 'usage' (消费返佣, 可选)
    source_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
    commission_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- status: 'pending', 'settled' (已结算到余额)
    settled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agent_commissions_agent_id ON agent_commissions(agent_id);
CREATE INDEX idx_agent_commissions_user_id ON agent_commissions(user_id);
CREATE INDEX idx_agent_commissions_order_id ON agent_commissions(order_id);
CREATE INDEX idx_agent_commissions_status ON agent_commissions(status);
```

---

## Phase 2：后端 - Ent Schema & Repository

### 2.1 更新 User schema
- 新增字段: `is_agent`, `agent_status`, `agent_commission_rate`, `agent_note`, `agent_approved_at`
- 新增 edge: `agent_commissions`

### 2.2 新建 AgentCommission schema
- `ent/schema/agent_commission.go`
- 字段: agent_id, user_id, order_id, source_type, source_amount, commission_rate, commission_amount, status, settled_at

### 2.3 新建 Repository
- `internal/repository/agent_repo.go` — 实现代理相关查询：
  - `ListSubUsers(agentID)` — 通过 referrals 表查找下级用户
  - `GetSubUserPaymentOrders(agentID)` — 下级用户的充值/购买记录
  - `GetSubUserUsageSummary(agentID)` — 下级用户的消费汇总
  - `CreateCommission()` / `ListCommissions()` / `SettleCommissions()`

---

## Phase 3：后端 - Service & Handler

### 3.1 AgentService (`internal/service/agent_service.go`)
- `ApplyForAgent(userID, note)` — 提交代理申请 (status → pending)
- `GetAgentStatus(userID)` — 获取代理状态
- `GetAgentDashboard(agentID)` — 代理仪表盘数据（下级总数、总充值、总消费、佣金统计）
- `ListSubUsers(agentID, pagination)` — 列出下级用户
- `ListSubUserFinancialLogs(agentID, pagination, filters)` — 下级用户资金变动（整合 payment_orders + usage_logs）
- `ListCommissions(agentID, pagination)` — 佣金记录
- `GetInviteLink(agentID)` — 生成代理推广链接（复用 invite_code）

### 3.2 AdminAgentService — 管理员操作
- `ListAgentApplications(pagination, status)` — 代理申请列表
- `ApproveAgent(userID, commissionRate)` — 批准代理
- `RejectAgent(userID)` — 拒绝代理
- `UpdateCommissionRate(userID, rate)` — 调整佣金比例
- `ListAllAgents(pagination)` — 所有代理列表
- `GetAgentDetail(agentID)` — 代理详情（含下级统计）

### 3.3 Handler
- `internal/handler/agent_handler.go` — 代理用户端接口
- `internal/handler/admin/agent_handler.go` — 管理员端接口

### 3.4 佣金触发
- 在 `payment_service.go` 的支付成功回调中，检查付款用户是否有上级代理（通过 referrals 查找），如有则自动创建佣金记录

---

## Phase 4：后端 - 路由注册

### 4.1 代理用户路由 (`routes/user.go`)
```
GET    /api/v1/agent/status          — 代理状态
POST   /api/v1/agent/apply           — 申请代理
GET    /api/v1/agent/dashboard       — 代理仪表盘
GET    /api/v1/agent/link            — 推广链接
GET    /api/v1/agent/sub-users       — 下级用户列表
GET    /api/v1/agent/financial-logs  — 下级资金变动
GET    /api/v1/agent/commissions     — 佣金记录
```

### 4.2 管理员路由 (`routes/admin.go`)
```
GET    /api/v1/admin/agents              — 代理列表
GET    /api/v1/admin/agents/applications — 申请列表
POST   /api/v1/admin/agents/:id/approve  — 批准
POST   /api/v1/admin/agents/:id/reject   — 拒绝
PUT    /api/v1/admin/agents/:id          — 更新佣金率
GET    /api/v1/admin/agents/:id          — 代理详情
```

---

## Phase 5：前端 - 代理后台

### 5.1 API 层
- `frontend/src/api/agent.ts` — 代理 API
- `frontend/src/api/admin/agents.ts` — 管理员代理 API

### 5.2 路由
在 router/index.ts 中新增：
- `/agent/dashboard` — 代理仪表盘（requiresAuth + requiresAgent）
- `/agent/sub-users` — 下级用户管理
- `/agent/commissions` — 佣金记录
- `/admin/agents` — 管理员代理管理

### 5.3 代理后台页面
- `views/agent/AgentDashboardView.vue` — 仪表盘（统计卡片 + 推广链接 + 快速概览）
- `views/agent/SubUsersView.vue` — 下级用户列表（邮箱、注册时间、余额、总消费、总充值）
- `views/agent/FinancialLogsView.vue` — 资金变动明细（充值记录、套餐购买、消费记录）
- `views/agent/CommissionsView.vue` — 佣金记录

### 5.4 管理员代理管理页面
- `views/admin/AgentsView.vue` — 代理列表 & 审核

### 5.5 用户端入口
- 在用户侧边栏新增"代理中心"入口
- 非代理用户显示"申请成为代理"按钮
- 已批准代理显示完整代理后台

---

## Phase 6：Wire 依赖注入更新

更新以下 wire.go 文件注册新的 service/repo/handler：
- `internal/repository/wire.go`
- `internal/service/wire.go`
- `internal/handler/wire.go`
- `internal/handler/admin/wire.go` (如有)
- 重新运行 `go generate ./...` 生成 wire_gen.go

---

## 实现顺序

1. 数据库迁移 (Phase 1)
2. Ent Schema (Phase 2.1, 2.2)
3. Repository (Phase 2.3)
4. Service (Phase 3.1, 3.2)
5. Handler (Phase 3.3)
6. 佣金触发集成 (Phase 3.4)
7. 路由注册 (Phase 4)
8. Wire DI 更新 (Phase 6)
9. 前端 API 层 (Phase 5.1)
10. 前端路由 (Phase 5.2)
11. 前端页面 (Phase 5.3, 5.4, 5.5)

## 文件变更清单（预估）

**新建文件** (~15个)：
- `backend/migrations/070_add_agent_system.sql`
- `backend/ent/schema/agent_commission.go`
- `backend/internal/repository/agent_repo.go`
- `backend/internal/service/agent.go` (domain model)
- `backend/internal/service/agent_service.go`
- `backend/internal/handler/agent_handler.go`
- `backend/internal/handler/admin/agent_handler.go`
- `frontend/src/api/agent.ts`
- `frontend/src/api/admin/agents.ts`
- `frontend/src/views/agent/AgentDashboardView.vue`
- `frontend/src/views/agent/SubUsersView.vue`
- `frontend/src/views/agent/FinancialLogsView.vue`
- `frontend/src/views/agent/CommissionsView.vue`
- `frontend/src/views/admin/AgentsView.vue`

**修改文件** (~10个)：
- `backend/ent/schema/user.go` — 新增代理字段
- `backend/internal/service/domain_constants.go` — 新增常量
- `backend/internal/service/payment_service.go` — 支付回调中触发佣金
- `backend/internal/handler/handler.go` — 新增 Agent handler
- `backend/internal/server/routes/user.go` — 代理路由
- `backend/internal/server/routes/admin.go` — 管理员代理路由
- `backend/internal/repository/wire.go`
- `backend/internal/service/wire.go`
- `backend/internal/handler/wire.go`
- `frontend/src/router/index.ts` — 新增路由
- `frontend/src/types/index.ts` — 新增类型
