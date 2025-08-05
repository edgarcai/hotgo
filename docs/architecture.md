# HotGo Architecture Document

## 文档创建设置

**模式**: 交互式 (Interactive)
**输出文件**: docs/architecture.md
**项目名称**: HotGo
**模板版本**: architecture-template-v2

---

## Introduction

本文档概述了 HotGo 项目的整体架构，包括后端系统、共享服务和非UI特定关注点。其主要目标是作为AI驱动开发的指导性架构蓝图，确保一致性并遵循选定的模式和技术。

**与前端架构的关系:**
如果项目包含重要的用户界面，将有单独的前端架构文档详细说明前端特定的设计，并且必须与本文档结合使用。本文档中记录的核心技术栈选择（见"技术栈"部分）对整个项目具有决定性意义，包括任何前端组件。

### 启动模板或现有项目

在进一步进行架构设计之前，我需要检查项目是否基于启动模板或现有代码库。

**分析结果:**
HotGo 是一个现有的成熟项目，基于以下技术栈：
- 后端：Go 1.24.4 + GoFrame v2.9.1 框架
- 前端：Vue 3 + TypeScript + Naive UI
- 数据库：MySQL + Redis
- 架构模式：前后端分离的单体应用

项目已有完整的目录结构和功能实现，当前正在开发双因素认证功能。

### 变更日志

| 日期 | 版本 | 描述 | 作者 |
|------|------|------|------|
| 2025-01-27 | 1.0 | 初始架构文档创建，基于现有HotGo项目 | Architect |

---

## High Level Architecture

### Technical Summary

HotGo 是一个现代化的企业级全栈应用框架，采用前后端分离的单体架构模式。后端基于 Go 语言和 GoFrame v2.9.1 框架构建，提供 RESTful API 服务；前端使用 Vue 3 + TypeScript + Naive UI 技术栈，通过 HTTP 请求与后端通信。系统采用分层架构模式，包含控制器层、业务逻辑层、数据访问层和模型层，支持插件化扩展和多租户架构。当前正在增强双因素认证功能，以提升系统安全性。

### High Level Overview

**架构风格**: 前后端分离的单体应用 (Monolithic with Frontend Separation)
- 后端：单一 Go 服务，包含所有业务逻辑和 API
- 前端：独立的 Vue 3 SPA 应用
- 通信：RESTful API，JSON 数据格式

**仓库结构**: 单仓库多模块 (Monorepo)
- 后端代码位于 `server/` 目录
- 前端代码位于 `web/` 目录
- 共享文档位于 `docs/` 目录

**服务架构**: 分层单体架构
- API 层：HTTP 路由和控制器
- 业务逻辑层：核心业务处理
- 数据访问层：数据库操作抽象
- 基础设施层：缓存、消息队列、外部服务集成

**用户交互流程**:
1. 用户通过浏览器访问 Vue 前端应用
2. 前端通过 Axios 发送 HTTP 请求到后端 API
3. 后端验证 JWT 令牌，处理业务逻辑
4. 数据持久化到 MySQL，缓存到 Redis
5. 返回 JSON 响应给前端渲染

**关键架构决策**:
- 选择单体架构以简化部署和开发，适合中小型团队
- 前后端分离支持独立开发和部署
- 基于 GoFrame 框架提供完整的企业级功能
- 插件化设计支持功能扩展

### High Level Project Diagram

```mermaid
graph TB
    User[用户] --> Web[Vue 3 前端应用]
    Web --> API[GoFrame API 服务]
    
    API --> Auth[认证中间件]
    API --> Router[路由层]
    Router --> Controller[控制器层]
    Controller --> Logic[业务逻辑层]
    Logic --> Service[服务接口层]
    Service --> DAO[数据访问层]
    
    DAO --> MySQL[(MySQL 数据库)]
    Logic --> Redis[(Redis 缓存)]
    
    API --> External[外部服务]
    External --> SMS[短信服务]
    External --> Payment[支付服务]
    External --> Storage[对象存储]
    
    subgraph "插件系统"
        Addons[插件模块]
    end
    
    Logic --> Addons
    
    subgraph "开发中功能"
        TwoFA[双因素认证]
    end
    
    Auth --> TwoFA
```

### Architectural and Design Patterns

基于 HotGo 项目的实际实现，以下是采用的关键架构和设计模式：

- **分层架构模式 (Layered Architecture)**: 采用经典的四层架构 - API层、业务逻辑层、数据访问层、数据模型层。_理由_: 提供清晰的关注点分离，便于维护和测试，符合 GoFrame 框架的最佳实践。

- **依赖注入模式 (Dependency Injection)**: 通过 GoFrame 的服务容器管理依赖关系。_理由_: 提高代码的可测试性和模块化程度，支持接口驱动开发。

- **仓储模式 (Repository Pattern)**: 数据访问层抽象数据库操作，通过 DAO 层实现。_理由_: 将业务逻辑与数据访问逻辑分离，支持数据库切换和单元测试。

- **中间件模式 (Middleware Pattern)**: 认证、权限、日志等横切关注点通过中间件实现。_理由_: 提供可重用的功能组件，支持请求处理管道。

- **插件架构模式 (Plugin Architecture)**: 通过 addons 目录支持功能扩展。_理由_: 支持第三方功能集成，保持核心系统的简洁性。

- **RESTful API 模式**: 统一的 HTTP API 设计，遵循 REST 原则。_理由_: 提供标准化的接口设计，便于前端集成和第三方调用。

- **JWT 认证模式**: 基于令牌的无状态认证机制。_理由_: 支持分布式部署，减少服务器端状态管理复杂度。

---

**需要用户确认**: 以上高层架构分析是否准确反映了您对 HotGo 项目的理解？是否需要调整或补充任何架构决策？

请选择以下选项：

1. 继续下一部分
2. 架构风格深入分析 - 探讨单体 vs 微服务的权衡
3. 技术债务评估 - 分析当前架构的局限性
4. 扩展性规划 - 讨论未来架构演进路径
5. 性能考量分析 - 评估当前架构的性能特征
6. 安全架构审查 - 深入分析安全设计模式
7. 集成模式探讨 - 分析外部服务集成策略
8. 数据流分析 - 详细梳理数据在系统中的流转
9. 架构决策记录 - 文档化关键架构决策的背景和理由

