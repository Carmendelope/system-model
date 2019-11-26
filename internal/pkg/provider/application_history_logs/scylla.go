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

package application_history_logs

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const ServiceInstanceHistoryTable = "Service_Instance_History"

var (
	ServiceInstanceHistoryColumns = []string{
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
	ServiceInstanceHistoryColumnsNoPK = []string{
		"created",
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

func NewScyllaApplicationHistoryLogsProvider(address string, port int, keyspace string) *ScyllaApplicationHistoryLogsProvider {
	provider := ScyllaApplicationHistoryLogsProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	_ = provider.Connect()
	return &provider
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Add(addLogRequest *entities.AddLogRequest) derrors.Error {
	sahlp.Lock()
	defer sahlp.Unlock()

	toAdd := entities.ServiceInstanceLog{
		OrganizationId:         addLogRequest.OrganizationId,
		AppDescriptorId:        addLogRequest.AppDescriptorId,
		AppInstanceId:          addLogRequest.AppInstanceId,
		ServiceGroupId:         addLogRequest.ServiceGroupId,
		ServiceGroupInstanceId: addLogRequest.ServiceGroupInstanceId,
		ServiceId:              addLogRequest.ServiceId,
		ServiceInstanceId:      addLogRequest.ServiceInstanceId,
		Created:                addLogRequest.Created,
		Terminated:             0,
	}

	pkComposite := sahlp.createServiceInstanceHistoryPKMap(addLogRequest.OrganizationId, addLogRequest.AppInstanceId, addLogRequest.ServiceInstanceId)
	return sahlp.UnsafeCompositeAdd(ServiceInstanceHistoryTable, pkComposite, ServiceInstanceHistoryColumns, toAdd)
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Update(updateLogRequest *entities.UpdateLogRequest) derrors.Error {
	sahlp.Lock()
	defer sahlp.Unlock()

	columns := []string{
		"terminated",
	}

	toUpdate := entities.ServiceInstanceLog{
		OrganizationId:    updateLogRequest.OrganizationId,
		AppInstanceId:     updateLogRequest.AppInstanceId,
		ServiceInstanceId: updateLogRequest.ServiceInstanceId,
		Terminated:        updateLogRequest.Terminated,
	}

	pkComposite := sahlp.createServiceInstanceHistoryPKMap(updateLogRequest.OrganizationId, updateLogRequest.AppInstanceId, updateLogRequest.ServiceInstanceId)
	return sahlp.UnsafeCompositeUpdate(ServiceInstanceHistoryTable, pkComposite, columns, toUpdate)
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Search(searchLogsRequest *entities.SearchLogsRequest) (*entities.LogResponse, derrors.Error) {
	sahlp.Lock()
	defer sahlp.Unlock()

	OrganizationIdMap := map[string]interface{}{
		"organization_id": searchLogsRequest.OrganizationId,
		"created":         searchLogsRequest.To,
	}

	result := make([]entities.ServiceInstanceLog, 0)

	// TODO: We should be able to perform this query without allowing filtering. It will involve changing the database design and probably adding an additional table
	sb := qb.Select(ServiceInstanceHistoryTable).Columns(ServiceInstanceHistoryColumns...).Where(qb.Eq("organization_id")).Where(qb.LtOrEq("created")).AllowFiltering()
	stmt, names := sb.ToCql()
	q := gocqlx.Query(sahlp.Session.Query(stmt), names).BindMap(OrganizationIdMap)
	qErr := q.SelectRelease(&result)
	if qErr != nil {
		return nil, derrors.NewGenericError("could not query database", qErr)
	}

	events := make([]entities.ServiceInstanceLog, 0)
	found := false
	for _, serviceInstanceLog := range result {
		if serviceInstanceLog.Terminated >= searchLogsRequest.From || serviceInstanceLog.Terminated == 0 {
			events = append(events, serviceInstanceLog)
			found = true
		}
	}

	if found {
		return &entities.LogResponse{
			OrganizationId: searchLogsRequest.OrganizationId,
			From:           searchLogsRequest.From,
			To:             searchLogsRequest.To,
			Events:         events,
		}, nil
	} else {
		return nil, derrors.NewNotFoundError("search log request").WithParams(searchLogsRequest)
	}
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Remove(removeLogRequest *entities.RemoveLogRequest) derrors.Error {
	sahlp.Lock()
	defer sahlp.Unlock()
	pkComposite := sahlp.createServiceInstanceHistoryAuxMap(removeLogRequest.OrganizationId, removeLogRequest.AppInstanceId)
	return sahlp.UnsafeCompositeRemove(ServiceInstanceHistoryTable, pkComposite)
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) ExistsServiceInstanceLog(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstanceId string) (bool, derrors.Error) {
	sahlp.Lock()
	defer sahlp.Unlock()
	pkComposite := sahlp.createServiceInstanceHistoryPKMap(organizationId, appInstanceId, serviceInstanceId)
	return sahlp.UnsafeGenericCompositeExist(ServiceInstanceHistoryTable, pkComposite)
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) Clear() derrors.Error {
	sahlp.Lock()
	defer sahlp.Unlock()

	if err := sahlp.UnsafeClear([]string{ServiceInstanceHistoryTable}); err != nil {
		return err
	}
	return nil
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) createServiceInstanceHistoryPKMap(organizationId string, appInstanceId string, serviceInstanceId string) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":     organizationId,
		"app_instance_id":     appInstanceId,
		"service_instance_id": serviceInstanceId,
	}
}

func (sahlp *ScyllaApplicationHistoryLogsProvider) createServiceInstanceHistoryAuxMap(organizationId string, appInstanceId string) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": organizationId,
		"app_instance_id": appInstanceId,
	}
}
