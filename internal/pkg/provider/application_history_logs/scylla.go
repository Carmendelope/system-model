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

package application_history_logs

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

var (
	ServiceInstanceLogsColumns = []string{
		"organization_id",
		"app_instance_id",
		"created",
		"service_instance_id",
		"terminated",
		"app_descriptor_id",
		"service_group_id",
		"service_group_instance_id",
		"service_id",
	}
)

type ScyllaApplicationHistoryLogsProvider struct {
	sync.Mutex
	scylladb.ScyllaDB
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Add (addLogRequest entities.AddLogRequest) derrors.Error {
	sahlp.Lock()
	defer sahlp.Unlock()



	return nil
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Update (addLogRequest entities.AddLogRequest) derrors.Error {
	return nil
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Search (addLogRequest entities.AddLogRequest) (derrors.Error, *entities.LogResponse) {
 	return nil, nil
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Remove (addLogRequest entities.AddLogRequest) derrors.Error {
	return nil
}