**YOLO 模式已启用** - 将快速完成剩余架构文档部分

## Data Models

基于 HotGo 项目的数据库结构分析，以下是核心数据模型：

### 核心实体关系

```mermaid
erDiagram
    User ||--o{ UserRole : has
    Role ||--o{ RoleMenu : contains
    Menu ||--o{ RoleMenu : assigned_to
    User ||--o{ LoginLog : generates
    User ||--o{ OperateLog : performs
    User ||--o{ TwoFactorAuth : enables
    User ||--o{ TwoFactorBackupCode : has
    
    User {
        int id PK
        string username
        string email
        string mobile
        string password_hash
        int status
        datetime created_at
        datetime updated_at
        boolean two_factor_enabled
    }
    
    TwoFactorAuth {
        int id PK
        int user_id FK
        string secret_key
        string qr_code_url
        boolean is_active
        datetime created_at
        datetime verified_at
    }
    
    TwoFactorBackupCode {
        int id PK
        int user_id FK
        string code_hash
        boolean is_used
        datetime created_at
        datetime used_at
    }
    
    Role {
        int id PK
        string name
        string key
        int status
        datetime created_at
    }
    
    Menu {
        int id PK
        int pid
        string title
        string name
        string path
        string component
        int type
        int status
    }
```

### 关键数据模型说明

**用户认证模型**:
- `User`: 核心用户实体，包含基本信息和双因素认证状态
- `TwoFactorAuth`: 双因素认证配置，存储TOTP密钥和状态
- `TwoFactorBackupCode`: 备用恢复代码，用于紧急访问

**权限管理模型**:
- `Role`: 角色定义，支持多角色分配
- `Menu`: 菜单权限，支持层级结构
- `RoleMenu`: 角色菜单关联表

**审计日志模型**:
- `LoginLog`: 登录日志，记录认证事件
- `OperateLog`: 操作日志，记录用户行为

## Components

### 后端组件架构

#### 1. 认证组件 (Authentication)
**职责**: 用户身份验证和会话管理
**接口**: 
- `Login(username, password) -> JWT Token`
- `VerifyTwoFactor(token, code) -> AuthResult`
- `RefreshToken(token) -> New Token`

**依赖**: 
- Database (用户数据)
- Redis (会话缓存)
- JWT Library (令牌生成)

**技术栈**: GoFrame Auth + JWT + Redis

#### 2. 双因素认证组件 (TwoFactor)
**职责**: TOTP双因素认证管理
**接口**:
- `EnableTwoFactor(userID) -> QRCode, Secret`
- `VerifyTOTP(userID, code) -> boolean`
- `GenerateBackupCodes(userID) -> []string`

**依赖**:
- TOTP Library (otp生成)
- QR Code Generator
- Crypto (密钥生成)

**技术栈**: pquerna/otp + skip2/go-qrcode

#### 3. 权限管理组件 (Authorization)
**职责**: 基于角色的访问控制
**接口**:
- `CheckPermission(userID, resource) -> boolean`
- `GetUserRoles(userID) -> []Role`
- `AssignRole(userID, roleID) -> error`

**依赖**:
- Database (权限数据)
- Cache (权限缓存)

**技术栈**: GoFrame ORM + Redis Cache

#### 4. API网关组件 (Gateway)
**职责**: 请求路由、中间件处理、限流
**接口**:
- HTTP路由处理
- 中间件链管理
- 请求/响应拦截

**依赖**:
- GoFrame Router
- 认证组件
- 日志组件

**技术栈**: GoFrame HTTP Server + Middleware

### 前端组件架构

#### 1. 认证模块 (Auth Module)
**职责**: 用户登录、双因素认证界面
**组件**:
- `LoginForm.vue`: 登录表单
- `TwoFactorSetup.vue`: 双因素认证设置
- `TwoFactorVerify.vue`: 双因素认证验证

**依赖**: Pinia Store, Axios, Naive UI

#### 2. 路由守卫 (Route Guards)
**职责**: 页面访问权限控制
**功能**:
- 登录状态检查
- 权限验证
- 双因素认证状态检查

**依赖**: Vue Router, Auth Store

#### 3. 状态管理 (State Management)
**职责**: 全局状态管理
**Store模块**:
- `authStore`: 认证状态
- `userStore`: 用户信息
- `permissionStore`: 权限数据

**技术栈**: Pinia + TypeScript

## External APIs

### 短信服务集成
**目的**: 双因素认证短信发送
**提供商**: 阿里云SMS
**文档**: https://help.aliyun.com/product/44282.html
**认证**: AccessKey + AccessSecret
**速率限制**: 1000条/小时
**关键端点**:
- `SendSms`: 发送验证码短信
- `QuerySendDetails`: 查询发送状态

### 支付服务集成
**目的**: 在线支付处理
**提供商**: 微信支付
**文档**: https://pay.weixin.qq.com/wiki/doc/api/index.html
**认证**: 商户号 + API密钥
**速率限制**: 600次/分钟
**关键端点**:
- `UnifiedOrder`: 统一下单
- `OrderQuery`: 订单查询
- `Refund`: 申请退款

### 对象存储集成
**目的**: 文件上传和存储
**提供商**: MinIO (私有部署)
**文档**: https://docs.min.io/
**认证**: AccessKey + SecretKey
**速率限制**: 无限制
**关键端点**:
- `PutObject`: 上传文件
- `GetObject`: 下载文件
- `DeleteObject`: 删除文件

## Core Workflows

### 用户登录工作流

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as Auth API
    participant D as Database
    participant R as Redis
    
    U->>F: 输入用户名密码
    F->>A: POST /api/auth/login
    A->>D: 验证用户凭据
    D-->>A: 用户信息
    
    alt 双因素认证已启用
        A-->>F: 需要双因素认证
        F->>U: 显示TOTP输入界面
        U->>F: 输入TOTP代码
        F->>A: POST /api/auth/verify-2fa
        A->>A: 验证TOTP代码
    end
    
    A->>A: 生成JWT Token
    A->>R: 缓存会话信息
    A-->>F: 返回Token和用户信息
    F->>F: 存储Token到localStorage
    F-->>U: 跳转到主页面
