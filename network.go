package goharv

import (
	"encoding/json"
	"net/http"

	cniv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/rancher/apiserver/pkg/types"
)

type Network cniv1.NetworkAttachmentDefinition

type NetworkList struct {
	types.Collection
	Data []*Network `json:"data"`
}

type NetworksClient struct {
	*apiClient
}

func newNetworksClient(c *Client) *NetworksClient {
	return &NetworksClient{
		apiClient: newAPIClient(c, "k8s.cni.cncf.io.network-attachment-definitions"),
	}
}

func (s *NetworksClient) List() (*NetworkList, error) {
	var collection NetworkList
	respCode, respBody, err := s.apiClient.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (s *NetworksClient) Create(obj *Network) (*Network, error) {
	var created *Network
	respCode, respBody, err := s.apiClient.Create(obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &created); err != nil {
		return nil, err
	}
	return created, nil
}

func (s *NetworksClient) Update(namespace, name string, obj *Network) (*Network, error) {
	var created *Network
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.apiClient.Update(namespacedName, obj)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusCreated {
		return nil, NewResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &created); err != nil {
		return nil, err
	}
	return created, nil
}

func (s *NetworksClient) Get(namespace, name string) (*Network, error) {
	var obj *Network
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.apiClient.Get(namespacedName)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *NetworksClient) Delete(namespace, name string) (*Network, error) {
	var obj *Network
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.apiClient.Delete(namespacedName)
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	if err = json.Unmarshal(respBody, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}
