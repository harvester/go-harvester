package apis

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	cdiv1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
)

type DataVolumesAPI struct {
	*Resource
}

type DataVolume cdiv1beta1.DataVolume

type DataVolumeList struct {
	types.Collection
	Data []*DataVolume `json:"data"`
}

func (s *DataVolumesAPI) List() (*DataVolumeList, error) {
	var collection DataVolumeList
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

func (s *DataVolumesAPI) Create(obj *DataVolume) (*DataVolume, error) {
	var created *DataVolume
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

func (s *DataVolumesAPI) Update(namespace, name string, obj *DataVolume) (*DataVolume, error) {
	var created *DataVolume
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.Resource.Update(namespacedName, obj)
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

func (s *DataVolumesAPI) Get(namespace, name string) (*DataVolume, error) {
	var obj *DataVolume
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.Resource.Get(namespacedName)
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

func (s *DataVolumesAPI) Delete(namespace, name string) (*DataVolume, error) {
	var obj *DataVolume
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
