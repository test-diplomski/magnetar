package proto

import (
	"errors"

	"github.com/c12s/magnetar/internal/domain"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/golang/protobuf/proto"
)

func LabelFromDomain(label domain.Label) (*api.Label, error) {
	value, err := ValueFromDomain(label.Value())
	if err != nil {
		return nil, err
	}
	return &api.Label{
		Key:   label.Key(),
		Value: value,
	}, nil
}

func LabelToDomain(l *api.Label) (domain.Label, error) {
	var label domain.Label
	var err error
	switch l.Value.Type {
	case api.Value_Bool:
		protoValue := &api.BoolValue{}
		err = proto.Unmarshal(l.Value.Marshalled, protoValue)
		if err == nil {
			label = domain.NewBoolLabel(l.Key, protoValue.Value)
		}
	case api.Value_Float64:
		protoValue := &api.Float64Value{}
		err = proto.Unmarshal(l.Value.Marshalled, protoValue)
		if err == nil {
			label = domain.NewFloat64Label(l.Key, protoValue.Value)
		}
	case api.Value_String:
		protoValue := &api.StringValue{}
		err = proto.Unmarshal(l.Value.Marshalled, protoValue)
		if err == nil {
			label = domain.NewStringLabel(l.Key, protoValue.Value)
		}
	default:
		err = errors.New("unsupported data type")
	}
	return label, err
}

func ValueFromDomain(value interface{}) (*api.Value, error) {
	var marshalled []byte
	var valueType api.Value_ValueTYpe
	var err error
	switch value := value.(type) {
	case bool:
		marshalled, err = proto.Marshal(&api.BoolValue{Value: value})
		valueType = api.Value_Bool
	case float64:
		marshalled, err = proto.Marshal(&api.Float64Value{Value: value})
		valueType = api.Value_Float64
	case string:
		marshalled, err = proto.Marshal(&api.StringValue{Value: value})
		valueType = api.Value_String
	default:
		err = errors.New("unsupported data type")
	}
	return &api.Value{
		Marshalled: marshalled,
		Type:       valueType,
	}, err
}

func NodeStringifiedFromDomain(node domain.Node) (*api.NodeStringified, error) {
	labels := make([]*api.LabelStringified, len(node.Labels))
	for i, label := range node.Labels {
		labelProto := &api.LabelStringified{}
		labelProto, err := LabelStringifiedFromDomain(label)
		if err != nil {
			return nil, err
		}
		labels[i] = labelProto
	}
	return &api.NodeStringified{
		Id:        node.Id.Value,
		Org:       node.Org,
		Labels:    labels,
		Resources: node.Resources,
	}, nil
}

func LabelStringifiedFromDomain(label domain.Label) (*api.LabelStringified, error) {
	return &api.LabelStringified{
		Key:   label.Key(),
		Value: label.StringValue(),
	}, nil
}

func NodeFromDomain(node domain.Node) (*api.Node, error) {
	resp := &api.Node{
		Id:          node.Id.Value,
		Org:         node.Org,
		Labels:      make([]*api.Label, len(node.Labels)),
		Resources:   node.Resources,
		BindAddress: node.BindAddress,
	}
	for i, label := range node.Labels {
		protoLabel, err := LabelFromDomain(label)
		if err != nil {
			return nil, err
		}
		resp.Labels[i] = protoLabel
	}
	return resp, nil
}

func NodeToDomain(node *api.Node) (*domain.Node, error) {
	resp := &domain.Node{
		Id: domain.NodeId{
			Value: node.Id,
		},
		Org:       node.Org,
		Labels:    make([]domain.Label, len(node.Labels)),
		Resources: node.Resources,
		BindAddress: node.BindAddress,
	}
	for i, protoLabel := range node.Labels {
		label, err := LabelToDomain(protoLabel)
		if err != nil {
			return nil, err
		}
		resp.Labels[i] = label
	}
	return resp, nil
}
