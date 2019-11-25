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
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for the application networking instances.
type Provider interface {
	// Add a new entry to the service instance history table
	Add(*entities.AddLogRequest) derrors.Error
	// Update an entry of the service instance history table
	Update(*entities.UpdateLogRequest) derrors.Error
	// Search for instances that were alive during a period defined in the request
	Search(*entities.SearchLogsRequest) (*entities.LogResponse, derrors.Error)
	// Remove an entry from the service instance history table
	Remove(*entities.RemoveLogRequest) derrors.Error

	// ExistsServiceInstanceLog checks if a ServiceInstanceLog exists
	ExistsServiceInstanceLog(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstanceId string) (bool, derrors.Error)

<<<<<<< HEAD
	// clear all application history logs
=======
	// clear all application history logs.
>>>>>>> a86c0292eaba4378d92c3f8e954f56c3995506a2
	Clear() derrors.Error
}
