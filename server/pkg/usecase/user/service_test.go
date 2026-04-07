package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	session "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user_session"
	auditdomain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	userservice "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/user"
)

//////////////////// MOCK REPO ////////////////////

type mockRepo struct {
	user *domain.User
}

type mockAuditRecorder struct {
	actions []string
}
type mockBlacklistBlacklisted struct{}
type expiredSessionRepo struct{}
type expiredResetRepo struct {
	mockRepo
}

var _ domain.Repository = (*mockRepo)(nil) // ✅ ensures full implementation

func (m *mockRepo) Create(u *domain.User) error {
	m.user = u
	return nil
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.user == nil {
		return nil, domain.ErrUserNotFound
	}

	if m.user.Email != email {
		return nil, domain.ErrUserNotFound
	}

	return m.user, nil
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.user == nil {
		return nil, errors.New("not found")
	}
	return m.user, nil
}

func (m *mockRepo) FindAll(page, limit int) ([]*domain.User, int, error) {
	return []*domain.User{m.user}, 1, nil
}

func (m *mockRepo) Update(u *domain.User) error {
	m.user = u
	return nil
}

func (m *mockRepo) WithTx(fn func(repo domain.Repository) error) error {
	return fn(m)
}



// -------- verification --------
func (m *mockRepo) SaveVerificationToken(userID uint, token string, expiresAt time.Time) error {
	return nil
}

func (m *mockRepo) FindVerificationByToken(token string) (*domain.EmailVerification, error) {
	return nil, nil
}

func (m *mockRepo) DeleteVerificationToken(id uint) error { return nil }

// -------- refresh tokens (REQUIRED FIX) --------
func (m *mockRepo) SaveRefreshToken(userID uint, token string, expiresAt time.Time) error {
	return nil
}

func (m *mockRepo) FindRefreshToken(token string) (*domain.RefreshToken, error) {
	return nil, nil
}

func (m *mockRepo) DeleteRefreshToken(id uint) error {
	return nil
}

func (m *mockRepo) ReplaceRefreshToken(id uint, token string, expiresAt time.Time) error {
	return nil
}

// -------- password reset --------
func (m *mockRepo) SavePasswordResetToken(userID uint, token string, expiresAt time.Time) error {
	return nil
}

