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
 *
 */

package eic

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
	"time"
)

func CreateTestCPU() []*entities.CPUInfo {
	cpus := make([]*entities.CPUInfo, 0)
	size := rand.Intn(10) + 1
	for i := 0; i < size; i++ {
		cpus = append(cpus, &entities.CPUInfo{
			Manufacturer: fmt.Sprintf("manufacturer_%d", i),
			Model:        fmt.Sprintf("model_%d", i),
			Architecture: fmt.Sprintf("architecture_%d", i),
			NumCores:     2,
		})
	}
	return cpus
}

func CreateTestNetInterfaces() []*entities.NetworkingHardwareInfo {
	netCards := make([]*entities.NetworkingHardwareInfo, 0)
	size := rand.Intn(10) + 1
	for i := 0; i < size; i++ {
		netCards = append(netCards, &entities.NetworkingHardwareInfo{
			Type:         fmt.Sprintf("type_%d", i),
			LinkCapacity: 100,
		})
	}
	return netCards
}

func CreateTestEdgeController() *entities.EdgeController {
	id := rand.Intn(200)
	labels := make(map[string]string, 0)
	size := rand.Intn(10) + 1
	for i := 0; i < size; i++ {
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

	operationSummary := entities.ECOpSummary{
		OperationId: entities.GenerateUUID(),
		Timestamp:   time.Now().Unix(),
		Status:      entities.OpStatusInProgress,
		Info:        "operation summary info",
	}

	return &entities.EdgeController{
		OrganizationId:   fmt.Sprintf("organization_%d", id),
		EdgeControllerId: entities.GenerateUUID(),
		Show:             true,
		Created:          time.Now().Unix(),
		Name:             fmt.Sprintf("name_%d", id),
		Labels:           labels,
		Location: &entities.InventoryLocation{
			Geolocation: "geolocation",
			Geohash:     "geohash",
		},
		Os:           os,
		Hardware:     hardware,
		Storage:      []*entities.StorageHardwareInfo{&storage},
		LastOpResult: &operationSummary,
	}
}
