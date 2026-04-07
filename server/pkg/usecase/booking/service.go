package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	goredis "github.com/redis/go-redis/v9"
)

type Service struct {
	repo          booking.Repository
	tracer        trace.Tracer
	roomRepo      room.Repository
	redisLocker   *redis.RedisLocker
	Client        *goredis.Client
}

func NewService(repo booking.Repository, roomRepo room.Repository, locker *redis.RedisLocker, client *goredis.Client) *Service {
	return &Service{
		repo:        repo,
		tracer:      otel.Tracer("booking-service"),
		roomRepo:    roomRepo,
		redisLocker: locker,
		Client:      client,
	}
}

func (s *Service) AvailabilityKey(roomID uint, date time.Time) string {
	return fmt.Sprintf("availability:%d:%s", roomID, date.Format("2006-01-02"))
}

func (s *Service) CreateBooking(ctx context.Context, b *booking.Booking) error {

	ctx, span := s.tracer.Start(ctx, "CreateBooking")
	defer span.End()

	if err := b.ValidateDates(); err != nil {
		span.RecordError(err)
		return err
	}

	b.Reference = GenerateBookingReference()
	b.Status = booking.BookingStatusPending
	b.PaymentStatus = booking.PaymentStatusPending
	b.ExpiresAt = time.Now().Add(15 * time.Minute)

	return s.withRoomLock(ctx, b.RoomID, func() error {

		err := s.repo.WithTransaction(func(tx booking.Repository) error {

			existing, err := tx.FindOverlappingBookings(
				b.RoomID,
				b.CheckInDate,
				b.CheckOutDate,
			)
			if err != nil {
				return err
			}

			if len(existing) > 0 {
				return errors.New("room already booked for selected dates")
			}

			return tx.Create(b)
		})

		if err != nil {
			return err
		}

		pipe := s.Client.Pipeline()

		for d := b.CheckInDate; d.Before(b.CheckOutDate); d = d.AddDate(0, 0, 1) {
			key := s.AvailabilityKey(b.RoomID, d)
			pipe.Set(ctx, key, "1", 24*time.Hour)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ConfirmBooking(ctx context.Context, bookingID uint) error {

	ctx, span := s.tracer.Start(ctx, "ConfirmBooking")
	defer span.End()

	b, err := s.repo.FindByID(bookingID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if b.Status == booking.BookingStatusCancelled {
		return errors.New("cannot confirm cancelled booking")
	}

	if b.PaymentStatus != booking.PaymentStatusCompleted {
		return errors.New("cannot confirm unpaid booking")
	}

	if time.Now().After(b.ExpiresAt) {
		return errors.New("booking expired")
	}

	b.Status = booking.BookingStatusConfirmed

	return s.repo.Update(b)
}

func (s *Service) CancelBooking(ctx context.Context, bookingID uint) error {

	ctx, span := s.tracer.Start(ctx, "CancelBooking")
	defer span.End()

	b, err := s.repo.FindByID(bookingID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return s.withRoomLock(ctx, b.RoomID, func() error {

		if b.Status == booking.BookingStatusCancelled {
			return errors.New("booking already cancelled")
		}

		b.Cancel()

		if err := s.repo.Update(b); err != nil {
			return err
		}

		pipe := s.Client.Pipeline()

		for d := b.CheckInDate; d.Before(b.CheckOutDate); d = d.AddDate(0, 0, 1) {
			key := s.AvailabilityKey(b.RoomID, d)
			pipe.Set(ctx, key, "0", 24*time.Hour)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ListBookings(ctx context.Context, page, limit int) ([]booking.Booking, int64, error) {
	ctx, span := s.tracer.Start(ctx, "ListBookings")
	defer span.End()

	return s.repo.List(ctx, page, limit)
}

func (s *Service) GetBookingByID(ctx context.Context, id uint) (*booking.Booking, error) {
	ctx, span := s.tracer.Start(ctx, "GetBookingByID")
	defer span.End()

	return s.repo.FindByID(id)
}

func (s *Service) GetBookingsByUserID(ctx context.Context, userID uint, page, limit int) ([]booking.Booking, int64, error) {
	ctx, span := s.tracer.Start(ctx, "GetBookingsByUserID")
	defer span.End()

	return s.repo.FindByUser(ctx, userID, page, limit)
}

func (s *Service) ConfirmPaymentAndBooking(ctx context.Context, bookingID uint) error {

	ctx, span := s.tracer.Start(ctx, "ConfirmPaymentAndBooking")
	defer span.End()

	// get booking first to know roomID
	b, err := s.repo.FindByID(bookingID)
	if err != nil {
		return err
	}

	return s.withRoomLock(ctx, b.RoomID, func() error {
		return s.repo.WithTransaction(func(txRepo booking.Repository) error {

			b, err := txRepo.FindByID(bookingID)
			if err != nil {
				return err
			}

			if b.Status == booking.BookingStatusCancelled {
				return errors.New("cannot confirm cancelled booking")
			}

			if b.PaymentStatus == booking.PaymentStatusCompleted {
				return errors.New("payment already completed")
			}

			b.PaymentStatus = booking.PaymentStatusCompleted
			b.Status = booking.BookingStatusConfirmed

			return txRepo.Update(b)
		})
	})
}

func (s *Service) CheckRoomAvailability(ctx context.Context, roomID uint, checkIn, checkOut time.Time) (bool, error) {

	if !checkIn.Before(checkOut) {
		return false, errors.New("invalid booking dates")
	}

	ctx, span := s.tracer.Start(ctx, "CheckRoomAvailability")
	defer span.End()

	pipe := s.Client.Pipeline()
	cmds := make([]*goredis.StringCmd, 0)

	for d := checkIn; d.Before(checkOut); d = d.AddDate(0, 0, 1) {
		cmds = append(cmds, pipe.Get(ctx, s.AvailabilityKey(roomID, d)))
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != goredis.Nil {
		return false, err
	}

	for _, cmd := range cmds {
		if cmd.Err() != nil && cmd.Err() != goredis.Nil {
			return false, cmd.Err()
		}
		if cmd.Err() == nil && cmd.Val() == "1" {
			return false, nil
		}
	}

	bookings, err := s.repo.FindOverlappingBookings(roomID, checkIn, checkOut)
	if err != nil {
		return false, err
	}

	available := len(bookings) == 0

	val := "0"
	if !available {
		val = "1"
	}

	for d := checkIn; d.Before(checkOut); d = d.AddDate(0, 0, 1) {
		s.Client.Set(ctx, s.AvailabilityKey(roomID, d), val, 10*time.Minute)
	}

	return available, nil
}

func (s *Service) ExpireBookings(ctx context.Context) error {

	ctx, span := s.tracer.Start(ctx, "ExpireBookings")
	defer span.End()

	now := time.Now()

	bookings, err := s.repo.FindExpiredBookings(now)
	if err != nil {
		span.RecordError(err)
		return err
	}

	for i := range bookings {
		b := &bookings[i]

		if b.Status == booking.BookingStatusPending {
			b.Status = booking.BookingStatusCancelled

			if err := s.repo.Update(b); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) withRoomLock(ctx context.Context, roomID uint, fn func() error) error {

	var unlock func()
	var err error

	for i := 0; i < 3; i++ {
		unlock, err = s.redisLocker.LockResource(ctx, "room", roomID, 5*time.Second)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
	}

	if err != nil {
		return err
	}
	defer unlock()

	return fn()
}