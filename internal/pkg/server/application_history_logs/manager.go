package application_history_logs

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-history-logs-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application_history_logs"
)

// Manager structure with the required providers for application history logs operations.
type Manager struct {
	AppHistoryLogsProvider application_history_logs.Provider
}

// NewManager creates a Manager using a provider.
func NewManager(appHistoryLogsProvider application_history_logs.Provider) Manager {
	return Manager{appHistoryLogsProvider}
}

func (m *Manager) Add(addLogRequest *entities.AddLogRequest) derrors.Error {
	aErr := m.AppHistoryLogsProvider.Add(addLogRequest)
	if aErr != nil {
		return aErr
	}
	return nil
}

func (m *Manager) Update(updateLogRequest *entities.UpdateLogRequest) derrors.Error {
	uErr := m.AppHistoryLogsProvider.Update(updateLogRequest)
	if uErr != nil {
		return uErr
	}
	return nil
}

func (m *Manager) Search(searchLogRequest *entities.SearchLogsRequest) (*entities.LogResponse, derrors.Error) {
	sErr, logResponse := m.AppHistoryLogsProvider.Search(searchLogRequest)
	if sErr != nil {
		return nil, sErr
	}
	return logResponse, nil
}

func (m *Manager) Remove(removeLogRequest *entities.RemoveLogRequest) derrors.Error {
	rErr := m.AppHistoryLogsProvider.Remove(removeLogRequest)
	if rErr != nil {
		return rErr
	}
	return nil
}

func (m *Manager) ExistServiceInstanceLog (addLogRequest *grpc_application_history_logs_go.AddLogRequest) (bool, derrors.Error) {
	return m.AppHistoryLogsProvider.ExistsServiceInstanceLog(addLogRequest.OrganizationId, addLogRequest.AppInstanceId, addLogRequest.Created, addLogRequest.ServiceInstanceId)
}