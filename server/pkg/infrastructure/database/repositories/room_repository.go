package repositories

import (
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"gorm.io/gorm"
)




type RoomRepository struct {
	db *gorm.DB
	}	

	func NewRoomRepository(db *gorm.DB) *RoomRepository {
		return &RoomRepository{db: db}
	}
	
	func (r *RoomRepository) Create(roomModel *room.Room) error{
		return r.db.Create(roomModel).Error
	}

	func (r *RoomRepository) Update(roomModel *room.Room) error {
		return r.db.Save(roomModel).Error
	}

	func (r *RoomRepository) FindByID(id uint) (*room.Room, error) {
		var roomModel room.Room
		err := r.db.First(&roomModel, id).Error
		if err != nil {
			return nil, err
		}
		return &roomModel, nil
	}

	func (r *RoomRepository) FindAll(status *room.RoomStatus) ([]*room.Room, error) {
		var rooms []*room.Room
		query := r.db.Model(&room.Room{})
		if status != nil {
			query = query.Where("status = ?", *status)
		}
		err := query.Find(&rooms).Error
		if err != nil {
			return nil, err
		}
		return rooms, nil
	} 

	func (r *RoomRepository) Delete(id uint) error {
		return r.db.Delete(&room.Room{}, id).Error
	}