func (m *mockRepo) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {
	return &domain.PasswordReset{
		ID:        1,
		UserID:    m.user.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *mockRepo) DeletePasswordResetToken(id uint) error { return nil }

// -------- cleanup --------
func (m *mockRepo) DeleteAllUserSessions(userID uint) error { return nil }
func (m *mockRepo) DeleteAllRefreshTokens(userID uint) error { return nil }

//////////////////// SESSION MOCK ////////////////////

type mockSessionRepo struct{}

func (m *mockSessionRepo) CreateUserSession(userID uint, deviceID, hash, ip, ua string, exp time.Time) error {
	return nil
}

func (m *mockSessionRepo) DeleteSession(userID uint, deviceID string) error { return nil }
func (m *mockSessionRepo) DeleteAllUserSessions(userID uint) error          { return nil }

func (m *mockSessionRepo) FindByToken(hash string) (*session.Session, error) {
	return &session.Session{
		ID:        1,
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}
func (m *mockAuditRecorder) LogEvent(ctx context.Context, event *auditdomain.AuditLog) error {
	m.actions = append(m.actions, event.Action)
	return nil
}
	

func (m *mockSessionRepo) GetUserSessions(userID uint) ([]*session.Session, error) {
	return nil, nil
}

func (m *mockSessionRepo) UpdateRefreshToken(sessionID uint, hash string, exp time.Time) error {
	return nil
}

//////////////////// RATE LIMITER ////////////////////

type mockRateLimiter struct{}

func (m *mockRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	return true, nil
}

type deniedRateLimiter struct{}

func (d *deniedRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	return false, nil
}

//////////////////// AUDIT ////////////////////

type mockAudit struct{}

func (m *mockAudit) LogEvent(ctx context.Context, event *auditdomain.AuditLog) error {
	return nil
}

//////////////////// BLACKLIST ////////////////////

type mockBlacklist struct{}

func (m *mockBlacklist) Blacklist(token string, d time.Duration) error { return nil }
func (m *mockBlacklist) IsBlacklisted(token string) bool               { return false }
func (m *mockBlacklistBlacklisted) Blacklist(string, time.Duration) error { return nil }
func (m *mockBlacklistBlacklisted) IsBlacklisted(string) bool  {return true}


//////////////////// SETUP ////////////////////

func setupService() *userservice.Service {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	u := &domain.User{
		ID:              1,
		Email:           "test@mail.com",
		Password:        string(hash),
		IsActive:        true,
		IsEmailVerified: true,
		Role:            common.RoleUser,
	}

	return userservice.NewService(
		&mockRepo{user: u},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)
}

func setupUser() *domain.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	return &domain.User{
		ID:              1,
		Name:            "Test User",
		Email:           "test@mail.com",
		Password:        string(hash),
		IsActive:        true,
		IsEmailVerified: true,
		Role:            common.RoleUser,
	}
}

//////////////////// TESTS ////////////////////

func TestRegister(t *testing.T) {

	service := userservice.NewService(
		&mockRepo{user: nil}, // ✅ empty repo
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	u, err := service.Register(context.Background(), "Test", "new@mail.com", "password")

	assert.NoError(t, err)
	assert.NotNil(t, u)
}
func TestLogin_Success(t *testing.T) {
	service := setupService()

	res, err := service.Login(context.Background(), "test@mail.com", "device1", "password", "127.0.0.1", "Chrome")

	assert.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
}

func TestLogin_InvalidPassword(t *testing.T) {
	service := setupService()

	_, err := service.Login(context.Background(), "test@mail.com", "device1", "wrong", "127.0.0.1", "Chrome")

	assert.Error(t, err)
}

func TestLogin_RateLimited(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{},
		&deniedRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), "test@mail.com", "device1", "password", "127.0.0.1", "Chrome")

	assert.Equal(t, domain.ErrRateLimited, err)
}

func TestOAuthLogin(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{user: nil},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	res, err := service.OAuthLogin(
		context.Background(),
		"Google User",
		"google@mail.com",
		"127.0.0.1",
		"Chrome",
		"device123",
	)

	assert.NoError(t, err)
	assert.NotNil(t, res.User)
}

func TestRefreshToken(t *testing.T) {
	service := setupService()

	res, err := service.RefreshToken(context.Background(), "token")

	assert.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
}

func TestResetPassword(t *testing.T) {
	service := setupService()

	err := service.ResetPassword(context.Background(), "token", "newpass")

	assert.NoError(t, err)
}

func TestPromoteToAdmin(t *testing.T) {
	service := setupService()

	_, err := service.PromoteToAdmin(context.Background(), 1, 1)

	assert.Error(t, err)
}

func TestDeactivateUser(t *testing.T) {
	service := setupService()

	_, err := service.DeactivateUser(context.Background(), 1, 1)

	assert.Error(t, err)
}

func TestGetUserByID(t *testing.T) {
	service := setupService()

	u, err := service.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), u.ID)
}

func TestLogoutAllDevices(t *testing.T) {
	service := setupService()

	err := service.LogoutAllDevices(context.Background(), 1, "token")

	assert.NoError(t, err)
}
func TestLogin_AuditLogTriggered(t *testing.T) {

	user := setupUser()

	audit := &mockAuditRecorder{}

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		audit,
		&mockBlacklist{},
	)

	_, err := service.Login(
		context.Background(),
		user.Email,
		"device1",
		"password",
		"127.0.0.1",
		"Chrome",
	)

	assert.NoError(t, err)
	assert.Contains(t, audit.actions, "LOGIN_SUCCESS")
assert.NotContains(t, audit.actions, "SUSPICIOUS_LOGIN")
}

