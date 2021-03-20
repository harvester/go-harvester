package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

type Service corev1.Service

type SVCList struct {
	types.Collection
	Data []*Service `json:"data"`
}

type ServicesClient struct {
	*apiClient
}

func newServicesClient(c *Client) *ServicesClient {
	return &ServicesClient{
		apiClient: newAPIClient(c, "services"),
	}
}

func (s *ServicesClient) List() (*SVCList, error) {
	var collection SVCList
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

func (s *ServicesClient) Get(name string) (*Service, error) {
	var obj *Service
	respCode, respBody, err := s.apiClient.Get(name)
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

func (s *ServicesClient) Create(obj *Service) (*Service, error) {
	var created *Service
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

func (s *ServicesClient) Delete(namespace, name string) (*Service, error) {
	var obj *Service
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
