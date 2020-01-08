/*
 * Copyright 2020 Nalej
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

package cluster

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
)

func CreateTestCluster(clusterID string) *entities.Cluster {

	id := rand.Intn(200)

	labels := make(map[string]string, 0)
	tam := rand.Intn(10) + 1
	for i := 0; i < tam; i++ {
		labels[fmt.Sprintf("label-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	return &entities.Cluster{

		OrganizationId:             fmt.Sprintf("organization_%d", id),
		ClusterId:                  fmt.Sprintf("cluster_%s", clusterID),
		Name:                       fmt.Sprintf("name_%d", id),
		ClusterType:                entities.ClusterType(1),
		Hostname:                   fmt.Sprintf("host_%s", clusterID),
		ControlPlaneHostname:       fmt.Sprintf("cp_host_%s", clusterID),
		Multitenant:                entities.MultitenantSupport(2),
		Status:                     entities.ClusterStatus(1),
		Labels:                     labels,
		Cordon:                     true,
		State:                      entities.Provisioning,
		MillicoresConversionFactor: .5,
	}
}
