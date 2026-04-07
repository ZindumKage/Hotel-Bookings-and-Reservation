package application

import (
	"context"
	"log"
	"time"

	bookingservice "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/booking"
)

type ExpiryWorker struct {
	service  *bookingservice.Service
	interval time.Duration
}

func NewExpiryWorker(service *bookingservice.Service, interval time.Duration) *ExpiryWorker {
	return &ExpiryWorker{
		service:  service,
		interval: interval,
	}
}

func (w *ExpiryWorker) Start(ctx context.Context) {

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("Booking expiry worker started")

	for {
		select {

		case <-ctx.Done():
			log.Println("Booking expiry worker stopped")
			return

		case <-ticker.C:

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Recovered in expiry worker:", r)
					}
				}()

				if err := w.service.ExpireBookings(ctx); err != nil {
					log.Println("booking expiry worker error:", err)
				}
			}()
		}
	}
}