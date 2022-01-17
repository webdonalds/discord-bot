package repositories

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//------------------------------------------------------------------------------

const MolluNotFound = MolluError("mollu not found")

type MolluError string

func (e MolluError) Error() string { return string(e) }

//------------------------------------------------------------------------------

const molluInfoListKey = "items.mollu"

type MolluInfo struct {
	ID            string     `json:"id"`
	NotifySetting bool       `json:"notify_setting"`
	CafeLastVisit *time.Time `json:"cafe_last_visit"`
	IsNotified    bool       `json:"is_notified"`
}

type MolluRepository interface {
	GetByID(ctx context.Context, id string) (MolluInfo, error)
	ListAll(ctx context.Context) ([]MolluInfo, error)
	Save(ctx context.Context, molluInfo MolluInfo) error
	SaveAll(ctx context.Context, molluInfoList []MolluInfo) error
}

type RedisMolluRepository struct {
	BaseRedisRepository
}

func NewRedisMolluRepository(rdb *redis.Client) (MolluRepository, error) {
	repo := &RedisMolluRepository{}
	repo.rdb = rdb
	return repo, nil
}

// TODO
func (repo *RedisMolluRepository) GetByID(ctx context.Context, id string) (MolluInfo, error) {
	return MolluInfo{}, nil
}

func (repo *RedisMolluRepository) ListAll(ctx context.Context) ([]MolluInfo, error) {
	return nil, nil
}

func (repo *RedisMolluRepository) Save(ctx context.Context, molluInfo MolluInfo) error {
	return nil
}

func (repo *RedisMolluRepository) SaveAll(ctx context.Context, molluInfoList []MolluInfo) error {
	return nil
}
