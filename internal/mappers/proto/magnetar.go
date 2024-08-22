package proto

import (
	"log"

	"github.com/c12s/magnetar/internal/domain"
	"github.com/c12s/magnetar/pkg/api"
)

func GetFromNodePoolReqToDomain(req *api.GetFromNodePoolReq) (*domain.GetFromNodePoolReq, error) {
	return &domain.GetFromNodePoolReq{
		Id: domain.NodeId{
			Value: req.NodeId,
		},
	}, nil
}

func GetFromNodePoolRespFromDomain(resp domain.GetFromNodePoolResp) (*api.GetFromNodePoolResp, error) {
	nodeProto, err := NodeStringifiedFromDomain(resp.Node)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrServerSide
	}
	return &api.GetFromNodePoolResp{
		Node: nodeProto,
	}, nil
}

func GetFromOrgReqToDomain(req *api.GetFromOrgReq) (*domain.GetFromOrgReq, error) {
	return &domain.GetFromOrgReq{
		Id: domain.NodeId{
			Value: req.NodeId,
		},
		Org: req.Org,
	}, nil
}

func GetFromOrgRespFromDomain(resp domain.GetFromOrgResp) (*api.GetFromOrgResp, error) {
	nodeProto, err := NodeStringifiedFromDomain(resp.Node)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrServerSide
	}
	return &api.GetFromOrgResp{
		Node: nodeProto,
	}, nil
}

func ClaimOwnershipReqToDomain(req *api.ClaimOwnershipReq) (*domain.ClaimOwnershipReq, error) {
	query, err := queryToDomain(req.Query)
	if err != nil {
		return nil, err
	}
	return &domain.ClaimOwnershipReq{
		Query: query,
		Org:   req.Org,
	}, nil
}

func ClaimOwnershipRespFromDomain(resp domain.ClaimOwnershipResp) (*api.ClaimOwnershipResp, error) {
	nodesProto := make([]*api.NodeStringified, 0)
	for _, node := range resp.Nodes {
		nodeProto, err := NodeStringifiedFromDomain(node)
		if err != nil {
			log.Println(err)
			return nil, domain.ErrServerSide
		}
		nodesProto = append(nodesProto, nodeProto)
	}
	return &api.ClaimOwnershipResp{
		Node: nodesProto,
	}, nil
}

func ListNodePoolReqToDomain(req *api.ListNodePoolReq) (*domain.ListNodePoolReq, error) {
	return &domain.ListNodePoolReq{}, nil
}

func ListNodePoolRespFromDomain(resp domain.ListNodePoolResp) (*api.ListNodePoolResp, error) {
	nodesProto := make([]*api.NodeStringified, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodeProto, err := NodeStringifiedFromDomain(node)
		if err != nil {
			log.Println(err)
			return nil, domain.ErrServerSide
		}
		nodesProto[i] = nodeProto
	}
	return &api.ListNodePoolResp{
		Nodes: nodesProto,
	}, nil
}

func ListOrgOwnedReqToDomain(req *api.ListOrgOwnedNodesReq) (*domain.ListOrgOwnedNodesReq, error) {
	return &domain.ListOrgOwnedNodesReq{
		Org: req.Org,
	}, nil
}

func ListOrgOwnedNodesRespFromDomain(resp domain.ListOrgOwnedNodesResp) (*api.ListOrgOwnedNodesResp, error) {
	nodesProto := make([]*api.NodeStringified, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodeProto, err := NodeStringifiedFromDomain(node)
		if err != nil {
			log.Println(err)
			return nil, domain.ErrServerSide
		}
		nodesProto[i] = nodeProto
	}
	return &api.ListOrgOwnedNodesResp{
		Nodes: nodesProto,
	}, nil
}

func ListAlldNodesRespFromDomain(nodes []domain.Node) (*api.ListAllNodesResp, error) {
	nodesProto := make([]*api.NodeStringified, len(nodes))
	for i, node := range nodes {
		nodeProto, err := NodeStringifiedFromDomain(node)
		if err != nil {
			log.Println(err)
			return nil, domain.ErrServerSide
		}
		nodesProto[i] = nodeProto
	}
	return &api.ListAllNodesResp{
		Nodes: nodesProto,
	}, nil
}

func QueryNodePoolReqToDomain(req *api.QueryNodePoolReq) (*domain.QueryNodePoolReq, error) {
	query, err := queryToDomain(req.Query)
	if err != nil {
		return nil, err
	}
	return &domain.QueryNodePoolReq{
		Query: query,
	}, nil
}

