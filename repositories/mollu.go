package repositories

import (
	"context"
	"encoding/json"
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

func (repo *RedisMolluRepository) GetByID(ctx context.Context, id string) (MolluInfo, error) {
	infoList, err := repo.ListAll(ctx)
	if err != nil {
		return MolluInfo{}, err
	}

	for _, info := range infoList {
		if info.ID == id {
			return info, nil
		}
	}
	return MolluInfo{}, MolluNotFound
}

func (repo *RedisMolluRepository) ListAll(ctx context.Context) ([]MolluInfo, error) {
	infoListJson, err := repo.rdb.Get(ctx, molluInfoListKey).Result()
	if err == redis.Nil {
		return []MolluInfo{}, nil
	} else if err != nil {
		return nil, err
	}

	var infoList []MolluInfo
	err = json.Unmarshal([]byte(infoListJson), &infoList)
	if err != nil {
		return nil, err
	}
	return infoList, nil
}

func (repo *RedisMolluRepository) Save(ctx context.Context, molluInfo MolluInfo) error {
	infoList, err := repo.ListAll(ctx)
	if err != nil {
		return err
	}

	changed := false
	for i, info := range infoList {
		if info.ID == molluInfo.ID {
			infoList[i] = molluInfo
			changed = true
			break
		}
	}
	if !changed {
		infoList = append(infoList, molluInfo)
	}
	return repo.SaveAll(ctx, infoList)
}

func (repo *RedisMolluRepository) SaveAll(ctx context.Context, molluInfoList []MolluInfo) error {
	infoListJson, err := json.Marshal(molluInfoList)
	if err != nil {
		return err
	}

	return repo.rdb.Set(ctx, molluInfoListKey, string(infoListJson), -1).Err()
}
