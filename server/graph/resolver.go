package graph

import (
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/audit_logs"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// valueOrEmpty is  a helper function to convert a nullable *string to string
func valueOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
	}	

type Resolver struct{
	UserService user.Service
	AuditLogService audit_logs.Service
	RoomService room.Service
	BookingService booking.Service
	
}
