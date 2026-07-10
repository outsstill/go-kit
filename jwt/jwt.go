package jwt

import (
	"errors"
	"fmt"
	"github.com/outsstill/go-kit/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

type JWT_TYPE int

const (
	USER_TOKEN_TYPE JWT_TYPE = iota + 1
	ADMIN_TOKEN_TYPE
	MERCHANT_TOKEN_TYPE
	AGENT_TOKEN_TYPE
)

// JWT 定义一个jwt对象
type JWT struct {

	// 秘钥，用以加密 JWT，读取配置信息 app.key
	SignKey []byte

	// 刷新 Token 的最大过期时间
	MaxRefresh time.Duration

	cfg Config
}

// JWTCustomClaims 自定义载荷
type JWTCustomClaims struct {
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	ExpireAtTime int64    `json:"expire_time"`
	Type         JWT_TYPE `json:"type"`

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.RegisteredClaims
}

type Config struct {
	Name       string   `mapstructure:"name" json:"name"`
	Key        string   `mapstructure:"key" json:"key"`
	MaxRefresh int64    `mapstructure:"max_refresh" json:"max_refresh"`
	Timezone   string   `mapstructure:"timezone" json:"timezone"`
	Expires    int64    `mapstructure:"expires" json:"expires"`
	Type       JWT_TYPE `mapstructure:"type" json:"type"`
}

func NewJWT(cfg Config) (*JWT, error) {

	if cfg.Key == "" {
		return nil, errors.New("jwt key is empty")
	}

	if cfg.Expires <= 0 {
		cfg.Expires = 120
	}

	if cfg.MaxRefresh <= 0 {
		cfg.MaxRefresh = 10080
	}

	if cfg.Timezone == "" {
		cfg.Timezone = "UTC"
	}

	return &JWT{
		SignKey:    []byte(cfg.Key),
		MaxRefresh: time.Duration(cfg.MaxRefresh) * time.Minute,
		cfg:        cfg,
	}, nil
}

func (jwt *JWT) timenowInTimezone() time.Time {
	chinaTimezone, err := time.LoadLocation(jwt.cfg.Timezone)
	if err != nil {
		logger.Error(fmt.Sprintf("timenowInTimezone ERROR: %s", err.Error()))
		return time.Now().UTC()
	}
	return time.Now().In(chinaTimezone)
}

// IssueToken 生成  Token，在登录成功时调用
func (jwt *JWT) IssueToken(userID string, userName string, types ...JWT_TYPE) (string, error) {

	ty := jwt.cfg.Type

	if len(types) > 0 {
		ty = types[0]
	}

	// 1. 构造用户 claims 信息(负荷)
	expireAt := jwt.expireAtTime()
	now := jwt.timenowInTimezone()
	// 2. 根据 claims 生成token对象
	token, err := jwt.GenerateJwt(jwt.SignKey, jwtpkg.SigningMethodHS256, JWTCustomClaims{
		userID,
		userName,
		expireAt.Unix(),
		JWT_TYPE(ty),
		jwtpkg.RegisteredClaims{
			ExpiresAt: jwtpkg.NewNumericDate(expireAt), // 过期时间
			IssuedAt:  jwtpkg.NewNumericDate(now),      // 签发时间
			NotBefore: jwtpkg.NewNumericDate(now),      // 生效时间
			Issuer:    jwt.cfg.Name,                    // 签发者
		},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// expireAtTime 过期时间
func (jwt *JWT) expireAtTime() time.Time {
	timenow := jwt.timenowInTimezone()
	expire := time.Duration(jwt.cfg.Expires) * time.Minute
	return timenow.Add(expire)
}

func (jwt *JWT) GenerateJwt(key any, method jwtpkg.SigningMethod, claims jwtpkg.Claims) (string, error) {
	token := jwtpkg.NewWithClaims(method, claims)
	return token.SignedString(key)
}

// ParserToken 解析 Token，中间件中调用
func (jwt *JWT) ParserToken(tokenString string) (*JWTCustomClaims, error) {

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	if err != nil {
		return nil, err
	}

	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 3. 将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ParserToken 解析 Token，中间件中调用
func (jwt *JWT) ParserTokenGin(c *gin.Context) (*JWTCustomClaims, error) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	if err != nil {
		return nil, err
	}

	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 3. 将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// 刷新 Token
func (jwt *JWT) RefreshTokenGin(c *gin.Context) (string, error) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		return "", err
	}

	// 验证 Token 是否有效
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	// 获取 Claims
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// 检查 Token 是否即将过期（例如剩余时间小于 5 分钟）
	if time.Until(claims.ExpiresAt.Time) > 5*time.Minute {
		return "", fmt.Errorf("token is not expired yet")
	}

	// 生成新的 Token
	claims.ExpiresAt = jwtpkg.NewNumericDate(jwt.expireAtTime())
	claims.IssuedAt = jwtpkg.NewNumericDate(time.Now())
	claims.NotBefore = jwtpkg.NewNumericDate(time.Now())
	newTokenString, err := jwt.createToken(*claims)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %v", err)
	}

	return newTokenString, nil
}

// 刷新 Token
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		return "", err
	}

	// 验证 Token 是否有效
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	// 获取 Claims
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// 检查 Token 是否即将过期（例如剩余时间小于 5 分钟）
	if time.Until(claims.ExpiresAt.Time) > 5*time.Minute {
		return "", fmt.Errorf("token is not expired yet")
	}

	// 生成新的 Token
	claims.ExpiresAt = jwtpkg.NewNumericDate(jwt.expireAtTime())
	claims.IssuedAt = jwtpkg.NewNumericDate(time.Now())
	claims.NotBefore = jwtpkg.NewNumericDate(time.Now())
	newTokenString, err := jwt.createToken(*claims)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %v", err)
	}

	return newTokenString, nil
}

// parseTokenString 使用 jwtpkg.ParseWithClaims 解析 Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// getTokenFromHeader 使用 jwtpkg.ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	// 使用HS256算法进行token生成
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}
