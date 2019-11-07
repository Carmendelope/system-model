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

package project

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const ProjectTable = "Project"

var allProjectColumns = []string{"owner_account_id", "project_id", "name", "created", "state", "state_info"}
var allProjectColumnsNoPK = []string{"name", "created", "state", "state_info"}

type ScyllaProjectProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaProjectProvider(address string, port int, keyspace string) *ScyllaProjectProvider {
	provider := ScyllaProjectProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

// disconnect from the database
func (sp *ScyllaProjectProvider) Disconnect() {
	sp.Lock()
	defer sp.Unlock()
	sp.ScyllaDB.Disconnect()
}

//
func (sp *ScyllaProjectProvider) createPKMap(accountID string, projectID string) map[string]interface{} {

	res := map[string]interface{}{
		"owner_account_id": accountID,
		"project_id":       projectID,
	}

	return res
}

// ------------------------------------------------------------------------------------------------
// This provider is for a table that has a PK with two fields, we use scylladb composite methods
// ------------------------------------------------------------------------------------------------
// Add a new project to the system.
func (sp *ScyllaProjectProvider) Add(project entities.Project) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createPKMap(project.OwnerAccountId, project.ProjectId)

	return sp.UnsafeCompositeAdd(ProjectTable, pkColumn, allProjectColumns, project)
}

// Update the information of a project.
func (sp *ScyllaProjectProvider) Update(project entities.Project) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createPKMap(project.OwnerAccountId, project.ProjectId)

	return sp.UnsafeCompositeUpdate(ProjectTable, pkColumn, allProjectColumnsNoPK, project)
}

// Exists checks if a project exists on the system.
func (sp *ScyllaProjectProvider) Exists(accountID string, projectID string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createPKMap(accountID, projectID)

	return sp.UnsafeGenericCompositeExist(ProjectTable, pkColumn)
}

// check if there is a project in the account with the received name
func (sp *ScyllaProjectProvider) ExistsByName(accountID string, name string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	indexMap := map[string]interface{}{
		"owner_account_id": accountID,
		"name":             name,
	}

	return sp.UnsafeGenericCompositeExist(ProjectTable, indexMap)

	return true, nil
}

// Get a project.
func (sp *ScyllaProjectProvider) Get(accountID string, projectID string) (*entities.Project, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createPKMap(accountID, projectID)

	var project interface{} = &entities.Project{}

	err := sp.UnsafeCompositeGet(ProjectTable, pkColumn, allProjectColumns, &project)
	if err != nil {
		return nil, err
	}
	return project.(*entities.Project), nil
}

// Remove a project
func (sp *ScyllaProjectProvider) Remove(accountID string, projectID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createPKMap(accountID, projectID)

	return sp.UnsafeCompositeRemove(ProjectTable, pkColumn)

}

// List all the projects of an account
func (sp *ScyllaProjectProvider) ListAccountProjects(accountID string) ([]entities.Project, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ProjectTable).Columns(allProjectColumns...).Where(qb.Eq("owner_account_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"owner_account_id": accountID,
	})

	projects := make([]entities.Project, 0)
	cqlErr := gocqlx.Select(&projects, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list projects")
	}

	return projects, nil
}

// Clear all projects
func (sp *ScyllaProjectProvider) Clear() derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{ProjectTable})
}
