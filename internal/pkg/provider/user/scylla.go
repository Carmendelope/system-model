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
 */

package user

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

const userTable = "users"
const userTablePK = "email"

const rowNotFound = "not found"

// TODO: ask to Dani if we need cluster.Consistency = gocql.Quorum
type ScyllaUserProvider struct {
	Address  string
	Port     int
	Keyspace string
	sync.Mutex
	Session *gocql.Session
}

func NewScyllaUserProvider(address string, port int, keyspace string) *ScyllaUserProvider {
	provider := ScyllaUserProvider{Address: address, Port: port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider
}

func (sp *ScyllaUserProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaUserProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}

func (sp *ScyllaUserProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}

}

func (sp *ScyllaUserProvider) unsafeExists(email string) (bool, derrors.Error) {

	var returnedEmail string

	stmt, names := qb.Select(userTable).Columns(userTablePK).Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		userTablePK: email})

	err := q.GetRelease(&returnedEmail)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if user exists")
		}
	}

	return true, nil
}

// check if the session is created
func (sp *ScyllaUserProvider) checkConnection() derrors.Error {
	if sp.Session == nil {
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaUserProvider) checkAndConnect() derrors.Error {

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

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) Add(user entities.User) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(user.Email)

	if err != nil {
		return derrors.AsError(err, "cannot add user")
	}
	if exists {
		return derrors.NewAlreadyExistsError(user.Email)
	}

	// insert a user

	stmt, names := qb.Insert(userTable).Columns("organization_id", "email", "name", "photo_base64", "member_since", "last_name", "title", "phone", "location").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(user)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add user")
	}

	return nil
}

// Update an existing user in the system
func (sp *ScyllaUserProvider) Update(user entities.User) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(user.Email)

	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError(user.Email)
	}

	// update a user
	stmt, names := qb.Update(userTable).Set("organization_id", "name", "photo_base64", "member_since", "last_name", "title", "phone", "location").Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(user)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update user")
	}

	return nil
}

// Exists checks if a user exists on the system.
func (sp *ScyllaUserProvider) Exists(email string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var returnedEmail string

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(userTable).Columns(userTablePK).Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		userTablePK: email})

	err := q.GetRelease(&returnedEmail)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if user exists")
		}
	}

	return true, nil
}

// Get a user.
func (sp *ScyllaUserProvider) Get(email string) (*entities.User, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	var user entities.User
	stmt, names := qb.Select(userTable).Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		userTablePK: email,
	})

	err := q.GetRelease(&user)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError(email)
		} else {
			return nil, derrors.AsError(err, "cannot get user")
		}
	}

	return &user, nil

}

// Remove a user.
func (sp *ScyllaUserProvider) Remove(email string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.unsafeExists(email)

	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("user").WithParams(email)
	}

	// remove a user
	stmt, _ := qb.Delete(userTable).Where(qb.Eq(userTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, email).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot remove user")
	}

	return nil
}

func (sp *ScyllaUserProvider) Clear() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE USERS").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return derrors.AsError(err, "cannot truncate users table")
	}

	return nil
}
