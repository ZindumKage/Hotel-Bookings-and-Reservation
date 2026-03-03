package booking

import (
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
)




type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(bookingModel *domain.Booking) error {
	return r.db.Create(toModel(bookingModel)).Error
}

func (r *BookingRepository) Update(bookingModel *domain.Booking) error {
	return r.db.Save(toModel(bookingModel)).Error
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

func (r *BookingRepository) List(page, limit int) ([]domain.Booking, int64, error) {
	var modelList []*models.Booking
	var total int64

	if err := r.db.Model(&models.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Offset((page-1)*limit).
		Limit(limit).
		Find(&modelList).Error; err != nil {
		return nil, 0, err
	}
var result []domain.Booking
for i := range modelList {
	result = append(result, *toDomain(modelList[i]))
}
return result, total, nil
}

func (r *BookingRepository) FindAll(status *domain.PaymentStatus, page, limit int) ([]domain.Booking, int64, error) {
    var modelList []models.Booking
    var total int64
    
    db := r.db.Model(&models.Booking{})

    // Apply filter if status is provided
    if status != nil {
        db = db.Where("payment_status = ?", string(*status))
    }

    // Get total count for pagination
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination and fetch
    offset := (page - 1) * limit
    err := db.Offset(offset).Limit(limit).Find(&modelList).Error
    if err != nil {
        return nil, 0, err
    }

    var result []domain.Booking
    for _, m := range modelList {
        result = append(result, *toDomain(&m))
    }
    
    return result, total, nil
}

func (r *BookingRepository) FindByUser(userID uint) ([]domain.Booking, error) {
	var modelList []models.Booking
	err := r.db.Where("user_id = ?", userID).Find(&modelList).Error
	if err != nil {
		return nil, err
	}
	
var result []domain.Booking
for i := range modelList {
	result = append(result, *toDomain(&modelList[i]))	
}
return result, nil
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




