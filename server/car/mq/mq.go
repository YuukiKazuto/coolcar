package mq

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
)

type Publisher interface {
	Publish(context.Context, *carpb.CarEntity) error
}

type Subscriber interface {
	Subscribe(context.Context) (chan *carpb.CarEntity, func(), error)
}