func TestLogin_NormalAudit(t *testing.T) {
	user := setupUser()

	audit := &mockAuditRecorder{}

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		audit,
		&mockBlacklist{},
	)

	_, _ = service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Contains(t, audit.actions, "LOGIN_SUCCESS")
	assert.NotContains(t, audit.actions, "SUSPICIOUS_LOGIN")
}

func TestLogin_SuspiciousAudit(t *testing.T) {
	user := setupUser()
	user.LastLoginIP = "1.1.1.1"
	user.LastUserAgent = "Old"

	audit := &mockAuditRecorder{}

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		audit,
		&mockBlacklist{},
	)

	_, _ = service.Login(context.Background(), user.Email, "dev", "password", "2.2.2.2", "New")

	assert.Contains(t, audit.actions, "SUSPICIOUS_LOGIN")
}

type mockSessionRepoWithTracking struct {
	calls int
}

func (m *mockSessionRepoWithTracking) CreateUserSession(userID uint, deviceID, hash, ip, ua string, exp time.Time) error {
	m.calls++
	return nil
}

func (m *mockSessionRepoWithTracking) DeleteSession(userID uint, deviceID string) error { return nil }
func (m *mockSessionRepoWithTracking) DeleteAllUserSessions(userID uint) error          { return nil }

func (m *mockSessionRepoWithTracking) FindByToken(hash string) (*session.Session, error) {
	return &session.Session{
		ID:        1,
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *mockSessionRepoWithTracking) GetUserSessions(userID uint) ([]*session.Session, error) {
	return nil, nil
}

func (m *mockSessionRepoWithTracking) UpdateRefreshToken(sessionID uint, hash string, exp time.Time) error {
	return nil
}

func TestLogin_MultipleDevices(t *testing.T) {

	user := setupUser()
	sessionRepo := &mockSessionRepoWithTracking{}

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		sessionRepo,
		&mockAudit{},
		&mockBlacklist{},
	)

	// first login
	_, err1 := service.Login(context.Background(), user.Email, "device1", "password", "ip", "ua")
	assert.NoError(t, err1)

	//  FIX: restore password (because service wipes it)
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user.Password = string(hash)

	// second login
	_, err2 := service.Login(context.Background(), user.Email, "device2", "password", "ip", "ua")
	assert.NoError(t, err2)

	assert.Equal(t, 2, sessionRepo.calls)
}

func TestLogin_EmailNotVerified(t *testing.T) {
	user := setupUser()
	user.IsEmailVerified = false

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Equal(t, domain.ErrEmailNotVerified, err)
}

func TestLogin_UserInactive(t *testing.T) {
	user := setupUser()
	user.IsActive = false

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Equal(t, domain.ErrUserInactive, err)
}

func TestLogin_AccountLocked(t *testing.T) {
	user := setupUser()
	lock := time.Now().Add(10 * time.Minute)
	user.AccountLockedUntil = &lock

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Equal(t, domain.ErrAccountLocked, err)
}


func TestLogin_LockAfterFailedAttempts(t *testing.T) {
	user := setupUser()
	user.FailedLoginAttempts = 4

	repo := &mockRepo{user: user}

	service := userservice.NewService(
		repo,
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "wrong", "ip", "ua")

	assert.Error(t, err)
	assert.NotNil(t, user.AccountLockedUntil)
}

func TestRefreshToken_Blacklisted(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklistBlacklisted{},
	)

	_, err := service.RefreshToken(context.Background(), "token")

	assert.Equal(t, domain.ErrInvalidOrExpiredToken, err)
}



func (e *expiredSessionRepo) CreateUserSession(
	userID uint,
	deviceID, hash, ip, ua string,
	exp time.Time,
) error {
	return nil
}

func (e *expiredSessionRepo) DeleteSession(userID uint, deviceID string) error {
	return nil
}

func (e *expiredSessionRepo) DeleteAllUserSessions(userID uint) error {
	return nil
}

func (e *expiredSessionRepo) FindByToken(hash string) (*session.Session, error) {
	return &session.Session{
		ID:        1,
		UserID:    1,
		ExpiresAt: time.Now().Add(-time.Hour), // ❌ expired
	}, nil
}

