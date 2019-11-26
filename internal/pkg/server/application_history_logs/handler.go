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
	"context"
	"github.com/nalej/grpc-application-history-logs-go"
	"github.com/nalej/grpc-common-go"
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

	entALR := entities.ToAddLogRequest(*addLogRequest)
	aErr := h.Manager.Add(&entALR)
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

	entULR := entities.ToUpdateLogRequest(*updateLogRequest)
	uErr := h.Manager.Update(&entULR)
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

	entSLR := entities.ToSearchLogsRequest(*searchLogRequest)
	logResponse, sErr := h.Manager.Search(&entSLR)
	if sErr != nil {
		log.Error().Str("search error", sErr.DebugReport()).Msg("cannot search log")
		return nil, conversions.ToGRPCError(sErr)
	}
	grpcLR := entities.ToGRPCLogRequest(*logResponse)
	return &grpcLR, nil
}

func (h *Handler) Remove(ctx context.Context, removeLogRequest *grpc_application_history_logs_go.RemoveLogsRequest) error {
	vErr := entities.ValidRemoveLogRequest(removeLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid remove log request")
		return conversions.ToGRPCError(vErr)
	}

	entRLR := entities.ToRemoveLogRequest(*removeLogRequest)
	rErr := h.Manager.Remove(&entRLR)
	if rErr != nil {
		log.Error().Str("remove error", rErr.DebugReport()).Msg("cannot remove log")
		return conversions.ToGRPCError(rErr)
	}
	return nil
}

func (h *Handler) ExistServiceInstanceLog(ctx context.Context, addLogRequest *grpc_application_history_logs_go.AddLogRequest) (*grpc_common_go.Exists, error) {
	vErr := entities.ValidAddLogRequest(addLogRequest)
	if vErr != nil {
		log.Error().Str("validation error", vErr.DebugReport()).Msg("invalid request")
		return &grpc_common_go.Exists{Exists: false}, conversions.ToGRPCError(vErr)
	}
	exists, err := h.Manager.ExistServiceInstanceLog(addLogRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot determine if the service instance log exists")
		return &grpc_common_go.Exists{Exists: false}, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Exists{Exists: exists}, nil
}
