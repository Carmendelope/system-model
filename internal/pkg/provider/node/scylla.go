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
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const nodeTable = "nodes"
const nodeTablePK = "node_id"
const rowNotFound = "not found"

type ScyllaNodeProvider struct {
	Address  string
	Port     int
	Keyspace string
	Session  *gocql.Session
	sync.Mutex
}

func NewScyllaNodeProvider(address string, port int, keyspace string) *ScyllaNodeProvider {
	provider := ScyllaNodeProvider{Address: address, Port: port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider

}

// connect to the database
func (sp *ScyllaNodeProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaNodeProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}

// disconnect from the database
func (sp *ScyllaNodeProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

func (sp *ScyllaNodeProvider) unsafeExists(nodeID string) (bool, derrors.Error) {

	var returnedId string

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(nodeTable).Columns(nodeTablePK).Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		nodeTablePK: nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if node exists")
		}
	}

	return true, nil
}

// check that the session is created
func (sp *ScyllaNodeProvider) checkConnection() derrors.Error {
	if sp.Session == nil {
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaNodeProvider) checkAndConnect() derrors.Error {

	err := sp.checkConnection()
	if err != nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err = sp.connect()
		if err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// Add a new node to the system.
func (sp *ScyllaNodeProvider) Add(node entities.Node) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(node.NodeId)

	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(node.NodeId)
	}

	// insert a user

	stmt, names := qb.Insert(nodeTable).Columns("organization_id", "cluster_id", "node_id", "ip", "labels", "status", "state").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(node)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add node")
	}

	return nil
}

// Update an existing node in the system
func (sp *ScyllaNodeProvider) Update(node entities.Node) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(node.NodeId)

	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError(node.NodeId)
	}

	// update a user
	stmt, names := qb.Update(nodeTable).Set("organization_id", "cluster_id", "ip", "labels", "status", "state").Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(node)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update node")
	}

	return nil
}

// Exists checks if a node exists on the system.
func (sp *ScyllaNodeProvider) Exists(nodeID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var returnedId string

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(nodeTable).Columns(nodeTablePK).Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		nodeTablePK: nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if node exists")
		}
	}

	return true, nil
}

// Get a node.
func (sp *ScyllaNodeProvider) Get(nodeID string) (*entities.Node, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	var node entities.Node
	stmt, names := qb.Select(nodeTable).Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		nodeTablePK: nodeID,
	})

	err := q.GetRelease(&node)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError(nodeID)
		} else {
			return nil, derrors.AsError(err, "cannot get node")
		}
	}

	return &node, nil
}

// Remove a node
func (sp *ScyllaNodeProvider) Remove(nodeID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(nodeID)

	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("node").WithParams(nodeID)
	}

	// remove a user
	stmt, _ := qb.Delete(nodeTable).Where(qb.Eq(nodeTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, nodeID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot remove node")
	}

	return nil
}

func (sp *ScyllaNodeProvider) Clear() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE Nodes").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return derrors.AsError(err, "cannot truncate node table")
	}

	return nil
}
