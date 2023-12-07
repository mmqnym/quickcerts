package data

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client = nil

// Connect to the redis database.
func ConnectRDB() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "qcs-cache:6379",
		Password: "",
		DB:       0,
	})
}

// Disconnect from the redis database.
func DisconnectRDB() {
	if rdb == nil {
		return
	}

	rdb.Close()
}

// Set the key cache corresponding to the device.
func SetDeviceKeyCache(key string, value interface{}) error {
	if rdb == nil {
		return errors.New("currently not connecting the redis database")
	}

	err := rdb.Set(ctx, key, value, time.Hour * 24 * 7).Err()
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
		return "", err
	}

	return key, nil
}