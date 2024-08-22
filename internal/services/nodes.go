package services

import (
	"context"
	"log"
	"strings"

	gravity_api "github.com/c12s/agent_queue/pkg/api"
	"github.com/c12s/magnetar/internal/domain"
	meridian_api "github.com/c12s/meridian/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
)

type NodeService struct {
	nodeRepo      domain.NodeRepo
	administrator *oortapi.AdministrationAsyncClient
	authorizer    AuthZService
	meridian      meridian_api.MeridianClient
	gravity       gravity_api.AgentQueueClient
}

func NewNodeService(nodeRepo domain.NodeRepo, evaluator oortapi.OortEvaluatorClient, administrator *oortapi.AdministrationAsyncClient, authorizer AuthZService, meridian meridian_api.MeridianClient, gravity gravity_api.AgentQueueClient) (*NodeService, error) {
	return &NodeService{
		nodeRepo:      nodeRepo,
		administrator: administrator,
		authorizer:    authorizer,
		meridian:      meridian,
		gravity:       gravity,
	}, nil
}

func (n *NodeService) GetFromNodePool(ctx context.Context, req domain.GetFromNodePoolReq) (*domain.GetFromNodePoolResp, error) {
	node, err := n.nodeRepo.Get(req.Id, "")
	if err != nil {
		return nil, err
	}
	return &domain.GetFromNodePoolResp{
		Node: *node,
	}, nil
}

func (n *NodeService) GetFromOrg(ctx context.Context, req domain.GetFromOrgReq) (*domain.GetFromOrgResp, error) {
	if !n.authorizer.Authorize(ctx, "node.get", "node", req.Id.Value) {
		return nil, domain.ErrForbidden
	}
	node, err := n.nodeRepo.Get(req.Id, req.Org)
	if err != nil {
		return nil, err
	}
	return &domain.GetFromOrgResp{
		Node: *node,
	}, nil
}

func (n *NodeService) ClaimOwnership(ctx context.Context, req domain.ClaimOwnershipReq) (*domain.ClaimOwnershipResp, error) {
	if !n.authorizer.Authorize(ctx, "node.put", "org", req.Org) {
		return nil, domain.ErrForbidden
	}
	cluster, err := n.nodeRepo.ListOrgOwnedNodes(req.Org)
	if err != nil {
		return nil, err
	}
	nodes, err := n.nodeRepo.QueryNodePool(req.Query)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		err = n.nodeRepo.Delete(node)
		if err != nil {
			log.Println(err)
			continue
		}
		node.Org = req.Org
		err = n.nodeRepo.Put(node)
		if err != nil {
			log.Println(err)
			continue
		}
		err = n.administrator.SendRequest(&oortapi.CreateInheritanceRelReq{
			From: &oortapi.Resource{
				Id:   req.Org,
				Kind: "org",
			},
			To: &oortapi.Resource{
				Id:   node.Id.Value,
				Kind: "node",
			},
		}, func(resp *oortapi.AdministrationAsyncResp) {
			if resp.Error != "" {
				log.Println(resp.Error)
			}
		})
		if err != nil {
			log.Println(err)
		}
	}
	// upsert ns
	listNodesResp, err := n.ListOrgOwnedNodes(ctx, domain.ListOrgOwnedNodesReq{
		Org: req.Org,
	})
	if err != nil {
		return nil, err
	}
	resources := make(map[string]float64)
	for _, node := range listNodesResp.Nodes {
		for resource, quota := range node.Resources {
			resources[resource] = resources[resource] + quota
		}
	}
	_, err = n.meridian.GetNamespace(ctx, &meridian_api.GetNamespaceReq{
		OrgId: req.Org,
		Name:  "default",
	})
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "not found") {
			_, err = n.meridian.AddNamespace(ctx, &meridian_api.AddNamespaceReq{
				OrgId:                     req.Org,
				Name:                      "default",
				Labels:                    make(map[string]string),
				Quotas:                    resources,
				SeccompDefinitionStrategy: "redefine",
				Profile: &meridian_api.SeccompProfile{
					Version:       "v1.0.0",
					DefaultAction: "ALLOW",
					Syscalls:      make([]*meridian_api.SyscallRule, 0),
				},
			})
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		_, err = n.meridian.SetNamespaceResources(ctx, &meridian_api.SetNamespaceResourcesReq{
			OrgId:  req.Org,
			Name:   "default",
			Quotas: resources,
		})
		if err != nil {
			log.Println(err)
		}
	}
	// join cluster
	if len(nodes) == 0 {
		return &domain.ClaimOwnershipResp{
			Nodes: nodes,
		}, nil
	}
	joinAddress := nodes[0].BindAddress
	if len(cluster) > 0 {
		joinAddress = cluster[0].BindAddress
	}
	log.Println("join address: " + joinAddress)
	for _, node := range nodes {
		_, err = n.gravity.JoinCluster(ctx, &gravity_api.JoinClusterRequest{
			NodeId:      node.Id.Value,
			JoinAddress: joinAddress,
			ClusterId:   req.Org,
		})
		if err != nil {
			log.Println(err)
		}
	}
	return &domain.ClaimOwnershipResp{
		Nodes: nodes,
	}, nil
}

func (n *NodeService) ListNodePool(ctx context.Context, req domain.ListNodePoolReq) (*domain.ListNodePoolResp, error) {
	nodes, err := n.nodeRepo.ListNodePool()
	if err != nil {
		return nil, err
	}
	return &domain.ListNodePoolResp{
		Nodes: nodes,
	}, nil
}

func (n *NodeService) ListOrgOwnedNodes(ctx context.Context, req domain.ListOrgOwnedNodesReq) (*domain.ListOrgOwnedNodesResp, error) {
	// if !n.authorizer.Authorize(ctx, "node.get", "org", req.Org) {
	// 	return nil, domain.ErrForbidden
	// }
	nodes, err := n.nodeRepo.ListOrgOwnedNodes(req.Org)
	if err != nil {
		return nil, err
	}
	return &domain.ListOrgOwnedNodesResp{
		Nodes: nodes,
	}, nil
}

func (n *NodeService) ListAllNodes(ctx context.Context) ([]domain.Node, error) {
	return n.nodeRepo.ListAllNodes()
}

func (n *NodeService) QueryNodePool(ctx context.Context, req domain.QueryNodePoolReq) (*domain.QueryNodePoolResp, error) {
	nodes, err := n.nodeRepo.QueryNodePool(req.Query)
	if err != nil {
		return nil, err
	}
	return &domain.QueryNodePoolResp{
		Nodes: nodes,
	}, nil
}

func (n *NodeService) QueryOrgOwnedNodes(ctx context.Context, req domain.QueryOrgOwnedNodesReq) (*domain.QueryOrgOwnedNodesResp, error) {
	if !n.authorizer.Authorize(ctx, "node.get", "org", req.Org) {
		return nil, domain.ErrForbidden
	}
	nodes, err := n.nodeRepo.QueryOrgOwnedNodes(req.Query, req.Org)
	if err != nil {
		return nil, err
	}
	return &domain.QueryOrgOwnedNodesResp{
		Nodes: nodes,
	}, nil
}
