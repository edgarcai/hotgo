# HotGo 棕地架构文档

## 引言

本文档记录了 HotGo 项目的**当前真实状态**，包括技术债务、实际实现模式和开发约束。本文档专为 AI 代理和开发者提供准确的系统理解，重点关注双因素认证功能的增强开发。

### 文档范围

重点关注与双因素认证功能增强相关的系统区域，包括认证系统、用户管理、前端组件和数据库设计。

### 变更日志

| 日期 | 版本 | 描述 | 作者 |
|------|------|------|------|
| 2025-01-27 | 1.0 | 初始棕地分析，聚焦2FA功能开发 | Architect |

## 快速参考 - 关键文件和入口点

### 理解系统的关键文件

- **后端主入口**: `server/main.go`
- **前端主入口**: `web/src/main.ts`
- **配置文件**: `server/hack/config.yaml`, `web/.env`
- **核心业务逻辑**: `server/internal/logic/`, `server/internal/service/`
- **API定义**: `server/api/`, `server/internal/controller/`
- **数据库模型**: `server/internal/model/entity/`, `server/internal/dao/`
- **前端组件**: `web/src/components/`, `web/src/views/`
- **认证相关**: `server/internal/library/token/`, `server/internal/logic/middleware/`

### 双因素认证功能影响区域

**当前状态**: 正在开发中，已有完整的PRD和前端规格文档

**需要修改的文件**:
- `server/internal/model/entity/` - 新增2FA数据模型
- `server/internal/logic/admin/site.go` - 登录流程集成
- `server/internal/library/token/token.go` - 认证流程增强
- `web/src/components/TwoFactor/` - 前端2FA组件（待创建）
- `web/src/views/auth/` - 登录页面修改
- `web/src/views/settings/` - 用户设置页面

## 高层架构

### 技术概要

HotGo 是一个现代化的全栈企业级应用框架，采用前后端分离架构。后端基于 Go 语言和 GoFrame 框架，前端使用 Vue 3 + TypeScript + Naive UI。系统设计注重模块化、可扩展性和开发效率。

### 实际技术栈（基于 package.json/go.mod）

| 类别 | 技术 | 版本 | 备注 |
|------|------|------|------|
| **后端运行时** | Go | 1.24.4 | 最新版本，性能优异 |
| **后端框架** | GoFrame | v2.9.1 | 企业级框架，功能完整 |
| **数据库** | MySQL | - | 主数据库，通过配置连接 |
| **缓存** | Redis | - | 会话存储和缓存 |
| **前端框架** | Vue | 3.4.38 | 组合式API |
| **UI库** | Naive UI | 2.42.0 | 现代化组件库 |
| **构建工具** | Vite | - | 快速构建和热重载 |
| **状态管理** | Pinia | 2.2.2 | Vue 3 推荐状态管理 |
| **HTTP客户端** | Axios | 1.7.7 | API请求处理 |
| **认证** | JWT | golang-jwt/jwt/v5 | 基于令牌的认证 |
| **包管理** | pnpm | - | 前端包管理器 |

### 仓库结构现状

- **类型**: 单仓库多模块（Monorepo）
- **包管理器**: Go Modules + pnpm
- **特殊说明**: 前后端代码分离在不同目录，但共享文档和配置

## 源码树和模块组织

### 项目结构（实际）

```text
hotgo/
├── server/                  # Go 后端服务
│   ├── main.go             # 应用入口点
│   ├── go.mod              # Go 模块依赖
│   ├── hack/config.yaml    # 开发配置（数据库连接等）
│   ├── internal/           # 内部业务逻辑
│   │   ├── cmd/            # 命令行工具
│   │   ├── controller/     # HTTP控制器
│   │   ├── logic/          # 业务逻辑层
│   │   ├── service/        # 服务接口定义
│   │   ├── model/          # 数据模型
│   │   ├── dao/            # 数据访问对象
│   │   ├── library/        # 公共库（认证、缓存等）
│   │   └── global/         # 全局配置和初始化
│   ├── api/                # API接口定义
│   ├── addons/             # 插件系统
│   └── utility/            # 工具函数
├── web/                    # Vue 前端应用
│   ├── src/
│   │   ├── main.ts         # 前端入口点
│   │   ├── App.vue         # 根组件
│   │   ├── components/     # 可复用组件
│   │   ├── views/          # 页面组件
│   │   ├── store/          # Pinia状态管理
│   │   ├── router/         # Vue Router配置
│   │   ├── api/            # API请求封装
│   │   └── utils/          # 工具函数
│   ├── package.json        # 前端依赖配置
│   ├── vite.config.ts      # Vite构建配置
│   └── tsconfig.json       # TypeScript配置
├── docs/                   # 项目文档
│   ├── prd.md             # 产品需求文档
│   ├── front-end-spec.md  # 前端规格文档
│   └── fullstack-architecture.md # 全栈架构文档
└── .bmad-core/            # BMAD开发工具配置
```

