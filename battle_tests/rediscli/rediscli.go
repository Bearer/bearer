package rediscli

import (
	"errors"

	"github.com/bearer/bearer/battle_tests/build"
	"github.com/bearer/bearer/battle_tests/config"
	"github.com/go-redis/redis"
)

var cli *redis.Client

func Setup() {
	cli = redis.NewClient(&redis.Options{
		Addr:     config.Runtime.Redis.Address,
		Password: config.Runtime.Redis.Password,
		DB:       0,
	})
}

func WorkerOnline() (int64, error) {
	key := build.Version + "_workers_online_" + build.Language + build.Attempt
	cmd := cli.Incr(key)
	return cmd.Result()
}

func WorkerOffline() (int64, error) {
	key := build.Version + "_workers_online_" + build.Language + build.Attempt
	cmd := cli.Decr(key)
	return cmd.Result()
}

func SetDocument(documentID string) error {
	key := build.Version + "_sheet_document_" + build.Attempt
	cmd := cli.Set(key, documentID, 0)
	result, err := cmd.Result()
	if result != "OK" {
		return errors.New("failed to set document")
	}
	return err
}

func GetDocument() (string, error) {
	key := build.Version + "_sheet_document_" + build.Attempt
	cmd := cli.Get(key)
	return cmd.Result()
}

func Init() error {
	if !config.Runtime.Redis.Init {
		return nil
	}

	key := build.Version + "_work_assigned_" + build.Language + build.Attempt
	cmd := cli.Set(key, 0, 0)
	_, err := cmd.Result()

	if err != nil {
		return err
	}

	key = build.Version + "_workers_online_" + build.Language + build.Attempt
	cmd = cli.Set(key, 0, 0)
	_, err = cmd.Result()

	if err != nil {
		return err
	}

	return nil
}

func PickUpWork() (int, error) {
	key := build.Version + "_work_assigned_" + build.Language + build.Attempt
	cmd := cli.Incr(key)
	counter, err := cmd.Result()
	return int(counter), err
}
