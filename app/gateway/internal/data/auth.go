package data

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"rag/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// authRepo implements biz.AuthRepo interface
type authRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuthRepo creates a new authentication repository
func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ValidateCredentials validates user credentials
func (r *authRepo) ValidateCredentials(ctx context.Context, username, password string) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("Validating credentials for user: %s", username)

	// TODO: 实际实现应该查询用户数据库
	// 这里是示例实现，实际应该从数据库获取用户信息并验证密码

	// 模拟用户数据
	mockUsers := map[string]struct {
		userID     string
		hashedPwd  string
		email      string
		roles      []string
		attributes map[string]string
	}{
		"admin": {
			userID:    "user-1",
			hashedPwd: "$2a$10$N9qo8uLOickgx2ZMRZoMye8QjqBF.JMhtMEHXOJZmDzjyqTb6bDhO", // "password"\n			email:     "admin@example.com",
			roles:     []string{"admin", "user"},
			attributes: map[string]string{
				"department": "IT",
				"level":      "senior",
			},
		},
		"user": {
			userID:    "user-2",
			hashedPwd: "$2a$10$N9qo8uLOickgx2ZMRZoMye8QjqBF.JMhtMEHXOJZmDzjyqTb6bDhO", // "password"
			email:     "user@example.com",
			roles:     []string{"user"},
			attributes: map[string]string{
				"department": "Business",
				"level":      "junior",
			},
		},
	}

	userData, exists := mockUsers[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(userData.hashedPwd), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	user := &biz.User{
		UserID:      userData.userID,
		Username:    username,
		Email:       userData.email,
		Roles:       userData.roles,
		Attributes:  userData.attributes,
		CreatedAt:   time.Now().AddDate(-1, 0, 0), // 假设一年前创建
		LastLoginAt: time.Now(),
	}

	r.log.WithContext(ctx).Infof("User validated successfully: %s", username)
	return user, nil
}

// GenerateAccessToken generates access and refresh tokens
func (r *authRepo) GenerateAccessToken(ctx context.Context, user *biz.User) (*biz.AuthToken, error) {
	r.log.WithContext(ctx).Infof("Generating token for user: %s", user.UserID)

	// JWT 密钥（实际应该从配置文件读取）
	jwtSecret := []byte("your-secret-key")

	// 生成访问令牌
	accessTokenClaims := jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"roles":    user.Roles,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // 1小时过期
		"iat":      time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
		"iat":     time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	authToken := &biz.AuthToken{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小时
		Scope:        "read write",
	}

	r.log.WithContext(ctx).Infof("Token generated successfully for user: %s", user.UserID)
	return authToken, nil
}

// ValidateAccessToken validates access token and returns user info
func (r *authRepo) ValidateAccessToken(ctx context.Context, tokenString string) (*biz.User, error) {
	r.log.WithContext(ctx).Info("Validating access token")

	jwtSecret := []byte("your-secret-key")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	rolesInterface, ok := claims["roles"]
	if !ok {
		return nil, fmt.Errorf("invalid roles in token")
	}

	var roles []string
	if roleSlice, ok := rolesInterface.([]interface{}); ok {
		for _, role := range roleSlice {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	user := &biz.User{
		UserID:   userID,
		Username: username,
		Roles:    roles,
	}

	r.log.WithContext(ctx).Infof("Token validated successfully for user: %s", userID)
	return user, nil
}

// RefreshAccessToken refreshes access token using refresh token
func (r *authRepo) RefreshAccessToken(ctx context.Context, refreshToken string) (*biz.AuthToken, error) {
	r.log.WithContext(ctx).Info("Refreshing access token")

	jwtSecret := []byte("your-secret-key")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	// 生成新的访问令牌
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// 生成新的刷新令牌
	newRefreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims)
	newRefreshTokenString, err := newRefreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	authToken := &biz.AuthToken{
		AccessToken:  accessTokenString,
		RefreshToken: newRefreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read write",
	}

	r.log.WithContext(ctx).Info("Access token refreshed successfully")
	return authToken, nil
}

// RevokeToken revokes a token
func (r *authRepo) RevokeToken(ctx context.Context, token string) error {
	r.log.WithContext(ctx).Info("Revoking token")
	// TODO: 实现令牌撤销逻辑，通常是将令牌加入黑名单
	// 这里是示例实现
	return nil
}

// UpdateLastLogin updates user's last login time
func (r *authRepo) UpdateLastLogin(ctx context.Context, userID string) error {
	r.log.WithContext(ctx).Infof("Updating last login for user: %s", userID)
	// TODO: 实现更新用户最后登录时间的逻辑
	// 这里是示例实现
	return nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
