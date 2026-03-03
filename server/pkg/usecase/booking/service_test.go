


package booking_test

import (
	"errors"
	"testing"
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
)


// ---------------- MOCK REPOSITORY ----------------

type mockRepo struct {
	bookings map[uint]*domain.Booking
	overlap bool
	repoError bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		bookings: make(map[uint]*domain.Booking),
	}
}

func (m *mockRepo) Create(b *domain.Booking) error {
	b.ID = 1
	m.bookings[b.ID] = b
	return nil
}

func (m *mockRepo) Update(b *domain.Booking) error {
	m.bookings[b.ID] = b
	return nil
}

func (m *mockRepo) Delete(id uint) error { return nil }

func (m *mockRepo) FindByID(id uint) (*domain.Booking, error) {
	b, ok := m.bookings[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (m *mockRepo) FindByUser(userID uint) ([]domain.Booking, error) {
	var result []domain.Booking
	for _, b := range m.bookings {
		if b.UserID == userID {
			result = append(result, *b)
		}
	}
	return result, nil
}



func (m *mockRepo) List(page, limit int) ([]domain.Booking, int64, error) {
	var result []domain.Booking
	for _, b := range m.bookings {
		result = append(result, *b)
	}
	return result, int64(len(result)), nil
}

func (m *mockRepo) UpdatePaymentStatus(id uint, status domain.PaymentStatus) error {
	return nil
}

func (m *mockRepo) FindByRoomID(roomID uint) ([]domain.Booking, error) { return nil, nil }
func (m *mockRepo) FindAll(status *domain.PaymentStatus, page, limit int) ([]domain.Booking, int64, error) {
	return nil, 0, nil
}

func (m *mockRepo) WithTransaction(fn func(repo domain.Repository) error) error {
	return fn(m)
}

func (m *mockRepo) FindOverlappingBookings(
	roomID uint,
	checkIn,
	checkOut time.Time,
) ([]domain.Booking, error) {
	
	if m.repoError {
		return nil, errors.New("Db errors")
	}

	if m.overlap {
		return []domain.Booking{
			{ID: 99},
		}, nil
	}

	return nil, nil
}

// ---------------- TESTS ----------------


// ✅ CreateBooking success
func TestCreateBooking_Success(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		UserID:       1,
		RoomID:       1,
		CheckInDate:  time.Now().Add(24 * time.Hour),
		CheckOutDate: time.Now().Add(48 * time.Hour),
	}

	err := service.CreateBooking(b)
	if err != nil {
		t.Fatal(err)
	}

	if b.Status != domain.BookingStatusPending {
		t.Fatal("expected booking status PENDING")
	}

	if b.PaymentStatus != domain.PaymentStatusPending {
		t.Fatal("expected payment status PENDING")
	}
}


// ❌ CreateBooking invalid dates
func TestCreateBooking_InvalidDates(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		UserID:       1,
		RoomID:       1,
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(-24 * time.Hour),
	}

	err := service.CreateBooking(b)
	if err == nil {
		t.Fatal("expected error for invalid dates")
	}
}


// ✅ ConfirmBooking success
func TestConfirmBooking_Success(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:            1,
		UserID:        1,
		Status:        domain.BookingStatusPending,
		PaymentStatus: domain.PaymentStatusCompleted,
	}

	repo.bookings[1] = b

	err := service.ConfirmBooking(1)
	if err != nil {
		t.Fatal(err)
	}

	if b.Status != domain.BookingStatusConfirmed {
		t.Fatal("expected booking to be CONFIRMED")
	}
}


// ❌ ConfirmBooking unpaid
func TestConfirmBooking_Unpaid(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:            1,
		Status:        domain.BookingStatusPending,
		PaymentStatus: domain.PaymentStatusPending,
	}

	repo.bookings[1] = b

	err := service.ConfirmBooking(1)
	if err == nil {
		t.Fatal("expected error for unpaid booking")
	}
}


// ✅ CancelBooking success
func TestCancelBooking_Success(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:     1,
		Status: domain.BookingStatusPending,
	}

	repo.bookings[1] = b

	err := service.CancelBooking(1)
	if err != nil {
		t.Fatal(err)
	}

	if b.Status != domain.BookingStatusCancelled {
		t.Fatal("expected booking to be CANCELLED")
	}
}


