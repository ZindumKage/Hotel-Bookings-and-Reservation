package graph

import (
	"strconv"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph/model"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
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