package services

import (
	"github.com/c12s/magnetar/internal/domain"
	"github.com/google/uuid"
)

type RegistrationService struct {
	nodeRepo domain.NodeRepo
}

func NewRegistrationService(nodeRepo domain.NodeRepo) (*RegistrationService, error) {
	return &RegistrationService{
		nodeRepo: nodeRepo,
	}, nil
}

func (r *RegistrationService) Register(req domain.RegistrationReq) (*domain.RegistrationResp, error) {
	node := domain.Node{
		Id: domain.NodeId{
			Value: generateNodeId(),
		},
		Labels:    req.Labels,
		Resources: req.Resources,
		BindAddress: req.BindAddress,
	}

	err := r.nodeRepo.Put(node)
	if err != nil {
		return nil, err
	}

	return &domain.RegistrationResp{
		NodeId: node.Id.Value,
	}, nil
}

func generateNodeId() string {
	return uuid.NewString()
}
