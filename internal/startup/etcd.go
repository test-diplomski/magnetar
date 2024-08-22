package startup

import (
	"fmt"
	etcd "go.etcd.io/etcd/client/v3"
)

func newEtcdClient(address string) (*etcd.Client, error) {
	return etcd.New(etcd.Config{
		Endpoints: []string{fmt.Sprintf("http://%s", address)},
	})
}
