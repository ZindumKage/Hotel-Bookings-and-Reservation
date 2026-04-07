package repositories

import (
	"sync"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *domain.Room) error {
	model := ToModel(room)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	room.ID = model.ID
	return nil
}

func (r *RoomRepository) Update(room *domain.Room) error {
	model := ToModel(room)
	return r.db.Save(&model).Error
}

func (r *RoomRepository) FindByID(id uint) (*domain.Room, error) {
	var model models.Room

	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	room := ToDomain(model)
	return &room, nil
}

func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}


func (r *RoomRepository) queryRooms(
	status *domain.RoomStatus,
	page int,
	limit int,
) ([]domain.Room, int64, error) {

	var roomModels []models.Room
	var total int64

	query := r.db.Model(&models.Room{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&roomModels).Error; err != nil {
		return nil, 0, err
	}

	return ToDomains(roomModels), total, nil
}


func (r *RoomRepository) List(page, limit int) ([]domain.Room, int64, error) {
	return r.queryRooms(nil, page, limit)
}

func (r *RoomRepository) FindAll(
	status *domain.RoomStatus,
	page int,
	limit int,
) ([]domain.Room, int64, error) {
	return r.queryRooms(status, page, limit)
}

func (r *RoomRepository) GetStatusByID(id uint) (domain.RoomStatus, error) {
	var model models.Room

	if err := r.db.Select("status").First(&model, id).Error; err != nil {
		return "", err
	}

	return domain.RoomStatus(model.Status), nil
}

func (r *RoomRepository) UpdateStatus(id uint, status domain.RoomStatus) error {
	return r.db.Model(&models.Room{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *RoomRepository) LockRoom(roomID uint) (*sync.Mutex, error) {
	var room models.Room

	err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", roomID).
		First(&room).Error

	return nil, err // no mutex in DB
}

func (r *RoomRepository) WithTransaction(fn func(repo domain.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &RoomRepository{db: tx}
		return fn(txRepo)
	})
}