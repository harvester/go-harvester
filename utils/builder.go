package utils

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/utils/pointer"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	cdiv1alpha1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"

	"github.com/futuretea/go-harvester/apis"
)

const (
	defaultVMGenerateName = "harv-"
	defaultVMNamespace    = "default"

	defaultVMCPUCores = 1
	defaultVMMemory   = "256Mi"

	defaultVMManagementNetworkName   = "default"
	defaultVMManagementInterfaceName = "default"
	defaultVMInterfaceModel          = "virtio"

	defaultVMCloudInitUserDataPasswordTemplate = `
#cloud-config
user: %s
password: %s
chpasswd: { expire: False }
ssh_pwauth: True`

	defaultVMCloudInitUserDataSSHKeyTemplate = `
#cloud-config
ssh_authorized_keys:
- >-
  %s`

	defaultVMCloudInitNetworkDataTemplate = `
network:
  version: 1
  config:
  - type: physical
    name: eth0
    subnets:
    - type: static
      address: %s
      gateway: %s`
)

type VMCloudInit struct {
	UserName  string
	Password  string
	PublicKey string
	Address   string
	Gateway   string
}

type VMBuilder struct {
	vm           *apis.VirtualMachine
	diskIndex    int
	networkIndex int
}

func NewVMBuilder(creator string) *VMBuilder {
	vmLabels := map[string]string{
		"harvester.cattle.io/creator": creator,
	}
	objectMeta := metav1.ObjectMeta{
		Namespace:    defaultVMNamespace,
		GenerateName: defaultVMGenerateName,
		Labels:       vmLabels,
	}
	running := pointer.BoolPtr(false)
	cpu := &kubevirtv1.CPU{
		Cores: defaultVMCPUCores,
	}
	resources := kubevirtv1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(defaultVMMemory),
		},
	}
	interfaces := []kubevirtv1.Interface{
		{
			Name:  defaultVMManagementInterfaceName,
			Model: defaultVMInterfaceModel,
			InterfaceBindingMethod: kubevirtv1.InterfaceBindingMethod{
				Masquerade: &kubevirtv1.InterfaceMasquerade{},
			},
		},
	}
	networks := []kubevirtv1.Network{
		{
			Name: defaultVMManagementNetworkName,
			NetworkSource: kubevirtv1.NetworkSource{
				Pod: &kubevirtv1.PodNetwork{},
			},
		},
	}
	template := &kubevirtv1.VirtualMachineInstanceTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: vmLabels,
		},
		Spec: kubevirtv1.VirtualMachineInstanceSpec{
			Domain: kubevirtv1.DomainSpec{
				CPU: cpu,
				Devices: kubevirtv1.Devices{
					Disks:      []kubevirtv1.Disk{},
					Interfaces: interfaces,
				},
				Resources: resources,
			},
			Networks: networks,
			Volumes:  []kubevirtv1.Volume{},
		},
	}

	vm := &apis.VirtualMachine{
		ObjectMeta: objectMeta,
		Spec: kubevirtv1.VirtualMachineSpec{
			Running:             running,
			Template:            template,
			DataVolumeTemplates: []kubevirtv1.DataVolumeTemplateSpec{},
		},
	}
	return &VMBuilder{
		vm: vm,
	}
}

func (v *VMBuilder) Name(name string) *VMBuilder {
	v.vm.ObjectMeta.Name = name
	v.vm.ObjectMeta.GenerateName = ""
	v.vm.Spec.Template.ObjectMeta.Labels["harvester.cattle.io/vmName"] = name
	return v
}

func (v *VMBuilder) Namespace(namespace string) *VMBuilder {
	v.vm.ObjectMeta.Namespace = namespace
	return v
}

func (v *VMBuilder) Memory(memory string) *VMBuilder {
	v.vm.Spec.Template.Spec.Domain.Resources.Requests = corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse(memory),
	}
	return v
}

func (v *VMBuilder) CPU(cores uint32) *VMBuilder {
	v.vm.Spec.Template.Spec.Domain.CPU.Cores = cores
	return v
}

func (v *VMBuilder) generateDiskName() string {
	return fmt.Sprintf("disk-%d", v.diskIndex)
}

func (v *VMBuilder) generateNetworkName() string {
	return fmt.Sprintf("network-%d", v.networkIndex)
}

func (v *VMBuilder) Blank(diskSize, diskBus string) *VMBuilder {
	return v.DataVolume(diskSize, diskBus)
}

func (v *VMBuilder) Image(diskSize, diskBus, sourceHTTPURL string) *VMBuilder {
	return v.DataVolume(diskSize, diskBus, sourceHTTPURL)
}

