package messaging

type Subscriber interface {
	Subscribe(handler func(msg []byte, replySubject string)) error
	Unsubscribe() error
}

type Publisher interface {
	Publish(msg []byte, subject string) error
	Request(msg []byte, subject, replySubject string) error
	GenerateReplySubject() string
}
