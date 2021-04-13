package commons

import (
	"github.com/go-redis/redis/v7"
)

//GetGoRedisClient 获取redis client
func GetGoRedisClient(opt *redis.Options) *redis.Client {
	client := redis.NewClient(opt)
	return client
}

func GetGoRedisConn(opt *redis.Options) *redis.Conn {
	client := redis.NewClient(opt)
	return client.Conn()
}

//redis server联通性校验
func CheckRedisClientConnect(r *redis.Client) error {
	_, err := r.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func CheckRedisClusterClientConnect(r *redis.ClusterClient) error {
	_, err := r.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
