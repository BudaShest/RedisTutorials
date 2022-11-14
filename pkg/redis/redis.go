package redis

import (
	"context"
	"errors"
	redis "github.com/go-redis/redis/v9"
	"time"
)

type Redis struct {
	connection *redis.Client
}

// New - новый экземпляр структуры Redis
func New(addr, password string, db int) *Redis {
	connection := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{
		connection: connection,
	}
}

// Set - установка значения
func (r *Redis) Set(key string, value any, ttl int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	timeToLive := time.Duration(ttl) * time.Second
	return r.connection.Set(ctx, key, value, timeToLive)
}

// Get - получение значения
func (r *Redis) Get(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.Get(ctx, key)
}

// MGet - Получить множество значений для множества ключей
func (r *Redis) MGet(keys ...string) *redis.SliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.MGet(ctx, keys...)
}

// MSet - Установить множество значений для множества ключей
func (r *Redis) MSet(keys []string, values []string) (*redis.StatusCmd, error) {
	if len(keys) != len(values) {
		return nil, errors.New("slice of keys is not fits to slice of values")
	}

	var keysWithValues map[string]any = make(map[string]any, len(keys))

	for idx, key := range keys {
		keysWithValues[key] = values[idx]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.MSet(ctx, keysWithValues), nil
}

// Exec - для выполнения кастомных комман
func (r *Redis) Exec(command, key string, value any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var (
		result any
		err    error
	)

	if value != nil {
		result, err = r.connection.Do(ctx, command, key, value).Result()
	} else {
		result, err = r.connection.Do(ctx, command, key).Result()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}

func (r *Redis) Lpush(key string, value any) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.LPush(ctx, key, value)
}

func (r *Redis) Rpush(key string, value any) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.RPush(ctx, key, value)
}

func (r *Redis) Lpop(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.LPop(ctx, key)
}

func (r *Redis) Rpop(key string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.RPop(ctx, key)
}

func (r *Redis) Lrange(key string, from, until int) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.LRange(ctx, key, int64(from), int64(until))
}

func (r *Redis) Ltrim(key string, from, until int) *redis.StatusCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.LTrim(ctx, key, int64(from), int64(until))
}

// todo comments in this file
func (r *Redis) BRPop(timeout int, keys []string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.BRPop(ctx, time.Duration(timeout)*time.Second, keys...)
}

func (r *Redis) BLPop(timeout int, keys []string) *redis.StringSliceCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.BLPop(ctx, time.Duration(timeout)*time.Second, keys...)
}

// Hashes
func (r *Redis) HGet(key, field string) *redis.StringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.HGet(ctx, key, field)
}

func (r *Redis) HSet(key string, values ...any) *redis.IntCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.HSet(ctx, key, values)
}

func (r *Redis) HGetAll(key string) *redis.MapStringStringCmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.connection.HGetAll(ctx, key)
}
