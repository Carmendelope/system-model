package role

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

/*
create table nalej.Roles (organization_id text, role_id text, name text, description text, created int, PRIMARY KEY (role_id));
*/

const addRole = "INSERT INTO Roles (organization_id, role_id, name, description , created) VALUES (?, ?, ?, ?, ?)"
const updateRole = "UPDATE Roles SET organization_id = ?, name = ?, description = ?, created = ? WHERE role_id = ?"
const exitsRole = "SELECT role_id from Roles where role_id = ?"
const selectRole = "SELECT organization_id, name, description , created from Roles where role_id = ?"
const deleteRole = "delete  from Roles where role_id = ?"

type ScyllaRoleProvider struct {
	Address string
	Keyspace string
	Session *gocql.Session
}

func NewSScyllaRoleProvider (address string, keyspace string) * ScyllaRoleProvider {
	return &ScyllaRoleProvider{address, keyspace, nil}
}

func (sp *ScyllaRoleProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace

	session, err := conf.CreateSession()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return conversions.ToDerror(err)
	}

	sp.Session = session

	return nil
}

func (sp *ScyllaRoleProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}

// Add a new role to the system.
func (sp *ScyllaRoleProvider) Add(role entities.Role) derrors.Error {

	// check if the role exists
	exists, err := sp.Exists(role.RoleId)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to add the role")
		return conversions.ToDerror(err)
	}
	if  exists {
		log.Info().Str("role_id", role.RoleId).Msg("unable to add the role, it alredy exists")
		return derrors.NewInvalidArgumentError("Role alredy exists")
	}

	// insert a user
	cqlErr := sp.Session.Query(addRole,role.OrganizationId, role.RoleId, role.Name, role.Description, role.Created).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to add the role")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Update an existing role in the system
func (sp *ScyllaRoleProvider) Update(role entities.Role) derrors.Error{
	// check if the user exists
	exists, err := sp.Exists(role.RoleId)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Str("role", role.RoleId).Msg("unable to update the role")
		return conversions.ToDerror(err)
	}
	if ! exists {
		log.Info().Str("role", role.RoleId).Msg("unable to update the role, not exists")
		return derrors.NewInvalidArgumentError("Role does not exit")
	}

	// insert a user
	cqlErr := sp.Session.Query(updateRole, role.OrganizationId, role.Name, role.Description, role.Created, role.RoleId).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to update the role")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Exists checks if a role exists on the system.
func (sp *ScyllaRoleProvider) Exists(roleID string) (bool, derrors.Error){

	// check if exists
	var recoveredEmail string
	err := sp.Session.Query(exitsRole, roleID).Scan(&recoveredEmail)

	if err == gocql.ErrNotFound{
		return false, nil
	}

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed role exists")
		return false, conversions.ToDerror(err)
	}

	return true, nil
}
// Get a role.
func (sp *ScyllaRoleProvider) Get(roleID string) (* entities.Role, derrors.Error) {
	// check if exists
	var organizationId, name, description string
	var created int64
	err := sp.Session.Query(selectRole, roleID).Scan(&organizationId,  &name, &description, &created)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed getting user")
		return nil, conversions.ToDerror(err)
	}

	return &entities.Role{OrganizationId:organizationId, Name:name, Description:description, RoleId:roleID, Created:created}, nil

}
// Remove a role
func (sp *ScyllaRoleProvider) Remove(roleID string) derrors.Error {

	// check if the user exists
	exists, err := sp.Exists(roleID)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Str("role", roleID).Msg("unable to remove the role")
		return conversions.ToDerror(err)
	}
	if ! exists {
		log.Info().Str("role", roleID).Msg("unable to remove the role, it not exists")
		return derrors.NewInvalidArgumentError("Role does not exit")
	}

	// insert a user
	cqlErr := sp.Session.Query(deleteRole, roleID).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to delete the role")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Truncate the table
func (sp *ScyllaRoleProvider) ClearTable() derrors.Error{
	err := sp.Session.Query("TRUNCATE TABLE ROLES").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return conversions.ToDerror(err)
	}

	return nil
}