func (e *expiredSessionRepo) GetUserSessions(userID uint) ([]*session.Session, error) {
	return nil, nil
}

func (e *expiredSessionRepo) UpdateRefreshToken(sessionID uint, hash string, exp time.Time) error {
	return nil
}

func TestRefreshToken_Expired(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{user: setupUser()},
		&mockRateLimiter{},
		&expiredSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.RefreshToken(context.Background(), "token")

	assert.Equal(t, domain.ErrInvalidOrExpiredToken, err)
}

func (e *expiredResetRepo) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {
	return &domain.PasswordReset{
		ID:        1,
		UserID:    1,
		ExpiresAt: time.Now().Add(-time.Hour),
	}, nil
}

func TestResetPassword_Expired(t *testing.T) {
	service := userservice.NewService(
		&expiredResetRepo{mockRepo{user: setupUser()}},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	err := service.ResetPassword(context.Background(), "token", "newpass")

	assert.Equal(t, domain.ErrInvalidOrExpiredToken, err)
}

func TestPromoteToAdmin_NotSuperAdmin(t *testing.T) {
	user := setupUser()

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.PromoteToAdmin(context.Background(), 1, 2)

	assert.Error(t, err)
}

func TestOAuthLogin_ExistingUser(t *testing.T) {
	user := setupUser()

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	res, err := service.OAuthLogin(
		context.Background(),
		user.Name,
		user.Email,
		"ip",
		"ua",
		"device",
	)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, res.User.Email)
}

//////////////////// EXTRA COVERAGE TESTS ////////////////////

// ---------- REGISTER ERROR ----------

type errorRepo struct {
	mockRepo
}

func (e *errorRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, errors.New("db error")
}

