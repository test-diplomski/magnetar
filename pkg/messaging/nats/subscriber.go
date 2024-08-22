package nats

import (
	"errors"
	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/nats-io/nats.go"
)

type subscriber struct {
	conn         *nats.Conn
	subscription *nats.Subscription
	subject      string
	queue        string
}

func NewSubscriber(conn *nats.Conn, subject, queue string) (messaging.Subscriber, error) {
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	return &subscriber{
		conn:    conn,
		subject: subject,
		queue:   queue,
	}, nil
}

func (s *subscriber) Subscribe(handler func(msg []byte, replySubject string)) error {
	if s.subscription != nil {
		return errors.New("already subscribed")
	}
	subscription, err := s.conn.QueueSubscribe(s.subject, s.queue, func(msg *nats.Msg) {
		handler(msg.Data, msg.Reply)
	})
	if err != nil {
		return err
	}
	s.subscription = subscription
	return nil
}

func (s *subscriber) Unsubscribe() error {
	if s.subscription != nil && s.subscription.IsValid() {
		return s.subscription.Drain()
	}
	return nil
}
