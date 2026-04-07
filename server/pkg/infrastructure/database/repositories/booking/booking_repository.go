package repositories

import (
	"context"

	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)




type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *domain.Booking) error {
	model := toModel(booking)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	// copy generated fields back
	booking.ID = model.ID
	booking.CreatedAt = model.CreatedAt
	booking.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *BookingRepository) Update(b *domain.Booking) error {

	return r.db.Model(&models.Booking{}).
		Where("id = ?", b.ID).
		Updates(map[string]interface{}{
			"status":         string(b.Status),
			"payment_status": string(b.PaymentStatus),
			"cancelled_at":   b.CancelledAt,
			"expires_at":     b.ExpiresAt,
		}).Error
}

func (r *BookingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

func (r *BookingRepository) FindByID(id uint) (*domain.Booking, error) {
	var m models.Booking
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *BookingRepository) queryBookings(
	ctx context.Context,
	query *gorm.DB,
	page int,
	limit int,
) ([]domain.Booking, int64, error) {

	var modelList []models.Booking
	var total int64

	// Attach context to DB operations
	db := query.WithContext(ctx)

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Fetch paginated records
	if err := db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&modelList).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain objects
	result := make([]domain.Booking, 0, len(modelList))

	for i := range modelList {
		domainBooking := toDomain(&modelList[i])
		if domainBooking != nil {
			result = append(result, *domainBooking)
		}
	}

	return result, total, nil
}
func (r *BookingRepository) List(ctx context.Context, page, limit int) ([]domain.Booking, int64, error) {
	query := r.db.Model(&models.Booking{})
	return r.queryBookings(ctx, query, page, limit)
}
func (r *BookingRepository) FindAll(
	status *domain.PaymentStatus,
	ctx context.Context,
	page, limit int,
) ([]domain.Booking, int64, error) {

	query := r.db.Model(&models.Booking{})

	if status != nil {
		query = query.Where("payment_status = ?", string(*status))
	}

	return r.queryBookings(ctx, query, page, limit)
}
func (r *BookingRepository) FindByUser(
	ctx context.Context,
	userID uint,
	page, limit int,
) ([]domain.Booking, int64, error) {

	query := r.db.
		Model(&models.Booking{}).
		Where("user_id = ?", userID)

	return r.queryBookings(ctx, query, page, limit)
}


func (r *BookingRepository) FindByRoomID(roomID uint) ([]domain.Booking, error) {
	var modelList []models.Booking
	err := r.db.Where("room_id = ?", roomID).Find(&modelList).Error
	if err != nil {
		return nil, err
	}
	
var result []domain.Booking
for i := range modelList {
	result = append(result, *toDomain(&modelList[i]))	
}
return result, nil 
}

func (r *BookingRepository) FindOverlappingBookings(roomID uint, checkIn, checkOut time.Time) ([]domain.Booking, error) {
	var modelList []models.Booking
	err := r.db.Where("room_id = ? AND check_in_date < ? AND check_out_date > ?", roomID, checkOut, checkIn).Find(&modelList).Error
	if err != nil {
		return nil, err
	}
	
var result []domain.Booking
for i := range modelList {
	result = append(result, *toDomain(&modelList[i]))	
}
return result, nil
}

func (r *BookingRepository) UpdatePaymentStatusTx(tx *gorm.DB, id uint, status domain.PaymentStatus) error {
	return tx.Model(&models.Booking{}).Where("id = ?", id).Update("payment_status", string(status)).Error
}

func (r *BookingRepository) UpdatePaymentStatus(id uint, status domain.PaymentStatus) error {
	return r.db.Transaction(func (tx *gorm.DB) error {
		return r.UpdatePaymentStatusTx(tx, id, status)
	})

}

func (r *BookingRepository) WithTransaction(fn func(tx domain.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &BookingRepository{db: tx}
		return fn(txRepo)
	})
}  

func (r *BookingRepository) LockRoom(roomID uint) error {
	
	var room room.Room

	return  r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", roomID).First(&room).Error
}


func (r *BookingRepository) FindExpiredBookings(now time.Time) ([]domain.Booking, error) {

	var modelsList []models.Booking

	err := r.db.
		Where("status = ? AND expires_at < ?", "PENDING", now).
		Find(&modelsList).Error

	if err != nil {
		return nil, err
	}

	result := make([]domain.Booking, 0, len(modelsList))

	for i := range modelsList {
		result = append(result, *toDomain(&modelsList[i]))
	}

	return result, nil
}


