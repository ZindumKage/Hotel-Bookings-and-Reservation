package review

import (
	"errors"
	"time"
)



type Review struct{
	ID    uint
	UserID uint
	
	RoomID  uint
	Comment string
	Rating int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

}

func (r *Review) Validate()error{
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5 stars")
	}
	if r.RoomID == 0 {
		return errors.New("room id required")
	}
	if r.UserID == 0{
		return errors.New("user id required")
	}
	return nil
}