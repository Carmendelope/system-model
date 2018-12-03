package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const userTable = "users"
const userTablePK = "email"

const rowNotFound = "not found"

// TODO: ask to Dani if we need cluster.Consistency = gocql.Quorum
type ScyllaUserProvider struct {
	Address string
	Port int
	Keyspace string
	Session *gocql.Session
}

func NewScyllaUserProvider (address string, port int, keyspace string) * ScyllaUserProvider {
	provider:= ScyllaUserProvider{address, port, keyspace, nil}
	provider.Connect()
	return &provider
}

func (sp *ScyllaUserProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaUserProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return conversions.ToDerror(err)
	}

	sp.Session = session

	return nil
}

func (sp *ScyllaUserProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}

// check if the session is created
func (sp *ScyllaUserProvider) CheckConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) Add(user entities.User) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(user.Email)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if  exists {
		return derrors.NewAlreadyExistsError(user.Email)
	}

	// insert a user

	stmt, names := qb.Insert(userTable).Columns("organization_id", "email", "name", "photo_url","member_since").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(user)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Update an existing user in the system
func (sp *ScyllaUserProvider) Update(user entities.User) derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(user.Email)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError(user.Email)
	}

	// update a user
	stmt, names := qb.Update(userTable).Set("organization_id", "name", "photo_url","member_since").Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(user)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Exists checks if a user exists on the system.
func (sp *ScyllaUserProvider) Exists(email string) (bool, derrors.Error) {

	var returnedEmail string

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(userTable).Columns(userTablePK).Where(qb.Eq(userTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		userTablePK: email })

	err := q.GetRelease(&returnedEmail)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil
}
// Get a user.
func (sp *ScyllaUserProvider) Get(email string) (* entities.User, derrors.Error) {

	// check connection
	if err := sp.CheckConnection(); err != nil {
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
			return nil, conversions.ToDerror(err)
		}else{
			return nil, derrors.NewNotFoundError(email)
		}
	}

	return &user, nil

}
// Remove a user.
func (sp *ScyllaUserProvider) Remove(email string) derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(email)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("user").WithParams(email)
	}

	// remove a user
	stmt, _ := qb.Delete(userTable).Where(qb.Eq(userTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, email).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

func (sp *ScyllaUserProvider) Clear() derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE USERS").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return conversions.ToDerror(err)
	}

	return nil
}