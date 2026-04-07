package user

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/auth"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	session "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user_session"
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	auditdomain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenBlacklist interface {
	Blacklist(token string, duration time.Duration) error
	IsBlacklisted(token string) bool
}

type AuditLogger interface {
	LogEvent(ctx context.Context, event *auditdomain.AuditLog) error
}

type Service struct {
	repo         domain.Repository
	rateLimiter  domain.RedisRateLimiter
	sessionRepo  session.Repository
	audit        AuditLogger
	blacklist    TokenBlacklist
}

func NewService(
	repo domain.Repository,
	rl domain.RedisRateLimiter,
	sessionRepo session.Repository,
	audit AuditLogger,
	blacklist TokenBlacklist,
) *Service {
	return &Service{
		repo:        repo,
		rateLimiter: rl,
		sessionRepo: sessionRepo,
		audit:       audit,
		blacklist:   blacklist,
	}
}

//////////////////// REGISTER ////////////////////

func (s *Service) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:            name,
		Email:           email,
		Password:        string(hashedPassword),
		Role:            common.RoleUser,
		IsActive:        true,
		IsEmailVerified: false,
	}

	err = s.repo.WithTx(func(txRepo domain.Repository) error {
		if err := txRepo.Create(user); err != nil {
			return err
		}

		token := uuid.New().String()
		hashedToken := auth.HashToken(token)
		expiresAt := time.Now().Add(10 * time.Minute)

		return txRepo.SaveVerificationToken(user.ID, hashedToken, expiresAt)
	})

	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

//////////////////// LOGIN ////////////////////

