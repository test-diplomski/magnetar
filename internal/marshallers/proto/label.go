package proto

import (
	"github.com/c12s/magnetar/internal/domain"
	mapper "github.com/c12s/magnetar/internal/mappers/proto"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/golang/protobuf/proto"
)

type protoLabelMarshaller struct {
}

func NewProtoLabelMarshaller() domain.LabelMarshaller {
	return &protoLabelMarshaller{}
}

func (p protoLabelMarshaller) Marshal(label domain.Label) ([]byte, error) {
	protoLabel, err := mapper.LabelFromDomain(label)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(protoLabel)
}

func (p protoLabelMarshaller) Unmarshal(labelMarshalled []byte) (domain.Label, error) {
	protoLabel := &api.Label{}
	err := proto.Unmarshal(labelMarshalled, protoLabel)
	if err != nil {
		return nil, err
	}
	return mapper.LabelToDomain(protoLabel)
}
