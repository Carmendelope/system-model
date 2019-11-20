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
	"context"
	"github.com/nalej/grpc-application-history-logs-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the application requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) Add(ctx context.Context, addLogRequest *grpc_application_history_logs_go.AddLogRequest) error {
	vErr := entities.ValidAddLogRequest(addLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid add log request")
		return conversions.ToGRPCError(vErr)
	}

	aErr := h.Manager.Add(addLogRequest)
	if aErr != nil {
		log.Error().Str("add error", aErr.DebugReport()).Msg("cannot add log")
		return conversions.ToGRPCError(aErr)
	}
	return nil
}

func (h *Handler) Update(ctx context.Context, updateLogRequest *grpc_application_history_logs_go.UpdateLogRequest) error {
	vErr := entities.ValidUpdateLogRequest(updateLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid update log request")
		return conversions.ToGRPCError(vErr)
	}

	uErr := h.Manager.Update(updateLogRequest)
	if uErr != nil {
		log.Error().Str("update error", uErr.DebugReport()).Msg("cannot update log")
		return conversions.ToGRPCError(uErr)
	}
	return nil
}

func (h *Handler) Search(ctx context.Context, searchLogRequest *grpc_application_history_logs_go.SearchLogRequest) (*grpc_application_history_logs_go.LogResponse, error) {
	vErr := entities.ValidSearchLogRequest(searchLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid add log request")
		return nil, conversions.ToGRPCError(vErr)
	}

	logResponse, sErr := h.Manager.Search(searchLogRequest)
	if sErr != nil {
		log.Error().Str("add error", sErr.DebugReport()).Msg("cannot add log")
		return nil, conversions.ToGRPCError(sErr)
	}
	return logResponse, nil
}

func (h *Handler) Remove(ctx context.Context, removeLogRequest *grpc_application_history_logs_go.RemoveLogsRequest) error {
	vErr := entities.ValidRemoveLogRequest(removeLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid add log request")
		return conversions.ToGRPCError(vErr)
	}

	rErr := h.Manager.Remove(removeLogRequest)
	if rErr != nil {
		log.Error().Str("add error", rErr.DebugReport()).Msg("cannot add log")
		return conversions.ToGRPCError(rErr)
	}
	return nil
}
