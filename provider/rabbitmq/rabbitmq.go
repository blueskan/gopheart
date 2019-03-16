package rabbitmq

import (
	"github.com/streadway/amqp"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type rabbitmqProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewRabbitmqProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &rabbitmqProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (rmqp rabbitmqProvider) GetName() string {
	return rmqp.name
}

func (rmqp rabbitmqProvider) GetInterval() time.Duration {
	return rmqp.interval
}

func (rmqp rabbitmqProvider) GetDownThreshold() int64 {
	return rmqp.downThreshold
}

func (rmqp rabbitmqProvider) GetUpThreshold() int64 {
	return rmqp.upThreshold
}

func (rmqp rabbitmqProvider) Heartbeat() bool {
	_, err := amqp.Dial(rmqp.connectionString)

	if err != nil {
		return false
	}

	return true
}
