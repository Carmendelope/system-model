/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package entities

// This file contains all the common structures of the edge-controller and the agents
// OpSummary
// Hw
// Storage
// OS info

import "github.com/nalej/grpc-inventory-go"

// Enumerate with the type of instances we can deploy in the system.
type OpStatus int32

const (
	OpStatusScheduled OpStatus = iota + 1
	OpStatusInProgress
	OpStatusSuccess
	OpStatusFail
	OpStatusCanceled
)

var OpStatusToGRPC = map[OpStatus]grpc_inventory_go.OpStatus{
	OpStatusScheduled:  grpc_inventory_go.OpStatus_SCHEDULED,
	OpStatusInProgress: grpc_inventory_go.OpStatus_INPROGRESS,
	OpStatusSuccess:    grpc_inventory_go.OpStatus_SUCCESS,
	OpStatusFail:       grpc_inventory_go.OpStatus_FAIL,
	OpStatusCanceled:   grpc_inventory_go.OpStatus_CANCELED,
}

var OpStatusFromGRPC = map[grpc_inventory_go.OpStatus]OpStatus{
	grpc_inventory_go.OpStatus_SCHEDULED:  OpStatusScheduled,
	grpc_inventory_go.OpStatus_INPROGRESS: OpStatusInProgress,
	grpc_inventory_go.OpStatus_SUCCESS:    OpStatusSuccess,
	grpc_inventory_go.OpStatus_FAIL:       OpStatusFail,
	grpc_inventory_go.OpStatus_CANCELED:   OpStatusCanceled,
}

type OperatingSystemClass int

const (
	LINUX = iota + 1
	WINDOWS
	DARWIN
)

var OperatingSystemClassToGRPC = map[OperatingSystemClass]grpc_inventory_go.OperatingSystemClass{
	LINUX:   grpc_inventory_go.OperatingSystemClass_LINUX,
	WINDOWS: grpc_inventory_go.OperatingSystemClass_WINDOWS,
	DARWIN:  grpc_inventory_go.OperatingSystemClass_DARWIN,
}
var OperatingSystemClassFromGRPC = map[grpc_inventory_go.OperatingSystemClass]OperatingSystemClass{
	grpc_inventory_go.OperatingSystemClass_LINUX:   LINUX,
	grpc_inventory_go.OperatingSystemClass_WINDOWS: WINDOWS,
	grpc_inventory_go.OperatingSystemClass_DARWIN:  DARWIN,
}

// OperatingSystemInfo contains information about the operating system of an asset. Notice that no
// enums have been used for either the name or the version as to no constraint the elements of the
// inventory even if we do not have an agent running supporting those.
type OperatingSystemInfo struct {
	// Name of the operating system. Expecting full name.
	Name string `json:"name,omitempty" cql:"name"`
	// Version installed.
	Version string `json:"version,omitempty" cql:"version"`
	// Class of the operating system - determines the binary format together with architecture
	Class OperatingSystemClass `json:"class,omitempty" cql:"op_class"`
	// Architecture of the OS.
	Architecture string `json:"architecture,omitempty" cql:"architecture"`
}

func NewOperatingSystemInfoFromGRPC(osInfo *grpc_inventory_go.OperatingSystemInfo) *OperatingSystemInfo {
	if osInfo == nil {
		return nil
	}
	return &OperatingSystemInfo{
		Name:         osInfo.Name,
		Version:      osInfo.Version,
		Class:        OperatingSystemClassFromGRPC[osInfo.Class],
		Architecture: osInfo.Architecture,
	}
}

func (os *OperatingSystemInfo) ToGRPC() *grpc_inventory_go.OperatingSystemInfo {
	if os == nil {
		return nil
	}
	return &grpc_inventory_go.OperatingSystemInfo{
		Name:         os.Name,
		Version:      os.Version,
		Class:        OperatingSystemClassToGRPC[os.Class],
		Architecture: os.Architecture,
	}
}

// HardareInfo contains information about the hardware of an asset.
type HardwareInfo struct {
	// CPUs contains the list of CPU available in a given asset.
	Cpus []*CPUInfo `json:"cpus,omitempty" cql:"cpus"`
	// InstalledRam contains the total RAM available in MB.
	InstalledRam int64 `json:"installed_ram,omitempty" cql:"installed_ram"`
	// NetInterfaces with the list of networking cards.
	NetInterfaces []*NetworkingHardwareInfo `json:"net_interfaces,omitempty" cql:"net_interfaces"`
}

