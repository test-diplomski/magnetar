package services

import (
	"context"
	"github.com/c12s/magnetar/internal/domain"
	oortapi "github.com/c12s/oort/pkg/api"
)

type LabelService struct {
	nodeRepo   domain.NodeRepo
	authorizer AuthZService
}

func NewLabelService(nodeRepo domain.NodeRepo, evaluator oortapi.OortEvaluatorClient, authorizer AuthZService) (*LabelService, error) {
	return &LabelService{
		nodeRepo:   nodeRepo,
		authorizer: authorizer,
	}, nil
}

func (l *LabelService) PutLabel(ctx context.Context, req domain.PutLabelReq) (*domain.PutLabelResp, error) {
	if !l.authorizer.Authorize(ctx, "node.label.put", "node", req.NodeId.Value) {
		return nil, domain.ErrForbidden
	}
	node, err := l.nodeRepo.Get(req.NodeId, req.Org)
	if err != nil {
		return nil, err
	}
	node, err = l.nodeRepo.PutLabel(*node, req.Label)
	if err != nil {
		return nil, err
	}
	return &domain.PutLabelResp{
		Node: *node,
	}, nil
}

func (l *LabelService) DeleteLabel(ctx context.Context, req domain.DeleteLabelReq) (*domain.DeleteLabelResp, error) {
	if !l.authorizer.Authorize(ctx, "node.label.delete", "node", req.NodeId.Value) {
		return nil, domain.ErrForbidden
	}
	node, err := l.nodeRepo.Get(req.NodeId, req.Org)
	if err != nil {
		return nil, err
	}
	node, err = l.nodeRepo.DeleteLabel(*node, req.LabelKey)
	if err != nil {
		return nil, err
	}
	return &domain.DeleteLabelResp{
		Node: *node,
	}, nil
}
