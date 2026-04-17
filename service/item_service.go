package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"go-day3/jobs"
	"go-day3/models"
	"go-day3/repository"
)

const itemsCacheKey = "items:all"

func itemCacheKey(id int) string {
	return "items:" + strconv.Itoa(id)
}

type ItemService interface {
	GetItems() ([]models.Item, error)
	GetItemByID(id int) (models.Item, error)
	CreateItem(name string) (models.Item, error)
	UpdateItem(id int, name string) (models.Item, error)
	DeleteItem(id int) error
}

type DefaultItemService struct {
	repo      repository.ItemRepository
	cache     *redis.Client
	jobQueue  *jobs.Queue
	jobWG     *sync.WaitGroup
	nextJobID atomic.Int64
}

func NewItemService(
	repo repository.ItemRepository,
	cache *redis.Client,
	jobQueue *jobs.Queue,
	jobWG *sync.WaitGroup,
) *DefaultItemService {
	return &DefaultItemService{
		repo:     repo,
		cache:    cache,
		jobQueue: jobQueue,
		jobWG:    jobWG,
	}
}

func (s *DefaultItemService) GetItems() ([]models.Item, error) {
	ctx := context.Background()

	cached, err := s.cache.Get(ctx, itemsCacheKey).Result()
	if err == nil {
		var items []models.Item
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			return items, nil
		}
	}

	items, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(items)
	if err == nil {
		_ = s.cache.Set(ctx, itemsCacheKey, bytes, 5*time.Minute).Err()
	}

	return items, nil
}

func (s *DefaultItemService) GetItemByID(id int) (models.Item, error) {
	ctx := context.Background()
	key := itemCacheKey(id)

	cached, err := s.cache.Get(ctx, key).Result()
	if err == nil {
		var item models.Item
		if err := json.Unmarshal([]byte(cached), &item); err == nil {
			return item, nil
		}
	}

	item, err := s.repo.GetByID(id)
	if err != nil {
		return models.Item{}, err
	}

	bytes, err := json.Marshal(item)
	if err == nil {
		_ = s.cache.Set(ctx, key, bytes, 5*time.Minute).Err()
	}

	return item, nil
}

func (s *DefaultItemService) CreateItem(name string) (models.Item, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Item{}, errors.New("name is required")
	}

	ctx := context.Background()

	item, err := s.repo.Create(name)
	if err != nil {
		return models.Item{}, err
	}

	_ = s.cache.Del(ctx, itemsCacheKey).Err()

	bytes, err := json.Marshal(item)
	if err == nil {
		_ = s.cache.Set(ctx, itemCacheKey(item.ID), bytes, 5*time.Minute).Err()
	}

	s.enqueueJob(item.ID, "reindex-item")

	return item, nil
}

func (s *DefaultItemService) UpdateItem(id int, name string) (models.Item, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Item{}, errors.New("name is required")
	}

	ctx := context.Background()

	item, err := s.repo.Update(id, name)
	if err != nil {
		return models.Item{}, err
	}

	_ = s.cache.Del(ctx, itemsCacheKey).Err()

	bytes, err := json.Marshal(item)
	if err == nil {
		_ = s.cache.Set(ctx, itemCacheKey(item.ID), bytes, 5*time.Minute).Err()
	} else {
		_ = s.cache.Del(ctx, itemCacheKey(item.ID)).Err()
	}

	s.enqueueJob(item.ID, "reindex-item")

	return item, nil
}

func (s *DefaultItemService) DeleteItem(id int) error {
	ctx := context.Background()

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	_ = s.cache.Del(ctx, itemsCacheKey).Err()
	_ = s.cache.Del(ctx, itemCacheKey(id)).Err()

	s.enqueueJob(id, "delete-from-index")

	return nil
}

func (s *DefaultItemService) enqueueJob(itemID int, jobType string) {
	if s.jobQueue == nil || s.jobWG == nil {
		return
	}

	jobID := int(s.nextJobID.Add(1))
	job := jobs.Job{
		ID:     jobID,
		ItemID: itemID,
		Type:   jobType,
	}

	s.jobWG.Add(1)

	select {
	case s.jobQueue.Jobs <- job:
	default:
		s.jobWG.Done()
	}
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
