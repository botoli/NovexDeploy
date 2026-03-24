package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"` // "deploy", "build", etc.
	Payload       json.RawMessage   `json:"payload"`
	CreatedAt     time.Time         `json:"created_at"`
	Status        string            `json:"status"`
}

type Queue interface {
	Enqueue(ctx context.Context, job *Job) error
	Dequeue(ctx context.Context) (*Job, error)
}

type RedisQueue struct {
	client *redis.Client
	queueKey string
}

func NewRedisQueue(addr string, password string, db int, queueKey string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,  // use default DB
	})

	return &RedisQueue{
		client:   rdb,
		queueKey: queueKey,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.client.LPush(ctx, q.queueKey, data).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (*Job, error) {
	// BRPOP blocks until an item is available
	result, err := q.client.BRPop(ctx, 0, q.queueKey).Result()
	if err != nil {
		return nil, err
	}
	
	// result[0] is the key, result[1] is the value
	var job Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}
