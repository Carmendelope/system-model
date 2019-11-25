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

package node

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new node to the system.
	Add(node entities.Node) derrors.Error
	// Update an existing node in the system
	Update(node entities.Node) derrors.Error
	// Exists checks if a node exists on the system.
	Exists(nodeID string) (bool, derrors.Error)
	// Get a node.
	Get(nodeID string) (*entities.Node, derrors.Error)
	// Remove a node
	Remove(nodeID string) derrors.Error
	// Clear nodes
	Clear() derrors.Error
}
