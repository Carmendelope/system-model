/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-inventory-go"
	"time"
)

// OperatingSystemInfo contains information about the operating system of an asset. Notice that no
// enums have been used for either the name or the version as to no constraint the elements of the
// inventory even if we do not have an agent running supporting those.
type OperatingSystemInfo struct {
	// Name of the operating system. Expecting full name.
	Name string `json:"name,omitempty"`
	// Version installed.
	Version              string   `json:"version,omitempty"`
}

func NewOperatingSystemInfoFromGRPC(osInfo *grpc_inventory_go.OperatingSystemInfo) *OperatingSystemInfo{
	if osInfo == nil{
		return nil
	}
	return &OperatingSystemInfo{
		Name:    osInfo.Name,
		Version: osInfo.Version,
	}
}

func (os * OperatingSystemInfo) ToGRPC() *grpc_inventory_go.OperatingSystemInfo{
	if os == nil{
		return nil
	}
	return &grpc_inventory_go.OperatingSystemInfo{
		Name:                 os.Name,
		Version:              os.Version,
	}
}

// HardareInfo contains information about the hardware of an asset.
type HardwareInfo struct {
	// CPUs contains the list of CPU available in a given asset.
	Cpus []*CPUInfo `json:"cpus,omitempty"`
	// InstalledRam contains the total RAM available in MB.
	InstalledRam int64 `json:"installed_ram,omitempty"`
	// NetInterfaces with the list of networking cards.
	NetInterfaces        []*NetworkingHardwareInfo `json:"net_interfaces,omitempty"`
}

func NewHardwareInfoFromGRPC(hardwareInfo * grpc_inventory_go.HardwareInfo) * HardwareInfo{
	if hardwareInfo == nil{
		return nil
	}
	cpus := make([]*CPUInfo, 0)
	for _, info := range hardwareInfo.Cpus{
		cpus = append(cpus, NewCPUInfoFromGRPC(info))
	}
	netCards := make([]*NetworkingHardwareInfo, 0)
	for _, net:= range hardwareInfo.NetInterfaces{
		netCards = append(netCards, NewNetworkingHardwareInfoFromGRPC(net))
	}
	return &HardwareInfo{
		Cpus:          cpus,
		InstalledRam:  hardwareInfo.InstalledRam,
		NetInterfaces: netCards,
	}
}

func (hi * HardwareInfo) ToGRPC() *grpc_inventory_go.HardwareInfo{
	if hi == nil{
		return nil
	}
	cpus := make([]*grpc_inventory_go.CPUInfo, 0)
	for _, info := range hi.Cpus{
		cpus = append(cpus, info.ToGRPC())
	}
	netCards := make([]*grpc_inventory_go.NetworkingHardwareInfo, 0)
	for _, net := range hi.NetInterfaces{
		netCards = append(netCards, net.ToGRPC())
	}
	return &grpc_inventory_go.HardwareInfo{
		Cpus:                 cpus,
		InstalledRam:         hi.InstalledRam,
		NetInterfaces:        netCards,
	}
}

// CPUInfo contains information about a CPU.
type CPUInfo struct {
	// Manufacturer of the CPU.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Model of the CPU.
	Model string `json:"model,omitempty"`
	// Architecture of the CPU.
	Architecture string `json:"architecture,omitempty"`
	// NumCores with the number of cores.
	NumCores             int32    `json:"num_cores,omitempty"`
}

func NewCPUInfoFromGRPC(cpu * grpc_inventory_go.CPUInfo) * CPUInfo{
	if cpu == nil{
		return nil
	}
	return &CPUInfo{
		Manufacturer: cpu.Manufacturer,
		Model:       cpu.Model,
		Architecture: cpu.Architecture,
		NumCores:     cpu.NumCores,
	}
}

func (ci * CPUInfo) ToGRPC() *grpc_inventory_go.CPUInfo{
	if ci == nil{
		return nil
	}
	return &grpc_inventory_go.CPUInfo{
		Manufacturer:         ci.Manufacturer,
		Model:                ci.Model,
		Architecture:         ci.Architecture,
		NumCores:             ci.NumCores,
	}
}

