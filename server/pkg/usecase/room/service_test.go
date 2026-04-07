package room_test

import (
	"errors"
	"fmt"
	"testing"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/room"
)

//////////////////////////////////////////////////////
//// MOCK REPOSITORY
//////////////////////////////////////////////////////

type mockRepo struct {
	rooms  map[uint]*domain.Room
	nextID uint
}

func (m *mockRepo) Create(r *domain.Room) error {
	m.nextID++
	r.ID = m.nextID
	m.rooms[r.ID] = r
	return nil
}

func (m *mockRepo) Update(r *domain.Room) error {
	if _, ok := m.rooms[r.ID]; !ok {
		return errors.New("room not found")
	}
	m.rooms[r.ID] = r
	return nil
}

func (m *mockRepo) FindByID(id uint) (*domain.Room, error) {
	r, ok := m.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return r, nil
}

func (m *mockRepo) GetByID(id uint) (*domain.Room, error) {
	return m.FindByID(id)
}

func (m *mockRepo) Delete(id uint) error {
	if _, ok := m.rooms[id]; !ok {
		return errors.New("room not found")
	}
	delete(m.rooms, id)
	return nil
}

func (m *mockRepo) UpdateStatus(id uint, status domain.RoomStatus) error {
	r, ok := m.rooms[id]
	if !ok {
		return errors.New("room not found")
	}
	r.Status = status
	return nil
}

func (m *mockRepo) GetStatusByID(id uint) (domain.RoomStatus, error) {
	r, ok := m.rooms[id]
	if !ok {
		return "", errors.New("room not found")
	}
	return r.Status, nil
}

func (m *mockRepo) FindAll(status *domain.RoomStatus, page int, limit int) ([]domain.Room, int64, error) {
	var filtered []domain.Room

	for _, r := range m.rooms {
		if status == nil || r.Status == *status {
			filtered = append(filtered, *r)
		}
	}

	total := int64(len(filtered))

	start := (page - 1) * limit
	if start >= len(filtered) {
		return []domain.Room{}, total, nil
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (m *mockRepo) List(page, limit int) ([]domain.Room, int64, error) {
	return m.FindAll(nil, page, limit)
}

//////////////////////////////////////////////////////
//// TEST HELPERS
//////////////////////////////////////////////////////

func newTestService() (*mockRepo, *usecase.Service) {
	repo := &mockRepo{
		rooms:  make(map[uint]*domain.Room),
		nextID: 0,
	}
	return repo, usecase.NewService(repo)
}

func newRoom(name string, price int64, number string, status domain.RoomStatus) *domain.Room {
	return &domain.Room{
		Name:       name,
		Price:      price,
		RoomNumber: number,
		Status:     status,
	}
}

func seedRooms(t *testing.T, svc *usecase.Service, count int, statusFn func(i int) domain.RoomStatus) {
	t.Helper()

	for i := 1; i <= count; i++ {
		room := newRoom(
			fmt.Sprintf("Room %d", i),
			int64(100*i),
			fmt.Sprintf("%d", 100+i),
			statusFn(i),
		)

		if err := svc.CreateRoom(room); err != nil {
			t.Fatalf("seed failed: %v", err)
		}
	}
}

func assertEqual[T comparable](t *testing.T, got, want T, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v want %v", msg, got, want)
	}
}

//////////////////////////////////////////////////////
//// TESTS
//////////////////////////////////////////////////////

func TestCreateRoom(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Deluxe", 200, "101", domain.Available)

	if err := svc.CreateRoom(room); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, room.ID, uint(1), "room id")
}

func TestGetRoomByID(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 300, "102", domain.Available)
	svc.CreateRoom(room)

	got, err := svc.GetRoomByID(room.ID)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, got.Name, room.Name, "room name")
}

func TestUpdateRoom(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 300, "102", domain.Available)
	svc.CreateRoom(room)

	room.Name = "Updated Suite"

	updated, err := svc.UpdateRoom(room) 
	if err != nil {
		t.Fatal(err)
	}


	assertEqual(t, updated.Name, "Updated Suite", "updated name")
}

func TestDeleteRoom(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 300, "102", domain.Available)
	svc.CreateRoom(room)

	if err := svc.DeleteRoom(room.ID); err != nil {
		t.Fatal(err)
	}

	_, err := svc.GetRoomByID(room.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

//////////////////////////////////////////////////////
//// STATUS TEST
//////////////////////////////////////////////////////

func TestUpdateRoomStatus(t *testing.T) {
	tests := []struct {
		name string
		from domain.RoomStatus
		to   domain.RoomStatus
	}{
		{"available to booked", domain.Available, domain.Booked},
		{"booked to available", domain.Booked, domain.Available},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, svc := newTestService()

			room := newRoom("Suite", 200, "101", tt.from)
			svc.CreateRoom(room)

			updated, err := svc.UpdateRoomStatus(room.ID, tt.to)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, updated.Status, tt.to, "room status")
		})
	}
}

//////////////////////////////////////////////////////
//// PAGINATION TESTS
//////////////////////////////////////////////////////

