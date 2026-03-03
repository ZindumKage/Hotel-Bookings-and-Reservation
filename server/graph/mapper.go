package graph

import (
	"fmt"
	"strconv"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph/model"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
)

// mapToGraphQLModel converts a slice of domain AuditLog to a slice of GraphQL AuditLog
func convertAuditLogs(logs []audit_logs.AuditLog) []*model.AuditLog {
	result := make([]*model.AuditLog, len(logs))
	for i, log := range logs {
		result[i] = &model.AuditLog{
			ID:         strconv.Itoa(int(log.ID)),
			UserID:     strconv.Itoa(int(log.UserID)),
			UserEmail:  &log.UserEmail,
			UserRole:   &log.UserRole,
			Action:     log.Action,
			Entity:     log.Entity,
			EntityID:   &log.EntityID,
			RiskLevel:  (*string)(&log.RiskLevel),
			Suspicious: &log.Suspicious,
			Reason:     &log.Reason,
			CreatedAt:  log.CreatedAt,
		}
	}
	return result
}

//

// func mapRoomsToGraphQL(domainRooms []room.Room) []*model.Room {
// 	gqlRooms := make([]*model.Room, len(domainRooms))
// 	for i, r := range domainRooms {
// 		gqlRooms[i] = &model.Room{
// 			ID:          strconv.Itoa(int(r.ID)),
// 			Name:        r.Name,
// 			Description: &r.Description,
// 			 Price:      r.Price,
// 			Status:      model.RoomStatus(r.Status),
// 			RoomNumber:  r.RoomNumber,
// 			Amenities:  r.Amenities, // assuming this is a []string in GraphQL as well
// 		}
// 	}
// 	return gqlRooms
// }

func MapToGraphQLRoom(r *room.Room) *model.Room {
	return &model.Room{
		ID:          strconv.Itoa(int(r.ID)),
		Name:        r.Name,
		Description: &r.Description,
		Price:       r.Price,
		Status:      model.RoomStatus(r.Status),
		RoomNumber:  r.RoomNumber,
		Amenities:   r.Amenities,
	}
}

func MapToGraphQLBookings(bookings []booking.Booking) []*model.Booking {
	var result []*model.Booking
	for i := range bookings {
		b := bookings[i]
		result = append(result, MapToGraphQLBooking(&b))
	}
	return result
}

func MapToGraphQLBooking(b *booking.Booking) *model.Booking {

	return &model.Booking{
		ID:            fmt.Sprintf("%d", b.ID),
		Reference:     b.Reference,
		User:          fmt.Sprintf("%d", b.UserID),
		RoomID:        fmt.Sprintf("%d", b.RoomID),
		RoomNumber:    b.RoomNumber,
		CheckInDate:   b.CheckInDate,
		CheckOutDate:  b.CheckOutDate,
		NightCount:    int64(b.NightCount),
		UnitPrice:     &b.UnitPrice,
		TaxAmount:     &b.TaxAmount,
		Discount:      &b.Discount,
		TotalPrice:    int64(b.TotalPrice),
		Status:        model.BookingStatus(b.Status),
		PaymentStatus: model.PaymentStatus(b.PaymentStatus),
		ExpiresAt:     b.ExpiresAt,
		CancelledAt:   b.CancelledAt,
	}
}
