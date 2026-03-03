package review

type Repository interface {
	Create(review *Review)error
	Update(review *Review) error
	Delete(id uint)error
	FindById(id uint) (*Review, error)
	FindByRoomID(roomID uint, page, limit int) ([]Review, int64, error)
	ExistsByUserAndRoom(userID, roomID uint)(bool, error)
	GetAverageRating(roomID uint) (float64, error)
}