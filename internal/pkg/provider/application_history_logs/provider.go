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
	grpc_application_history_logs_go "github.com/nalej/grpc-application-history-logs-go"
)

// Provider for the application networking instances.
type Provider interface {
	Add(*grpc_application_history_logs_go.AddLogRequest) derrors.Error
	Update(*grpc_application_history_logs_go.UpdateLogRequest) derrors.Error
	Search(*grpc_application_history_logs_go.SearchLogRequest) (*grpc_application_history_logs_go.LogResponse, derrors.Error)
	Remove(request *grpc_application_history_logs_go.RemoveLogsRequest) derrors.Error
}
