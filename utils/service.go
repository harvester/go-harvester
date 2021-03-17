package utils

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/futuretea/go-harvester/apis"
)

type ServiceBuilder struct {
	vm       *apis.VirtualMachine
	services map[string]*apis.Service
}

func NewServiceBuilder(vm *apis.VirtualMachine) *ServiceBuilder {
	return &ServiceBuilder{
		vm:       vm,
		services: make(map[string]*apis.Service),
	}
}

func (s *ServiceBuilder) Expose(name string, serviceType corev1.ServiceType, ports ...int32) *ServiceBuilder {
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
	servicePorts := make([]corev1.ServicePort, 0, len(ports))
	for _, port := range ports {
		servicePort := corev1.ServicePort{
			Name: fmt.Sprintf("%s-%d", name, port),
			Port: port,
			TargetPort: intstr.IntOrString{
				IntVal: port,
			},
		}
		servicePorts = append(servicePorts, servicePort)
	}
	svc := &apis.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Ports:    servicePorts,
			Selector: vm.Spec.Template.ObjectMeta.Labels,
		},
	}
	s.services[name] = svc
	return s
}

func (s *ServiceBuilder) Services() map[string]*apis.Service {
	return s.services
}
