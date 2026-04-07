package audit_logs_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	bookingDomain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"

	repositories "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/audit_logs"
	bookingRepo "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories/booking"

	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/audit_logs"
	bookingService "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockPublisher struct{}

func (m *mockPublisher) Publish(topic string, msg []byte) error {
	return nil
}

func setupDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(
		&models.AuditLog{},
		&models.Booking{},
	)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestAuditAndBookingIntegration(t *testing.T) {

	ctx := context.Background()
	db := setupDB(t)

	auditRepo := repositories.NewAuditRepository(db)
	pub := &mockPublisher{}
	auditSvc := usecase.NewService(auditRepo, pub)

	bookingRepository := bookingRepo.NewBookingRepository(db)
	bookingSvc := bookingService.NewService(bookingRepository)

	//-----------------------------------------
	// Test 1: Create audit log
	//-----------------------------------------

	before := json.RawMessage(`{"status":"PENDING"}`)
	after := json.RawMessage(`{"status":"CONFIRMED"}`)

	logEvent := &domain.AuditLog{
		UserID:      1,
		UserEmail:   "admin@test.com",
		UserRole:    "ADMIN",
		Action:      "UPDATE_BOOKING",
		Entity:      "BOOKING",
		EntityID:    "123",
		BeforeState: before,
		AfterState:  after,
	}

	err := auditSvc.LogEvent(ctx, logEvent)
	if err != nil {
		t.Fatal(err)
	}

	filter := domain.AuditFilter{
		UserID: &logEvent.UserID,
	}

	logs, total, err := auditSvc.GetAuditLogs(ctx, filter, 1, 10)
	if len(logs) == 0 {t.Fatal("expected logs")}
	if err != nil {
		t.Fatal(err)
	}

	if total != 1 {
		t.Fatalf("expected 1 log got %d", total)
	}

	//-----------------------------------------
	// Test 2: High risk detection
	//-----------------------------------------

	deleteLog := &domain.AuditLog{
		UserID:      2,
		UserEmail:   "user@test.com",
		UserRole:    "USER",
		Action:      "DELETE_USER",
		Entity:      "USER",
		EntityID:    "456",
		BeforeState: json.RawMessage(`{"active":true}`),
		AfterState:  json.RawMessage(`{}`),
	}

	err = auditSvc.LogEvent(ctx, deleteLog)
	if err != nil {
		t.Fatal(err)
	}

	if !deleteLog.Suspicious {
		t.Fatal("expected suspicious log")
	}

	if deleteLog.RiskLevel != "HIGH" {
		t.Fatalf("expected HIGH risk got %s", deleteLog.RiskLevel)
	}

	//-----------------------------------------
	// Test 3: Create booking
	//-----------------------------------------

	b := &bookingDomain.Booking{
		UserID:       1,
		RoomID:       101,
		CheckInDate:  time.Now().Add(24 * time.Hour),
		CheckOutDate: time.Now().Add(48 * time.Hour),
		Status:       bookingDomain.BookingStatusPending,
	}

	err = bookingSvc.CreateBooking(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	if b.ID == 0 {
		t.Fatal("booking ID should be set")
	}

	//-----------------------------------------
	// Test 4: Confirm payment
	//-----------------------------------------

	err = bookingSvc.ConfirmPaymentAndBooking(ctx, b.ID)
	if err != nil {
		t.Fatal(err)
	}

	//-----------------------------------------
	// Test 5: Overlapping booking
	//-----------------------------------------

	available, err := bookingSvc.CheckRoomAvailability(
		ctx,
		b.RoomID,
		b.CheckInDate,
		b.CheckOutDate,
	)

	if err != nil {
		t.Fatal(err)
	}

	if available {
		t.Fatal("room should NOT be available")
	}

	t.Log("Integration tests passed")
}