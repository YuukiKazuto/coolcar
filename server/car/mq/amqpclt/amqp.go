package amqpclt

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	err = declareExchange(ch, exchange)
	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, car *carpb.CarEntity) error {
	b, err := json.Marshal(car)
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}
	return p.ch.Publish(
		p.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Body: b,
		},
	)
}

type Subscriber struct {
	conn     *amqp.Connection
	exchange string
	logger   *zap.Logger
}

func NewSubscriber(conn *amqp.Connection, exchange string, logger *zap.Logger) (*Subscriber, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	defer ch.Close()
	err = declareExchange(ch, exchange)
	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Subscriber{
		conn:     conn,
		exchange: exchange,
		logger:   logger,
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context) (chan *carpb.CarEntity, func(), error) {
	msgCh, cleanUp, err := s.SubscribeRaw(ctx)
	if err != nil {
		return nil, cleanUp, err
	}
	carCh := make(chan *carpb.CarEntity)
	go func() {
		for msg := range msgCh {
			var car carpb.CarEntity
			err := json.Unmarshal(msg.Body, &car)
			if err != nil {
				s.logger.Error("cannot unmarshal", zap.Error(err))
			}
			carCh <- &car
		}
		close(carCh)
	}()
	return carCh, cleanUp, nil
}

func (s *Subscriber) SubscribeRaw(ctx context.Context) (<-chan amqp.Delivery, func(), error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot allocate channel: %v", err)
	}
	closeCh := func() {
		err := ch.Close()
		if err != nil {
			s.logger.Error("cannot close channel", zap.Error(err))
		}
	}
	q, err := ch.QueueDeclare(
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, closeCh, fmt.Errorf("cannot declare queue: %v", err)
	}
	cleanUp := func() {
		_, err := ch.QueueDelete(
			q.Name,
			false, // ifUnused
			false, // ifEmpty
			false, // noWait
		)
		if err != nil {
			s.logger.Error("cannot delete queue", zap.String("name", q.Name), zap.Error(err))
		}
		closeCh()
	}
	err = ch.QueueBind(
		q.Name,
		"",
		s.exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot bind: %v", err)
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot comsume queue: %v", err)
	}
	return msgs, cleanUp, nil

}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}
