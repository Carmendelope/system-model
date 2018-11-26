/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
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
	Get(nodeID string) (* entities.Node, derrors.Error)
	// Remove a node
	Remove(nodeID string) derrors.Error
	// Clear nodes
	Clear() derrors.Error
}