func (s *Service) Login(
	ctx context.Context,
	email, deviceID, password, ipAddress, userAgent string,
) (*domain.AuthPayload, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	allowed, err := s.rateLimiter.Allow("login:"+email+":"+ipAddress, 10, time.Minute)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, domain.ErrRateLimited
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsEmailVerified {
		return nil, domain.ErrEmailNotVerified
	}
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}
	if user.AccountLockedUntil != nil && user.AccountLockedUntil.After(time.Now()) {
		return nil, domain.ErrAccountLocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		user.FailedLoginAttempts++

		if user.FailedLoginAttempts >= 5 {
			lock := time.Now().Add(15 * time.Minute)
			user.AccountLockedUntil = &lock
			user.FailedLoginAttempts = 0
		}

		_ = s.repo.Update(user)
		return nil, domain.ErrInvalidCredentials
	}

	// suspicious login detection
	if DetectSuspiciousLogin(user.LastLoginIP, ipAddress, user.LastUserAgent, userAgent) {
		_ = s.audit.LogEvent(ctx, &auditdomain.AuditLog{
			UserID:    user.ID,
			UserEmail: user.Email,
			UserRole:  string(user.Role),
			Action:    "SUSPICIOUS_LOGIN",
			Entity:    "USER",
			EntityID:  strconv.Itoa(int(user.ID)),
			IPAddress: ipAddress,
			UserAgent: userAgent,
		})
	}

	user.FailedLoginAttempts = 0
	user.AccountLockedUntil = nil
	user.LastLoginIP = ipAddress
	user.LastUserAgent = userAgent

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashed := auth.HashToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	err = s.sessionRepo.CreateUserSession(
		user.ID,
		deviceID,
		hashed,
		ipAddress,
		userAgent,
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	// audit success login
	_ = s.audit.LogEvent(ctx, &auditdomain.AuditLog{
		UserID:    user.ID,
		Action:    "LOGIN_SUCCESS",
		Entity:    "AUTH",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	user.Password = ""

	return &domain.AuthPayload{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceID:     deviceID,
	}, nil
}

//////////////////// REFRESH TOKEN ////////////////////

func (s *Service) RefreshToken(ctx context.Context, oldToken string) (*domain.AuthPayload, error) {

	if s.blacklist.IsBlacklisted(oldToken) {
		return nil, domain.ErrInvalidOrExpiredToken
	}

	hashed := auth.HashToken(oldToken)

	sessionData, err := s.sessionRepo.FindByToken(hashed)
	if err != nil {
		return nil, domain.ErrInvalidOrExpiredToken
	}

	if sessionData.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrInvalidOrExpiredToken
	}

	user, err := s.repo.FindByID(ctx, sessionData.UserID)
	if err != nil {
		return nil, err
	}

	newAccess, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	newRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedNew := auth.HashToken(newRefresh)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// rotate refresh token
	if err := s.sessionRepo.UpdateRefreshToken(sessionData.ID, hashedNew, expiresAt); err != nil {
		return nil, err
	}

	return &domain.AuthPayload{
		User:         user,
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

//////////////////// LOGOUT ////////////////////

func (s *Service) LogoutAllDevices(ctx context.Context, userID uint, accessToken string) error {

	if err := s.repo.DeleteAllRefreshTokens(userID); err != nil {
		return err
	}

	if err := s.sessionRepo.DeleteAllUserSessions(userID); err != nil {
		return err
	}

	if accessToken != "" {
		_ = s.blacklist.Blacklist(accessToken, 15*time.Minute)
	}

	return s.audit.LogEvent(ctx, &auditdomain.AuditLog{
		UserID: userID,
		Action: "LOGOUT_ALL_DEVICES",
		Entity: "AUTH",
	})
}

//////////////////// PASSWORD RESET ////////////////////

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {

	hashedToken := auth.HashToken(token)

	record, err := s.repo.FindPasswordResetByToken(hashedToken)
	if err != nil || record.ExpiresAt.Before(time.Now()) {
		return domain.ErrInvalidOrExpiredToken
	}

	user, err := s.repo.FindByID(ctx, record.UserID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.FailedLoginAttempts = 0
	user.AccountLockedUntil = nil

	if err := s.repo.Update(user); err != nil {
		return err
	}

	_ = s.repo.DeletePasswordResetToken(record.ID)
	_ = s.repo.DeleteAllRefreshTokens(user.ID)
	_ = s.sessionRepo.DeleteAllUserSessions(user.ID)

	_ = s.audit.LogEvent(ctx, &auditdomain.AuditLog{
		UserID: user.ID,
		Action: "PASSWORD_RESET",
		Entity: "AUTH",
	})

	return nil
}

//////////////////// HELPERS ////////////////////

func DetectSuspiciousLogin(lastIP, newIP, lastUA, newUA string) bool {
	return lastIP != "" && (lastIP != newIP || lastUA != newUA)
}

func (s *Service) PromoteToAdmin(ctx context.Context, actorID, targetUserID uint) (*domain.User, error) {

	actor, err := s.repo.FindByID(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if actor.Role != common.RoleSuperAdmin {
		return nil, errors.New("only super admins can promote users")
	}

	if actorID == targetUserID {
		return nil, errors.New("cannot promote yourself")
	}

	target, err := s.repo.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	if target.Role == common.RoleSuperAdmin {
		return nil, errors.New("cannot modify super admin")
	}

	if target.Role == common.RoleAdmin {
		return nil, errors.New("user is already an admin")
	}

	target.Role = common.RoleAdmin

	if err := s.repo.Update(target); err != nil {
		return nil, err
	}

	target.Password = ""
	return target, nil
}


func (s *Service) DeactivateUser(ctx context.Context, actorID uint, targetUserID uint) (*domain.User, error) {

	actor, err := s.repo.FindByID(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if actor.Role != common.RoleSuperAdmin {
		return nil, errors.New("only super admins can deactivate users")
	}

	if actorID == targetUserID {
		return nil, errors.New("cannot deactivate yourself")
	}

	target, err := s.repo.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	if target.Role == common.RoleSuperAdmin {
		return nil, errors.New("cannot deactivate another super admin")
	}

	target.IsActive = false

	if err := s.repo.Update(target); err != nil {
		return nil, err
	}

	target.Password = ""
	return target, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetAllUsers(page, limit int) ([]*domain.User, int, error) {
	return s.repo.FindAll(page, limit)
} 

func (s *Service) OAuthLogin(
	ctx context.Context,
	name string,
	email string,
	ip string,
	userAgent string,
	deviceID string,
) (*domain.AuthPayload, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repo.FindByEmail(ctx, email)

	if err != nil {

		if errors.Is(err, domain.ErrUserNotFound) {

			user = &domain.User{
				Name:            name,
				Email:           email,
				Role:            common.RoleUser,
				IsActive:        true,
				IsEmailVerified: true,
			}

			if err := s.repo.Create(user); err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashed := auth.HashToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Create session
	err = s.sessionRepo.CreateUserSession(
		user.ID,
		deviceID,
		hashed,
		ip,
		userAgent,
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	// Update last login info
	user.LastLoginIP = ip
	user.LastUserAgent = userAgent
	_ = s.repo.Update(user)

	// Audit log
	_ = s.audit.LogEvent(ctx, &auditdomain.AuditLog{
		UserID:    user.ID,
		UserEmail: user.Email,
		UserRole:  string(user.Role),
		Action:    "OAUTH_LOGIN",
		Entity:    "USER",
		EntityID:  strconv.Itoa(int(user.ID)),
		IPAddress: ip,
		UserAgent: userAgent,
	})

	user.Password = ""

	return &domain.AuthPayload{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		DeviceID:     deviceID,
	}, nil
}