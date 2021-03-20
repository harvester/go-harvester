package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvester.cattle.io/v1alpha1"
)

type Image harv1.VirtualMachineImage

type ImageList struct {
	types.Collection
	Data []*Image `json:"data"`
}

type ImagesClient struct {
	*apiClient
}

func newImagesClient(c *Client) *ImagesClient {
	return &ImagesClient{
		apiClient: newAPIClient(c, "harvester.cattle.io.virtualmachineimages"),
	}
}

func (s *ImagesClient) List() (*ImageList, error) {
	var collection ImageList
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

func (s *ImagesClient) Create(obj *Image) (*Image, error) {
	var created *Image
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

func (s *ImagesClient) Update(namespace, name string, obj *Image) (*Image, error) {
	var created *Image
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

func (s *ImagesClient) Get(namespace, name string) (*Image, error) {
	var obj *Image
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

func (s *ImagesClient) Delete(namespace, name string) (*Image, error) {
	var obj *Image
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
