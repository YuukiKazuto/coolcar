package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/rental/ai"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profClient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	coolenvpb "coolcar/shared/coolenv"
	"coolcar/shared/server"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	addr              = flag.String("addr", ":8082", "address to listen")
	mongoURI          = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri")
	blobAddr          = flag.String("blob_addr", "localhost:8083", "address for blob service")
	aiAddr            = flag.String("ai_addr", "localhost:18001", "address for ai service")
	carAddr           = flag.String("car_addr", "localhost:8084", "address for car service")
	authPublicKeyFile = flag.String("auth_public_key_file", "shared/auth/public.key", "public key file for auth")
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
	ac, err := grpc.Dial(*aiAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect aiservice", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")
	blobConn, err := grpc.Dial(*blobAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect blob service", zap.Error(err))
	}
	aiClient := &ai.Client{
		AIClient:  coolenvpb.NewAIServiceClient(ac),
		UseRealAI: false,
	}
	profileService := &profile.Service{
		IdentityResolver:  aiClient,
		PhotoGetExpire:    5 * time.Second,
		PhotoUploadExpire: 10 * time.Second,
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		Mongo:             profiledao.NewMongo(db),
		Logger:            logger,
	}
	carConn, err := grpc.Dial(*carAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              *addr,
		AuthPublicKeyFile: *authPublicKeyFile,
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				DistanceCalc: aiClient,
				ProfileManager: &profClient.Manager{
					Fetcher: profileService,
				},
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profileService)
		},
	}))
}
