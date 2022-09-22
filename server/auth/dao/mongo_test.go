package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("62f7080a7688d6e7651b1046")),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("62f7080a7688d6e7651b1047")),
			openIDField:        "openid_2",
		},
	})
	if err != nil {
		panic(err)
	}
	mgutil.NewObjIDWithValue(id.AccountID("62f7080a7688d6e7651b1048"))
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "62f7080a7688d6e7651b1046",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "62f7080a7688d6e7651b1047",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "62f7080a7688d6e7651b1048",
		},
	}
	for _, cc := range cases {
		t.Run(
			cc.name,
			func(t *testing.T) {
				id, err := m.ResolveAccountID(context.Background(), cc.openID)
				if err != nil {
					t.Errorf("faild resolve account id for  %q: %v", cc.openID, err)
				}
				if id.String() != cc.want {
					t.Errorf("resolve account id: want: %q, got: %q", cc.want, id)
				}
			},
		)
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
