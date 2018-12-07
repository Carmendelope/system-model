package role

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const roleTable = "roles"
const roleTablePK = "role_id"

const rowNotFound = "not found"


type ScyllaRoleProvider struct {
	Address string
	Port int
	Keyspace string
	Session *gocql.Session
}

func NewSScyllaRoleProvider (address string, port int, keyspace string) * ScyllaRoleProvider {
	provider := ScyllaRoleProvider{address, port,keyspace, nil}
	provider.Connect()
	return &provider

}

func (sp *ScyllaRoleProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaRoleProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
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

func (sp *ScyllaRoleProvider) CheckConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaRoleProvider) CheckAndConnect () derrors.Error{

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err = sp.Connect()
		if err != nil  {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// Add a new role to the system.
func (sp *ScyllaRoleProvider) Add(role entities.Role) derrors.Error {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// check if the role exists
	exists, err := sp.Exists(role.RoleId)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if  exists {
		return derrors.NewAlreadyExistsError(role.RoleId)
	}

	// insert a role
	stmt, names := qb.Insert(roleTable).Columns("organization_id","role_id","name","description","internal","created").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(role)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// Update an existing role in the system
func (sp *ScyllaRoleProvider) Update(role entities.Role) derrors.Error{

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(role.RoleId)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError(role.RoleId)
	}

	// update the role
	stmt, names := qb.Update(roleTable).Set("organization_id","name","description","internal","created").Where(qb.Eq(roleTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(role)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// Exists checks if a role exists on the system.
func (sp *ScyllaRoleProvider) Exists(roleID string) (bool, derrors.Error){

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	// check if exists
	var recoveredRoleID string
	stmt, names := qb.Select(roleTable).Columns(roleTablePK).Where(qb.Eq(roleTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		roleTablePK: roleID })

	err := q.GetRelease(&recoveredRoleID)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}
	return true, nil
}

// Get a role.
func (sp *ScyllaRoleProvider) Get(roleID string) (* entities.Role, derrors.Error) {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	var role entities.Role
	stmt, names := qb.Select(roleTable).Where(qb.Eq(roleTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		roleTablePK: roleID,
	})

	err := q.GetRelease(&role)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, conversions.ToDerror(err)
		}else{
			return nil, derrors.NewNotFoundError(roleID)
		}
	}

	return &role, nil

}

// Remove a role
func (sp *ScyllaRoleProvider) Remove(roleID string) derrors.Error {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// check if the role exists
	exists, err := sp.Exists(roleID)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("role").WithParams(roleID)
	}

	// remove the role
	stmt, _ := qb.Delete(roleTable).Where(qb.Eq(roleTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, roleID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Truncate the table

func (sp *ScyllaRoleProvider) Clear() derrors.Error{

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE ROLES").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	return nil
}