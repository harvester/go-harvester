package apis

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

type NodesAPI struct {
	*Resource
}

type Node corev1.Node

type NodeList struct {
	types.Collection
	Data []*Node `json:"data"`
}

func (s *NodesAPI) List() (*NodeList, error) {
	var collection NodeList
	respCode, respBody, err := s.Resource.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, ResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (s *NodesAPI) Get(name string) (*Node, error) {
	var obj *Node
	respCode, respBody, err := s.Resource.Get(name)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, ResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}
