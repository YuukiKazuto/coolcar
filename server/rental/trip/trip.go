package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"coolcar/shared/mongo/objid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

type Service struct {
	DistanceCalc   DistanceCalc
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

// ProfileManager defines the ACL (Anti Corruption Layer)
// for profile verification logic.
type ProfileManager interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

// CarManager defines the ACL for car management.
type CarManager interface {
	Verify(context.Context, id.CarID, *rentalpb.Location) error
	Unlock(context.Context, id.CarID, id.AccountID, id.TripID, string) error
	Lock(context.Context, id.CarID) error
}

// POIManager resolves POI(Point Of Interest).
type POIManager interface {
	Resolve(ctx context.Context, location *rentalpb.Location) (string, error)
}

type DistanceCalc interface {
	DistanceKm(context.Context, *rentalpb.Location, *rentalpb.Location) (float64, error)
}

var nowFuc = func() int64 {
	return time.Now().Unix()
}

const centsPerSec = 0.7

func (s *Service) CreateTrip(ctx context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	aid, err := auth.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if request.CarId == "" || request.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	// 验证驾驶者身份
	iID, err := s.ProfileManager.Verify(ctx, aid)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	//检查车辆状态
	carID := id.CarID(request.CarId)
	err = s.CarManager.Verify(ctx, carID, request.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	ls := s.calcCurrentStatus(ctx, &rentalpb.LocationStatus{
		Location:     request.Start,
		TimestampSec: nowFuc(),
	}, request.Start)
	tr, err := s.Mongo.CreateTrip(ctx, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		Start:      ls,
		Current:    ls,
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		IdentityId: iID.String(),
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}
	go func() {
		err = s.CarManager.Unlock(context.Background(), carID, aid, objid.ToTripID(tr.ID), request.AvatarUrl)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()
	return &rentalpb.TripEntity{
		Id:   tr.ID.Hex(),
		Trip: tr.Trip,
	}, nil
}

func (s *Service) GetTrip(ctx context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	aid, err := auth.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tr, err := s.Mongo.GetTrip(ctx, id.TripID(request.Id), aid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	return tr.Trip, nil
}

func (s *Service) GetTrips(ctx context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	aid, err := auth.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	trips, err := s.Mongo.GetTrips(ctx, aid, request.Status)
	if err != nil {
		s.Logger.Error("cannot get trips", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	res := &rentalpb.GetTripsResponse{}
	for _, tr := range trips {
		res.Trips = append(res.Trips, &rentalpb.TripEntity{
			Id:   tr.ID.Hex(),
			Trip: tr.Trip,
		})
	}
	return res, nil
}

func (s *Service) UpdateTrip(ctx context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	accountID, err := auth.AccountIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tid := id.TripID(request.Id)
	tr, err := s.Mongo.GetTrip(ctx, id.TripID(request.Id), accountID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}
	if tr.Trip.Status == rentalpb.TripStatus_FINISHED {
		return nil, status.Error(codes.FailedPrecondition, "cannot update a finished trip")
	}
	if tr.Trip.Current == nil {
		s.Logger.Error("trip without current set", zap.String("id", tid.String()))
	}
	cur := tr.Trip.Current.Location
	if request.Current != nil {
		cur = request.Current
	}
	tr.Trip.Current = s.calcCurrentStatus(ctx, tr.Trip.Current, cur)
	if request.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
		s.Logger.Debug(tr.Trip.CarId)
		err := s.CarManager.Lock(ctx, id.CarID(tr.Trip.CarId))
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "cannot lock car: %v", err)
		}
	}
	err = s.Mongo.UpdateTrip(ctx, tid, accountID, tr.UpdatedAt, tr.Trip)
	if err != nil {
		return nil, status.Error(codes.Aborted, "")
	}
	return tr.Trip, nil
}

func (s *Service) calcCurrentStatus(ctx context.Context, last *rentalpb.LocationStatus, cur *rentalpb.Location) *rentalpb.LocationStatus {
	now := nowFuc()
	elapsedSec := float64(now - last.TimestampSec)

	dist, err := s.DistanceCalc.DistanceKm(ctx, last.Location, cur)
	if err != nil {
		s.Logger.Warn("cannot calculate distance", zap.Error(err))
	}

	poi, err := s.POIManager.Resolve(ctx, cur)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", cur), zap.Error(err))
	}

	return &rentalpb.LocationStatus{
		Location:     cur,
		FeeCent:      last.FeeCent + int32(centsPerSec*elapsedSec*2*rand.Float64()),
		KmDriven:     last.KmDriven + dist,
		PoiName:      poi,
		TimestampSec: now,
	}
}
