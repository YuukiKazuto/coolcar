package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"github.com/golang/protobuf/proto"
	"hash/fnv"
)

var poi = []string{
	"双子塔",
	"厦门大学",
	"环岛路",
	"椰风寨",
	"狐尾山气象台",
	"山海健步道",
}

type Manager struct {
}

func (*Manager) Resolve(ctx context.Context, location *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(location)
	if err != nil {
		return "", err
	}
	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
