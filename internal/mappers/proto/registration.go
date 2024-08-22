package proto

import (
	"github.com/c12s/magnetar/internal/domain"
	"github.com/c12s/magnetar/pkg/api"
)

func RegistrationReqToDomain(req *api.RegistrationReq) (*domain.RegistrationReq, error) {
	var labels []domain.Label
	for _, protoLabel := range req.Labels {
		label, err := LabelToDomain(protoLabel)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	return &domain.RegistrationReq{
		Labels:      labels,
		Resources:   req.Resources,
		BindAddress: req.BindAddress,
	}, nil
}

func RegistrationRespFromDomain(resp domain.RegistrationResp) (*api.RegistrationResp, error) {
	return &api.RegistrationResp{
		NodeId: resp.NodeId,
	}, nil
}