### 关键模块及其用途

**后端核心模块**:
- **认证系统**: `internal/library/token/` - JWT令牌管理，支持多端登录控制
- **用户管理**: `internal/logic/admin/` - 管理员用户相关业务逻辑
- **中间件**: `internal/logic/middleware/` - 认证、权限、日志等中间件
- **数据访问**: `internal/dao/` - 自动生成的数据访问层
- **业务服务**: `internal/service/` - 服务接口定义

**前端核心模块**:
- **认证组件**: `src/views/auth/` - 登录、注册等认证页面
- **用户设置**: `src/views/settings/` - 用户个人设置页面
- **API封装**: `src/api/` - 后端API请求封装
- **状态管理**: `src/store/` - 全局状态管理

## 数据模型和API

### 数据模型

**现有核心模型**:
- **用户模型**: `internal/model/entity/admin_member.go` - 管理员用户信息
- **登录日志**: `internal/model/entity/sys_login_log.go` - 登录记录
- **系统日志**: `internal/model/entity/sys_log.go` - 操作日志

**双因素认证相关模型（规划中）**:
- **2FA设置**: 用户双因素认证配置（待实现）
- **备用码**: 恢复码存储（待实现）
- **2FA日志**: 认证操作记录（待实现）

### API规范

**现有认证API**:
- `POST /api/admin/site/login` - 用户登录
- `POST /api/admin/site/logout` - 用户登出
- `GET /api/admin/site/profile` - 获取用户信息

**双因素认证API（规划中）**:
- `GET /api/admin/auth/2fa/status` - 获取2FA状态
- `POST /api/admin/auth/2fa/setup` - 初始化2FA设置
- `POST /api/admin/auth/2fa/verify` - 验证2FA代码
- `POST /api/admin/auth/2fa/disable` - 禁用2FA

## 技术债务和已知问题

### 关键技术债务

1. **认证系统**: 当前JWT实现较为基础，缺少刷新令牌机制的完整实现
2. **前端组件**: TwoFactor组件目录存在但为空，需要完整实现
3. **数据库迁移**: 缺少正式的数据库版本管理工具，依赖手动SQL脚本
4. **错误处理**: 前后端错误处理机制不够统一，需要标准化
5. **测试覆盖**: 缺少完整的单元测试和集成测试

### 开发约束和注意事项

- **数据库连接**: 开发环境硬编码密码 `123456789a`，生产环境需要环境变量
- **端口配置**: 后端默认3000端口，前端开发服务器需要避免冲突
- **代码生成**: 使用 `gf gen dao` 自动生成数据访问层，不要手动修改生成的文件
- **插件系统**: addons目录包含插件机制，新功能可以考虑插件化实现
- **多语言支持**: 系统支持中文，注释和错误信息需要中文化

## 集成点和外部依赖

### 外部服务

| 服务 | 用途 | 集成类型 | 关键文件 |
|------|------|----------|----------|
| MySQL | 主数据库 | 直连 | `hack/config.yaml` |
| Redis | 缓存和会话 | 直连 | `internal/library/cache/` |
| 阿里云SMS | 短信服务 | SDK | `internal/library/sms/` |
| 微信支付 | 支付集成 | API | `internal/library/payment/` |
| 对象存储 | 文件存储 | SDK | `internal/library/storager/` |

### 内部集成点

- **前后端通信**: RESTful API，JSON格式，统一响应结构
- **认证机制**: JWT Bearer Token，存储在localStorage
- **权限控制**: 基于角色的访问控制（RBAC）
- **日志系统**: 统一日志格式，支持链路追踪

