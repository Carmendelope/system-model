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
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupApplicationHistoryLogsProvider struct {
	sync.Mutex
	// serviceInstanceLogs indexed by organizationId
	serviceInstanceLogs map[string][]*entities.ServiceInstanceLog
}

func (m *MockupApplicationHistoryLogsProvider) ExistsServiceInstanceLog(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstanceId string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	return m.unsafeExistsServiceInstanceLog(organizationId, appInstanceId, serviceGroupInstanceId, serviceInstanceId)
}

func (m *MockupApplicationHistoryLogsProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.serviceInstanceLogs = make(map[string][]*entities.ServiceInstanceLog, 0)
	return nil
}

func NewMockupApplicationHistoryLogsProvider() *MockupApplicationHistoryLogsProvider {
	return &MockupApplicationHistoryLogsProvider{
		serviceInstanceLogs: make(map[string][]*entities.ServiceInstanceLog, 0),
	}
}

func (m *MockupApplicationHistoryLogsProvider) Add(addLogRequest *entities.AddLogRequest) derrors.Error {
	m.Lock()
	defer m.Unlock()

	toAdd := AddLogRequestToServiceInstanceLog(*addLogRequest)
	_, err := m.unsafeExistsServiceInstanceLog(addLogRequest.OrganizationId, addLogRequest.AppInstanceId, addLogRequest.ServiceGroupInstanceId, addLogRequest.ServiceInstanceId)
	if err != nil {
		m.serviceInstanceLogs[addLogRequest.OrganizationId] = append(m.serviceInstanceLogs[addLogRequest.OrganizationId], &toAdd)
	} else {
		return err
	}

	return nil
}

func (m *MockupApplicationHistoryLogsProvider) Update(updateLogRequest *entities.UpdateLogRequest) derrors.Error {
	m.Lock()
	defer m.Unlock()

	list, exists := m.serviceInstanceLogs[updateLogRequest.OrganizationId]
	if exists {
		found := false
		newLogs := make([]*entities.ServiceInstanceLog, len(list))
		for i, serviceInstanceLog := range list {
			if serviceInstanceLog.AppInstanceId == updateLogRequest.AppInstanceId {
				toAdd := serviceInstanceLog
				toAdd.Terminated = updateLogRequest.Terminated
				newLogs[i] = toAdd
				found = true
			} else {
				newLogs[i] = serviceInstanceLog
			}
		}
		if !found {
			return derrors.NewNotFoundError("app instance id").WithParams(updateLogRequest.AppInstanceId)
		}
		m.serviceInstanceLogs[updateLogRequest.OrganizationId] = newLogs
	} else {
		return derrors.NewNotFoundError("organization id").WithParams(updateLogRequest.OrganizationId)
	}

	return nil
}

func (m *MockupApplicationHistoryLogsProvider) Search(searchLogsRequest *entities.SearchLogsRequest) (derrors.Error, *entities.LogResponse) {
	m.Lock()
	defer m.Unlock()

	var events *[]entities.ServiceInstanceLog = nil
	list, exists := m.serviceInstanceLogs[searchLogsRequest.OrganizationId]
	if !exists {
		return derrors.NewNotFoundError("organization id").WithParams(searchLogsRequest.OrganizationId), nil
	}

	for _, serviceInstanceLog := range list {
		if (serviceInstanceLog.OrganizationId == searchLogsRequest.OrganizationId && serviceInstanceLog.Created >= searchLogsRequest.From) || (serviceInstanceLog.OrganizationId == searchLogsRequest.OrganizationId && serviceInstanceLog.Terminated >= searchLogsRequest.To) {
			*events = append(*events, *serviceInstanceLog)
		}
	}

	return nil, &entities.LogResponse{
		OrganizationId: searchLogsRequest.OrganizationId,
		From:           searchLogsRequest.From,
		To:             searchLogsRequest.To,
		Events:         *events,
	}
}

func (m *MockupApplicationHistoryLogsProvider) Remove(removeLogRequest *entities.RemoveLogRequest) derrors.Error {
	m.Lock()
	defer m.Unlock()

	list, exists := m.serviceInstanceLogs[removeLogRequest.OrganizationId]
	if exists {
		found := false
		newLogs := make([]*entities.ServiceInstanceLog, len(list))
		for i, serviceInstanceLog := range list {
			if serviceInstanceLog.AppInstanceId == removeLogRequest.AppInstanceId {
				found = true
			} else {
				newLogs[i] = serviceInstanceLog
			}
		}
		if !found {
			return derrors.NewNotFoundError("app instance id").WithParams(removeLogRequest.AppInstanceId)
		}
		m.serviceInstanceLogs[removeLogRequest.OrganizationId] = newLogs
	} else {
		return derrors.NewNotFoundError("organization id").WithParams(removeLogRequest.OrganizationId)
	}

	return nil
}

func (m *MockupApplicationHistoryLogsProvider) unsafeExistsServiceInstanceLog(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstanceId string) (bool, derrors.Error) {
	list, exists := m.serviceInstanceLogs[organizationId]
	if exists {
		for _, serviceInstanceLog := range list {
			if serviceInstanceLog.AppInstanceId == appInstanceId && serviceInstanceId == serviceInstanceLog.ServiceInstanceId && serviceGroupInstanceId == serviceInstanceLog.ServiceGroupInstanceId {
				return true, nil
			}
		}
	} else {
		return false, derrors.NewNotFoundError("organization id").WithParams(organizationId)
	}
	return false, nil
}

func AddLogRequestToServiceInstanceLog(addLogRequest entities.AddLogRequest) entities.ServiceInstanceLog {
	return entities.ServiceInstanceLog{
		OrganizationId:         addLogRequest.OrganizationId,
		AppDescriptorId:        addLogRequest.AppDescriptorId,
		AppInstanceId:          addLogRequest.AppInstanceId,
		ServiceGroupId:         addLogRequest.ServiceGroupId,
		ServiceGroupInstanceId: addLogRequest.ServiceGroupInstanceId,
		ServiceId:              addLogRequest.ServiceId,
		ServiceInstanceId:      addLogRequest.ServiceInstanceId,
		Created:                addLogRequest.Created,
	}
}