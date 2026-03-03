package booking

import (
	"log"
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
)


type ExpiryWorker struct {
	repo domain.Repository
}

func NewExpiryWorker(repo domain.Repository) *ExpiryWorker {
	return &ExpiryWorker{repo: repo}
}

func (w *ExpiryWorker) Start(interval time.Duration){
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			w.expireBookings()
		}
	}()
}

func (w *ExpiryWorker) expireBookings() {
	bookings,_,err := w.repo.FindAll(nil, 1, 1000)
	if err != nil {
		log.Println("expiry worker error:", err)
		return
	}
	now := time.Now()
	for _, b := range bookings{
		if b.PaymentStatus == domain.PaymentStatusPending && now.After(b.ExpiresAt) {
			b.Cancel()
			_= w.repo.Update(&b)
		}	
	}}