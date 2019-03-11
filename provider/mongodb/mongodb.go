package mongodb

import (
	"context"
	"time"

	"github.com/blueskan/gopheart/provider"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDbProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int
	upThreshold      int
}

func NewMongoDbProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int,
) provider.Provider {
	return &mongoDbProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (mdp mongoDbProvider) GetName() string {
	return mdp.name
}

func (mdp mongoDbProvider) GetInterval() time.Duration {
	return mdp.interval
}

func (mdp mongoDbProvider) GetDownThreshold() int {
	return mdp.downThreshold
}

func (mdp mongoDbProvider) GetUpThreshold() int {
	return mdp.upThreshold
}

func (mdp mongoDbProvider) Heartbeat() bool {
	ctx, _ := context.WithTimeout(context.Background(), mdp.timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mdp.connectionString))
	if err != nil {
		return false
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}
