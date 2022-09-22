package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	ctx := auth.ContextWithAccountID(context.Background(), id.AccountID("account1"))
	pm := &profileManager{}
	cm := &carManager{}
	s := newService(ctx, t, pm, cm)
	req := &rentalpb.CreateTripRequest{
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
		CarId: "car1",
	}
	nowFuc = func() int64 {
		return 1661685003
	}
	pm.iID = "identity1"
	golden := `{"account_id":"account1","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"fee_cent":1546005263,"poi_name":"厦门大学","timestamp_sec":1661685003},"current":{"location":{"latitude":32.123,"longitude":114.2525},"fee_cent":1546005263,"poi_name":"厦门大学","timestamp_sec":1661685003},"identity_id":"identity1"}`
	cases := []struct {
		name         string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:   "normal_create",
			tripID: "62fd24b25218a12088c70042",
			want:   golden,
		},
		{
			name:       "profile_err",
			tripID:     "62fd24b25218a12088c70043",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			tripID:       "62fd24b25218a12088c70044",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			tripID:       "62fd24b25218a12088c70045",
			carUnlockErr: fmt.Errorf("unlock"),
			want:         golden,
		},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.verifyErr = cc.carVerifyErr
			cm.unlockErr = cc.carUnlockErr
			res, err := s.CreateTrip(ctx, req)
			if cc.wantErr {
				if err == nil {
					t.Error("want error got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id: want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response: %v", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response: want %s, got %s", cc.want, got)
			}
		})
	}
}

func TestTripLifecycle(t *testing.T) {
	ctx := auth.ContextWithAccountID(context.Background(), id.AccountID("account_for_lifecycle"))
	s := newService(ctx, t, &profileManager{}, &carManager{})
	tid := id.TripID("62fe24b25218a12088c70044")
	mgutil.NewObjIDWithValue(tid)
	cases := []struct {
		name    string
		now     int64
		op      func() (*rentalpb.Trip, error)
		want    string
		wantErr bool
	}{
		{
			name: "create_trip",
			now:  10000,
			op: func() (*rentalpb.Trip, error) {
				e, err := s.CreateTrip(ctx, &rentalpb.CreateTripRequest{
					Start: &rentalpb.Location{
						Latitude:  32.123,
						Longitude: 114.2525,
					},
					CarId: "car1",
				})
				if err != nil {
					return nil, err
				}
				return e.Trip, nil
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"厦门大学","timestamp_sec":10000},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"厦门大学","timestamp_sec":10000},"status":1}`,
		},
		{
			name: "update_trip",
			now:  20000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(ctx, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
					Current: &rentalpb.Location{
						Latitude:  28.132,
						Longitude: 123.12345,
					},
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"厦门大学","timestamp_sec":10000},"current":{"location":{"latitude":28.132,"longitude":123.12345},"fee_cent":10677,"km_driven":100,"poi_name":"环岛路","timestamp_sec":20000},"status":1}`,
		},
		{
			name: "finish_trip",
			now:  30000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(ctx, &rentalpb.UpdateTripRequest{
					Id:      tid.String(),
					EndTrip: true,
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"厦门大学","timestamp_sec":10000},"current":{"location":{"latitude":28.132,"longitude":123.12345},"fee_cent":24674,"km_driven":100,"poi_name":"环岛路","timestamp_sec":30000},"end":{"location":{"latitude":28.132,"longitude":123.12345},"fee_cent":24674,"km_driven":100,"poi_name":"环岛路","timestamp_sec":30000},"status":2}`,
		},
		{
			name: "query_trip",
			now:  40000,
			op: func() (*rentalpb.Trip, error) {
				return s.GetTrip(ctx, &rentalpb.GetTripRequest{
					Id: tid.String(),
				})
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"厦门大学","timestamp_sec":10000},"current":{"location":{"latitude":28.132,"longitude":123.12345},"fee_cent":24674,"km_driven":100,"poi_name":"环岛路","timestamp_sec":30000},"end":{"location":{"latitude":28.132,"longitude":123.12345},"fee_cent":24674,"km_driven":100,"poi_name":"环岛路","timestamp_sec":30000},"status":2}`,
		},
		{
			name: "update_after_finished",
			now:  50000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(ctx, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
				})
			},
			wantErr: true,
		},
	}
	rand.Seed(1345)
	for _, cc := range cases {
		nowFuc = func() int64 {
			return cc.now
		}
		trip, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: operation failed: %v", cc.name, err)
			continue
		}
		b, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s: failed marshall response: %v", cc.name, err)
		}
		got := string(b)
		if cc.want != got {
			t.Errorf("%s: incorrect response; want: %s, got: %s", cc.name, cc.want, got)
		}
	}
}

func newService(ctx context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {
	mc, err := mongotesting.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	database := mc.Database("coolcar")
	mongotesting.SetupIndexes(ctx, database)
	return &Service{
		DistanceCalc:   &distCalc{},
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(database),
		Logger:         logger,
	}
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(ctx context.Context, accountID id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (c *carManager) Unlock(ctx context.Context, carID id.CarID, accountID id.AccountID, tripID id.TripID, s string) error {
	return c.unlockErr
}

func (c *carManager) Lock(ctx context.Context, carID id.CarID) error {
	return nil
}

func (c *carManager) Verify(ctx context.Context, carID id.CarID, location *rentalpb.Location) error {
	return c.verifyErr
}

type distCalc struct {
}

func (d *distCalc) DistanceKm(ctx context.Context, from *rentalpb.Location, to *rentalpb.Location) (float64, error) {
	if from.Latitude == to.Latitude && from.Longitude == to.Longitude {
		return 0, nil
	}
	return 100, nil
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
