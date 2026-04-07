// package booking_test

// import (
// 	"context"
// 	"errors"
// 	"sync"
// 	"sync/atomic"
// 	"testing"
// 	"time"

// 	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
// 	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
// 	"gorm.io/gorm"
// )

// type mockRepo struct {
// 	bookings  map[uint]*domain.Booking
// 	repoError bool
// 	overlap   bool
// 	lockedRoomID *uint
// 	locks  map[uint]*sync.Mutex
// 	global   sync.Mutex
// }




// func newMockRepo() *mockRepo {
// 	return &mockRepo{
// 		bookings: make(map[uint]*domain.Booking),
// 		locks:    make(map[uint]*sync.Mutex),
// 	}
// }

// func (m *mockRepo) Create(b *domain.Booking) error {
// 	m.global.Lock()
// 	defer m.global.Unlock()
	
// 	b.ID = uint(len(m.bookings) + 1)
// 	m.bookings[b.ID] = b
// 	return nil
// }

// func (m *mockRepo) Update(b *domain.Booking) error {
// 	m.bookings[b.ID] = b
// 	return nil
// }

// func (m *mockRepo) Delete(id uint) error { return nil }

// func (m *mockRepo) FindByID(id uint) (*domain.Booking, error) {
// 	b, ok := m.bookings[id]
// 	if !ok {
// 		return nil, errors.New("not found")
// 	}
// 	return b, nil
// }

// func (m *mockRepo) List(ctx context.Context, page, limit int) ([]domain.Booking, int64, error) {
// 	var result []domain.Booking
// 	for _, b := range m.bookings {
// 		result = append(result, *b)
// 	}
// 	return result, int64(len(result)), nil
// }

// func (m *mockRepo) FindByUser(ctx context.Context, userID uint, page, limit int) ([]domain.Booking, int64, error) {
// 	var result []domain.Booking
// 	for _, b := range m.bookings {
// 		if b.UserID == userID {
// 			result = append(result, *b)
// 		}
// 	}
// 	return result, int64(len(result)), nil
// }

// func (m *mockRepo) FindOverlappingBookings(roomID uint, checkIn, checkOut time.Time) ([]domain.Booking, error) {

// 	if m.repoError {
// 		return nil, errors.New("db error")
// 	}

// 	var result []domain.Booking

// 	for _, b := range m.bookings {
// 		if b.RoomID == roomID &&
// 			checkIn.Before(b.CheckOutDate) &&
// 			checkOut.After(b.CheckInDate) {

// 			result = append(result, *b)
// 		}
// 	}

	
// 	if m.overlap {
// 		return []domain.Booking{{ID: 999}}, nil
// 	}

// 	return result, nil
// }

// func (m *mockRepo) FindExpiredBookings(now time.Time) ([]domain.Booking, error) {
// 	return []domain.Booking{
// 		{
// 			ID:     1,
// 			Status: domain.BookingStatusPending,
// 		},
// 	}, nil
// }


// func (m *mockRepo) FindByRoomID(roomID uint) ([]domain.Booking, error) { return nil, nil }

// func (m *mockRepo) FindAll(status *domain.PaymentStatus, ctx context.Context, page, limit int) ([]domain.Booking, int64, error) {
// 	return nil, 0, nil
// }

// func (m *mockRepo) UpdatePaymentStatus(id uint, status domain.PaymentStatus) error {
// 	return nil
// }


// func (m *mockRepo) UpdatePaymentStatusTx(tx *gorm.DB, id uint, status domain.PaymentStatus) error {
// 	return nil
// }

// var ctx = context.Background()

// func (m *mockRepo) LockRoom(roomID uint) error {
// 	m.global.Lock()

// 	defer m.global.Unlock()
// 	return nil
// }

// func (m *mockRepo) WithTransaction(fn func(repo domain.Repository) error) error {
	
// 	return fn(m)
// }

// // ---------------- TESTS ----------------

// func TestCreateBooking_Success(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		UserID:       1,
// 		RoomID:       1,
// 		CheckInDate:  time.Now().Add(24 * time.Hour),
// 		CheckOutDate: time.Now().Add(48 * time.Hour),
// 	}

// 	err := service.CreateBooking(ctx, b)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if b.Status != domain.BookingStatusPending {
// 		t.Fatal("expected PENDING booking")
// 	}

