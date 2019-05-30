/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
	"time"
)

func CreateTestCPU() []*entities.CPUInfo{
	cpus := make([]*entities.CPUInfo, 0)
	size := rand.Intn(10) +1
	for i:=0; i<size; i++{
		cpus = append(cpus, &entities.CPUInfo{
			Manufacturer: fmt.Sprintf("manufacturer_%d", i),
			Model:        fmt.Sprintf("model_%d", i),
			Architecture: fmt.Sprintf("architecture_%d", i),
			NumCores:     2,
		})
	}
	return cpus
}

func CreateTestNetInterfaces() []*entities.NetworkingHardwareInfo{
	netCards := make([]*entities.NetworkingHardwareInfo, 0)
	size := rand.Intn(10) +1
	for i:=0; i<size; i++{
		netCards = append(netCards, &entities.NetworkingHardwareInfo{
			Type:         fmt.Sprintf("type_%d", i),
			LinkCapacity: 100,
		})
	}
	return netCards
}

func CreateTestAsset() * entities.Asset{
	id:= rand.Intn(200)
	labels := make (map[string]string, 0)
	size := rand.Intn(10) +1
	for i:=0; i<size; i++{
		labels[fmt.Sprintf("label-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	os := &entities.OperatingSystemInfo{
		Name:    "FakeOS",
		Version: "1.0",
	}
	hardware := &entities.HardwareInfo{
		Cpus:          CreateTestCPU(),
		InstalledRam:  100,
		NetInterfaces: CreateTestNetInterfaces(),
	}
	storage := entities.StorageHardwareInfo{
		Type:          "FakeStorage",
		TotalCapacity: 100,
	}
	return &entities.Asset{
		OrganizationId: fmt.Sprintf("organization_%d", id),
		EdgeControllerId: entities.GenerateUUID(),
		AssetId:        entities.GenerateUUID(),
		AgentId:        fmt.Sprintf("agent_%d", id),
		Show:           true,
		Created:        time.Now().Unix(),
		Labels:         labels,
		Os:             os,
		Hardware:       hardware,
		Storage:        []entities.StorageHardwareInfo{storage},
		EicNetIp:       "1.1.1.1",
	}
}