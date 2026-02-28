package room

import "time"

type RoomStatus string

const (
	Available   RoomStatus = "AVAILABLE"
	Booked      RoomStatus = "BOOKED"
	Maintenance RoomStatus = "MAINTENANCE"
	Reserved    RoomStatus = "RESERVED"
)

type Room struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      RoomStatus `json:"status"`
	Price       int64    `json:"price"`
	Amenities   []string   `json:"amenities,omitempty"`
	RoomNumber  string     `json:"roomNumber"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