// 	if b.PaymentStatus != domain.PaymentStatusPending {
// 		t.Fatal("expected payment PENDING")
// 	}
// }

// func TestCreateBooking_InvalidDates(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		UserID:       1,
// 		RoomID:       1,
// 		CheckInDate:  time.Now(),
// 		CheckOutDate: time.Now().Add(-24 * time.Hour),
// 	}

// 	err := service.CreateBooking(ctx, b)

// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// }

// func TestConfirmBooking_Success(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusPending,
// 		PaymentStatus: domain.PaymentStatusCompleted,
// 		ExpiresAt:     time.Now().Add(time.Hour),
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmBooking(ctx, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if b.Status != domain.BookingStatusConfirmed {
// 		t.Fatal("expected confirmed")
// 	}
// }

// func TestConfirmBooking_Unpaid(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusPending,
// 		PaymentStatus: domain.PaymentStatusPending,
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected unpaid error")
// 	}
// }

// func TestCancelBooking_Success(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:     1,
// 		Status: domain.BookingStatusPending,
// 	}

// 	repo.bookings[1] = b

// 	err := service.CancelBooking(ctx, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if b.Status != domain.BookingStatusCancelled {
// 		t.Fatal("expected cancelled")
// 	}
// }

// func TestConfirmPaymentAndBooking_Success(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusPending,
// 		PaymentStatus: domain.PaymentStatusPending,
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmPaymentAndBooking(ctx, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if b.Status != domain.BookingStatusConfirmed {
// 		t.Fatal("expected confirmed")
// 	}

// 	if b.PaymentStatus != domain.PaymentStatusCompleted {
// 		t.Fatal("expected payment completed")
// 	}
// }

// func TestCreateBooking_Overlapping(t *testing.T) {

// 	repo := newMockRepo()

// 	// ✅ Seed existing booking (already booked)
// 	repo.bookings[1] = &domain.Booking{
// 		RoomID:       1,
// 		CheckInDate:  time.Now().Add(1 * time.Hour),
// 		CheckOutDate: time.Now().Add(5 * time.Hour),
// 	}

// 	service := usecase.NewService(repo)

// 	// ❗ New booking that overlaps
// 	b := &domain.Booking{
// 		UserID:       1,
// 		RoomID:       1,
// 		CheckInDate:  time.Now().Add(2 * time.Hour),
// 		CheckOutDate: time.Now().Add(4 * time.Hour),
// 	}

// 	err := service.CreateBooking(ctx, b)

// 	if err == nil {
// 		t.Fatal("expected overlap error")
// 	}
// }

// func TestGetBookingByID(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	repo.bookings[1] = &domain.Booking{ID: 1}

// 	b, err := service.GetBookingByID(ctx, 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if b.ID != 1 {
// 		t.Fatal("wrong booking")
// 	}
// }

// func TestGetBookingsByUserID(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	repo.bookings[1] = &domain.Booking{ID: 1, UserID: 5}
// 	repo.bookings[2] = &domain.Booking{ID: 2, UserID: 5}
// 	repo.bookings[3] = &domain.Booking{ID: 3, UserID: 7}

// 	list, _, err := service.GetBookingsByUserID(ctx, 5, 1, 10)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(list) != 2 {
// 		t.Fatal("expected 2 bookings")
// 	}
// }

// func TestConfirmBooking_Expired(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusPending,
// 		PaymentStatus: domain.PaymentStatusCompleted,
// 		ExpiresAt:     time.Now().Add(-time.Hour), // expired
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected expired error")
// 	}
// }

// func TestConfirmBooking_Cancelled(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusCancelled,
// 		PaymentStatus: domain.PaymentStatusCompleted,
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected cancelled error")
// 	}
// }

// func TestCancelBooking_AlreadyCancelled(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:     1,
// 		Status: domain.BookingStatusCancelled,
// 	}

// 	repo.bookings[1] = b

// 	err := service.CancelBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected already cancelled error")
// 	}
// }

// func TestConfirmPaymentAndBooking_AlreadyPaid(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusPending,
// 		PaymentStatus: domain.PaymentStatusCompleted,
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmPaymentAndBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected already paid error")
// 	}
// }

// func TestConfirmPaymentAndBooking_Cancelled(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		ID:            1,
// 		Status:        domain.BookingStatusCancelled,
// 		PaymentStatus: domain.PaymentStatusPending,
// 	}

// 	repo.bookings[1] = b

// 	err := service.ConfirmPaymentAndBooking(ctx, 1)

// 	if err == nil {
// 		t.Fatal("expected cancelled error")
// 	}
// }

// func TestCreateBooking_RepoError(t *testing.T) {
// 	repo := newMockRepo()
// 	repo.repoError = true

// 	service := usecase.NewService(repo)

// 	b := &domain.Booking{
// 		UserID:       1,
// 		RoomID:       1,
// 		CheckInDate:  time.Now().Add(time.Hour),
// 		CheckOutDate: time.Now().Add(48 * time.Hour),
// 	}

// 	err := service.CreateBooking(ctx, b)

// 	if err == nil {
// 		t.Fatal("expected repo error")
// 	}
// }

// func TestCheckRoomAvailability_Success(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	ok, err := service.CheckRoomAvailability(
// 		ctx,
// 		1,
// 		time.Now().Add(time.Hour),
// 		time.Now().Add(48*time.Hour),
// 	)

// 	if err != nil || !ok {
// 		t.Fatal("expected available room")
// 	}
// }

// func TestCheckRoomAvailability_InvalidDates(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	ok, err := service.CheckRoomAvailability(
// 		ctx,
// 		1,
// 		time.Now(),
// 		time.Now().Add(-time.Hour),
// 	)

// 	if err == nil || ok {
// 		t.Fatal("expected invalid date error")
// 	}
// }

// func TestCheckRoomAvailability_Overlap(t *testing.T) {
// 	repo := newMockRepo()

// 	//  Seed existing booking
// 	repo.bookings[1] = &domain.Booking{
// 		RoomID:       1,
// 		CheckInDate:  time.Now().Add(1 * time.Hour),
// 		CheckOutDate: time.Now().Add(5 * time.Hour),
// 	}

// 	service := usecase.NewService(repo)

// 	ok, err := service.CheckRoomAvailability(
// 		ctx,
// 		1,
// 		time.Now().Add(2*time.Hour), // overlaps
// 		time.Now().Add(4*time.Hour),
// 	)

// 	if err != nil || ok {
// 		t.Fatal("expected unavailable room")
// 	}
// }

// func TestExpireBookings_Success(t *testing.T) {
// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	err := service.ExpireBookings(ctx)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// type expireErrorRepo struct {
// 	mockRepo
// }

// func (e *expireErrorRepo) FindExpiredBookings(now time.Time) ([]domain.Booking, error) {
// 	return nil, errors.New("fail")
// }

// func TestExpireBookings_Error(t *testing.T) {
// 	repo := &expireErrorRepo{}
// 	service := usecase.NewService(repo)

// 	err := service.ExpireBookings(ctx)

// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// }

// func TestCreateBooking_Concurrent(t *testing.T) {

// 	repo := newMockRepo()
// 	service := usecase.NewService(repo)

// 	var successCount int32 // ✅ shared
// 	total := 10

// 	var wg sync.WaitGroup

// 	for i := 0; i < total; i++ {
// 		wg.Add(1)

// 		go func() {
// 			defer wg.Done()

// 			b := &domain.Booking{
// 				UserID:       1,
// 				RoomID:       1,
// 				CheckInDate:  time.Now().Add(time.Hour),
// 				CheckOutDate: time.Now().Add(48 * time.Hour),
// 			}

// 			err := service.CreateBooking(ctx, b)

// 			if err == nil {
// 				atomic.AddInt32(&successCount, 1)
// 			}
// 		}()
// 	}

// 	wg.Wait()

// 	if atomic.LoadInt32(&successCount) != 1 {
// 		t.Fatalf("expected only 1 success, got %d", successCount)
// 	}
// }package booking_test


package booking_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/redis"

	goredis "github.com/redis/go-redis/v9"
)

/* =========================
   MOCK REPO
========================= */

type mockRepo struct {
	bookings map[uint]*domain.Booking
	mu       sync.Mutex
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		bookings: make(map[uint]*domain.Booking),
	}
}

func (m *mockRepo) Create(b *domain.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b.ID = uint(len(m.bookings) + 1)
	m.bookings[b.ID] = b
	return nil
}

func (m *mockRepo) Update(b *domain.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.bookings[b.ID] = b
	return nil
}

func (m *mockRepo) Delete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.bookings, id)
	return nil
}