## 开发和部署

### 本地开发设置

**后端启动**:
```bash
cd server
go mod tidy
gf run main.go
```

**前端启动**:
```bash
cd web
pnpm install
pnpm dev
```

**数据库设置**:
- 本地MySQL，用户名: `root`，密码: `123456789a`
- 数据库名: `hotgo`
- 执行初始化SQL脚本（如果有）

### 构建和部署流程

**后端构建**:
```bash
gf build
# 或使用Docker
gf docker
```

**前端构建**:
```bash
pnpm build
```

**部署说明**:
- 后端: 编译为单一可执行文件，包含静态资源
- 前端: 静态文件部署到CDN或Web服务器
- 配置: 生产环境使用环境变量覆盖默认配置

## 测试现状

### 当前测试覆盖

- **单元测试**: 基础覆盖，主要在核心业务逻辑
- **集成测试**: 有限，主要针对API端点
- **前端测试**: 基本没有，依赖手动测试
- **E2E测试**: 无

### 运行测试

```bash
# 后端测试
cd server
go test ./...

# 前端测试（如果有）
cd web
pnpm test
```

## 双因素认证功能开发状态

### 当前进展

**已完成**:
- ✅ 产品需求文档 (PRD)
- ✅ 前端UI/UX规格文档
- ✅ 全栈架构设计文档
- ✅ 数据库表结构设计

**进行中**:
- 🔄 后端API实现
- 🔄 前端组件开发
- 🔄 数据库迁移脚本

**待开始**:
- ⏳ 集成测试
- ⏳ 用户验收测试
- ⏳ 安全审计

### 需要修改的文件

**后端新增文件**:
- `internal/model/entity/user_two_factor.go` - 2FA数据模型
- `internal/dao/user_two_factor.go` - 2FA数据访问
- `internal/logic/admin/two_factor.go` - 2FA业务逻辑
- `internal/controller/admin/two_factor.go` - 2FA控制器
- `api/admin/two_factor.go` - 2FA API定义

**后端修改文件**:
- `internal/logic/admin/site.go` - 登录流程集成2FA验证
- `internal/library/token/token.go` - 认证流程增强
- `internal/logic/middleware/admin_auth.go` - 中间件支持2FA

**前端新增文件**:
- `src/components/TwoFactor/` - 2FA组件库
- `src/views/auth/TwoFactorAuth.vue` - 2FA验证页面
- `src/api/twoFactor.ts` - 2FA API封装
- `src/store/modules/twoFactor.ts` - 2FA状态管理

**前端修改文件**:
- `src/views/auth/Login.vue` - 登录页面集成2FA
- `src/views/settings/Security.vue` - 安全设置页面
- `src/store/modules/auth.ts` - 认证状态管理增强

### 集成考虑事项

- 必须与现有JWT认证中间件兼容
- 遵循现有API响应格式 `{code, message, data}`
- 前端组件需要与Naive UI设计系统保持一致
- 数据库变更需要支持平滑升级，不影响现有用户
- 安全日志需要集成到现有日志系统

## 附录 - 常用命令和脚本

### 常用开发命令

```bash
# 后端相关
gf run main.go          # 启动开发服务器
gf gen dao              # 生成数据访问层
gf build                # 构建生产版本
gf docker               # Docker构建

# 前端相关
pnpm dev                # 启动开发服务器
pnpm build              # 构建生产版本
pnpm lint               # 代码检查
pnpm type-check         # 类型检查

# 数据库相关
mysql -u root -p123456789a hotgo  # 连接数据库
```

### 调试和故障排除

- **日志位置**: 控制台输出，可配置文件输出
- **调试模式**: 设置环境变量 `GF_GCFG_FILE=hack/config.yaml`
- **常见问题**: 
  - 端口冲突: 检查3000端口占用
  - 数据库连接: 确认MySQL服务运行和密码正确
  - 前端代理: 检查vite.config.ts中的proxy配置

### 代码生成工具

```bash
# 生成数据访问层（基于数据库表）
gf gen dao

# 生成控制器（基于API定义）
gf gen ctrl

# 生成服务接口
gf gen service
```

---

**文档维护说明**: 本文档应随着代码变更及时更新，特别是在双因素认证功能开发完成后，需要更新实现状态和实际文件路径。