func TestListRoomsPagination(t *testing.T) {
	tests := []struct {
		name      string
		totalSeed int
		page      int
		limit     int
		wantCount int
		wantTotal int64
	}{
		{"first page", 5, 1, 2, 2, 5},
		{"middle page", 5, 2, 2, 2, 5},
		{"last partial page", 5, 3, 2, 1, 5},
		{"page overflow", 5, 10, 2, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, svc := newTestService()

			seedRooms(t, svc, tt.totalSeed, func(i int) domain.RoomStatus {
				return domain.Available
			})

			rooms, total, err := svc.GetRooms(nil, tt.page, tt.limit)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, int(total), int(tt.wantTotal), "total count")
			assertEqual(t, len(rooms), tt.wantCount, "page size")
		})
	}
}

//////////////////////////////////////////////////////
//// FILTER TEST
//////////////////////////////////////////////////////

func TestListRoomsWithStatusFilter(t *testing.T) {
	_, svc := newTestService()

	seedRooms(t, svc, 5, func(i int) domain.RoomStatus {
		if i%2 == 0 {
			return domain.Booked
		}
		return domain.Available
	})

	filter := domain.Booked

	rooms, total, err := svc.GetRooms(&filter, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, int(total), 2, "filtered total")
	assertEqual(t, len(rooms), 2, "filtered page size")
}

//////////////////////////////////////////////////////
//// EXTRA EDGE CASE TESTS (SENIOR LEVEL)
//////////////////////////////////////////////////////

// ---------- CREATE VALIDATION ----------

func TestCreateRoom_InvalidInput(t *testing.T) {
	_, svc := newTestService()

	tests := []struct {
		name string
		room *domain.Room
	}{
		{"empty name", newRoom("", 100, "101", domain.Available)},
		{"zero price", newRoom("Room", 0, "101", domain.Available)},
		{"negative price", newRoom("Room", -10, "101", domain.Available)},
		{"empty number", newRoom("Room", 100, "", domain.Available)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreateRoom(tt.room)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

// ---------- CREATE DEFAULT STATUS ----------

func TestCreateRoom_DefaultStatus(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Standard", 100, "101", "")
	err := svc.CreateRoom(room)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, room.Status, domain.Available, "default status")
}

// ---------- UPDATE NOT FOUND ----------

func TestUpdateRoom_NotFound(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Ghost", 100, "999", domain.Available)
	room.ID = 99

	_, err := svc.UpdateRoom(room)
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---------- DELETE NOT FOUND ----------

func TestDeleteRoom_NotFound(t *testing.T) {
	_, svc := newTestService()

	err := svc.DeleteRoom(999)
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---------- DOUBLE DELETE (IDEMPOTENCY SIGNAL) ----------

func TestDeleteRoom_DoubleDelete(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 200, "101", domain.Available)
	svc.CreateRoom(room)

	err1 := svc.DeleteRoom(room.ID)
	err2 := svc.DeleteRoom(room.ID)

	if err1 != nil {
		t.Fatal(err1)
	}
	if err2 == nil {
		t.Fatal("expected error on second delete")
	}
}

// ---------- STATUS UPDATE NOT FOUND ----------

func TestUpdateRoomStatus_NotFound(t *testing.T) {
	_, svc := newTestService()

	_, err := svc.UpdateRoomStatus(999, domain.Booked)
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---------- PAGINATION EDGE CASES ----------

func TestListRooms_InvalidPagination(t *testing.T) {
	_, svc := newTestService()

	

	seedRooms(t, svc, 5, func(i int) domain.RoomStatus {
		return domain.Available
	})

	tests := []struct {
		name  string
		page  int
		limit int
	}{
		{"zero page", 0, 2},
		{"zero limit", 1, 0},
		{"negative page", -1, 2},
		{"negative limit", 1, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		_, _, err := svc.GetRooms(nil, tt.page, tt.limit)

if err == nil {
	t.Fatal("expected error for invalid pagination")
}
		})
	}
}

// ---------- FILTER NO MATCH ----------

func TestListRooms_FilterNoMatch(t *testing.T) {
	_, svc := newTestService()

	seedRooms(t, svc, 5, func(i int) domain.RoomStatus {
		return domain.Available
	})

	filter := domain.Booked

	rooms, total, err := svc.GetRooms(&filter, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, int(total), 0, "total")
	assertEqual(t, len(rooms), 0, "rooms length")
}

// ---------- UPDATE PRESERVES ID ----------

func TestUpdateRoom_PreservesID(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 200, "101", domain.Available)
	svc.CreateRoom(room)

	originalID := room.ID

	room.Name = "Updated"
	updated, err := svc.UpdateRoom(room)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, updated.ID, originalID, "id should not change")
}

// ---------- CONCURRENT UPDATE (BASIC SAFETY SIGNAL) ----------

func TestRoom_ConcurrentUpdate(t *testing.T) {
	_, svc := newTestService()

	room := newRoom("Suite", 200, "101", domain.Available)
	svc.CreateRoom(room)

	done := make(chan bool)

	go func() {
		svc.UpdateRoomStatus(room.ID, domain.Booked)
		done <- true
	}()

	go func() {
		svc.UpdateRoomStatus(room.ID, domain.Available)
		done <- true
	}()

	<-done
	<-done

	// no panic = pass (basic concurrency signal)
}