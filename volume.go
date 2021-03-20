package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	cdiv1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type DataVolume cdiv1beta1.DataVolume

type DataVolumeList struct {
	types.Collection
	Data []*DataVolume `json:"data"`
}

type DataVolumesClient struct {
	*apiClient
}

func newDataVolumesClient(c *Client) *DataVolumesClient {
	return &DataVolumesClient{
		apiClient: newAPIClient(c, "cdi.kubevirt.io.datavolumes"),
	}
}

func (s *DataVolumesClient) List() (*DataVolumeList, error) {
	var collection DataVolumeList
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

func (s *DataVolumesClient) Create(obj *DataVolume) (*DataVolume, error) {
	var created *DataVolume
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

func (s *DataVolumesClient) Update(namespace, name string, obj *DataVolume) (*DataVolume, error) {
	var created *DataVolume
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

func (s *DataVolumesClient) Get(namespace, name string) (*DataVolume, error) {
	var obj *DataVolume
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

func (s *DataVolumesClient) Delete(namespace, name string) (*DataVolume, error) {
	var obj *DataVolume
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
