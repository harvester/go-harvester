package goharv

import (
	"encoding/json"
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

type VirtualMachineInstance kubevirtv1.VirtualMachineInstance

type VirtualMachineInstanceList struct {
	types.Collection
	Data []*VirtualMachineInstance `json:"data"`
}

type VirtualMachineInstanceClient struct {
	*apiClient
}

func newVirtualMachineInstanceClient(c *Client) *VirtualMachineInstanceClient {
	return &VirtualMachineInstanceClient{
		apiClient: newAPIClient(c, "kubevirt.io.virtualmachineinstance"),
	}
}

func (s *VirtualMachineInstanceClient) List() (*VirtualMachineInstanceList, error) {
	var collection VirtualMachineInstanceList
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

func (s *VirtualMachineInstanceClient) Get(namespace, name string) (*VirtualMachineInstance, error) {
	var obj *VirtualMachineInstance
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