```

### 双因素认证设置工作流

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as Auth API
    participant T as TOTP Service
    participant D as Database
    
    U->>F: 点击启用双因素认证
    F->>A: POST /api/user/enable-2fa
    A->>T: 生成TOTP密钥
    T-->>A: 密钥和二维码
    A->>D: 保存密钥(未激活状态)
    A-->>F: 返回二维码URL
    F->>U: 显示二维码
    
    U->>U: 扫描二维码到认证器
    U->>F: 输入验证码
    F->>A: POST /api/user/verify-2fa-setup
    A->>T: 验证TOTP代码
    
    alt 验证成功
        A->>D: 激活双因素认证
        A->>A: 生成备用代码
        A-->>F: 返回备用代码
        F->>U: 显示备用代码
    else 验证失败
        A-->>F: 返回错误
        F->>U: 显示错误信息
    end
```

## REST API Specification

```yaml
openapi: 3.0.3
info:
  title: HotGo Backend API
  description: HotGo项目后端API规范
  version: 1.0.0
  contact:
    name: HotGo Team
    email: support@hotgo.com

servers:
  - url: http://localhost:8000/api
    description: 开发环境
  - url: https://api.hotgo.com
    description: 生产环境

paths:
  /auth/login:
    post:
      summary: 用户登录
      tags: [Authentication]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                username:
                  type: string
                  example: "admin"
                password:
                  type: string
                  example: "password123"
                captcha:
                  type: string
                  example: "abc123"
              required: [username, password]
      responses:
        '200':
          description: 登录成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 0
                  message:
                    type: string
                    example: "success"
                  data:
                    type: object
                    properties:
                      token:
                        type: string
                      user:
                        $ref: '#/components/schemas/User'
                      needTwoFactor:
                        type: boolean
        '401':
          description: 认证失败
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/verify-2fa:
    post:
      summary: 双因素认证验证
      tags: [Authentication]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                token:
                  type: string
                  description: 临时认证令牌
                code:
                  type: string
                  description: TOTP验证码
                  example: "123456"
              required: [token, code]
      responses:
        '200':
          description: 验证成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 0
                  data:
                    type: object
                    properties:
                      token:
                        type: string
                      user:
                        $ref: '#/components/schemas/User'

  /user/enable-2fa:
    post:
      summary: 启用双因素认证
      tags: [User]
      security:
        - BearerAuth: []
      responses:
        '200':
          description: 成功生成双因素认证配置
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 0
                  data:
                    type: object
                    properties:
                      qrCodeUrl:
                        type: string
                        description: 二维码图片URL
                      secret:
                        type: string
                        description: TOTP密钥(仅用于手动输入)

  /user/verify-2fa-setup:
    post:
      summary: 验证双因素认证设置
      tags: [User]
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: string
                  description: TOTP验证码
                  example: "123456"
              required: [code]
      responses:
        '200':
          description: 验证成功，双因素认证已启用
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 0
                  data:
                    type: object
                    properties:
                      backupCodes:
                        type: array
                        items:
                          type: string
                        description: 备用恢复代码

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
          example: 1
        username:
          type: string
          example: "admin"
        email:
          type: string
          example: "admin@hotgo.com"
        mobile:
          type: string
          example: "13800138000"
        status:
          type: integer
          example: 1
        twoFactorEnabled:
          type: boolean
          example: false
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time
    
    Error:
      type: object
      properties:
        code:
          type: integer
          example: 1
        message:
          type: string
          example: "错误信息"
        data:
          type: object
          nullable: true

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

## Database Schema

基于 MySQL 8.0+ 的数据库设计，针对性能和可扩展性进行优化：

```sql
-- 用户表
CREATE TABLE `hg_admin_member` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态:1正常,2禁用',
  `two_factor_enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用双因素认证',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员用户表';

-- 双因素认证表
CREATE TABLE `hg_admin_two_factor_auth` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `secret_key` varchar(255) NOT NULL COMMENT 'TOTP密钥',
  `qr_code_url` varchar(500) DEFAULT NULL COMMENT '二维码URL',
  `is_active` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否激活',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `verified_at` datetime DEFAULT NULL COMMENT '验证时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  CONSTRAINT `fk_two_factor_user` FOREIGN KEY (`user_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='双因素认证表';

