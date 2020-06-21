package redis

import (
	"fmt"
	"strconv"

	"github.com/MrWormHole/url-shortener/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repository := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repository.client = client

	return repository, nil
}

func (r *redisRepository) generateKey(hash string) string {
	return fmt.Sprintf("redirect: %s", hash)
}

func (r *redisRepository) Find(hash string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(hash)

	data, err := r.client.HGetAll(r.client.Context(), key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect.Hash = data["hash"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Hash)
	data := map[string]interface{}{
		"hash":       redirect.Hash,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(r.client.Context(), key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
