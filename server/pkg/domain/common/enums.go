package common

type Role string

const (
	RoleUser       Role = "USER"
	RoleAdmin      Role = "ADMIN"
	RoleSuperAdmin Role = "SUPER_ADMIN"
)

type BookingStatus string

const (
	BookingPending    BookingStatus = "PENDING"
	BookingConfirmed  BookingStatus = "CONFIRMED"
	BookingCancelled  BookingStatus = "CANCELLED"
	BookingCheckedIn  BookingStatus = "CHECKED_IN"
	BookingCheckedOut BookingStatus = "CHECKED_OUT"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentRefunded  PaymentStatus = "REFUNDED"
)

type RoomStatus string

const (
	RoomAvailable   RoomStatus = "AVAILABLE"
	RoomOccupied    RoomStatus = "OCCUPIED"
	RoomMaintenance RoomStatus = "MAINTENANCE"
	RoomReserved    RoomStatus = "RESERVED"
)

type User string

const (
	UserActive   User = "ACTIVE"
	UserInactive User = "INACTIVE"
)

type ReviewStatus string

const (
	ReviewPending ReviewStatus = "PENDING"

	ReviewApproved ReviewStatus = "APPROVED"

	ReviewRejected ReviewStatus = "REJECTED"

	ReviewFlagged ReviewStatus = "FLAGGED"

	ReviewArchived ReviewStatus = "ARCHIVED"
)
