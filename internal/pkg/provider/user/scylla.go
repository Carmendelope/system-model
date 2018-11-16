package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

const addUser = "INSERT INTO Users (organization_id, email, name, photo_url, member_since) VALUES (?, ?, ?, ?, ?)"
const updateUser = "UPDATE Users SET organization_id = ?, name = ?, photo_url = ?, member_since = ? WHERE email = ?"
const exitsUser = "SELECT email from Users where email = ?"
const selectUser = "SELECT organization_id, name, photo_url, member_since from Users where email = ?"
const deleteUser = "delete  from nalej.Users where email = ?"

// TODO: ask to Dani if we need cluster.Consistency = gocql.Quorum
type ScyllaUserProvider struct {
	Address string
	Keyspace string
	Session *gocql.Session
}

func NewScyllaUserProvider (address string, keyspace string) * ScyllaUserProvider {
	return &ScyllaUserProvider{address, keyspace, nil}
}

func (sp *ScyllaUserProvider) Connect() derrors.Error {

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

func (sp *ScyllaUserProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}


func (sp *ScyllaUserProvider) Add(user entities.User) derrors.Error{


	// check if the user exists
	exists, err := sp.Exists(user.Email)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to add the user")
		return conversions.ToDerror(err)
	}
	if  exists {
		log.Info().Str("email", user.Email).Str("user", user.Email).Msg("unable to update the user, user alredy exists")
		return derrors.NewInvalidArgumentError("User alredy exists")
	}

	// insert a user
	cqlErr := sp.Session.Query(addUser,user.OrganizationId, user.Email, user.Name, user.PhotoUrl, user.MemberSince).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to add the user")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Update an existing user in the system
func (sp *ScyllaUserProvider) Update(user entities.User) derrors.Error {

	// check if the user exists
	exists, err := sp.Exists(user.Email)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Str("user", user.Email).Msg("unable to update the user")
		return conversions.ToDerror(err)
	}
	if ! exists {
		log.Info().Str("email", user.Email).Msg("unable to update the user, not exists")
		return derrors.NewInvalidArgumentError("User does not exit")
	}

	// insert a user
	cqlErr := sp.Session.Query(updateUser, user.OrganizationId, user.Name, user.PhotoUrl, user.MemberSince, user.Email).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to update the user")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Exists checks if a user exists on the system.
func (sp *ScyllaUserProvider) Exists(email string) (bool, derrors.Error) {

	// check if exists
	var recoveredEmail string
	err := sp.Session.Query(exitsUser, email).Scan(&recoveredEmail)

	if err == gocql.ErrNotFound{
		return false, nil
	}

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed user exists")
		return false, conversions.ToDerror(err)
	}

	return true, nil
}
// Get a user.
func (sp *ScyllaUserProvider) Get(email string) (* entities.User, derrors.Error) {

	// check if exists
	var organizationId, name, photoUrl string
	var memberSince int64
	err := sp.Session.Query(selectUser, email).Scan(&organizationId,  &name, &photoUrl, &memberSince)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed getting user")
		return nil, conversions.ToDerror(err)
	}

	return &entities.User{OrganizationId:organizationId, Email:email, Name:name, MemberSince: memberSince, PhotoUrl:photoUrl}, nil

}
// Remove a user.
func (sp *ScyllaUserProvider) Remove(email string) derrors.Error {

	// check if the user exists
	exists, err := sp.Exists(email)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Str("user", email).Msg("unable to remove the user")
		return conversions.ToDerror(err)
	}
	if ! exists {
		log.Info().Str("email", email).Msg("unable to remove the user, not exists")
		return derrors.NewInvalidArgumentError("User does not exit")
	}

	// insert a user
	cqlErr := sp.Session.Query(deleteUser, email).Exec()

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to delete the user")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

func (sp *ScyllaUserProvider) ClearTable() derrors.Error{

	err := sp.Session.Query("TRUNCATE TABLE USERS").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return conversions.ToDerror(err)
	}

	return nil
}