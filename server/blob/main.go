package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	"coolcar/blob/blob"
	"coolcar/blob/cos"
	"coolcar/blob/dao"
	"coolcar/shared/server"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

var (
	addr      = flag.String("addr", ":8083", "address to listen")
	mongoURI  = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri")
	cosAddr   = flag.String("cos_addr", "<URL>", "cos address")
	cosSecID  = flag.String("cos_sec_id", "<SEC_ID>", "cos secret id")
	cosSecKey = flag.String("cos_sec_key", "<SEC_KEY>", "cos secret key")
)

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")

	st, err := cos.NewService(*cosAddr, *cosSecID, *cosSecKey)
	if err != nil {
		logger.Fatal("cannot create cos service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "blob",
		Addr:   *addr,
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			blobpb.RegisterBlobServiceServer(s, &blob.Service{
				Storage: st,
				Mongo:   dao.NewMongo(db),
				Logger:  logger,
			})
		},
	}))
}
