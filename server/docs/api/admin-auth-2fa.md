# 管理员双因子认证 API 文档

## 概述

本文档描述了管理员双因子认证相关的 API 端点。这些 API 允许管理员启用、禁用、验证双因子认证，以及管理备用恢复码。

## 基础信息

- **基础路径**: `/api/admin/auth`
- **认证方式**: Bearer Token (JWT)
- **内容类型**: `application/json`

## API 端点

### 1. 启用双因子认证

**端点**: `POST /api/admin/auth/enable-2fa`

**描述**: 为当前登录的管理员启用双因子认证，生成 TOTP 密钥和备用恢复码。

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**: 无需请求体（用户ID从认证token中获取）

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "secret": "JBSW Y3DP EHPK 3PXP",
    "qrCodeURL": "otpauth://totp/HotGo:admin_1?secret=JBSWY3DPEHPK3PXP&issuer=HotGo",
    "backupCodes": [
      "12345678",
      "87654321",
      "11223344",
      "44332211",
      "55667788",
      "88776655",
      "99001122",
      "22110099"
    ]
  }
}
```

**错误响应**:
- `401`: 用户未登录
- `500`: 服务器内部错误

### 2. 验证双因子认证设置

**端点**: `POST /api/admin/auth/verify-2fa-setup`

**描述**: 验证 TOTP 验证码以确认启用双因子认证。

**请求体**:
```json
{
  "code": "123456"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "双因子认证已成功启用"
}
```

**错误响应**:
- `400`: 验证码无效或已过期
- `401`: 用户未登录
- `500`: 服务器内部错误

### 3. 禁用双因子认证

**端点**: `POST /api/admin/auth/disable-2fa`

**描述**: 禁用当前用户的双因子认证。

**请求体**:
```json
{
  "code": "123456"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "双因子认证已禁用"
}
```

**错误响应**:
- `400`: 验证码无效或双因子认证未启用
- `401`: 用户未登录
- `500`: 服务器内部错误

### 4. 获取双因子认证状态

**端点**: `GET /api/admin/auth/2fa-status`

**描述**: 获取当前用户的双因子认证状态信息。

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "isEnabled": true,
    "backupCodesCount": 6,
    "lastUsedAt": "2024-01-15T10:30:00Z"
  }
}
```

**错误响应**:
- `401`: 用户未登录
- `500`: 服务器内部错误

### 5. 重新生成备用恢复码

**端点**: `POST /api/admin/auth/regenerate-backup-codes`

**描述**: 重新生成备用恢复码（需要提供有效的 TOTP 验证码）。

**请求体**:
```json
{
  "code": "123456"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "备用恢复码已重新生成",
  "data": {
    "backupCodes": [
      "12345678",
      "87654321",
      "11223344",
      "44332211",
      "55667788",
      "88776655",
      "99001122",
      "22110099"
    ]
  }
}
```

**错误响应**:
- `400`: 验证码无效或双因子认证未启用
- `401`: 用户未登录
- `500`: 服务器内部错误

### 6. 验证双因子认证

**端点**: `POST /api/admin/auth/verify-2fa`

**描述**: 验证 TOTP 验证码或备用恢复码。

**请求体**:
```json
{
  "code": "123456",
  "backupCode": ""
}
```

或使用备用恢复码:
```json
{
  "code": "",
  "backupCode": "12345678"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "验证成功"
}
```

**错误响应**:
- `400`: 验证码无效或双因子认证未启用
- `401`: 用户未登录
- `500`: 服务器内部错误

## 错误代码说明

| 错误代码 | 描述 |
|---------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（用户未登录） |
| 403 | 禁止访问 |
| 500 | 服务器内部错误 |

## 安全注意事项

1. **密钥存储**: TOTP 密钥在数据库中加密存储
2. **备用恢复码**: 备用恢复码经过哈希处理后存储，使用后立即失效
3. **验证码时效**: TOTP 验证码有效期为 30 秒
4. **重放攻击防护**: 已使用的验证码在时间窗口内不能重复使用
5. **审计日志**: 所有双因子认证操作都会记录审计日志

## 使用流程

### 启用双因子认证流程

1. 调用 `POST /api/admin/auth/enable-2fa` 获取 TOTP 密钥和二维码
2. 用户使用认证器应用扫描二维码或手动输入密钥
3. 调用 `POST /api/admin/auth/verify-2fa-setup` 验证设置
4. 保存备用恢复码到安全位置

### 登录验证流程

1. 用户完成用户名/密码验证
2. 系统检查是否启用双因子认证
3. 如已启用，要求用户提供 TOTP 验证码
4. 调用 `POST /api/admin/auth/verify-2fa` 验证
5. 验证通过后完成登录

### 禁用双因子认证流程

1. 用户提供当前有效的 TOTP 验证码或备用恢复码
2. 调用 `POST /api/admin/auth/disable-2fa` 禁用
3. 系统清除相关的双因子认证数据

## 示例代码

### JavaScript (Axios)

```javascript
// 启用双因子认证
const enable2FA = async () => {
  try {
    const response = await axios.post('/api/admin/auth/enable-2fa', {}, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    console.log('2FA enabled:', response.data);
  } catch (error) {
    console.error('Error enabling 2FA:', error.response.data);
  }
};

// 验证设置
const verify2FASetup = async (code) => {
  try {
    const response = await axios.post('/api/admin/auth/verify-2fa-setup', {
      code: code
    }, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    console.log('2FA setup verified:', response.data);
  } catch (error) {
    console.error('Error verifying 2FA setup:', error.response.data);
  }
};
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Enable2FAResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Secret      string   `json:"secret"`
        QRCodeURL   string   `json:"qrCodeURL"`
        BackupCodes []string `json:"backupCodes"`
    } `json:"data"`
}

func enable2FA(token string) (*Enable2FAResponse, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", "/api/admin/auth/enable-2fa", nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result Enable2FAResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

## 更新日志

- **v1.0.0** (2024-01-15): 初始版本，包含所有基础双因子认证功能