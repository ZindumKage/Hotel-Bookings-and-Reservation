package room

import (
	"errors"
	"fmt"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
)


type Service struct {
	repo room.Repository
}

func NewService(repo room.Repository) *Service {
	return &Service{repo: repo}
}

// CreateRoom - creates a new room
func (s *Service) CreateRoom(r *room.Room) error{
	if r.Name == "" || r.Price <= 0 || r.RoomNumber == "" {
		return errors.New("Invalid Room Data")
	}
	if r.Status == "" {
		r.Status = room.Available
	}
	return s.repo.Create(r)
}
// UpdateRoom - updates an existing room
func (s *Service) UpdateRoom(r *room.Room) error {
	existing, err := s.repo.FindByID(r.ID)
	if err != nil {
		return errors.New("Room not found")
	}
	existing.Name = r.Name
	existing.Description = r.Description
	existing.Status = r.Status
	existing.Price = r.Price
	existing.Amenities = r.Amenities
	existing.RoomNumber = r.RoomNumber
	return s.repo.Update(existing)
}
// GetRoomByID - retrieves a room by its ID
func (s *Service) GetRoomByID(id uint) (*room.Room, error) {
	return s.repo.FindByID(id)
}
// ListRooms - lists rooms with pagination and optional status filter
 func (s *Service) GetRooms(status *room.RoomStatus, page, limit int) ([]*room.Room, int64, error) {	
	rooms, total, err := s.repo.FindAll(status, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*room.Room, len(rooms))
	for i := range rooms {
		result[i] = &rooms[i]
}
return result, total, nil
}
// DeleteRoom - deletes a room by its ID

func (s *Service) DeleteRoom(id uint) error {
	return s.repo.Delete(id)
}

// UpdateRoomStatus - updates the status of a room
func (s *Service) UpdateRoomStatus(id uint, status room.RoomStatus) (*room.Room, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("Room with ID %d not found")
	}

	room.Status = status
	room.UpdatedAt = time.Now()
	 
	if err := s.repo.Update(room); err != nil {
		return nil, fmt.Errorf("Failed to update room status: %w", err)	
	}
	return room, nil
}