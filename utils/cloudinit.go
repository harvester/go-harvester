package utils

import (
	"fmt"

	kubevirtv1 "kubevirt.io/client-go/api/v1"
)

const (
	defaultCloudInitUserDataBasic = `#cloud-config
package_update: true
packages:
- qemu-guest-agent
runcmd:
- [systemctl, enable, qemu-guest-agent]
- [systemctl, start, qemu-guest-agent]
`
	defaultCloudInitUserDataPasswordTemplate = `
user: %s
password: %s
chpasswd: { expire: False }
ssh_pwauth: True`

	defaultCloudInitUserDataSSHKeyTemplate = `
ssh_authorized_keys:
- >-
  %s`
	defaultCloudInitNetworkDataTemplate = `
network:
  version: 1
  config:
  - type: physical
    name: %s`
	defaultCloudInitNetworkDataDHCPTemplate = `
    subnets:
    - type: dhcp`
	defaultCloudInitNetworkDataStaticTemplate = `
    subnets:
    - type: static
      address: %s
      gateway: %s`
)

type VMCloudInit struct {
	UserName      string
	Password      string
	PublicKey     string
	InterfaceName string
	Address       string
	Gateway       string
}

func generateCloudInit(vmCloudInit *VMCloudInit) (userData string, networkData string) {
	// userData
	userData = defaultCloudInitUserDataBasic
	if vmCloudInit.Password != "" && vmCloudInit.UserName != "" {
		userData += fmt.Sprintf(defaultCloudInitUserDataPasswordTemplate, vmCloudInit.UserName, vmCloudInit.Password)
	}
	if vmCloudInit.PublicKey != "" {
		userData += fmt.Sprintf(defaultCloudInitUserDataSSHKeyTemplate, vmCloudInit.PublicKey)
	}
	// networkData
	if vmCloudInit.InterfaceName == "" {
		return
	}
	networkData = fmt.Sprintf(defaultCloudInitNetworkDataTemplate, vmCloudInit.InterfaceName)
	if vmCloudInit.Address != "" && vmCloudInit.Gateway != "" {
		networkData += fmt.Sprintf(defaultCloudInitNetworkDataStaticTemplate, vmCloudInit.Address, vmCloudInit.Gateway)
	} else {
		networkData += defaultCloudInitNetworkDataDHCPTemplate
	}
	return userData, networkData
}

func (v *VMBuilder) CloudInit(vmCloudInit *VMCloudInit) *VMBuilder {
	if vmCloudInit == nil {
		return v
	}
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
	userData, networkData := generateCloudInit(vmCloudInit)
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
