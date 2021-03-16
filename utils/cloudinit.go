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
      netmask: %s
      gateway: %s`
)

type VMCloudInit struct {
	UserName      string
	Password      string
	PublicKey     string
	InterfaceName string
	Address       string
	NetMask       string
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

	if vmCloudInit.Address != "" && vmCloudInit.Gateway != "" && vmCloudInit.NetMask != "" {
		networkData += fmt.Sprintf(defaultCloudInitNetworkDataStaticTemplate, vmCloudInit.Address, vmCloudInit.NetMask, vmCloudInit.Gateway)
	} else {
		networkData += defaultCloudInitNetworkDataDHCPTemplate
	}
	return userData, networkData
}

func (v *VMBuilder) CloudInit(vmCloudInit *VMCloudInit) *VMBuilder {
	if vmCloudInit == nil {
		return v
	}
	userData, networkData := generateCloudInit(vmCloudInit)
	diskName := "cloudinitdisk"
	diskBus := "virtio"
	// Disks
	var (
		diskExist bool
		diskIndex int
	)
	disks := v.vm.Spec.Template.Spec.Domain.Devices.Disks
	for i, disk := range disks {
		if disk.Name == diskName {
			diskExist = true
			diskIndex = i
			break
		}
	}

	disk := kubevirtv1.Disk{
		Name: diskName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: diskBus,
			},
		},
	}
	if diskExist {
		disks[diskIndex] = disk
	} else {
		disks = append(disks, disk)
	}

	v.vm.Spec.Template.Spec.Domain.Devices.Disks = disks

	// Volumes
	var (
		volumeExist bool
		volumeIndex int
	)
	volumes := v.vm.Spec.Template.Spec.Volumes
	for i, volume := range volumes {
		if volume.Name == diskName {
			volumeExist = true
			volumeIndex = i
			break
		}
	}
	volume := kubevirtv1.Volume{
		Name: diskName,
		VolumeSource: kubevirtv1.VolumeSource{
			CloudInitNoCloud: &kubevirtv1.CloudInitNoCloudSource{
				UserData:    userData,
				NetworkData: networkData,
			},
		},
	}
	if volumeExist {
		volumes[volumeIndex] = volume
	} else {
		volumes = append(volumes, volume)
	}
	v.vm.Spec.Template.Spec.Volumes = volumes
	return v
}
