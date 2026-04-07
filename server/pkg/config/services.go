package config

import (
	

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database"
	session "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/user_session"

	auditrepo "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/audit_logs"
	bookingrepo "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/booking"
	roomrepo "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/room"
	userrepo "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/user"

	redisinfra "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/redis"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/audit_logs"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/user"
)

var (
	UserService     *user.Service
	RoomService     *room.Service
	BookingService  *booking.Service
	AuditLogService *audit_logs.Service
)

func InitServices() {

	userRepo := userrepo.NewUserRepository(database.DB)
	sessionRepo := session.NewUserSessionRepository(database.DB)
	roomRepo := roomrepo.NewRoomRepository(database.DB)
	bookingRepo := bookingrepo.NewBookingRepository(database.DB)
	auditRepo := auditrepo.NewAuditRepository(database.DB)


	rateLimiter := redisinfra.NewRedisRateLimiter(Redis)
	tokenBlacklist := redisinfra.NewRedisBlacklist(Redis)

	publisher := audit_logs.NewLogPublisher()

	AuditLogService = audit_logs.NewService(auditRepo, publisher)

	UserService = user.NewService(userRepo, rateLimiter, sessionRepo , AuditLogService, tokenBlacklist )
	RoomService = room.NewService(roomRepo)
	BookingService = booking.NewService(bookingRepo)
	AuditLogService = audit_logs.NewService(auditRepo, publisher)
}