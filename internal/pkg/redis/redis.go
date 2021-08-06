package redis

import "github.com/go-redis/redis"

var client *redis.Client

func InitRedis(url string, pwd string) error {
	client = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pwd,
		DB:       0,
	})
	return client.Ping().Err()
}

func GetClient() *redis.Client {
	return client
}