// NetworkingHardwareInfo with the list of network interfaces that are available.
type NetworkingHardwareInfo struct {
	// Type of network interface. For example, ethernet, wifi, infiniband, etc.
	Type string `json:"type,omitempty"`
	// Link capacity in Mbps.
	LinkCapacity         int64    `json:"link_capacity,omitempty"`
}

func NewNetworkingHardwareInfoFromGRPC(netInfo * grpc_inventory_go.NetworkingHardwareInfo) * NetworkingHardwareInfo{
	if netInfo == nil{
		return nil
	}
	return &NetworkingHardwareInfo{
		Type:         netInfo.Type,
		LinkCapacity: netInfo.LinkCapacity,
	}
}

func (nhi * NetworkingHardwareInfo) ToGRPC() *grpc_inventory_go.NetworkingHardwareInfo{
	if nhi == nil{
		return nil
	}
	return &grpc_inventory_go.NetworkingHardwareInfo{
		Type:                 nhi.Type,
		LinkCapacity:         nhi.LinkCapacity,
	}
}

// StorageHardwareInfo with the storage related information.
type StorageHardwareInfo struct {
	// Type of storage.
	Type string `json:"type,omitempty"`
	// Total capacity in MB.
	TotalCapacity        int64    `json:"total_capacity,omitempty"`
}

func NewStorageHardwareInfoFromGRPC(storage * grpc_inventory_go.StorageHardwareInfo) * StorageHardwareInfo{
	if storage == nil{
		return nil
	}
	return &StorageHardwareInfo{
		Type:          storage.Type,
		TotalCapacity: storage.TotalCapacity,
	}
}

func (shi * StorageHardwareInfo) ToGRPC() *grpc_inventory_go.StorageHardwareInfo{
	if shi == nil{
		return nil
	}
	return &grpc_inventory_go.StorageHardwareInfo{
		Type:                 shi.Type,
		TotalCapacity:        shi.TotalCapacity,
	}
}

// Asset represents an element in the network from which we register some type of information. Example of
// assets could be workstations, nodes in a cluster, or other type of hardware.
type Asset struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AssetId with the asset identifier.
	AssetId string `json:"asset_id,omitempty"`
	// AgentId with the agent identifier that is monitoring this asset if any.
	AgentId string `json:"agent_id,omitempty"`
	// Show flag to determine if this asset should be shown on the UI. This flag is internally used
	// for the async uninstall/removal of the asset.
	Show bool `json:"show,omitempty"`
	// Created time
	Created int64 `json:"created,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// OS contains Operating System information.
	Os *OperatingSystemInfo `json:"os,omitempty"`
	// Hardware information.
	Hardware *HardwareInfo `json:"hardware,omitempty"`
	// Storage information.
	Storage *StorageHardwareInfo `json:"storage,omitempty"`
	// EicNetIp contains the current IP address that connects the asset to the EIC.
	EicNetIp             string   `json:"eic_net_ip,omitempty"`
}

func NewAssetFromGRPC(addRequest * grpc_inventory_go.AddAssetRequest) *Asset{
	return &Asset{
		OrganizationId: addRequest.OrganizationId,
		AssetId:        GenerateUUID(),
		AgentId:        addRequest.AgentId,
		Show:           true,
		Created:        time.Now().Unix(),
		Labels:         addRequest.Labels,
		Os:             NewOperatingSystemInfoFromGRPC(addRequest.Os),
		Hardware:       NewHardwareInfoFromGRPC(addRequest.Hardware),
		Storage:        NewStorageHardwareInfoFromGRPC(addRequest.Storage),
	}
}

func (a * Asset) ToGRPC() *grpc_inventory_go.Asset{
	return &grpc_inventory_go.Asset{
		OrganizationId:       a.OrganizationId,
		AssetId:              a.AssetId,
		AgentId:              a.AgentId,
		Show:                 a.Show,
		Created:              a.Created,
		Labels:               a.Labels,
		Os:                   a.Os.ToGRPC(),
		Hardware:             a.Hardware.ToGRPC(),
		Storage:              a.Storage.ToGRPC(),
		EicNetIp:             a.EicNetIp,
	}
}