## **1. 产品概述**

- **产品名称**：AniX
    
- **产品类型**：面向二次元群体的 **免费在线漫画阅读网站**
    
- **核心价值**：
    
    - 提供高质量、免费、高清的漫画资源。
        
    - 打造二次元爱好者的交流与互动社区。
        
    - 支持多平台（Web/H5/App）无障碍阅读。
        
    

---

## **2. 产品定位与目标用户**


### **2.1 产品定位**

- 面向 **二次元用户** 的 **漫画阅读 + 社交互动** 平台。
    
- 以 **免费 + 高速阅读体验** 为核心竞争力。
    
- 兼顾 **UGC（用户生成内容）**，逐步建立社区生态。
    

### **2.2 目标用户**

1. **核心用户群体**
    
    - 年龄：12-35 岁
        
    - 特征：二次元爱好者，习惯在线阅读，愿意分享与讨论
        
    
2. **次要用户群体**
    - 动漫新人，想要快速获取热门作品
    - 海外华人及日漫爱好者

### **2.3 用户需求**

- 免费、快速获取漫画资源
    
- 支持手机端、PC 端无障碍阅读
    
- 漫画分类清晰、更新及时
    
- 可以和其他漫迷互动、讨论
    
- 支持收藏、书签、历史记录

---

## **3. 核心功能需求**
### **3.1 用户端功能**

1. **注册与登录**
    - 支持邮箱、手机号注册
    - 游客模式可直接阅读内容
    - 收藏/书签需登录
    - 免广告需要会员
    
2. **漫画阅读**
    
    - 漫画列表：分类浏览（热血/恋爱/冒险/搞笑/同人…）
        
    - 漫画详情页：作品介绍、作者、更新状态、评分
        
    - 在线阅读器：支持竖读/横读模式、自动翻页、夜间模式
        
    - 收藏/书签：用户可标记阅读进度
        
    - 阅读历史：自动保存上次阅读进度
        
3. **搜索与推荐**
    
    - 关键词搜索：按标题/作者/标签
        
    - 智能推荐：根据历史阅读偏好推荐作品
        
    - 排行榜：热门、最新、用户评分
        
4. **社区互动**
    
    - 章节评论
        
    - 点赞/回复功能
        
    - 用户个人主页（头像、昵称、关注、粉丝）
        
    
5. **个人中心**
    
    - 我的收藏
        
    - 我的书签/历史
        
    - 设置（夜间模式、字体、语言）
        
    

---

### **3.2 管理后台功能**

1. **漫画管理**
    
    - 漫画上架/下架
        
    - 分类、标签管理
        
    - 章节上传、顺序管理
        
    - 漫画版权信息标注
        
    
2. **用户管理**
    
    - 用户账号管理（封禁/解封）
        
    - 评论审核
        
    - 举报处理
        
    
3. **数据统计**
    
    - PV/UV、用户留存率
        
    - 漫画热度排行
        
    - 用户行为数据（收藏、评论活跃度）
        

---

## **4. 业务流程**
### **4.1 阅读流程**

用户进入首页 → 选择分类/搜索漫画 → 打开详情页 → 开始在线阅读 → 可收藏/评论

### **4.2 用户互动流程**
用户阅读章节 → 评论区留言 → 其他用户点赞/回复 → 评论审核（后台可控）

### **4.3 内容运营流程**
编辑/爬虫收集漫画 → 后台上传漫画章节 → 设置分类与标签 → 上架到前台 → 用户可搜索/阅读

---

## **5. 非功能需求**

1. **性能**
    - 页面加载 < 2 秒
    - 支持高并发：单机支持 1 万 QPS，可通过负载均衡扩展
2. **兼容性**
    - 支持 Web、H5、小程序，后期支持移动端 APP
3. **安全性**
    - 用户数据加密存储（密码 hash）
    - 防止爬虫盗链（CDN 鉴权）
    - 评论敏感词过滤
    
4. **扩展性**

    - 模块化设计，支持后续接入 **付费阅读 / 广告系统 / 社区功能**
        
---
# **📌 数据库表设计（MySQL 8）**
## **1. 用户系统**

