package repositories

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const deliveryTrackListKey = "items.delivery_track"

type DeliveryTrack struct {
	Mention       string     `json:"mention"`
	CarrierID     string     `json:"carrier_id"`
	TrackID       string     `json:"track_id"`
	ItemName      string     `json:"item_name,omitempty"`
	LastTimestamp *time.Time `json:"last_timestamp"`
}

type DeliveryTrackRepository interface {
	ListAllShouldExecute(ctx context.Context) ([]DeliveryTrack, error)
	Append(ctx context.Context, track *DeliveryTrack, runAt *time.Time) error
}

type RedisDeliveryTrackRepository struct {
	BaseRedisRepository
}

func NewRedisDeliveryTrackRepository(rdb *redis.Client) (DeliveryTrackRepository, error) {
	repo := &RedisDeliveryTrackRepository{}
	repo.rdb = rdb
	return repo, nil
}

func (repo *RedisDeliveryTrackRepository) ListAllShouldExecute(ctx context.Context) ([]DeliveryTrack, error) {
	// Get all tracks that reached waiting time.
	timeScore := strconv.FormatInt(time.Now().Unix(), 10)
	out, err := repo.rdb.ZRangeByScore(ctx, deliveryTrackListKey, &redis.ZRangeBy{Min: "0", Max: timeScore}).Result()
	if err != nil {
		return nil, err
	} else if len(out) == 0 {
		return []DeliveryTrack{}, nil
	}

	// Remove all reached items.
	err = repo.rdb.ZRemRangeByScore(ctx, deliveryTrackListKey, "0", timeScore).Err()
	if err != nil {
		return nil, err
	}

	// Parse all items and return
	result := make([]DeliveryTrack, len(out))
	for idx, rowJson := range out {
		var row DeliveryTrack
		_ = json.Unmarshal([]byte(rowJson), &row)
		result[idx] = row
	}

	return result, nil
}

func (repo *RedisDeliveryTrackRepository) Append(ctx context.Context, track *DeliveryTrack, runAt *time.Time) error {
	row, _ := json.Marshal(track)
	return repo.rdb.ZAdd(ctx, deliveryTrackListKey, &redis.Z{Member: string(row), Score: float64(runAt.Unix())}).Err()
}
