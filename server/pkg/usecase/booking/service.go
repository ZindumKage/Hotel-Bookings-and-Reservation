package booking

import (
	"errors"
	"time"

	
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBooking(booking *domain.Booking) error{
	if err := booking.ValidateDates(); err != nil {
		return err
	}
	existing, err := s.repo.FindOverlappingBookings(
		booking.RoomID,
		booking.CheckInDate,
		booking.CheckOutDate,
	)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return  errors.New("Room already booked for the selected dates")
}
booking.Status = domain.BookingStatusPending
booking.PaymentStatus = domain.PaymentStatusPending
booking.ExpiresAt = time.Now().Add(15 * time.Minute)
return s.repo.Create(booking)
}

func (s *Service) ConfirmBooking(bookingID uint) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return err
	}
	if booking.PaymentStatus != domain.PaymentStatusCompleted{
		return errors.New("Cannot Confirm Unpaid Booking")
	}
	booking.Status = domain.BookingStatusConfirmed
	return s.repo.Update(booking)
}

func (s *Service) CancelBooking(bookingID uint) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return err
	}
	if booking.Status == domain.BookingStatusCancelled {
		return errors.New("Booking is already cancelled")
	}
	booking.Cancel()
	return s.repo.Update(booking)
}

func (s *Service) ListBookings(page, limit int) ([]domain.Booking, int64, error) {
	return s.repo.List(page, limit)
}

func (s *Service) GetBookingByID(id uint) (*domain.Booking, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetBookingsByUserID(userID uint) ([]domain.Booking, error) {
	return s.repo.FindByUser(userID)
}

func (s *Service) ConfirmPaymentAndBooking(bookingID uint) error {
	return s.repo.WithTransaction(func(txRepo domain.Repository) error {

		booking, err := txRepo.FindByID(bookingID)
		if err != nil {
			return err
		}

		if booking.Status == domain.BookingStatusCancelled {
			return errors.New("Cannot confirm cancelled booking")
		}

		if booking.PaymentStatus == domain.PaymentStatusCompleted {
			return errors.New("Payment already completed")
		}

		booking.PaymentStatus = domain.PaymentStatusCompleted
		booking.Status = domain.BookingStatusConfirmed

		if err := txRepo.Update(booking); err != nil {
			return err
		}

		return nil
	})
}