func TestRegister_RepoError(t *testing.T) {
	service := userservice.NewService(
		&errorRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Register(context.Background(), "Test", "fail@mail.com", "password")

	assert.Error(t, err)
}

// ---------- LOGIN UPDATE FAIL ----------

type updateErrorRepo struct {
	mockRepo
}

func (u *updateErrorRepo) Update(user *domain.User) error {
	return errors.New("update failed")
}

func TestLogin_UpdateFails(t *testing.T) {
	user := setupUser()

	service := userservice.NewService(
		&updateErrorRepo{mockRepo{user: user}},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Error(t, err)
}

// ---------- SESSION FAIL ----------

type sessionFailRepo struct {
	mockSessionRepo
}

func (s *sessionFailRepo) CreateUserSession(
	userID uint, deviceID, hash, ip, ua string, exp time.Time,
) error {
	return errors.New("session failed")
}

func TestLogin_SessionFails(t *testing.T) {
	user := setupUser()

	service := userservice.NewService(
		&mockRepo{user: user},
		&mockRateLimiter{},
		&sessionFailRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), user.Email, "dev", "password", "ip", "ua")

	assert.Error(t, err)
}

// ---------- REFRESH TOKEN UPDATE FAIL ----------

type updateRefreshFail struct {
	mockSessionRepo
}

func (u *updateRefreshFail) UpdateRefreshToken(sessionID uint, hash string, exp time.Time) error {
	return errors.New("update failed")
}

func TestRefreshToken_UpdateFails(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{user: setupUser()},
		&mockRateLimiter{},
		&updateRefreshFail{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.RefreshToken(context.Background(), "token")

	assert.Error(t, err)
}

// ---------- LOGOUT FAIL ----------

type logoutErrorRepo struct {
	mockRepo
}

func (l *logoutErrorRepo) DeleteAllRefreshTokens(userID uint) error {
	return errors.New("fail")
}

func TestLogout_Error(t *testing.T) {
	service := userservice.NewService(
		&logoutErrorRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	err := service.LogoutAllDevices(context.Background(), 1, "token")

	assert.Error(t, err)
}

// ---------- DETECT SUSPICIOUS LOGIN ----------

func TestDetectSuspiciousLogin(t *testing.T) {
	assert.False(t, userservice.DetectSuspiciousLogin("", "ip", "", "ua"))
	assert.False(t, userservice.DetectSuspiciousLogin("ip", "ip", "ua", "ua"))
	assert.True(t, userservice.DetectSuspiciousLogin("ip1", "ip2", "ua", "ua"))
	assert.True(t, userservice.DetectSuspiciousLogin("ip", "ip", "ua1", "ua2"))
}

//////////////////// MISSING EDGE CASES ////////////////////

// -------- Register TX FAIL --------

type txFailRepo struct {
	mockRepo
}

func (t *txFailRepo) WithTx(fn func(repo domain.Repository) error) error {
	return errors.New("tx failed")
}

func TestRegister_TxFails(t *testing.T) {
	service := userservice.NewService(
		&txFailRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Register(context.Background(), "Test", "fail@mail.com", "pass")

	assert.Error(t, err)
}

// -------- Rate limiter ERROR --------

type errorRateLimiter struct{}

func (e *errorRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	return false, errors.New("redis down")
}

func TestLogin_RateLimiterError(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{},
		&errorRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), "mail", "dev", "pass", "ip", "ua")

	assert.Error(t, err)
}

// -------- FindByEmail ERROR --------

type findEmailErrorRepo struct {
	mockRepo
}

func (f *findEmailErrorRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, errors.New("db error")
}

func TestLogin_FindByEmailError(t *testing.T) {
	service := userservice.NewService(
		&findEmailErrorRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.Login(context.Background(), "mail", "dev", "pass", "ip", "ua")

	assert.Equal(t, domain.ErrInvalidCredentials, err)
}

// -------- Refresh FindByToken ERROR --------

type findTokenErrorRepo struct {
	mockSessionRepo
}

func (f *findTokenErrorRepo) FindByToken(hash string) (*session.Session, error) {
	return nil, errors.New("not found")
}

func TestRefreshToken_FindError(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{},
		&mockRateLimiter{},
		&findTokenErrorRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.RefreshToken(context.Background(), "token")

	assert.Equal(t, domain.ErrInvalidOrExpiredToken, err)
}

// -------- Refresh FindByID ERROR --------

type findByIDErrorRepo struct {
	mockRepo
}

func (f *findByIDErrorRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	return nil, errors.New("fail")
}

func TestRefreshToken_UserFetchFails(t *testing.T) {
	service := userservice.NewService(
		&findByIDErrorRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.RefreshToken(context.Background(), "token")

	assert.Error(t, err)
}

// -------- ResetPassword update FAIL --------

type updateFailRepo struct {
	mockRepo
}

func (u *updateFailRepo) Update(user *domain.User) error {
	return errors.New("fail")
}

func TestResetPassword_UpdateFails(t *testing.T) {
	service := userservice.NewService(
		&updateFailRepo{mockRepo{user: setupUser()}},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	err := service.ResetPassword(context.Background(), "token", "pass")

	assert.Error(t, err)
}

// -------- OAuth Create FAIL --------

type oauthCreateFailRepo struct {
	mockRepo
}

func (o *oauthCreateFailRepo) Create(u *domain.User) error {
	return errors.New("fail")
}

func TestOAuthLogin_CreateFails(t *testing.T) {
	service := userservice.NewService(
		&oauthCreateFailRepo{},
		&mockRateLimiter{},
		&mockSessionRepo{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.OAuthLogin(context.Background(), "name", "mail", "ip", "ua", "dev")

	assert.Error(t, err)
}

// -------- OAuth Session FAIL --------

type oauthSessionFail struct {
	mockSessionRepo
}

func (o *oauthSessionFail) CreateUserSession(
	userID uint, deviceID, hash, ip, ua string, exp time.Time,
) error {
	return errors.New("fail")
}

func TestOAuthLogin_SessionFails(t *testing.T) {
	service := userservice.NewService(
		&mockRepo{user: nil},
		&mockRateLimiter{},
		&oauthSessionFail{},
		&mockAudit{},
		&mockBlacklist{},
	)

	_, err := service.OAuthLogin(context.Background(), "name", "mail", "ip", "ua", "dev")

	assert.Error(t, err)
}