-- 双因素认证备用代码表
CREATE TABLE `hg_admin_two_factor_backup_code` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `code_hash` varchar(255) NOT NULL COMMENT '备用代码哈希',
  `is_used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已使用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `used_at` datetime DEFAULT NULL COMMENT '使用时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_used` (`is_used`),
  CONSTRAINT `fk_backup_code_user` FOREIGN KEY (`user_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='双因素认证备用代码表';

-- 角色表
CREATE TABLE `hg_admin_role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `key` varchar(50) NOT NULL COMMENT '角色标识',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态:1正常,2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE `hg_admin_member_role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `member_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_member_role` (`member_id`,`role_id`),
  KEY `idx_role_id` (`role_id`),
  CONSTRAINT `fk_member_role_member` FOREIGN KEY (`member_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_member_role_role` FOREIGN KEY (`role_id`) REFERENCES `hg_admin_role` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 登录日志表
CREATE TABLE `hg_admin_login_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `member_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `login_ip` varchar(45) NOT NULL COMMENT '登录IP',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `status` tinyint(4) NOT NULL COMMENT '状态:1成功,2失败',
  `error_msg` varchar(255) DEFAULT NULL COMMENT '错误信息',
  `two_factor_used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否使用双因素认证',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_member_id` (`member_id`),
  KEY `idx_login_ip` (`login_ip`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';
```

### 性能优化策略

**索引策略**:
- 主键使用自增ID，提供最佳插入性能
- 为常用查询字段创建索引（用户名、邮箱、状态等）
- 复合索引用于多字段查询优化

**分区策略**:
- 日志表按月分区，提高查询性能
- 历史数据归档策略，保持表大小可控

**缓存策略**:
- 用户信息缓存（Redis，TTL: 30分钟）
- 权限数据缓存（Redis，TTL: 1小时）
- 会话数据缓存（Redis，TTL: 24小时）

## Source Code Tree

基于 GoFrame 框架的项目结构，体现分层架构和模块化设计：

```
server/
├── api/                          # API接口定义
│   ├── admin/                    # 管理后台API
│   │   ├── auth.go              # 认证相关API
│   │   ├── user.go              # 用户管理API
│   │   └── two_factor.go        # 双因素认证API
│   └── home/                     # 前台API
├── internal/                     # 内部业务逻辑
│   ├── cmd/                      # 命令行工具
│   │   └── cmd.go               # 主命令入口
│   ├── controller/               # 控制器层
│   │   ├── admin/               # 管理后台控制器
│   │   │   ├── auth.go          # 认证控制器
│   │   │   ├── user.go          # 用户控制器
│   │   │   └── two_factor.go    # 双因素认证控制器
│   │   └── home/                # 前台控制器
│   ├── dao/                      # 数据访问对象
│   │   ├── admin_member.go      # 用户DAO
│   │   ├── admin_role.go        # 角色DAO
│   │   ├── two_factor_auth.go   # 双因素认证DAO
│   │   └── login_log.go         # 登录日志DAO
│   ├── logic/                    # 业务逻辑层
│   │   ├── admin/               # 管理后台业务逻辑
│   │   │   ├── auth.go          # 认证业务逻辑
│   │   │   ├── user.go          # 用户业务逻辑
│   │   │   └── two_factor.go    # 双因素认证业务逻辑
│   │   └── middleware/          # 中间件
│   │       ├── auth.go          # 认证中间件
│   │       ├── cors.go          # 跨域中间件
│   │       └── rate_limit.go    # 限流中间件
│   ├── model/                    # 数据模型
│   │   ├── entity/              # 数据库实体
│   │   │   ├── admin_member.go  # 用户实体
│   │   │   ├── admin_role.go    # 角色实体
│   │   │   └── two_factor_auth.go # 双因素认证实体
│   │   └── do/                  # 数据对象
│   ├── service/                  # 服务接口定义
│   │   ├── auth.go              # 认证服务接口
│   │   ├── user.go              # 用户服务接口
│   │   └── two_factor.go        # 双因素认证服务接口
│   └── router/                   # 路由配置
│       ├── admin.go             # 管理后台路由
│       └── home.go              # 前台路由
├── manifest/                     # 配置文件
│   ├── config/                  # 应用配置
│   │   ├── config.yaml          # 主配置文件
│   │   ├── config.example.yaml  # 配置示例
│   │   └── database.yaml        # 数据库配置
│   └── deploy/                  # 部署配置
│       ├── docker-compose.yml   # Docker编排
│       └── nginx.conf           # Nginx配置
├── resource/                     # 静态资源
│   ├── template/                # 模板文件
│   └── public/                  # 公共资源
├── utility/                      # 工具库
│   ├── encrypt/                 # 加密工具
│   ├── validate/                # 验证工具
│   └── convert/                 # 转换工具
├── go.mod                        # Go模块定义
├── go.sum                        # 依赖校验
└── main.go                       # 应用入口

web/
├── src/
│   ├── api/                     # API请求封装
│   │   ├── auth.ts              # 认证API
│   │   ├── user.ts              # 用户API
│   │   └── two-factor.ts        # 双因素认证API
│   ├── components/              # 公共组件
│   │   ├── TwoFactor/           # 双因素认证组件
│   │   │   ├── Setup.vue        # 设置组件
│   │   │   ├── Verify.vue       # 验证组件
│   │   │   └── BackupCodes.vue  # 备用代码组件
│   │   └── Auth/                # 认证组件
│   │       ├── LoginForm.vue    # 登录表单
│   │       └── AuthGuard.vue    # 认证守卫
│   ├── store/                   # 状态管理
│   │   ├── modules/             # 状态模块
│   │   │   ├── auth.ts          # 认证状态
│   │   │   ├── user.ts          # 用户状态
│   │   │   └── permission.ts    # 权限状态
│   │   └── index.ts             # 状态入口
│   ├── router/                  # 路由配置
│   │   ├── index.ts             # 路由入口
│   │   ├── guards.ts            # 路由守卫
│   │   └── modules/             # 路由模块
│   ├── views/                   # 页面组件
│   │   ├── auth/                # 认证页面
│   │   │   ├── login.vue        # 登录页面
│   │   │   └── two-factor.vue   # 双因素认证页面
│   │   └── admin/               # 管理页面
│   ├── utils/                   # 工具函数
│   │   ├── auth.ts              # 认证工具
│   │   ├── request.ts           # 请求工具
│   │   └── crypto.ts            # 加密工具
│   └── types/                   # 类型定义
│       ├── auth.d.ts            # 认证类型
│       └── user.d.ts            # 用户类型
├── package.json                 # 前端依赖
├── vite.config.ts              # Vite配置
└── tsconfig.json               # TypeScript配置
```

### 架构设计原则

**后端架构**:
- **分层架构**: API -> Controller -> Logic -> Service -> DAO
- **依赖注入**: 通过GoFrame的依赖注入容器管理组件
- **接口隔离**: Service层定义接口，Logic层实现具体业务逻辑
- **单一职责**: 每个模块专注于特定功能领域

**前端架构**:
- **组件化**: 可复用的Vue组件设计
- **状态管理**: Pinia集中管理应用状态
- **类型安全**: TypeScript提供编译时类型检查
- **模块化**: 按功能模块组织代码结构

## Infrastructure and Deployment

### 部署架构

```mermaid
graph TB
    subgraph "Production Environment"
        LB[Load Balancer<br/>Nginx]
        
        subgraph "Application Tier"
            APP1[HotGo Instance 1<br/>:8000]
            APP2[HotGo Instance 2<br/>:8000]
        end
        
        subgraph "Data Tier"
            DB[(MySQL 8.0<br/>Master-Slave)]
            REDIS[(Redis Cluster<br/>Cache & Session)]
            MINIO[(MinIO<br/>Object Storage)]
        end
        
        subgraph "Monitoring"
            PROM[Prometheus]
            GRAF[Grafana]
            ALERT[AlertManager]
        end
    end
    
    LB --> APP1
    LB --> APP2
    APP1 --> DB
    APP2 --> DB
    APP1 --> REDIS
    APP2 --> REDIS
    APP1 --> MINIO
    APP2 --> MINIO
    
    PROM --> APP1
    PROM --> APP2
    PROM --> DB
    PROM --> REDIS
    GRAF --> PROM
    ALERT --> PROM
```

### 部署策略

**容器化部署**:
```yaml
# docker-compose.yml
version: '3.8'
services:
  hotgo-app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - GF_GERROR_BRIEF=true
      - GF_DATABASE_LINK=mysql:root:123456789a@tcp(mysql:3306)/hotgo
      - GF_REDIS_ADDRESS=redis:6379
    depends_on:
      - mysql
      - redis
    restart: unless-stopped
    
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 123456789a
      MYSQL_DATABASE: hotgo
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    restart: unless-stopped
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - hotgo-app
    restart: unless-stopped

volumes:
  mysql_data:
```

**环境管理**:
- **开发环境**: 本地Docker Compose
- **测试环境**: Kubernetes集群
- **生产环境**: Kubernetes + Helm Charts

**部署流程**:
1. **代码提交** → Git仓库
2. **CI/CD触发** → GitHub Actions/GitLab CI
3. **构建镜像** → Docker Build
4. **安全扫描** → 容器安全检查
5. **部署测试** → 测试环境验证
6. **生产部署** → 蓝绿部署/滚动更新
7. **健康检查** → 服务状态监控

**回滚策略**:
- **快速回滚**: Kubernetes Rollout Undo
- **数据库回滚**: 数据库备份恢复
- **配置回滚**: ConfigMap版本管理

## Error Handling Strategy

### 错误分类和处理

**错误类别**:
```go
// 业务错误代码定义
const (
    // 认证相关错误 (1000-1999)
    ErrInvalidCredentials    = 1001  // 用户名或密码错误
    ErrAccountDisabled       = 1002  // 账户已禁用
    ErrTwoFactorRequired     = 1003  // 需要双因素认证
    ErrInvalidTwoFactorCode  = 1004  // 双因素认证码无效
    ErrTwoFactorNotEnabled   = 1005  // 双因素认证未启用
    
    // 权限相关错误 (2000-2999)
    ErrPermissionDenied      = 2001  // 权限不足
    ErrResourceNotFound      = 2002  // 资源不存在
    ErrOperationNotAllowed   = 2003  // 操作不被允许
    
    // 系统相关错误 (9000-9999)
    ErrInternalServer        = 9001  // 内部服务器错误
    ErrDatabaseConnection    = 9002  // 数据库连接错误
    ErrRedisConnection       = 9003  // Redis连接错误
)
```

**错误处理模式**:
```go
// 统一错误响应结构
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
    TraceID string `json:"trace_id,omitempty"`
}

// 错误处理中间件
func ErrorHandler(r *ghttp.Request) {
    r.Middleware.Next()
    
    if err := r.GetError(); err != nil {
        var response ErrorResponse
        
        switch e := err.(type) {
        case *gerror.Error:
            response.Code = e.Code()
            response.Message = e.Error()
        default:
            response.Code = ErrInternalServer
            response.Message = "Internal server error"
        }
        
        response.TraceID = gtrace.GetTraceID(r.Context())
        
        // 记录错误日志
        g.Log().Error(r.Context(), "Request error", g.Map{
            "error":    err.Error(),
            "trace_id": response.TraceID,
            "path":     r.URL.Path,
            "method":   r.Method,
        })
        
        r.Response.WriteJsonExit(response)
    }
}
```

**日志标准**:
- **结构化日志**: 使用JSON格式，便于日志分析
- **日志级别**: DEBUG, INFO, WARN, ERROR, FATAL
- **上下文信息**: 包含TraceID、用户ID、请求路径等
- **敏感信息**: 密码、令牌等敏感信息不记录到日志

**监控和告警**:
- **错误率监控**: 通过Prometheus监控API错误率
- **响应时间**: 监控API响应时间分布
- **告警规则**: 错误率超过阈值时触发告警

## Coding Standards

### 后端编码标准 (Go)

**核心原则**:
- **简洁性**: 代码应该简洁明了，避免过度设计
- **可读性**: 代码应该自解释，必要时添加注释
- **一致性**: 遵循团队统一的编码风格
- **性能**: 关注性能，避免不必要的内存分配

**命名约定**:
```go
// 包名：小写，简短，有意义
package auth

// 常量：大写，下划线分隔
const (
    DEFAULT_TIMEOUT = 30 * time.Second
    MAX_RETRY_COUNT = 3
)

// 变量：驼峰命名
var userService IUserService

// 函数：驼峰命名，动词开头
func GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
    // 实现
}

