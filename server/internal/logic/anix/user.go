// Package anix AniX用户业务逻辑
package anix

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"

	v1 "hotgo/api/anix/user/v1"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/anixin"
	"hotgo/internal/service"
)

type sAnixUser struct{}

// Model 用户ORM模型
func (s *sAnixUser) Model(ctx context.Context) *gdb.Model {
	return dao.AnixUser.Ctx(ctx)
}

// Register 用户注册
func (s *sAnixUser) Register(ctx context.Context, in *anixin.RegisterInp) (err error) {
	// 检查用户名是否已存在
	exists, err := s.IsUsernameExists(ctx, in.Username)
	if err != nil {
		return err
	}
	if exists {
		return gerror.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.IsEmailExists(ctx, in.Email)
	if err != nil {
		return err
	}
	if exists {
		return gerror.New("邮箱已被注册")
	}

	// 验证用户协议
	if !in.AgreeTerms {
		return gerror.New("请同意用户协议")
	}

	// 密码加密
	hashedPassword, err := s.HashPassword(in.Password)
	if err != nil {
		return gerror.Wrap(err, "密码加密失败")
	}

	// 创建用户
	user := &entity.AnixUser{
		Username:  in.Username,
		Email:     in.Email,
		Nickname:  in.Nickname,
		Password:  hashedPassword,
		Status:    1, // 正常状态
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	_, err = s.Model(ctx).Data(user).Insert()
	if err != nil {
		return gerror.Wrap(err, "用户注册失败")
	}

	return nil
}

// Login 用户登录
func (s *sAnixUser) Login(ctx context.Context, in *anixin.LoginInp) (res *anixin.LoginModel, err error) {
	var user *entity.AnixUser

	// 判断是邮箱还是用户名登录
	if gstr.Contains(in.Account, "@") {
		user, err = s.GetUserByEmail(ctx, in.Account)
	} else {
		user, err = s.GetUserByUsername(ctx, in.Account)
	}

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, gerror.New("用户已被禁用")
	}

	// 验证密码
	if !s.VerifyPassword(in.Password, user.Password) {
		return nil, gerror.New("密码错误")
	}

	// 生成Token
	token, err := s.GenerateToken(ctx, user)
	if err != nil {
		return nil, gerror.Wrap(err, "生成Token失败")
	}

	// 更新最后登录信息
	err = s.UpdateLastLogin(ctx, int64(user.Id), in.Ip)
	if err != nil {
		g.Log().Warning(ctx, "更新最后登录信息失败:", err)
	}

	res = &anixin.LoginModel{
		User:  user,
		Token: token,
	}

	return res, nil
}

// GetProfile 获取用户信息
func (s *sAnixUser) GetProfile(ctx context.Context, userId int64) (res *v1.UserInfo, err error) {
	var user *entity.AnixUser
	err = s.Model(ctx).Where("id", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Wrap(err, "获取用户信息失败")
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	res = &v1.UserInfo{
		Id:          int64(user.Id),
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		Birthday:    gconv.String(user.Birthday),
		Mobile:      user.Phone, // 实体模型中使用phone字段
		Status:      user.Status,
		LastLoginAt: gconv.String(user.LastLoginAt),
		LastLoginIp: user.LastLoginIp,
		CreatedAt:   gconv.String(user.CreatedAt),
		UpdatedAt:   gconv.String(user.UpdatedAt),
	}

	return res, nil
}

// UpdateProfile 更新用户信息
func (s *sAnixUser) UpdateProfile(ctx context.Context, in *anixin.UpdateProfileInp) (err error) {
	updateData := g.Map{
		"updated_at": gtime.Now(),
	}

	if in.Nickname != "" {
		updateData["nickname"] = in.Nickname
	}
	if in.Gender >= 0 {
		updateData["gender"] = in.Gender
	}
	if in.Birthday != nil {
		updateData["birthday"] = in.Birthday
	}
	if in.Mobile != "" {
		updateData["phone"] = in.Mobile // 实体模型中使用phone字段
	}
	if in.Avatar != "" {
		updateData["avatar"] = in.Avatar
	}

	_, err = s.Model(ctx).Where("id", in.UserId).Data(updateData).Update()
	if err != nil {
		return gerror.Wrap(err, "更新用户信息失败")
	}

	return nil
}

// ChangePassword 修改密码
func (s *sAnixUser) ChangePassword(ctx context.Context, in *anixin.ChangePasswordInp) (err error) {
	// 获取用户信息
	var user *entity.AnixUser
	err = s.Model(ctx).Where("id", in.UserId).Scan(&user)
	if err != nil {
		return gerror.Wrap(err, "获取用户信息失败")
	}

	if user == nil {
		return gerror.New("用户不存在")
	}

	// 验证原密码
	if !s.VerifyPassword(in.OldPassword, user.Password) {
		return gerror.New("原密码错误")
	}

	// 加密新密码
	newHashedPassword, err := s.HashPassword(in.NewPassword)
	if err != nil {
		return gerror.Wrap(err, "密码加密失败")
	}

	// 更新密码
	_, err = s.Model(ctx).Where("id", in.UserId).Data(g.Map{
		"password":   newHashedPassword,
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrap(err, "修改密码失败")
	}

	return nil
}

// Logout 用户登出
func (s *sAnixUser) Logout(ctx context.Context, token string) (err error) {
	// TODO: 实现Token黑名单机制
	// 这里可以将Token加入黑名单，或者使用Redis存储已登出的Token
	g.Log().Info(ctx, "用户登出, Token:", token)
	return nil
}

// VerifyToken 验证Token
func (s *sAnixUser) VerifyToken(ctx context.Context, token string) (userId int64, err error) {
	// 解析JWT Token
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.Cfg().MustGet(ctx, "jwt.secret").String()), nil
	})

	if err != nil {
		return 0, gerror.Wrap(err, "Token解析失败")
	}

	if !jwtToken.Valid {
		return 0, gerror.New("Token无效")
	}

	// 获取用户ID
	userIdFloat, ok := claims["user_id"]
	if !ok {
		return 0, gerror.New("Token中缺少用户ID")
	}

	userId = gconv.Int64(userIdFloat)
	if userId <= 0 {
		return 0, gerror.New("用户ID无效")
	}

	return userId, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *sAnixUser) GetUserByUsername(ctx context.Context, username string) (user *entity.AnixUser, err error) {
	err = s.Model(ctx).Where("username", username).Scan(&user)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *sAnixUser) GetUserByEmail(ctx context.Context, email string) (user *entity.AnixUser, err error) {
	err = s.Model(ctx).Where("email", email).Scan(&user)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	return user, nil
}

