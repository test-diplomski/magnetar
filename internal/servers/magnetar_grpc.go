package servers

import (
	"context"
	"errors"

	"github.com/c12s/magnetar/internal/domain"
	"github.com/c12s/magnetar/internal/mappers/proto"
	"github.com/c12s/magnetar/internal/services"
	"github.com/c12s/magnetar/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MagnetarGrpcServer struct {
	api.UnimplementedMagnetarServer
	nodeService  services.NodeService
	labelService services.LabelService
}

func NewMagnetarGrpcServer(nodeService services.NodeService, labelService services.LabelService) (api.MagnetarServer, error) {
	return &MagnetarGrpcServer{
		nodeService:  nodeService,
		labelService: labelService,
	}, nil
}

func (m *MagnetarGrpcServer) GetFromNodePool(ctx context.Context, req *api.GetFromNodePoolReq) (*api.GetFromNodePoolResp, error) {
	domainReq, err := proto.GetFromNodePoolReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.GetFromNodePool(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.GetFromNodePoolRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) GetFromOrg(ctx context.Context, req *api.GetFromOrgReq) (*api.GetFromOrgResp, error) {
	domainReq, err := proto.GetFromOrgReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.GetFromOrg(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.GetFromOrgRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) ClaimOwnership(ctx context.Context, req *api.ClaimOwnershipReq) (*api.ClaimOwnershipResp, error) {
	domainReq, err := proto.ClaimOwnershipReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.ClaimOwnership(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.ClaimOwnershipRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) ListNodePool(ctx context.Context, req *api.ListNodePoolReq) (*api.ListNodePoolResp, error) {
	domainReq, err := proto.ListNodePoolReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.ListNodePool(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.ListNodePoolRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) ListOrgOwnedNodes(ctx context.Context, req *api.ListOrgOwnedNodesReq) (*api.ListOrgOwnedNodesResp, error) {
	domainReq, err := proto.ListOrgOwnedReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.ListOrgOwnedNodes(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.ListOrgOwnedNodesRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) QueryNodePool(ctx context.Context, req *api.QueryNodePoolReq) (*api.QueryNodePoolResp, error) {
	domainReq, err := proto.QueryNodePoolReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.QueryNodePool(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.QueryNodePoolRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) QueryOrgOwnedNodes(ctx context.Context, req *api.QueryOrgOwnedNodesReq) (*api.QueryOrgOwnedNodesResp, error) {
	domainReq, err := proto.QueryOrgOwnedNodesReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.nodeService.QueryOrgOwnedNodes(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.QueryOrgOwnedNodesRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) PutBoolLabel(ctx context.Context, req *api.PutBoolLabelReq) (*api.PutLabelResp, error) {
	domainReq, err := proto.PutBoolLabelReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.labelService.PutLabel(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.PutLabelRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) PutFloat64Label(ctx context.Context, req *api.PutFloat64LabelReq) (*api.PutLabelResp, error) {
	domainReq, err := proto.PutFloat64LabelReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.labelService.PutLabel(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.PutLabelRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) PutStringLabel(ctx context.Context, req *api.PutStringLabelReq) (*api.PutLabelResp, error) {
	domainReq, err := proto.PutStringLabelReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.labelService.PutLabel(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.PutLabelRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) DeleteLabel(ctx context.Context, req *api.DeleteLabelReq) (*api.DeleteLabelResp, error) {
	domainReq, err := proto.DeleteLabelReqToDomain(req)
	if err != nil {
		return nil, err
	}
	domainResp, err := m.labelService.DeleteLabel(ctx, *domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, err
	}
	return proto.DeleteLabelRespFromDomain(*domainResp)
}

func (m *MagnetarGrpcServer) ListAllNodes(ctx context.Context, req *api.ListAllNodesReq) (*api.ListAllNodesResp, error) {
	nodes, err := m.nodeService.ListAllNodes(ctx)
	if err != nil {
		return nil, err
	}
	return proto.ListAlldNodesRespFromDomain(nodes)
}

func GetAuthInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md.Get("authz-token")) > 0 {
			ctx = context.WithValue(ctx, "authz-token", md.Get("authz-token")[0])
		}
		// Calls the handler
		return handler(ctx, req)
	}
}
