package review

import (
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/review"
	"gorm.io/gorm"
)


type ReviewRepository struct{
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return  &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *domain.Review) error {
model := toModel(review)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	// Update domain object with DB values
	review.ID = model.ID
	review.CreatedAt = model.CreatedAt
	review.UpdatedAt = model.UpdatedAt

	return nil
}









