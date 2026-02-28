package user

import (
	"time"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	Role      common.Role
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}