package repos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/c12s/magnetar/internal/domain"
	"github.com/juliangruber/go-intersect"
	etcd "go.etcd.io/etcd/client/v3"
	"golang.org/x/exp/slices"
)

// data model
// for get operations
// key - nodes/pool/{nodeId} | nodes/orgs/{orgId}/{nodeId}
// value - protobuf node (id + org + labels)
// for query operations
// key - labels/pool/{labelKey}/{nodeId} | labels/orgs/{orgId}/{labelKey}/{nodeId}
// value - protobuf label (key + value)

type nodeEtcdRepo struct {
	etcd            *etcd.Client
	nodeMarshaller  domain.NodeMarshaller
	labelMarshaller domain.LabelMarshaller
}

func NewNodeEtcdRepo(etcd *etcd.Client, nodeMarshaller domain.NodeMarshaller, labelMarshaller domain.LabelMarshaller) (domain.NodeRepo, error) {
	return &nodeEtcdRepo{
		etcd:            etcd,
		nodeMarshaller:  nodeMarshaller,
		labelMarshaller: labelMarshaller,
	}, nil
}

func (n nodeEtcdRepo) Put(node domain.Node) error {
	err := n.putNodeGetModel(node)
	if err != nil {
		return err
	}
	return n.putNodeQueryModel(node)
}

func (n nodeEtcdRepo) Delete(node domain.Node) error {
	err := n.deleteNodeGetModel(node)
	if err != nil {
		return err
	}
	return n.deleteNodeQueryModel(node)
}

func (n nodeEtcdRepo) Get(nodeId domain.NodeId, org string) (*domain.Node, error) {
	key := getKey(domain.Node{Id: nodeId, Org: org})
	resp, err := n.etcd.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New("node not found")
	}
	return n.nodeMarshaller.Unmarshal(resp.Kvs[0].Value)
}

func (n nodeEtcdRepo) ListNodePool() ([]domain.Node, error) {
	keyPrefix := fmt.Sprintf("%s/pool", getKeyPrefix)
	return n.listNodes(keyPrefix)
}

func (n nodeEtcdRepo) ListOrgOwnedNodes(org string) ([]domain.Node, error) {
	keyPrefix := fmt.Sprintf("%s/orgs/%s", getKeyPrefix, org)
	return n.listNodes(keyPrefix)
}

func (n nodeEtcdRepo) ListAllNodes() ([]domain.Node, error) {
	return n.listNodes(getKeyPrefix)
}

func (n nodeEtcdRepo) QueryNodePool(query domain.Query) ([]domain.Node, error) {
	keyPrefix := fmt.Sprintf("%s/pool", queryKeyPrefix)
	if len(query) == 0 {
		return n.listNodes(keyPrefix)
	}
	nodeIds, err := n.queryNodes(query, keyPrefix)
	if err != nil {
		return nil, err
	}
	nodes := make([]domain.Node, 0)
	for _, nodeId := range nodeIds {
		node, err := n.Get(nodeId, "")
		if err != nil {
			log.Println(err)
			continue
		}
		nodes = append(nodes, *node)
	}
	return nodes, nil
}

func (n nodeEtcdRepo) QueryOrgOwnedNodes(query domain.Query, org string) ([]domain.Node, error) {
	keyPrefix := fmt.Sprintf("%s/orgs/%s", queryKeyPrefix, org)
	if len(query) == 0 {
		return n.listNodes(keyPrefix)
	}
	nodeIds, err := n.queryNodes(query, keyPrefix)
	if err != nil {
		return nil, err
	}
	nodes := make([]domain.Node, 0)
	for _, nodeId := range nodeIds {
		node, err := n.Get(nodeId, org)
		if err != nil {
			log.Println(err)
			continue
		}
		nodes = append(nodes, *node)
	}
	return nodes, nil
}

func (n nodeEtcdRepo) listNodes(keyPrefix string) ([]domain.Node, error) {
	nodes := make([]domain.Node, 0)
	resp, err := n.etcd.Get(context.TODO(), keyPrefix, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		node, err := n.nodeMarshaller.Unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	}
	return nodes, nil
}

func (n nodeEtcdRepo) PutLabel(node domain.Node, label domain.Label) (*domain.Node, error) {
	err := n.putLabelGetModel(node, label)
	if err != nil {
		return nil, err
	}
	err = n.putLabelQueryModel(node, label)
	if err != nil {
		return nil, err
	}
	return n.Get(node.Id, node.Org)
}

func (n nodeEtcdRepo) DeleteLabel(node domain.Node, labelKey string) (*domain.Node, error) {
	err := n.deleteLabelGetModel(node, labelKey)
	if err != nil {
		return nil, err
	}
	err = n.deleteLabelQueryModel(node, labelKey)
	if err != nil {
		return nil, err
	}
	return n.Get(node.Id, node.Org)
}