// 结构体：驼峰命名，名词
type UserLoginRequest struct {
    Username string `json:"username" v:"required|length:3,50"`
    Password string `json:"password" v:"required|length:6,50"`
}
```

**错误处理**:
```go
// 使用 GoFrame 的错误处理机制
func (s *sUser) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
    // 验证用户凭据
    user, err := s.dao.GetUserByUsername(ctx, req.Username)
    if err != nil {
        return nil, gerror.Wrap(err, "获取用户信息失败")
    }
    
    if user == nil {
        return nil, gerror.NewCode(gcode.CodeNotFound, "用户不存在")
    }
    
    // 验证密码
    if !s.verifyPassword(req.Password, user.PasswordHash) {
        return nil, gerror.NewCode(gcode.CodeInvalidParameter, "密码错误")
    }
    
    return &model.LoginResponse{
        Token: token,
        User:  user,
    }, nil
}
```

**数据库操作**:
```go
// 使用 GoFrame ORM，避免 SQL 注入
func (d *dao) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
    var user entity.User
    err := d.DB().Model("admin_member").Where("id", userID).Scan(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// 事务处理
func (s *sUser) CreateUserWithRole(ctx context.Context, req *model.CreateUserRequest) error {
    return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // 创建用户
        userID, err := tx.Model("admin_member").Data(req.User).InsertAndGetId()
        if err != nil {
            return err
        }
        
        // 分配角色
        _, err = tx.Model("admin_member_role").Data(g.Map{
            "member_id": userID,
            "role_id":   req.RoleID,
        }).Insert()
        return err
    })
}
```

### 前端编码标准 (TypeScript/Vue)

**命名约定**:
```typescript
// 接口：PascalCase，I前缀
interface IUserInfo {
  id: number
  username: string
  email: string
}

