package utils

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/futuretea/go-harvester/apis"
)

func NewServiceBuilder(vm *apis.VirtualMachine) *ServiceBuilder {
	return &ServiceBuilder{
		vm:       vm,
		services: make(map[string]*apis.Service),
	}
}

func (s *ServiceBuilder) Expose(name string, port int32, nodePort ...int32) *ServiceBuilder {
	vm := s.vm
	objectMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", vm.Name, name),
		Namespace: vm.Namespace,
		Labels:    vm.Spec.Template.ObjectMeta.Labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: vm.APIVersion,
				Kind:       vm.Kind,
				Name:       vm.Name,
				UID:        vm.UID,
			},
		},
	}
	svc := &apis.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: name,
					Port: port,
					TargetPort: intstr.IntOrString{
						IntVal: port,
					},
				},
			},
			Selector: vm.Spec.Template.ObjectMeta.Labels,
			Type:     corev1.ServiceTypeNodePort,
		},
	}
	if len(nodePort) != 0 {
		svc.Spec.Ports[0].NodePort = nodePort[0]
	}
	s.services[name] = svc
	return s
}

func (s *ServiceBuilder) Services() map[string]*apis.Service {
	return s.services
}
