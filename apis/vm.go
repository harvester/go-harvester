package apis

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rancher/apiserver/pkg/types"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

type VirtualMachine kubevirtv1.VirtualMachine

type VirtualMachineList struct {
	types.Collection
	Data []*VirtualMachine `json:"data"`
}

type VirtualMachinesClient struct {
	*Resource
}

func (s *VirtualMachinesClient) List() (*VirtualMachineList, error) {
	var collection VirtualMachineList
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

func (s *VirtualMachinesClient) Create(obj *VirtualMachine) (*VirtualMachine, error) {
	var created *VirtualMachine
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

func (s *VirtualMachinesClient) Update(namespace, name string, obj *VirtualMachine) (*VirtualMachine, error) {
	var created *VirtualMachine
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

func (s *VirtualMachinesClient) Get(namespace, name string) (*VirtualMachine, error) {
	var obj *VirtualMachine
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

func (s *VirtualMachinesClient) Delete(namespace, name string, removedDisks []string) (*VirtualMachine, error) {
	var obj *VirtualMachine
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.Resource.Delete(namespacedName, map[string]string{
		"removedDisks": strings.Join(removedDisks, ","),
	})
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

func (s *VirtualMachinesClient) Kill(namespace, name string) (*VirtualMachine, error) {
	vm, err := s.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	*vm.Spec.Running = false
	*vm.Spec.Template.Spec.TerminationGracePeriodSeconds = 0
	return s.Update(namespace, name, vm)
}

func (s *VirtualMachinesClient) Start(namespace, name string) error {
	return s.simpleAction(namespace, name, "start")
}

func (s *VirtualMachinesClient) Stop(namespace, name string) error {
	return s.simpleAction(namespace, name, "stop")
}

func (s *VirtualMachinesClient) Restart(namespace, name string) error {
	return s.simpleAction(namespace, name, "restart")
}

func (s *VirtualMachinesClient) simpleAction(namespace, name, action string) error {
	namespacedName := namespace + "/" + name
	respCode, respBody, err := s.Resource.Action(namespacedName, action, nil)
	if err != nil {
		return err
	}
	if respCode != http.StatusNoContent {
		return NewResponseError(respCode, respBody)
	}
	return nil
}
