package application_history_logs

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-history-logs-go"
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

func (m *Manager) Add(addLogRequest *grpc_application_history_logs_go.AddLogRequest) derrors.Error {
	aErr := m.AppHistoryLogsProvider.Add(addLogRequest)
	if aErr != nil {
		return aErr
	}
	return nil
}

func (m *Manager) Update(updateLogRequest *grpc_application_history_logs_go.UpdateLogRequest) derrors.Error {
	uErr := m.AppHistoryLogsProvider.Update(updateLogRequest)
	if uErr != nil {
		return uErr
	}
	return nil
}

func (m *Manager) Search(searchLogRequest *grpc_application_history_logs_go.SearchLogRequest) (*grpc_application_history_logs_go.LogResponse, derrors.Error) {
	logResponse, sErr := m.AppHistoryLogsProvider.Search(searchLogRequest)
	if sErr != nil {
		return nil, sErr
	}
	return logResponse, nil
}

func (m *Manager) Remove(removeLogRequest *grpc_application_history_logs_go.RemoveLogsRequest) derrors.Error {
	rErr := m.AppHistoryLogsProvider.Remove(removeLogRequest)
	if rErr != nil {
		return rErr
	}
	return nil
}