// ❌ CancelBooking already cancelled
func TestCancelBooking_AlreadyCancelled(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:     1,
		Status: domain.BookingStatusCancelled,
	}

	repo.bookings[1] = b

	err := service.CancelBooking(1)
	if err == nil {
		t.Fatal("expected error for already cancelled booking")
	}
}


// ✅ ConfirmPaymentAndBooking success
func TestConfirmPaymentAndBooking_Success(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:            1,
		Status:        domain.BookingStatusPending,
		PaymentStatus: domain.PaymentStatusPending,
	}

	repo.bookings[1] = b

	err := service.ConfirmPaymentAndBooking(1)
	if err != nil {
		t.Fatal(err)
	}

	if b.Status != domain.BookingStatusConfirmed {
		t.Fatal("expected booking confirmed")
	}

	if b.PaymentStatus != domain.PaymentStatusCompleted {
		t.Fatal("expected payment completed")
	}
}
func TestCreateBooking_Overlapping(t *testing.T) {
	repo := newMockRepo()
	repo.overlap = true // simulate overlap

	service := usecase.NewService(repo)

	b := &domain.Booking{
		UserID:       1,
		RoomID:       1,
		CheckInDate:  time.Now().Add(24 * time.Hour),
		CheckOutDate: time.Now().Add(48 * time.Hour),
	}

	err := service.CreateBooking(b)
	if err == nil {
		t.Fatal("expected overlap error")
	}
}

func TestConfirmPaymentAndBooking_Cancelled(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:            1,
		Status:        domain.BookingStatusCancelled,
		PaymentStatus: domain.PaymentStatusPending,
	}
	repo.bookings[1] = b

	err := service.ConfirmPaymentAndBooking(1)
	if err == nil {
		t.Fatal("expected error for cancelled booking")
	}
}

func TestConfirmPaymentAndBooking_AlreadyPaid(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	b := &domain.Booking{
		ID:            1,
		Status:        domain.BookingStatusPending,
		PaymentStatus: domain.PaymentStatusCompleted,
	}
	repo.bookings[1] = b

	err := service.ConfirmPaymentAndBooking(1)
	if err == nil {
		t.Fatal("expected error for already paid")
	}
}

func TestListBookings(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	repo.bookings[1] = &domain.Booking{ID: 1}
	repo.bookings[2] = &domain.Booking{ID: 2}

	list, total, err := service.ListBookings(1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if total != 2 || len(list) != 2 {
		t.Fatal("expected 2 bookings")
	}
}
func TestGetBookingByID(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	repo.bookings[1] = &domain.Booking{ID: 1}

	b, err := service.GetBookingByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if b.ID != 1 {
		t.Fatal("wrong booking returned")
	}
}

func TestGetBookingsByUserID(t *testing.T) {
	repo := newMockRepo()
	service := usecase.NewService(repo)

	repo.bookings[1] = &domain.Booking{ID: 1, UserID: 5}
	repo.bookings[2] = &domain.Booking{ID: 2, UserID: 5}
	repo.bookings[3] = &domain.Booking{ID: 3, UserID: 7}

	list, err := service.GetBookingsByUserID(5)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatal("expected 2 bookings for user 5")
	}
}

func TestCreateBooking_RepoError(t *testing.T) {
    repo := newMockRepo()
    repo.repoError = true

    service := usecase.NewService(repo)

    b := &domain.Booking{
        UserID:       1,
        RoomID:       1,
        CheckInDate:  time.Now().Add(24 * time.Hour),
        CheckOutDate: time.Now().Add(48 * time.Hour),
    }

    err := service.CreateBooking(b)
    if err == nil {
        t.Fatal("expected repo error")
    }
}

func TestConfirmBooking_NotFound(t *testing.T) {
    repo := newMockRepo()
    service := usecase.NewService(repo)

    err := service.ConfirmBooking(99)
    if err == nil {
        t.Fatal("expected not found error")
    }
}
func TestCancelBooking_NotFound(t *testing.T) {
    repo := newMockRepo()
    service := usecase.NewService(repo)

    err := service.CancelBooking(99)
    if err == nil {
        t.Fatal("expected not found error")
    }
}
func TestConfirmPaymentAndBooking_NotFound(t *testing.T) {
    repo := newMockRepo()
    service := usecase.NewService(repo)

    err := service.ConfirmPaymentAndBooking(99)
    if err == nil {
        t.Fatal("expected not found error")
    }
}