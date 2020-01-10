/*
 * Copyright 2020 Nalej
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

package organization_setting

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const organizationSettingTable = "OrganizationSetting"

var organizationSettingTableColumns = []string{"organization_id", "key", "value", "description"}
var organizationSettingTableColumnsNoPK = []string{"value", "description"}

type ScyllaOrganizationSettingProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaOrganizationSettingProvider(address string, port int, keyspace string) *ScyllaOrganizationSettingProvider {
	provider := ScyllaOrganizationSettingProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

func (s *ScyllaOrganizationSettingProvider) createPKMap(OrganizationID string, key string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"key":             key,
	}

	return res
}

// Add a new setting for an organization.
func (s *ScyllaOrganizationSettingProvider) Add(setting entities.OrganizationSetting) derrors.Error {

	s.Lock()
	defer s.Unlock()

	pk := s.createPKMap(setting.OrganizationId, setting.Key)

	return s.UnsafeCompositeAdd(organizationSettingTable, pk, organizationSettingTableColumns, setting)
}

// Check if a setting is defined for an organization
func (s *ScyllaOrganizationSettingProvider) Exists(organizationID string, key string) (bool, derrors.Error) {
	s.Lock()
	defer s.Unlock()

	pk := s.createPKMap(organizationID, key)

	return s.UnsafeGenericCompositeExist(organizationSettingTable, pk)

}

// Get a setting organization.
func (s *ScyllaOrganizationSettingProvider) Get(organizationID string, key string) (*entities.OrganizationSetting, derrors.Error) {
	s.Lock()
	defer s.Unlock()

	pk := s.createPKMap(organizationID, key)
	var setting interface{} = &entities.OrganizationSetting{}

	err := s.UnsafeCompositeGet(organizationSettingTable, pk, organizationSettingTableColumns, &setting)
	if err != nil {
		return nil, err
	}
	return setting.(*entities.OrganizationSetting), nil

}

// List all the settings of an organization.
func (s *ScyllaOrganizationSettingProvider) List(organizationID string) ([]entities.OrganizationSetting, derrors.Error) {
	s.Lock()
	defer s.Unlock()

	// check connection
	if err := s.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(organizationSettingTable).Columns(organizationSettingTableColumns...).ToCql()
	q := gocqlx.Query(s.Session.Query(stmt), names)

	settings := make([]entities.OrganizationSetting, 0)
	cqlErr := q.SelectRelease(&settings)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list settings")
	}

	return settings, nil

}

// Update a setting of an organization
func (s *ScyllaOrganizationSettingProvider) Update(setting entities.OrganizationSetting) derrors.Error {
	s.Lock()
	defer s.Unlock()

	pk := s.createPKMap(setting.OrganizationId, setting.Key)

	return s.UnsafeCompositeUpdate(organizationSettingTable, pk, organizationSettingTableColumnsNoPK, setting)

}

func (s *ScyllaOrganizationSettingProvider) Remove(organizationID string, key string) derrors.Error {
	s.Lock()
	defer s.Unlock()

	pk := s.createPKMap(organizationID, key)

	return s.UnsafeCompositeRemove(organizationSettingTable, pk)
}

func (s *ScyllaOrganizationSettingProvider) Clear() derrors.Error {
	s.Lock()
	defer s.Unlock()

	return s.UnsafeClear([]string{organizationSettingTable})

}
