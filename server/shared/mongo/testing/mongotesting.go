package mongotesting

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	image         = "mongo:latest"
	containerPort = "27017/tcp"
)

var mongoURI string

const defaultMongoURI = "mongodb://localhost:27017"

func RunWithMongoInDocker(m *testing.M) int {
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	resp, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			ExposedPorts: nat.PortSet{
				containerPort: {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				containerPort: []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: "0",
					},
				},
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()

	containerID := resp.ID
	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports["27017/tcp"][0]
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)
	return m.Run()
}
func NewClient(ctx context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo uri not set.Please run RunWithMongoInDocker in TestMain")
	}
	return mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
}
func NewDefaultClient(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(defaultMongoURI))
}
func SetupIndexes(ctx context.Context, database *mongo.Database) error {
	_, err := database.Collection("account").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "open_id",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	_, err = database.Collection("trip").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "trip.accountid",
				Value: 1,
			},
			{
				Key:   "trip.status",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
			"trip.status": 1,
		}),
	})
	if err != nil {
		return err
	}
	_, err = database.Collection("profile").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "accountid",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	return err
}
