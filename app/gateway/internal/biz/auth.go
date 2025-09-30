package biz

import (
	"context"
	"time"

	commmonv1 "rag/api/common/v1"
	v1 "rag/api/gateway/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.GatewayServiceErrorReason_UNKNOWN_ERROR.String(), "user not found")
)

// User represents user information
type User struct {
	UserID      string
	Username    string
	Email       string
	Roles       []string
	Attributes  map[string]string
	CreatedAt   time.Time
	LastLoginAt time.Time
}

// AuthToken represents authentication token information
type AuthToken struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	Scope        string
}

// AuthRepo defines the data access interface for authentication
type AuthRepo interface {
	// 验证用户凭据
	ValidateCredentials(ctx context.Context, username, password string) (*User, error)
	// 生成访问令牌
	GenerateAccessToken(ctx context.Context, user *User) (*AuthToken, error)
	// 验证访问令牌
	ValidateAccessToken(ctx context.Context, token string) (*User, error)
	// 刷新令牌
	RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthToken, error)
	// 撤销令牌
	RevokeToken(ctx context.Context, token string) error
	// 更新用户最后登录时间
	UpdateLastLogin(ctx context.Context, userID string) error
}

// AuthUsecase handles authentication business logic
type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

// NewAuthUsecase creates a new authentication usecase
func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Login handles user login
func (uc *AuthUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	uc.log.WithContext(ctx).Infof("User login attempt: %s", req.Username)
	// return nil, errors.New(400, "Bad Request", "Invalid username or password").WithMetadata(map[string]string{"username": "shenmingjie"})
	return nil, errors.New(400, commmonv1.ErrorCode_ERROR_CODE_BAD_GATEWAY.String(), "")
	// 验证用户凭据
	user, err := uc.repo.ValidateCredentials(ctx, req.Username, req.Password)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to validate credentials: %v", err)
		return nil, err
	}

	// 生成访问令牌
	token, err := uc.repo.GenerateAccessToken(ctx, user)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to generate token: %v", err)
		return nil, err
	}

	// 更新最后登录时间
	if err := uc.repo.UpdateLastLogin(ctx, user.UserID); err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to update last login: %v", err)
	}

	// 构建响应
	response := &v1.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		Scope:        token.Scope,
		UserInfo: &v1.UserInfo{
			UserId:     user.UserID,
			Username:   user.Username,
			Email:      user.Email,
			Roles:      user.Roles,
			Attributes: user.Attributes,
		},
	}

	uc.log.WithContext(ctx).Infof("User login successful: %s", req.Username)
	return response, nil
}

// RefreshToken handles token refresh
func (uc *AuthUsecase) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	uc.log.WithContext(ctx).Info("Token refresh attempt")

	// 刷新访问令牌
	token, err := uc.repo.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to refresh token: %v", err)
		return nil, err
	}

	response := &v1.RefreshTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
	}

	uc.log.WithContext(ctx).Info("Token refresh successful")
	return response, nil
}

// ValidateToken validates access token and returns user info
func (uc *AuthUsecase) ValidateToken(ctx context.Context, token string) (*User, error) {
	return uc.repo.ValidateAccessToken(ctx, token)
}
