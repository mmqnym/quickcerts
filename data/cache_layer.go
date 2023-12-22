package data

import (
	"context"
	"errors"
	"strconv"
	"time"

	cfg "github.com/mmq88/quickcerts/configs"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client = nil

// Connect to the redis database.
func ConnectRDB() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.CACHE_CONFIG.HOST + ":" + strconv.Itoa(cfg.CACHE_CONFIG.PORT),
		Password: cfg.CACHE_CONFIG.PASSWORD,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		rdb = nil
		return errors.New("failed to access the redis database")
	}

	return nil
}

// Disconnect from the redis database.
func DisconnectRDB() error {
	if rdb == nil {
		return errors.New("currently not connecting the redis database")
	}

	err := rdb.Close()
	rdb = nil
	return err
}

// Set the key cache corresponding to the device.
func SetDeviceKeyCache(key string, value interface{}) error {
	if rdb == nil {
		return errors.New("currently not connecting the redis database")
	}

	err := rdb.Set(ctx, key, value, time.Hour*24*7).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get the key cache corresponding to the device. if exists, or return "".
func GetDeviceKeyCache(deviceInfoBase string) (string, error) {
	if rdb == nil {
		return "", errors.New("currently not connecting the redis database")
	}

	key, err := rdb.Get(ctx, deviceInfoBase).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", errors.New("the key not exist in the cache")
		}
		return "", err
	}

	return key, nil
}

// Not a secure way to delete cache, only for testing.
func DeleteTestingCache(deviceInfoBase string) error {
	if rdb == nil {
		return errors.New("currently not connecting the redis database")
	}

	err := rdb.Del(ctx, deviceInfoBase).Err()
	if err != nil {
		return err
	}

	return nil
}
