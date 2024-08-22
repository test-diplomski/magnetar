package servers

import (
	"github.com/c12s/magnetar/internal/mappers/proto"
	"github.com/c12s/magnetar/internal/services"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/c12s/magnetar/pkg/messaging"
	"log"
)

type RegistrationAsyncServer struct {
	subscriber messaging.Subscriber
	publisher  messaging.Publisher
	service    services.RegistrationService
}

func NewRegistrationAsyncServer(subscriber messaging.Subscriber, publisher messaging.Publisher, service services.RegistrationService) (*RegistrationAsyncServer, error) {
	return &RegistrationAsyncServer{
		subscriber: subscriber,
		publisher:  publisher,
		service:    service,
	}, nil
}

func (n *RegistrationAsyncServer) Serve() error {
	return n.subscriber.Subscribe(n.register)
}

func (n *RegistrationAsyncServer) register(msg []byte, replySubject string) {
	reqProto := &api.RegistrationReq{}
	err := reqProto.Unmarshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := proto.RegistrationReqToDomain(reqProto)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := n.service.Register(*req)
	if err != nil {
		log.Println(err)
		return
	}
	respProto, err := proto.RegistrationRespFromDomain(*resp)
	if err != nil {
		log.Println(err)
		return
	}
	respMarshalled, err := respProto.Marshal()
	if err != nil {
		log.Println(err)
		return
	}
	err = n.publisher.Publish(respMarshalled, replySubject)
	if err != nil {
		log.Println(err)
	}
}

func (n *RegistrationAsyncServer) GracefulStop() {
	err := n.subscriber.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
}
