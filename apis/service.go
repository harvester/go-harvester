package apis

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

type ServicesAPI struct {
	*Resource
}

func (s *ServicesAPI) List() (*SVCList, error) {
	var collection SVCList
	respCode, respBody, err := s.Resource.List()
	if err != nil {
		return nil, err
	}
	if respCode != http.StatusOK {
		return nil, NewResponseError(respCode, respBody)
	}
	err = json.Unmarshal(respBody, &collection)
	return &collection, err
}

func (s *ServicesAPI) Get(name string) (*Service, error) {
	var obj *Service
	respCode, respBody, err := s.Resource.Get(name)
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

func (s *ServicesAPI) Create(obj *Service) (*Service, error) {
	var created *Service
	respCode, respBody, err := s.Resource.Create(obj)
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

func (s *ServicesAPI) Delete(namespace, name string) (*Service, error) {
	var obj *Service
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.Resource.Delete(namespacedName)
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