```
-- 用户表
CREATE TABLE users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    avatar_url VARCHAR(255) DEFAULT NULL COMMENT '头像',
    gender ENUM('male','female','secret') DEFAULT 'secret' COMMENT '性别',
    birthday DATE DEFAULT NULL COMMENT '生日',
    bio VARCHAR(255) DEFAULT NULL COMMENT '个人简介',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
);

-- 用户行为表（收藏、点赞、关注漫画）
CREATE TABLE user_actions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    action_type ENUM('favorite','like','follow') NOT NULL COMMENT '收藏/点赞/关注',
    target_type ENUM('comic','chapter','author') NOT NULL COMMENT '目标对象类型',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '目标对象ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_action(user_id, action_type, target_type, target_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## **2. 漫画 & 作者**

```
-- 漫画表
CREATE TABLE comics (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '漫画标题',
    description TEXT COMMENT '简介',
    cover_url VARCHAR(255) COMMENT '封面图',
    author_id BIGINT UNSIGNED COMMENT '作者ID',
    status ENUM('ongoing','completed','hiatus') DEFAULT 'ongoing' COMMENT '连载状态',
    tags JSON COMMENT '标签（如 热血, 校园, 魔幻）',
    views BIGINT UNSIGNED DEFAULT 0 COMMENT '总浏览量',
    likes BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    favorites BIGINT UNSIGNED DEFAULT 0 COMMENT '收藏数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 作者表
CREATE TABLE authors (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '作者名',
    bio TEXT COMMENT '作者简介',
    avatar_url VARCHAR(255) COMMENT '头像',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## **3. 章节 & 漫画内容**

```
-- 漫画章节表
CREATE TABLE chapters (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    comic_id BIGINT UNSIGNED NOT NULL COMMENT '所属漫画',
    title VARCHAR(200) NOT NULL COMMENT '章节标题',
    chapter_number INT UNSIGNED NOT NULL COMMENT '章节序号',
    views BIGINT UNSIGNED DEFAULT 0 COMMENT '浏览量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    UNIQUE KEY uq_comic_chapter(comic_id, chapter_number)
);

-- 漫画图片页
CREATE TABLE chapter_pages (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    chapter_id BIGINT UNSIGNED NOT NULL,
    page_number INT UNSIGNED NOT NULL COMMENT '页码',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    UNIQUE KEY uq_chapter_page(chapter_id, page_number)
);
```

---

## **4. 评论 & 弹幕**

```
-- 评论表
CREATE TABLE comments (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    target_type ENUM('comic','chapter') NOT NULL COMMENT '评论对象',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '对象ID',
    content TEXT NOT NULL COMMENT '评论内容',
    likes BIGINT UNSIGNED DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 弹幕表
CREATE TABLE danmu (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    chapter_id BIGINT UNSIGNED NOT NULL COMMENT '所属章节',
    user_id BIGINT UNSIGNED NOT NULL,
    content VARCHAR(255) NOT NULL COMMENT '弹幕内容',
    position FLOAT NOT NULL COMMENT '出现时间/位置（秒或百分比）',
    color VARCHAR(20) DEFAULT '#FFFFFF' COMMENT '弹幕颜色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## **5. 阅读记录**

```
CREATE TABLE reading_history (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    comic_id BIGINT UNSIGNED NOT NULL,
    chapter_id BIGINT UNSIGNED NOT NULL,
    last_page INT UNSIGNED DEFAULT 1 COMMENT '上次读到第几页',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_comic(user_id, comic_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
);
```

---

## **6. 广告 / 变现**

```
-- 广告表
CREATE TABLE ads (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '广告标题',
    image_url VARCHAR(255) NOT NULL COMMENT '广告图',
    link_url VARCHAR(255) NOT NULL COMMENT '跳转链接',
    position ENUM('banner','reader_top','reader_bottom','popup') NOT NULL COMMENT '广告位置',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 广告点击日志
CREATE TABLE ad_click_logs (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    ad_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED DEFAULT NULL,
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ad_id) REFERENCES ads(id) ON DELETE CASCADE
);
```

---

## **7. 数据分析 & 推荐**

- comic_ratings 用户评分表
    
- search_logs 搜索日志
    
- recommendations 个性化推荐