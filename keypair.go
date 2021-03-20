package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvester.cattle.io/v1alpha1"
)

type KeyPair harv1.KeyPair

type KeyPairList struct {
	types.Collection
	Data []*KeyPair `json:"data"`
}

type KeyPairsClient struct {
	*apiClient
}

func newKeyPairsClient(c *Client) *KeyPairsClient {
	return &KeyPairsClient{
		apiClient: newAPIClient(c, "harvester.cattle.io.keypairs"),
	}
}

func (s *KeyPairsClient) List() (*KeyPairList, error) {
	var collection KeyPairList
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

func (s *KeyPairsClient) Create(obj *KeyPair) (*KeyPair, error) {
	var created *KeyPair
	respCode, respBody, err := s.apiClient.Create(obj)
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

func (s *KeyPairsClient) Update(namespace, name string, obj *KeyPair) (*KeyPair, error) {
	var created *KeyPair
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.apiClient.Update(namespacedName, obj)
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

func (s *KeyPairsClient) Get(namespace, name string) (*KeyPair, error) {
	var obj *KeyPair
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

func (s *KeyPairsClient) Delete(namespace, name string) (*KeyPair, error) {
	var obj *KeyPair
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