func queryToDomain(query []*api.Selector) (domain.Query, error) {
	queryDomain := make([]domain.Selector, 0)
	for _, selector := range query {
		selectorDomain, err := selectorToDomain(selector)
		if err != nil {
			log.Println(err)
			return nil, domain.ErrServerSide
		}
		queryDomain = append(queryDomain, *selectorDomain)
	}
	return queryDomain, nil
}

func QueryNodePoolRespFromDomain(resp domain.QueryNodePoolResp) (*api.QueryNodePoolResp, error) {
	protoResp := &api.QueryNodePoolResp{
		Nodes: make([]*api.NodeStringified, 0),
	}
	for _, node := range resp.Nodes {
		protoNode := &api.NodeStringified{
			Id:        node.Id.Value,
			Labels:    make([]*api.LabelStringified, 0),
			Resources: node.Resources,
		}
		for _, label := range node.Labels {
			protoLabel := &api.LabelStringified{
				Key:   label.Key(),
				Value: label.StringValue(),
			}
			protoNode.Labels = append(protoNode.Labels, protoLabel)
		}
		protoResp.Nodes = append(protoResp.Nodes, protoNode)
	}
	return protoResp, nil
}

func QueryOrgOwnedNodesReqToDomain(req *api.QueryOrgOwnedNodesReq) (*domain.QueryOrgOwnedNodesReq, error) {
	query, err := queryToDomain(req.Query)
	if err != nil {
		return nil, err
	}
	return &domain.QueryOrgOwnedNodesReq{
		Org:   req.Org,
		Query: query,
	}, nil
}

func QueryOrgOwnedNodesRespFromDomain(resp domain.QueryOrgOwnedNodesResp) (*api.QueryOrgOwnedNodesResp, error) {
	protoResp := &api.QueryOrgOwnedNodesResp{
		Nodes: make([]*api.NodeStringified, 0),
	}
	for _, node := range resp.Nodes {
		protoNode := &api.NodeStringified{
			Id:        node.Id.Value,
			Labels:    make([]*api.LabelStringified, 0),
			Resources: node.Resources,
		}
		for _, label := range node.Labels {
			protoLabel := &api.LabelStringified{
				Key:   label.Key(),
				Value: label.StringValue(),
			}
			protoNode.Labels = append(protoNode.Labels, protoLabel)
		}
		protoResp.Nodes = append(protoResp.Nodes, protoNode)
	}
	return protoResp, nil
}

func PutBoolLabelReqToDomain(req *api.PutBoolLabelReq) (*domain.PutLabelReq, error) {
	return &domain.PutLabelReq{
		NodeId: domain.NodeId{
			Value: req.NodeId,
		},
		Label: domain.NewBoolLabel(req.Label.Key, req.Label.Value),
		Org:   req.Org,
	}, nil
}

func PutFloat64LabelReqToDomain(req *api.PutFloat64LabelReq) (*domain.PutLabelReq, error) {
	return &domain.PutLabelReq{
		NodeId: domain.NodeId{
			Value: req.NodeId,
		},
		Label: domain.NewFloat64Label(req.Label.Key, req.Label.Value),
		Org:   req.Org,
	}, nil
}

func PutStringLabelReqToDomain(req *api.PutStringLabelReq) (*domain.PutLabelReq, error) {
	return &domain.PutLabelReq{
		NodeId: domain.NodeId{
			Value: req.NodeId,
		},
		Label: domain.NewStringLabel(req.Label.Key, req.Label.Value),
		Org:   req.Org,
	}, nil
}

func PutLabelRespFromDomain(resp domain.PutLabelResp) (*api.PutLabelResp, error) {
	node, err := NodeStringifiedFromDomain(resp.Node)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrServerSide
	}
	return &api.PutLabelResp{
		Node: node,
	}, nil
}

func DeleteLabelReqToDomain(req *api.DeleteLabelReq) (*domain.DeleteLabelReq, error) {
	return &domain.DeleteLabelReq{
		NodeId: domain.NodeId{
			Value: req.NodeId,
		},
		LabelKey: req.LabelKey,
		Org:      req.Org,
	}, nil
}

func DeleteLabelRespFromDomain(resp domain.DeleteLabelResp) (*api.DeleteLabelResp, error) {
	node, err := NodeStringifiedFromDomain(resp.Node)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrServerSide
	}
	return &api.DeleteLabelResp{
		Node: node,
	}, nil
}

func selectorToDomain(query *api.Selector) (*domain.Selector, error) {
	shouldBe, err := domain.NewCompResultFromString(query.ShouldBe)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrServerSide
	}
	return &domain.Selector{
		LabelKey: query.LabelKey,
		ShouldBe: shouldBe,
		Value:    query.Value,
	}, nil
}