func NewHardwareInfoFromGRPC(hardwareInfo *grpc_inventory_go.HardwareInfo) *HardwareInfo {
	if hardwareInfo == nil {
		return nil
	}
	cpus := make([]*CPUInfo, 0)
	for _, info := range hardwareInfo.Cpus {
		cpus = append(cpus, NewCPUInfoFromGRPC(info))
	}
	netCards := make([]*NetworkingHardwareInfo, 0)
	for _, net := range hardwareInfo.NetInterfaces {
		netCards = append(netCards, NewNetworkingHardwareInfoFromGRPC(net))
	}
	return &HardwareInfo{
		Cpus:          cpus,
		InstalledRam:  hardwareInfo.InstalledRam,
		NetInterfaces: netCards,
	}
}

func (hi *HardwareInfo) ToGRPC() *grpc_inventory_go.HardwareInfo {
	if hi == nil {
		return nil
	}
	cpus := make([]*grpc_inventory_go.CPUInfo, 0)
	for _, info := range hi.Cpus {
		cpus = append(cpus, info.ToGRPC())
	}
	netCards := make([]*grpc_inventory_go.NetworkingHardwareInfo, 0)
	for _, net := range hi.NetInterfaces {
		netCards = append(netCards, net.ToGRPC())
	}
	return &grpc_inventory_go.HardwareInfo{
		Cpus:          cpus,
		InstalledRam:  hi.InstalledRam,
		NetInterfaces: netCards,
	}
}

// CPUInfo contains information about a CPU.
type CPUInfo struct {
	// Manufacturer of the CPU.
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	// Model of the CPU.
	Model string `json:"model,omitempty" cql:"model"`
	// Architecture of the CPU.
	Architecture string `json:"architecture,omitempty" cql:"architecture"`
	// NumCores with the number of cores.
	NumCores int32 `json:"num_cores,omitempty" cql:"num_cores"`
}

func NewCPUInfoFromGRPC(cpu *grpc_inventory_go.CPUInfo) *CPUInfo {
	if cpu == nil {
		return nil
	}
	return &CPUInfo{
		Manufacturer: cpu.Manufacturer,
		Model:        cpu.Model,
		Architecture: cpu.Architecture,
		NumCores:     cpu.NumCores,
	}
}

func (ci *CPUInfo) ToGRPC() *grpc_inventory_go.CPUInfo {
	if ci == nil {
		return nil
	}
	return &grpc_inventory_go.CPUInfo{
		Manufacturer: ci.Manufacturer,
		Model:        ci.Model,
		Architecture: ci.Architecture,
		NumCores:     ci.NumCores,
	}
}

// NetworkingHardwareInfo with the list of network interfaces that are available.
type NetworkingHardwareInfo struct {
	// Type of network interface. For example, ethernet, wifi, infiniband, etc.
	Type string `json:"type,omitempty" cql:"type"`
	// Link capacity in Mbps.
	LinkCapacity int64 `json:"link_capacity,omitempty" cql:"link_capacity"`
}

func NewNetworkingHardwareInfoFromGRPC(netInfo *grpc_inventory_go.NetworkingHardwareInfo) *NetworkingHardwareInfo {
	if netInfo == nil {
		return nil
	}
	return &NetworkingHardwareInfo{
		Type:         netInfo.Type,
		LinkCapacity: netInfo.LinkCapacity,
	}
}

func (nhi *NetworkingHardwareInfo) ToGRPC() *grpc_inventory_go.NetworkingHardwareInfo {
	if nhi == nil {
		return nil
	}
	return &grpc_inventory_go.NetworkingHardwareInfo{
		Type:         nhi.Type,
		LinkCapacity: nhi.LinkCapacity,
	}
}

// StorageHardwareInfo with the storage related information.
type StorageHardwareInfo struct {
	// Type of storage.
	Type string `json:"type,omitempty" cql:"type"`
	// Total capacity in MB.
	TotalCapacity int64 `json:"total_capacity,omitempty" cql:"total_capacity"`
}

func NewStorageHardwareInfoFromGRPC(storage *grpc_inventory_go.StorageHardwareInfo) *StorageHardwareInfo {
	if storage == nil {
		return nil
	}
	return &StorageHardwareInfo{
		Type:          storage.Type,
		TotalCapacity: storage.TotalCapacity,
	}
}

func (shi *StorageHardwareInfo) ToGRPC() *grpc_inventory_go.StorageHardwareInfo {
	if shi == nil {
		return nil
	}
	return &grpc_inventory_go.StorageHardwareInfo{
		Type:          shi.Type,
		TotalCapacity: shi.TotalCapacity,
	}
}
