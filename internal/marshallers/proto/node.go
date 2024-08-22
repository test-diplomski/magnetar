package proto

import (
	"github.com/c12s/magnetar/internal/domain"
	mapper "github.com/c12s/magnetar/internal/mappers/proto"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/golang/protobuf/proto"
)

type protoNodeMarshaller struct {
}

func NewProtoNodeMarshaller() domain.NodeMarshaller {
	return &protoNodeMarshaller{}
}

func (p protoNodeMarshaller) Marshal(node domain.Node) ([]byte, error) {
	protoNode, err := mapper.NodeFromDomain(node)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(protoNode)
}

func (p protoNodeMarshaller) Unmarshal(nodeMarshalled []byte) (*domain.Node, error) {
	protoNode := &api.Node{}
	err := proto.Unmarshal(nodeMarshalled, protoNode)
	if err != nil {
		return nil, err
	}
	return mapper.NodeToDomain(protoNode)
}