// IsEmailExists 检查邮箱是否存在
func (s *sAnixUser) IsEmailExists(ctx context.Context, email string) (exists bool, err error) {
	count, err := s.Model(ctx).Where("email", email).Count()
	if err != nil {
		return false, gerror.Wrap(err, "检查邮箱失败")
	}
	return count > 0, nil
}

// IsUsernameExists 检查用户名是否存在
func (s *sAnixUser) IsUsernameExists(ctx context.Context, username string) (exists bool, err error) {
	count, err := s.Model(ctx).Where("username", username).Count()
	if err != nil {
		return false, gerror.Wrap(err, "检查用户名失败")
	}
	return count > 0, nil
}

// GenerateToken 生成JWT Token
func (s *sAnixUser) GenerateToken(ctx context.Context, user *entity.AnixUser) (token string, err error) {
	// 获取JWT配置
	secret := g.Cfg().MustGet(ctx, "jwt.secret", "hotgo-anix-secret").String()
	expireTime := g.Cfg().MustGet(ctx, "jwt.expire", 7*24*3600).Int() // 默认7天

	// 创建Claims
	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}

	// 生成Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", gerror.Wrap(err, "生成Token失败")
	}

	return token, nil
}

// HashPassword 密码加密
func (s *sAnixUser) HashPassword(password string) (hashedPassword string, err error) {
	// 使用MD5加盐加密
	salt := "anix_salt_2024"
	hashedPassword = gmd5.MustEncrypt(password + salt)
	return hashedPassword, nil
}

// VerifyPassword 验证密码
func (s *sAnixUser) VerifyPassword(password, hashedPassword string) bool {
	salt := "anix_salt_2024"
	return gmd5.MustEncrypt(password+salt) == hashedPassword
}

// UpdateLastLogin 更新最后登录信息
func (s *sAnixUser) UpdateLastLogin(ctx context.Context, userId int64, ip string) (err error) {
	_, err = s.Model(ctx).Where("id", userId).Data(g.Map{
		"last_login_at": gtime.Now(),
		"last_login_ip": ip,
		"updated_at":    gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrap(err, "更新最后登录信息失败")
	}
	return nil
}

// GetUserFromContext 从上下文获取用户信息
func (s *sAnixUser) GetUserFromContext(ctx context.Context) (user *entity.AnixUser, err error) {
	userIdValue := ctx.Value("user_id")
	if userIdValue == nil {
		return nil, gerror.New("未找到用户信息")
	}

	userId := gconv.Int64(userIdValue)
	if userId <= 0 {
		return nil, gerror.New("用户ID无效")
	}

	err = s.Model(ctx).Where("id", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Wrap(err, "获取用户信息失败")
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	return user, nil
}

// GetUserFromRequest 从请求获取用户信息
func (s *sAnixUser) GetUserFromRequest(r *ghttp.Request) (user *entity.AnixUser, err error) {
	// 从Header获取Token
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, gerror.New("未提供访问令牌")
	}

	// 移除Bearer前缀
	if gstr.HasPrefix(token, "Bearer ") {
		token = gstr.SubStr(token, 7)
	}

	// 验证Token获取用户ID
	userId, err := s.VerifyToken(r.Context(), token)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	err = s.Model(r.Context()).Where("id", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Wrap(err, "获取用户信息失败")
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	return user, nil
}

func init() {
	service.RegisterAnixUser(&sAnixUser{})
}