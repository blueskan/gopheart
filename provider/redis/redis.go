package redis

import (
	"regexp"
	"time"

	"github.com/blueskan/gopheart/provider"
	"github.com/go-redis/redis"
)

type redisProvider struct {
	name          string
	redisHost     string
	redisUsername string
	redisPassword string
	redisPort     string
	timeout       time.Duration
	interval      time.Duration
	downThreshold int
	upThreshold   int
}

func NewRedisProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int,
) provider.Provider {
	regex := regexp.MustCompile(`redis:\/\/(.*?)@(.*?):([0-9]*)`)
	res := regex.FindStringSubmatch(connectionString)

	var host, port, username, password string

	if len(res) <= 0 {
		regex = regexp.MustCompile(`redis:\/\/(.*?):([0-9]*)`)
		res = regex.FindStringSubmatch(connectionString)

		host = res[1]
		port = res[2]
	} else {
		username = res[1]
		password = res[2]
		host = res[3]
		port = res[4]
	}

	return &redisProvider{
		name:          name,
		redisHost:     host,
		redisPort:     port,
		redisUsername: username,
		redisPassword: password,
		timeout:       timeout,
		interval:      interval,
		downThreshold: downThreshold,
		upThreshold:   upThreshold,
	}
}

func (rp redisProvider) GetName() string {
	return rp.name
}

func (rp redisProvider) GetInterval() time.Duration {
	return rp.interval
}

func (rp redisProvider) GetDownThreshold() int {
	return rp.downThreshold
}

func (rp redisProvider) GetUpThreshold() int {
	return rp.upThreshold
}

func (rp redisProvider) Heartbeat() bool {
	client := redis.NewClient(&redis.Options{
		Addr:        rp.redisHost + ":" + rp.redisPort,
		Password:    rp.redisPassword,
		DB:          0,
		DialTimeout: rp.timeout,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return false
	}

	if pong != "PONG" {
		return false
	}

	return true
}