// 类型：PascalCase
type LoginStatus = 'pending' | 'success' | 'failed'

// 变量：camelCase
const userInfo: IUserInfo = {
  id: 1,
  username: 'admin',
  email: 'admin@example.com'
}

// 函数：camelCase，动词开头
const getUserInfo = async (userId: number): Promise<IUserInfo> => {
  // 实现
}

// 常量：UPPER_SNAKE_CASE
const API_BASE_URL = '/api'
const MAX_RETRY_COUNT = 3
```

**Vue组件规范**:
```vue
<template>
  <!-- 使用 kebab-case 命名 -->
  <div class="user-profile">
    <user-avatar :src="userInfo.avatar" />
    <user-info :user="userInfo" @update="handleUserUpdate" />
  </div>
</template>

<script setup lang="ts">
// 组件名：PascalCase
defineOptions({
  name: 'UserProfile'
})

// Props 定义
interface Props {
  userId: number
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  readonly: false
})

// Emits 定义
interface Emits {
  update: [user: IUserInfo]
  delete: [userId: number]
}

const emit = defineEmits<Emits>()

// 响应式数据
const userInfo = ref<IUserInfo | null>(null)
const loading = ref(false)

// 计算属性
const displayName = computed(() => {
  return userInfo.value?.realName || userInfo.value?.username || '未知用户'
})

// 方法
const handleUserUpdate = (user: IUserInfo) => {
  userInfo.value = user
  emit('update', user)
}
</script>

<style scoped>
.user-profile {
  /* 使用 BEM 命名规范 */
}
</style>
```

**API调用规范**:
```typescript
// API 模块：统一错误处理
export const userApi = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await request.post<LoginResponse>('/auth/login', data)
      return response.data
    } catch (error) {
      // 统一错误处理
      throw new ApiError('登录失败', error)
    }
  },
  
  async getUserInfo(userId: number): Promise<IUserInfo> {
    const response = await request.get<IUserInfo>(`/user/${userId}`)
    return response.data
  }
}
```

## Testing Strategy and Standards

### 测试理念

**测试金字塔**:
- **单元测试 (70%)**: 测试单个函数/方法的逻辑
- **集成测试 (20%)**: 测试模块间的交互
- **端到端测试 (10%)**: 测试完整的用户流程

### 后端测试

**单元测试示例**:
```go
// internal/logic/auth/auth_test.go
func TestUserLogin(t *testing.T) {
    // 使用 GoFrame 的测试框架
    gtest.C(t, func(t *gtest.T) {
        // 准备测试数据
        ctx := context.Background()
        req := &model.LoginRequest{
            Username: "testuser",
            Password: "password123",
        }
        
        // Mock 依赖
        mockDao := &mockUserDao{
            user: &entity.User{
                ID:           1,
                Username:     "testuser",
                PasswordHash: "$2a$10$...", // bcrypt hash
                Status:       1,
            },
        }
        
        service := &sUser{dao: mockDao}
        
        // 执行测试
        resp, err := service.Login(ctx, req)
        
        // 断言结果
        t.AssertNil(err)
        t.AssertNE(resp, nil)
        t.AssertNE(resp.Token, "")
        t.AssertEQ(resp.User.Username, "testuser")
    })
}

// 双因素认证测试
func TestTwoFactorAuth(t *testing.T) {
    gtest.C(t, func(t *gtest.T) {
        ctx := context.Background()
        userID := int64(1)
        
        // 测试启用双因素认证
        secret, qrCode, err := twoFactorService.Enable(ctx, userID)
        t.AssertNil(err)
        t.AssertNE(secret, "")
        t.AssertNE(qrCode, "")
        
        // 测试验证 TOTP 代码
        code := generateTOTPCode(secret) // 生成有效的 TOTP 代码
        valid, err := twoFactorService.Verify(ctx, userID, code)
        t.AssertNil(err)
        t.AssertEQ(valid, true)
    })
}
```

**集成测试示例**:
```go
// test/integration/auth_test.go
func TestAuthIntegration(t *testing.T) {
    // 启动测试服务器
    s := g.Server()
    s.SetPort(8080)
    s.Start()
    defer s.Shutdown()
    
    // 测试登录 API
    client := g.Client()
    response := client.PostContent("/api/auth/login", g.Map{
        "username": "admin",
        "password": "123456",
    })
    
    var result map[string]interface{}
    json.Unmarshal([]byte(response), &result)
    
    assert.Equal(t, 0, result["code"])
    assert.NotEmpty(t, result["data"].(map[string]interface{})["token"])
}
```

### 前端测试

**组件测试示例**:
```typescript
// tests/components/LoginForm.test.ts
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import LoginForm from '@/components/Auth/LoginForm.vue'

describe('LoginForm', () => {
  it('should render login form correctly', () => {
    const wrapper = mount(LoginForm)
    
    expect(wrapper.find('[data-testid="username-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="password-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="login-button"]').exists()).toBe(true)
  })
  
  it('should emit login event with correct data', async () => {
    const wrapper = mount(LoginForm)
    
    // 填写表单
    await wrapper.find('[data-testid="username-input"]').setValue('admin')
    await wrapper.find('[data-testid="password-input"]').setValue('password')
    
    // 提交表单
    await wrapper.find('[data-testid="login-button"]').trigger('click')
    
    // 验证事件
    expect(wrapper.emitted('login')).toBeTruthy()
    expect(wrapper.emitted('login')[0]).toEqual([{
      username: 'admin',
      password: 'password'
    }])
  })
})
```

**E2E测试示例**:
```typescript
// tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Authentication Flow', () => {
  test('should login successfully with valid credentials', async ({ page }) => {
    // 访问登录页面
    await page.goto('/login')
    
    // 填写登录表单
    await page.fill('[data-testid="username-input"]', 'admin')
    await page.fill('[data-testid="password-input"]', 'password123')
    
    // 点击登录按钮
    await page.click('[data-testid="login-button"]')
    
    // 验证跳转到主页
    await expect(page).toHaveURL('/dashboard')
    
    // 验证用户信息显示
    await expect(page.locator('[data-testid="user-avatar"]')).toBeVisible()
  })
  
  test('should handle two-factor authentication', async ({ page }) => {
    // 登录需要双因素认证的用户
    await page.goto('/login')
    await page.fill('[data-testid="username-input"]', 'user-with-2fa')
    await page.fill('[data-testid="password-input"]', 'password123')
    await page.click('[data-testid="login-button"]')
    
    // 应该跳转到双因素认证页面
    await expect(page).toHaveURL('/auth/two-factor')
    
    // 输入 TOTP 代码
    await page.fill('[data-testid="totp-input"]', '123456')
    await page.click('[data-testid="verify-button"]')
    
    // 验证最终登录成功
    await expect(page).toHaveURL('/dashboard')
  })
})
```

### 测试数据管理

**测试数据库**:
```yaml
# test/config/database.yaml
database:
  default:
    link: "mysql:root:123456789a@tcp(localhost:3306)/hotgo_test"
    debug: true
