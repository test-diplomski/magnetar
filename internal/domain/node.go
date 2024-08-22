package domain

type Node struct {
	Id          NodeId
	Org         string
	Labels      []Label
	Resources   map[string]float64
	BindAddress string
}

func (n Node) Claimed() bool {
	return len(n.Org) > 0
}

type NodeId struct {
	Value string
}

type Query []Selector

type Selector struct {
	LabelKey string
	ShouldBe ComparisonResult
	Value    string
}

type NodeRepo interface {
	Put(node Node) error
	Get(nodeId NodeId, org string) (*Node, error)
	Delete(node Node) error
	ListNodePool() ([]Node, error)
	ListOrgOwnedNodes(org string) ([]Node, error)
	QueryNodePool(query Query) ([]Node, error)
	QueryOrgOwnedNodes(query Query, org string) ([]Node, error)
	PutLabel(node Node, label Label) (*Node, error)
	DeleteLabel(node Node, labelKey string) (*Node, error)
	ListAllNodes() ([]Node, error)
}

type NodeMarshaller interface {
	Marshal(node Node) ([]byte, error)
	Unmarshal(nodeMarshalled []byte) (*Node, error)
}

type GetFromNodePoolReq struct {
	Id NodeId
}

type GetFromNodePoolResp struct {
	Node Node
}

type GetFromOrgReq struct {
	Id  NodeId
	Org string
}

type GetFromOrgResp struct {
	Node Node
}

type ClaimOwnershipReq struct {
	Query Query
	Org   string
}

type ClaimOwnershipResp struct {
	Nodes []Node
}

type ListNodePoolReq struct {
}

type ListNodePoolResp struct {
	Nodes []Node
}

type ListOrgOwnedNodesReq struct {
	Org string
}

type ListOrgOwnedNodesResp struct {
	Nodes []Node
}

type QueryNodePoolReq struct {
	Query Query
}

type QueryNodePoolResp struct {
	Nodes []Node
}

type QueryOrgOwnedNodesReq struct {
	Query Query
	Org   string
}

type QueryOrgOwnedNodesResp struct {
	Nodes []Node
}