func (v *VMBuilder) DataVolume(diskSize, diskBus string, sourceHTTPURL ...string) *VMBuilder {
	diskName := v.generateDiskName()
	volumeMode := corev1.PersistentVolumeFilesystem
	dataVolumeName := fmt.Sprintf("%s-%s-%s", v.vm.Name, diskName, rand.String(5))
	// DataVolumeTemplates
	dataVolumeTemplates := v.vm.Spec.DataVolumeTemplates
	dataVolumeSpecSource := cdiv1alpha1.DataVolumeSource{
		Blank: &cdiv1alpha1.DataVolumeBlankImage{},
	}

	if len(sourceHTTPURL) > 0 && sourceHTTPURL[0] != "" {
		dataVolumeSpecSource = cdiv1alpha1.DataVolumeSource{
			HTTP: &cdiv1alpha1.DataVolumeSourceHTTP{
				URL: sourceHTTPURL[0],
			},
		}
	}
	dataVolumeTemplate := kubevirtv1.DataVolumeTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:        dataVolumeName,
			Labels:      nil,
			Annotations: nil,
		},
		Spec: cdiv1alpha1.DataVolumeSpec{
			Source: dataVolumeSpecSource,
			PVC: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(diskSize),
					},
				},
				VolumeMode: &volumeMode,
			},
		},
	}
	dataVolumeTemplates = append(dataVolumeTemplates, dataVolumeTemplate)
	v.vm.Spec.DataVolumeTemplates = dataVolumeTemplates
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	disks = append(disks, kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			DataVolume: &kubevirtv1.DataVolumeSource{
				Name: dataVolumeName,
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	// diskIndex
	v.diskIndex++
	return v
}

func (v *VMBuilder) ExistingDataVolume(dataVolumeName, diskBus string) *VMBuilder {
	diskName := v.generateDiskName()
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	disks = append(disks, kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			DataVolume: &kubevirtv1.DataVolumeSource{
				Name: dataVolumeName,
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	// diskIndex
	v.diskIndex++
	return v
}

func (v *VMBuilder) ContainerDisk(diskBus, imageName, ImagePullPolicy string, isCDRom bool) *VMBuilder {
	diskName := v.generateDiskName()
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	diskDevice := kubevirtv1.DiskDevice{
		Disk: &kubevirtv1.DiskTarget{
			Bus: diskBus,
		},
	}
	if isCDRom {
		diskDevice = kubevirtv1.DiskDevice{
			CDRom: &kubevirtv1.CDRomTarget{
				Bus: diskBus,
			},
		}
	}
	disks = append(disks, kubevirtv1.Disk{
		Name:       diskName,
		DiskDevice: diskDevice,
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			ContainerDisk: &kubevirtv1.ContainerDiskSource{
				Image:           imageName,
				ImagePullPolicy: corev1.PullPolicy(ImagePullPolicy),
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	// diskIndex
	v.diskIndex++
	return v
}

func (v *VMBuilder) Container(diskBus, imageName, ImagePullPolicy string) *VMBuilder {
	return v.ContainerDisk(diskBus, imageName, ImagePullPolicy, false)
}

func (v *VMBuilder) CDRom(diskBus, imageName, ImagePullPolicy string) *VMBuilder {
	return v.ContainerDisk(diskBus, imageName, ImagePullPolicy, true)
}

func (v *VMBuilder) CloudInit(vmCloudInit *VMCloudInit) *VMBuilder {
	diskName := "cloudinitdisk"
	diskBus := "virtio"
	// Disks
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	for _, disk := range disks {
		if disk.Name == diskName {
			return v
		}
	}

	disks = append(disks, kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks
	// Volumes
	var userData, networkData string
	if vmCloudInit != nil {
		if vmCloudInit.Password != "" {
			userData = fmt.Sprintf(defaultVMCloudInitUserDataPasswordTemplate, vmCloudInit.UserName, vmCloudInit.Password)
		}
		if vmCloudInit.PublicKey != "" {
			userData = fmt.Sprintf(defaultVMCloudInitUserDataSSHKeyTemplate, vmCloudInit.PublicKey)
		}
		if vmCloudInit.Address != "" && vmCloudInit.Gateway != "" {
			networkData = fmt.Sprintf(defaultVMCloudInitNetworkDataTemplate, vmCloudInit.Address, vmCloudInit.Gateway)
		}
	}
	volumes := v.vm.Spec.Template.Spec.Volumes
	volumes = append(volumes, kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			CloudInitNoCloud: &kubevirtv1.CloudInitNoCloudSource{
				UserData:    userData,
				NetworkData: networkData,
			},
		},
	})
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}

func (v *VMBuilder) Network(networkModel string) *VMBuilder {
	networkName := v.generateNetworkName()
	// Networks
	networks := v.vm.Spec.Template.Spec.Networks
	networks = append(networks, kubevirtv1.Network{
		Name: networkName,
		NetworkSource: kubevirtv1.NetworkSource{
			Multus: &kubevirtv1.MultusNetwork{
				NetworkName: networkName,
				Default:     false,
			},
		},
	})
	v.vm.Spec.Template.Spec.Networks = networks
	// Interfaces
	interfaces := v.vm.Spec.Template.Spec.Domain.Devices.Interfaces
	interfaces = append(interfaces, kubevirtv1.Interface{
		Name:  networkName,
		Model: networkModel,
		InterfaceBindingMethod: kubevirtv1.InterfaceBindingMethod{
			Bridge: &kubevirtv1.InterfaceBridge{},
		},
	})
	v.vm.Spec.Template.Spec.Domain.Devices.Interfaces = interfaces
	// networkIndex
	v.networkIndex++
	return v
}

func (v *VMBuilder) Run() *apis.VirtualMachine {
	v.vm.Spec.Running = pointer.BoolPtr(true)
	return v.vm
}

func (v *VMBuilder) VM() *apis.VirtualMachine {
	return v.vm
}

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

func (s *ServiceBuilder) Expose(name string, port int32) *ServiceBuilder {
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
	s.services[name] = svc
	return s
}

func (s *ServiceBuilder) Services() map[string]*apis.Service {
	return s.services
}