```

**数据清理策略**:
```go
// 每个测试后清理数据
func tearDown(t *testing.T) {
    // 清理测试数据
    g.DB().Exec("TRUNCATE TABLE hg_admin_member")
    g.DB().Exec("TRUNCATE TABLE hg_admin_two_factor_auth")
    g.DB().Exec("TRUNCATE TABLE hg_admin_login_log")
}
```

## Security

### 输入验证

**后端验证**:
```go
// 使用 GoFrame 的验证器
type LoginRequest struct {
    Username string `json:"username" v:"required|length:3,50|regex:^[a-zA-Z0-9_]+$"`
    Password string `json:"password" v:"required|length:6,50"`
    Captcha  string `json:"captcha" v:"required|length:4,6"`
}

// 自定义验证规则
func init() {
    gvalid.RegisterRule("strong-password", func(ctx context.Context, in gvalid.RuleInput) error {
        password := in.Value.String()
        if !isStrongPassword(password) {
            return gerror.New("密码强度不足，需包含大小写字母、数字和特殊字符")
        }
        return nil
    })
}
```

**前端验证**:
```typescript
// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ]
}
```

### 认证与授权

**JWT 安全配置**:
```go
// JWT 配置
type JWTConfig struct {
    SecretKey     string        `json:"secret_key"`
    ExpireTime    time.Duration `json:"expire_time"`
    RefreshTime   time.Duration `json:"refresh_time"`
    Issuer        string        `json:"issuer"`
    SigningMethod string        `json:"signing_method"`
}

// 生成安全的 JWT Token
func GenerateToken(userID int64, username string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "iat":      time.Now().Unix(),
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
        "iss":      "hotgo",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.JWT.SecretKey))
}
```

**权限检查中间件**:
```go
func AuthMiddleware(r *ghttp.Request) {
    token := r.Header.Get("Authorization")
    if token == "" {
        r.Response.WriteJsonExit(g.Map{
            "code":    401,
            "message": "未提供认证令牌",
        })
        return
    }
    
    // 验证 JWT Token
    claims, err := ValidateToken(strings.TrimPrefix(token, "Bearer "))
    if err != nil {
        r.Response.WriteJsonExit(g.Map{
            "code":    401,
            "message": "无效的认证令牌",
        })
        return
    }
    
    // 设置用户上下文
    r.SetCtxVar("user_id", claims["user_id"])
    r.SetCtxVar("username", claims["username"])
    
    r.Middleware.Next()
}
```

### 密钥管理

**环境变量配置**:
```bash
# .env
JWT_SECRET_KEY=your-super-secret-key-here
DB_PASSWORD=123456789a
REDIS_PASSWORD=redis-password
TWO_FACTOR_ISSUER=HotGo
```

**密钥轮换策略**:
```go
// 支持多个 JWT 密钥，实现平滑轮换
type KeyManager struct {
    currentKey  string
    previousKey string
    nextKey     string
}

