package room

import (
	"sync"

	
)



type Repository interface {
	Create(room *Room) error
	Update(room *Room) error
	Delete(id uint) error
	FindByID(id uint) (*Room, error)


	List(page, limit int) ([]Room, int64, error)
	FindAll(status *RoomStatus, page, limit int) ([]Room, int64, error)

	LockRoom(roomID uint) (*sync.Mutex, error)
	WithTransaction(fn func(repo Repository) error) error
	GetStatusByID(id uint) (RoomStatus, error)
	UpdateStatus(id uint, status RoomStatus) error
}

