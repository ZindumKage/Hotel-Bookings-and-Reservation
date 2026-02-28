package audit_logs_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/repositories"
	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/audit_logs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Dummy publisher
type mockPublisher struct{}

func (m *mockPublisher) Publish(topic string, msg []byte) error {
	return nil
}

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&audit_logs.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAuditLogIntegration(t *testing.T) {
	db := setupDB(t)

	repo := repositories.NewAuditRepository(db)
	pub := &mockPublisher{}
	service := usecase.NewService(repo, pub)

	before := json.RawMessage(`{"field1":"old"}`)
	after := json.RawMessage(`{"field1":"new"}`)

	log := &audit_logs.AuditLog{
		UserID:      1,
		UserEmail:   "stan@test.com",
		UserRole:    "ADMIN",
		Action:      "UPDATE_BOOKING",
		Entity:      "BOOKING",
		EntityID:    "123",
		BeforeState: before,
		AfterState:  after,
	}

	if err := service.LogEvent(log); err != nil {
		t.Fatal(err)
	}

	// Fetch logs from repo
	filter := audit_logs.AuditFilter{UserID: &log.UserID}
	logs, total, err := service.GetAuditLogs(context.Background(),filter, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if total != 1 || logs[0].Action != "UPDATE_BOOKING" {
		t.Fatalf("expected 1 log with action UPDATE_BOOKING, got %v", logs)
	}
	t.Logf("Integration test passed: %+v", logs[0])
}