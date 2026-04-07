package room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/room"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/redis"
	goredis "github.com/redis/go-redis/v9"
)

type Service struct {
	repo         room.Repository
	redisLocker  *redis.RedisLocker
	redisClient  *goredis.Client
}

func NewService(
	repo room.Repository,
	redisLocker *redis.RedisLocker,
	client *goredis.Client,
) *Service {
	return &Service{
		repo:        repo,
		redisLocker: redisLocker,
		redisClient: client,
	}
}

/* =========================
   CACHE KEYS
========================= */

func (s *Service) roomKey(id uint) string {
	return fmt.Sprintf("room:%d", id)
}

func (s *Service) roomsKey(
	version int64,
	status *room.RoomStatus,
	page, limit int,
) string {

	statusStr := "all"
	if status != nil {
		statusStr = string(*status)
	}

	return fmt.Sprintf("rooms:v%d:%s:%d:%d", version, statusStr, page, limit)
}

/* =========================
   CREATE
========================= */

func (s *Service) CreateRoom(ctx context.Context, r *room.Room) error {
	if r.Name == "" || r.Price <= 0 || r.RoomNumber == "" {
		return errors.New("invalid room data")
	}

	if r.Status == "" {
		r.Status = room.Available
	}

	if err := s.repo.Create(r); err != nil {
		return err
	}

	s.invalidateRoomsCache(ctx)
	return nil
}

/* =========================
   UPDATE
========================= */

func (s *Service) UpdateRoom(ctx context.Context, r *room.Room) (*room.Room, error) {
	existing, err := s.repo.FindByID(r.ID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	existing.Name = r.Name
	existing.Description = r.Description
	existing.Status = r.Status
	existing.Price = r.Price
	existing.Amenities = r.Amenities
	existing.RoomNumber = r.RoomNumber
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	s.invalidateRoomCache(ctx, r.ID)
	s.invalidateRoomsCache(ctx)

	return existing, nil
}

/* =========================
   GET ROOM BY ID (CACHED)
========================= */

func (s *Service) GetRoomByID(ctx context.Context, id uint) (*room.Room, error) {

	key := s.roomKey(id)

	val, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var r room.Room
		if err := json.Unmarshal([]byte(val), &r); err == nil {
			return &r, nil
		}
	}

	if err != nil && err != goredis.Nil {
		return nil, err
	}

	r, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(r)
	if err == nil {
		ttl := 10*time.Minute + time.Duration(rand.Intn(60))*time.Second
		_ = s.redisClient.Set(ctx, key, data, ttl).Err()
	}

	return r, nil
}

/* =========================
   LIST ROOMS (CACHED)
========================= */
func (s *Service) GetRooms(
	ctx context.Context,
	status *room.RoomStatus,
	page, limit int,
) ([]*room.Room, int64, error) {

	if page <= 0 || limit <= 0 {
		return nil, 0, errors.New("invalid pagination parameters")
	}

	version, err := s.getRoomsVersion(ctx)
	if err != nil {
		return nil, 0, err
	}

	key := s.roomsKey(version, status, page, limit)

	val, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var cached struct {
			Rooms []*room.Room
			Total int64
		}
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached.Rooms, cached.Total, nil
		}
	}

	if err != nil && err != goredis.Nil {
		return nil, 0, err
	}

	rooms, total, err := s.repo.FindAll(status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*room.Room, len(rooms))
	for i := range rooms {
		result[i] = &rooms[i]
	}

	cacheData := struct {
		Rooms []*room.Room
		Total int64
	}{
		Rooms: result,
		Total: total,
	}

	data, err := json.Marshal(cacheData)
	if err == nil {
		ttl := 5*time.Minute + time.Duration(rand.Intn(60))*time.Second
		_ = s.redisClient.Set(ctx, key, data, ttl).Err()
	}

	return result, total, nil
}

/* =========================
   DELETE
========================= */

func (s *Service) DeleteRoom(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateRoomCache(ctx, id)
	s.invalidateRoomsCache(ctx)

	return nil
}

/* =========================
   STATUS UPDATE (LOCKED)
========================= */

func (s *Service) UpdateRoomStatus(
	ctx context.Context,
	roomID uint,
	status string,
) (*room.Room, error) {

	if !isValidStatus(status) {
		return nil, errors.New("invalid room status")
	}

	roomStatus := room.RoomStatus(status)

	var updatedRoom *room.Room

	err := s.withRoomLock(ctx, roomID, func() error {
		return s.repo.WithTransaction(func(tx room.Repository) error {

			r, err := tx.FindByID(roomID)
			if err != nil {
				return fmt.Errorf("room with ID %d not found", roomID)
			}

			r.Status = roomStatus
			r.UpdatedAt = time.Now()

			if err := tx.Update(r); err != nil {
				return err
			}

			updatedRoom = r
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	s.invalidateRoomCache(ctx, roomID)
	s.invalidateRoomsCache(ctx)

	return updatedRoom, nil
}

/* =========================
   LOCK HELPER
========================= */

func (s *Service) withRoomLock(
	ctx context.Context,
	roomID uint,
	fn func() error,
) error {

	var unlock func()
	var err error

	for i := 0; i < 3; i++ {
		unlock, err = s.redisLocker.LockResource(ctx, "room", roomID, 3*time.Second)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
	}

	if err != nil {
		return err
	}
	defer unlock()

	return fn()
}

/* =========================
   CACHE INVALIDATION
========================= */

func (s *Service) invalidateRoomCache(ctx context.Context, id uint) {
	_ = s.redisClient.Del(ctx, s.roomKey(id)).Err()
}

func (s *Service) invalidateRoomsCache(ctx context.Context) {
	iter := s.redisClient.Scan(ctx, 0, "rooms:*", 0).Iterator()
	for iter.Next(ctx) {
		_ = s.redisClient.Incr(ctx, s.roomsVersionKey()).Err()
	}
}

/* =========================
   HELPERS
========================= */

func isValidStatus(s string) bool {
	switch room.RoomStatus(s) {
	case room.Available, room.Booked, room.Maintenance:
		return true
	default:
		return false
	}
}
func (s *Service) roomsVersionKey() string {
	return "rooms:version"
}

func (s *Service) getRoomsVersion(ctx context.Context) (int64, error) {
	v, err := s.redisClient.Get(ctx, s.roomsVersionKey()).Int64()

	if err == goredis.Nil {
		// initialize version = 1
		if err := s.redisClient.Set(ctx, s.roomsVersionKey(), 1, 0).Err(); err != nil {
			return 0, err
		}
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return v, nil
}