package repositories

import (
	"strings"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
)

func ToDomain(model models.Room) domain.Room {
	return domain.Room{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Status:      domain.RoomStatus(model.Status),
		RoomNumber:  model.RoomNumber,
		Amenities:   strings.Split(model.Amenities, ","),
	}
}

func ToDomains(models []models.Room) []domain.Room {
	result := make([]domain.Room, len(models))

	for i, m := range models {
		result[i] = ToDomain(m)
	}

	return result
}

func ToModel(domainRoom *domain.Room) models.Room {
	return models.Room{
		Model: gorm.Model{
			ID: domainRoom.ID,
		},
		Name:        domainRoom.Name,
		Description: domainRoom.Description,
		Price:       domainRoom.Price,
		Status:      string(domainRoom.Status),
		RoomNumber:  domainRoom.RoomNumber,
		Amenities:   strings.Join(domainRoom.Amenities, ","),
	}
}