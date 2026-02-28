package room



type Repository interface {
	Create(room *Room) error
	GetByID(id uint) (*Room, error)
	Update(room *Room) error
	Delete(id uint) error
	List(page, limit int) ([]Room, int64, error)
	FindByID(id uint) (*Room, error)
	FindAll(status *RoomStatus, page, limit int) ([]Room, int64, error)
	GetStatusByID(id uint) (RoomStatus, error)
	UpdateStatus(id uint, status RoomStatus) error
}