func (n nodeEtcdRepo) deleteNodeGetModel(node domain.Node) error {
	key := getKey(node)
	_, err := n.etcd.Delete(context.TODO(), key)
	return err
}

func (n nodeEtcdRepo) deleteNodeQueryModel(node domain.Node) (err error) {
	for _, label := range node.Labels {
		key := queryKey(node, label.Key())
		_, delErr := n.etcd.Delete(context.TODO(), key)
		err = errors.Join(err, delErr)
	}
	return err
}

func (n nodeEtcdRepo) putNodeGetModel(node domain.Node) error {
	nodeMarshalled, err := n.nodeMarshaller.Marshal(node)
	if err != nil {
		return err
	}
	key := getKey(node)
	_, err = n.etcd.Put(context.TODO(), key, string(nodeMarshalled))
	return err
}

func (n nodeEtcdRepo) putLabelGetModel(node domain.Node, label domain.Label) error {
	labelIndex := -1
	for i, nodeLabel := range node.Labels {
		if nodeLabel.Key() == label.Key() {
			labelIndex = i
		}
	}
	if labelIndex >= 0 {
		node.Labels[labelIndex] = label
	} else {
		node.Labels = append(node.Labels, label)
	}
	return n.putNodeGetModel(node)
}

func (n nodeEtcdRepo) deleteLabelGetModel(node domain.Node, labelKey string) error {
	labelIndex := -1
	for i, nodeLabel := range node.Labels {
		if nodeLabel.Key() == labelKey {
			labelIndex = i
		}
	}
	if labelIndex >= 0 {
		node.Labels = slices.Delete(node.Labels, labelIndex, labelIndex+1)
		return n.putNodeGetModel(node)
	}
	return nil
}

func (n nodeEtcdRepo) putNodeQueryModel(node domain.Node) error {
	for _, label := range node.Labels {
		err := n.putLabelQueryModel(node, label)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n nodeEtcdRepo) putLabelQueryModel(node domain.Node, label domain.Label) error {
	labelMarshalled, err := n.labelMarshaller.Marshal(label)
	if err != nil {
		return err
	}
	key := queryKey(node, label.Key())
	_, err = n.etcd.Put(context.TODO(), key, string(labelMarshalled))
	return err
}

func (n nodeEtcdRepo) deleteLabelQueryModel(node domain.Node, labelKey string) error {
	key := queryKey(node, labelKey)
	_, err := n.etcd.Delete(context.TODO(), key)
	return err
}

func (n nodeEtcdRepo) queryNodes(query domain.Query, keyPrefix string) ([]domain.NodeId, error) {
	nodeIds := make([]domain.NodeId, 0)
	for i, selector := range query {
		currNodes, err := n.selectNodes(selector, keyPrefix)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			nodeIds = currNodes
		} else {
			intersection := intersect.Simple(nodeIds, currNodes)
			nodeIds = make([]domain.NodeId, len(intersection))
			for i, node := range intersection {
				nodeIds[i] = node.(domain.NodeId)
			}
		}
	}
	return nodeIds, nil
}

func (n nodeEtcdRepo) selectNodes(selector domain.Selector, keyPrefix string) ([]domain.NodeId, error) {
	prefix := fmt.Sprintf("%s/%s/", keyPrefix, selector.LabelKey)
	resp, err := n.etcd.Get(context.TODO(), prefix, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}
	nodeIds := make([]domain.NodeId, 0)
	for _, kv := range resp.Kvs {
		nodeLabel, err := n.labelMarshaller.Unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		cmpResult, err := nodeLabel.Compare(selector.Value)
		if err != nil {
			log.Println(err)
			continue
		}
		if slices.Contains(cmpResult, selector.ShouldBe) {
			nodeId := extractNodeIdFromQueryKey(string(kv.Key))
			nodeIds = append(nodeIds, domain.NodeId{
				Value: nodeId,
			})
		}
	}
	return nodeIds, nil
}

const (
	getKeyPrefix   = "nodes"
	queryKeyPrefix = "labels"
)

func getKey(node domain.Node) string {
	if node.Claimed() {
		return fmt.Sprintf("%s/orgs/%s/%s", getKeyPrefix, node.Org, node.Id.Value)
	}
	return fmt.Sprintf("%s/pool/%s", getKeyPrefix, node.Id.Value)
}

func queryKey(node domain.Node, labelKey string) string {
	if node.Claimed() {
		return fmt.Sprintf("%s/orgs/%s/%s/%s", queryKeyPrefix, node.Org, labelKey, node.Id.Value)
	}
	return fmt.Sprintf("%s/pool/%s/%s", queryKeyPrefix, labelKey, node.Id.Value)
}

func extractNodeIdFromQueryKey(key string) string {
	if strings.HasPrefix(key, fmt.Sprintf("%s/pool", queryKeyPrefix)) {
		return strings.Split(key, "/")[3]
	}
	return strings.Split(key, "/")[4]
}