func (m *mockRepo) FindByID(id uint) (*domain.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.bookings[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (m *mockRepo) FindOverlappingBookings(roomID uint, checkIn, checkOut time.Time) ([]domain.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, b := range m.bookings {
		if b.RoomID == roomID &&
			checkIn.Before(b.CheckOutDate) &&
			checkOut.After(b.CheckInDate) {
			return []domain.Booking{*b}, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) List(ctx context.Context, page, limit int) ([]domain.Booking, int64, error) {
	return nil, 0, nil
}

func (m *mockRepo) FindByUser(ctx context.Context, userID uint, page, limit int) ([]domain.Booking, int64, error) {
	return nil, 0, nil
}

func (m *mockRepo) FindExpiredBookings(now time.Time) ([]domain.Booking, error) {
	return []domain.Booking{}, nil
}

func (m *mockRepo) WithTransaction(fn func(repo domain.Repository) error) error {
	return fn(m)
}

func (m *mockRepo) FindAll(
	status *domain.PaymentStatus,
	ctx context.Context,
	page,
	limit int,
) ([]domain.Booking, int64, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	var result []domain.Booking

	for _, b := range m.bookings {
		result = append(result, *b)
	}

	return result, int64(len(result)), nil
}

func (m *mockRepo) FindByRoomID(roomID uint) ([]domain.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []domain.Booking

	for _, b := range m.bookings {
		if b.RoomID == roomID {
			result = append(result, *b)
		}
	}

	return result, nil
}
func (m *mockRepo) LockRoom(roomID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}
func (m *mockRepo) UpdatePaymentStatus(id uint, status domain.PaymentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.bookings[id]
	if !ok {
		return errors.New("not found")
	}

	b.PaymentStatus = status
	return nil
}

/* =========================
   HELPER
========================= */

func newService(repo *mockRepo) *usecase.Service {
	redisClient := goredis.NewClient(&goredis.Options{})

	locker := redis.NewRedisLocker(redisClient)

	return usecase.NewService( repo, nil, locker, redisClient)
}

/* =========================
   TESTS
========================= */

func TestCreateBooking_Success(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	b := &domain.Booking{
		UserID:       1,
		RoomID:       1,
		CheckInDate:  time.Now().Add(time.Hour),
		CheckOutDate: time.Now().Add(48 * time.Hour),
	}

	if err := service.CreateBooking(context.Background(), b); err != nil {
		t.Fatal(err)
	}
}

func TestCreateBooking_Overlap(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	repo.bookings[1] = &domain.Booking{
		RoomID:       1,
		CheckInDate:  time.Now().Add(time.Hour),
		CheckOutDate: time.Now().Add(5 * time.Hour),
	}

	b := &domain.Booking{
		RoomID:       1,
		CheckInDate:  time.Now().Add(2 * time.Hour),
		CheckOutDate: time.Now().Add(4 * time.Hour),
	}

	if service.CreateBooking(context.Background(), b) == nil {
		t.Fatal("expected overlap error")
	}
}

func TestCancelBooking(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	repo.bookings[1] = &domain.Booking{
		ID:     1,
		Status: domain.BookingStatusPending,
	}

	if err := service.CancelBooking(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestCheckAvailability_DBFallback(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	ok, err := service.CheckRoomAvailability(
		context.Background(),
		1,
		time.Now().Add(time.Hour),
		time.Now().Add(48*time.Hour),
	)

	if err != nil || !ok {
		t.Fatal("expected available")
	}
}

func TestCheckAvailability_InvalidDates(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	ok, err := service.CheckRoomAvailability(
		context.Background(),
		1,
		time.Now(),
		time.Now().Add(-time.Hour),
	)

	if err == nil || ok {
		t.Fatal("expected invalid date error")
	}
}

func TestConcurrentBooking(t *testing.T) {
	repo := newMockRepo()
	service := newService(repo)

	var success int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			b := &domain.Booking{
				RoomID:       1,
				CheckInDate:  time.Now().Add(time.Hour),
				CheckOutDate: time.Now().Add(48 * time.Hour),
			}

			if service.CreateBooking(context.Background(), b) == nil {
				atomic.AddInt32(&success, 1)
			}
		}()
	}

	wg.Wait()

	if success != 1 {
		t.Fatalf("expected 1 success, got %d", success)
	}
}