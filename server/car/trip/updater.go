package trip

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/car/mq"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunUpdater(sub mq.Subscriber, client rentalpb.TripServiceClient, logger *zap.Logger) {
	ch, cleanUp, err := sub.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		logger.Fatal("cannot subscribe", zap.Error(err))
	}
	for carEntity := range ch {
		if carEntity.Car.Status == carpb.CarStatus_UNLOCKED && carEntity.Car.TripId != "" && carEntity.Car.Driver.Id != "" {
			_, err := client.UpdateTrip(context.Background(), &rentalpb.UpdateTripRequest{
				Id: carEntity.Car.TripId,
				Current: &rentalpb.Location{
					Latitude:  carEntity.Car.Position.Latitude,
					Longitude: carEntity.Car.Position.Longitude,
				},
			}, grpc.PerRPCCredentials(&impersonation{AccountID: id.AccountID(carEntity.Car.Driver.Id)}))
			if err != nil {
				logger.Error("cannot update trip", zap.String("trip_id", carEntity.Car.TripId), zap.Error(err))
			}
		}
	}
}

type impersonation struct {
	AccountID id.AccountID
}

func (i *impersonation) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		auth.ImpersonateAccountHeader: i.AccountID.String(),
	}, nil
}

func (i *impersonation) RequireTransportSecurity() bool {
	return false
}