func (km *KeyManager) ValidateToken(tokenString string) (*jwt.Token, error) {
    // 尝试使用当前密钥验证
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(km.currentKey), nil
    })
    
    if err != nil && km.previousKey != "" {
        // 尝试使用前一个密钥验证（支持密钥轮换）
        token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(km.previousKey), nil
        })
    }
    
    return token, err
}
```

### API 安全

**速率限制**:
```go
// 基于 IP 的速率限制
func RateLimitMiddleware(r *ghttp.Request) {
    clientIP := r.GetClientIp()
    key := fmt.Sprintf("rate_limit:%s", clientIP)
    
    // 使用 Redis 实现滑动窗口限流
    count, err := g.Redis().Incr(r.Context(), key)
    if err != nil {
        g.Log().Error(r.Context(), "Redis error in rate limit", err)
        r.Middleware.Next()
        return
    }
    
    if count == 1 {
        g.Redis().Expire(r.Context(), key, time.Minute)
    }
    
    if count > 100 { // 每分钟最多 100 次请求
        r.Response.WriteJsonExit(g.Map{
            "code":    429,
            "message": "请求过于频繁，请稍后再试",
        })
        return
    }
    
    r.Middleware.Next()
}
```

**CORS 配置**:
```go
func CORSMiddleware(r *ghttp.Request) {
    corsOptions := r.Response.DefaultCORSOptions()
    corsOptions.AllowDomain = []string{"https://admin.hotgo.com", "http://localhost:3000"}
    corsOptions.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    corsOptions.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
    corsOptions.AllowCredentials = true
    
    r.Response.CORS(corsOptions)
    r.Middleware.Next()
}
```

### 数据保护

**敏感数据加密**:
```go
// 密码哈希
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 双因素认证密钥加密存储
func EncryptTOTPSecret(secret string) (string, error) {
    key := []byte(config.Encryption.Key) // 32 字节密钥
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(secret), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

**数据库安全**:
```sql
-- 创建专用数据库用户
CREATE USER 'hotgo_app'@'%' IDENTIFIED BY 'strong-password-here';
GRANT SELECT, INSERT, UPDATE, DELETE ON hotgo.* TO 'hotgo_app'@'%';

-- 启用 SSL 连接
-- 在配置中使用: mysql://user:pass@tcp(host:port)/db?tls=true
```

### 安全测试

**安全扫描集成**:
```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Gosec Security Scanner
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: './...'
          
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
```

---

## 文档完成总结

**YOLO 模式已完成 HotGo 后端架构文档的创建**，包含以下完整章节：

✅ **引言** - 项目概述和文档目标  
✅ **高层架构** - 技术概要、架构图和设计模式  
✅ **技术栈** - 完整的技术选择和决策说明  
✅ **数据模型** - 实体关系图和核心数据结构  
✅ **组件架构** - 前后端组件设计和职责划分  
✅ **外部API** - 第三方服务集成说明  
✅ **核心工作流** - 关键业务流程的时序图  
✅ **REST API规范** - OpenAPI 3.0 格式的API文档  
✅ **数据库Schema** - 完整的MySQL表结构设计  
✅ **源代码树** - 项目文件组织和架构原则  
✅ **基础设施部署** - 容器化部署和CI/CD流程  
✅ **错误处理策略** - 统一的错误处理和日志规范  
✅ **编码标准** - Go和TypeScript的编码规范  
✅ **测试策略** - 单元测试、集成测试和E2E测试  
✅ **安全** - 全面的安全措施和最佳实践

该文档为 HotGo 项目提供了完整的后端架构指导，特别关注了双因素认证功能的实现细节，为开发团队和AI代理提供了准确的技术参考。

## Tech Stack

这是**决定性的**技术选择部分。基于对 HotGo 项目的分析，以下是当前使用的技术栈：

### Cloud Infrastructure

- **Provider:** 本地部署 / 私有云
- **Key Services:** MySQL 数据库、Redis 缓存、本地文件存储
- **Deployment Regions:** 根据部署需求配置

### Technology Stack Table

| 类别 | 技术 | 版本 | 用途 | 理由 |
|------|------|------|------|------|
| **语言** | Go | 1.24.4 | 后端开发语言 | 高性能、并发支持、强类型、丰富生态 |
| **后端框架** | GoFrame | v2.9.1 | 企业级Web框架 | 功能完整、文档齐全、社区活跃、适合企业开发 |
| **前端语言** | TypeScript | 5.3+ | 前端开发语言 | 强类型支持、更好的IDE体验、减少运行时错误 |
| **前端框架** | Vue | 3.4.38 | 前端UI框架 | 组合式API、性能优异、生态成熟 |
| **UI组件库** | Naive UI | 2.42.0 | Vue 3组件库 | 现代化设计、TypeScript支持、组件丰富 |
| **构建工具** | Vite | 最新 | 前端构建工具 | 快速热重载、ES模块支持、开发体验优秀 |
| **状态管理** | Pinia | 2.2.2 | Vue状态管理 | Vue 3官方推荐、TypeScript友好、简洁API |
| **HTTP客户端** | Axios | 1.7.7 | API请求库 | 功能完整、拦截器支持、广泛使用 |
| **数据库** | MySQL | 8.0+ | 关系型数据库 | 成熟稳定、性能优异、事务支持完整 |
| **缓存** | Redis | 6.0+ | 内存数据库 | 高性能缓存、会话存储、丰富数据结构 |
| **认证** | JWT | golang-jwt/jwt/v5 | 令牌认证 | 无状态、跨域支持、标准化 |
| **ORM** | GoFrame ORM | v2.9.1 | 数据库ORM | 框架内置、功能完整、代码生成支持 |
| **路由** | GoFrame Router | v2.9.1 | HTTP路由 | 中间件支持、RESTful设计、性能优异 |
| **日志** | GoFrame Log | v2.9.1 | 日志系统 | 结构化日志、多输出支持、性能优化 |
| **配置管理** | GoFrame Config | v2.9.1 | 配置管理 | 多格式支持、环境变量、热重载 |
| **包管理** | Go Modules | - | Go包管理 | 官方工具、版本管理、依赖解析 |
| **前端包管理** | pnpm | 最新 | 前端包管理 | 快速安装、磁盘空间优化、monorepo支持 |
| **代码检查** | ESLint | 最新 | 前端代码检查 | 代码质量保证、团队规范统一 |
| **代码格式化** | Prettier | 最新 | 代码格式化 | 统一代码风格、自动格式化 |
| **CSS预处理** | Sass/SCSS | 最新 | CSS预处理器 | 变量支持、嵌套语法、模块化 |
| **消息队列** | Apache RocketMQ | v2+ | 异步消息处理 | 高可靠性、分布式、事务消息支持 |
| **短信服务** | 阿里云SMS | 最新SDK | 短信发送 | 稳定可靠、覆盖面广、价格合理 |
| **支付集成** | 微信支付 | 最新SDK | 支付处理 | 用户基数大、集成简单、安全可靠 |
| **对象存储** | MinIO | v7+ | 文件存储 | S3兼容、私有部署、成本可控 |

### 关键技术决策说明

**后端技术选择**:
- **Go + GoFrame**: 选择Go语言是因为其出色的并发性能和简洁的语法，GoFrame提供了完整的企业级功能，减少了重复开发
- **MySQL**: 作为主数据库，提供ACID事务保证和成熟的生态系统
- **Redis**: 用于缓存和会话存储，显著提升系统性能

**前端技术选择**:
- **Vue 3 + TypeScript**: Vue 3的组合式API提供更好的逻辑复用，TypeScript增强代码可维护性
- **Naive UI**: 现代化的Vue 3组件库，与TypeScript深度集成
- **Vite**: 提供快速的开发体验和优化的生产构建

**架构模式选择**:
- **单体架构**: 适合当前团队规模和业务复杂度，简化部署和运维
- **前后端分离**: 支持独立开发和部署，提高开发效率
- **RESTful API**: 标准化的接口设计，便于集成和扩展

---

**需要用户确认**: 以上技术栈选择是否符合项目需求？是否需要调整或补充任何技术选择？

请选择以下选项：

1. 继续下一部分
2. 技术栈深度分析 - 详细评估每个技术选择的优缺点
3. 版本管理策略 - 讨论技术栈版本升级和兼容性
4. 性能基准测试 - 分析当前技术栈的性能特征
5. 安全性评估 - 评估技术栈的安全特性和风险
6. 替代方案分析 - 探讨其他可行的技术选择
7. 成本效益分析 - 评估技术栈的开发和运维成本
8. 团队技能匹配 - 分析团队技能与技术栈的匹配度
9. 生态系统评估 - 评估技术栈的社区支持和生态完整性

选择 1-9 或直接输入您的问题/反馈：