package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	harv1 "github.com/rancher/harvester/pkg/apis/harvester.cattle.io/v1alpha1"
)

type Setting harv1.Setting

type SettingList struct {
	types.Collection
	Data []*Setting `json:"data"`
}

type SettingsClient struct {
	*apiClient
}

func newSettingsClient(c *Client) *SettingsClient {
	return &SettingsClient{
		apiClient: newAPIClient(c, "harvester.cattle.io.settings"),
	}
}

func (s *SettingsClient) List() (*SettingList, error) {
	var collection SettingList
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

func (s *SettingsClient) Create(obj *Setting) (*Setting, error) {
	var created *Setting
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

func (s *SettingsClient) Update(name string, obj *Setting) (*Setting, error) {
	var created *Setting
	respCode, respBody, err := s.apiClient.Update(name, obj)
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

func (s *SettingsClient) Get(name string) (*Setting, error) {
	var obj *Setting
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

func (s *SettingsClient) Delete(name string) (*Setting, error) {
	var obj *Setting
	respCode, respBody, err := s.apiClient.Delete(name)
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
