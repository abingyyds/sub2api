# 支付宝集成完成说明

## 已完成的工作

### 1. 添加依赖
- 在 `go.mod` 中添加了支付宝 SDK: `github.com/smartwalle/alipay/v3 v3.2.23`
- 修正了 Go 版本号从 1.25.5 到 1.22.5

### 2. 数据库变更
- 创建迁移文件: `migrations/066_add_alipay_support.sql`
- 在 `payment_orders` 表中添加 `alipay_trade_no` 字段
- 更新 Schema: `ent/schema/payment_order.go` 添加支付宝交易号字段

### 3. 常量定义
- 在 `domain_constants.go` 中添加:
  - 支付方式常量: `PaymentMethodAlipayNative`
  - 支付宝配置键: `SettingKeyAlipayEnabled`, `SettingKeyAlipayAppID` 等

### 4. 配置管理
- 更新 `SystemSettings` 结构体，添加支付宝配置字段
- 在 `setting_service.go` 中:
  - `parseSettings()` 方法添加支付宝配置解析
  - `UpdateSettings()` 方法添加支付宝配置保存
  - `InitializeDefaultSettings()` 添加支付宝默认配置

### 5. 支付服务
- 更新 `PaymentOrder` 结构体，添加 `AlipayTradeNo` 字段
- 修改 `CreateOrder()` 方法，支持 `payMethod` 参数
- 修改 `CreateRechargeOrder()` 方法，支持 `payMethod` 参数
- 添加支付宝相关方法:
  - `initAlipayClient()`: 初始化支付宝客户端
  - `createAlipayNativeOrder()`: 创建支付宝扫码支付订单
  - `HandleAlipayNotify()`: 处理支付宝回调通知

### 6. Repository 层
- 更新 `UpdateStatus()` 方法，根据支付方式设置对应的交易号字段
- 更新 `toServicePaymentOrder()` 方法，映射 `AlipayTradeNo` 字段

### 7. Handler 层
- 更新请求 DTO:
  - `CreateOrderRequest` 添加 `PayMethod` 字段
  - `CreateRechargeRequest` 添加 `PayMethod` 字段
- 修改 `CreateOrder()` 和 `CreateRecharge()` 方法，传递支付方式参数
- 添加 `AlipayNotify()` 方法处理支付宝回调

### 8. 路由配置
- 在 `routes/user.go` 中添加支付宝回调路由: `POST /api/v1/payment/alipay/notify`

### 9. 管理后台
- 更新 `dto/settings.go`，添加支付宝配置字段
- 更新 `admin/setting_handler.go`:
  - `GetSettings()` 返回支付宝配置
  - `UpdateSettingsRequest` 添加支付宝字段
  - `UpdateSettings()` 保存支付宝配置

## API 使用说明

### 创建订单（支持支付宝）
```bash
POST /api/v1/payment/orders
{
  "plan_key": "group_1",
  "promo_code": "",
  "pay_method": "alipay"  // "wechat" | "alipay"
}
```

### 创建充值订单（支持支付宝）
```bash
POST /api/v1/payment/recharge
{
  "amount": 100.0,
  "promo_code": "",
  "pay_method": "alipay"  // "wechat" | "alipay"
}
```

### 支付宝回调
```
POST /api/v1/payment/alipay/notify
```

## 管理后台配置

在管理后台的系统设置中，需要配置以下支付宝参数：

1. **alipay_enabled**: 是否启用支付宝支付
2. **alipay_app_id**: 支付宝应用 ID
3. **alipay_private_key**: 应用私钥（RSA2）
4. **alipay_public_key**: 支付宝公钥
5. **alipay_notify_url**: 异步通知地址
6. **alipay_is_production**: 是否为生产环境

## 注意事项

1. 支付宝 SDK 需要 Go 版本支持，当前项目使用 Go 1.22.5
2. 支付宝私钥和公钥需要在支付宝开放平台配置
3. 回调地址必须是公网可访问的 HTTPS 地址
4. 支付宝使用当面付（扫码支付）方式
5. 回调验签由 SDK 自动处理

## 下一步工作

1. 运行数据库迁移: `go run cmd/migrate/main.go up`
2. 在管理后台配置支付宝参数
3. 测试支付宝支付流程
4. 前端添加支付方式选择界面
