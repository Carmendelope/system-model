package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

// Table constants
// ---------
// -- User
// ---------
const UserTable = "users"
const UserTablePK = "email"
var allUserColumns = []string{"organization_id", "email", "name", "photo_url","member_since", "contact_info"}
var allUserColumnsNoPK = []string{"organization_id", "name", "photo_url","member_since", "contact_info"}

// ------------------
// -- Account User
// ------------------
const AccountUserTable = "AccountUser"
var allAccountUserColumns = []string{"account_id", "email", "role_id", "internal", "status"}
var allAccountUserColumnsNoPK = []string{"role_id", "internal", "status"}

// ------------------------
// -- Account User Invite
// ------------------------
const AccountUserInviteTable = "AccountUserInvite"
var allAccountUserInviteColumns = []string{"account_id", "email", "role_id", "invited_by", "msg", "expires"}
var allAccountUserInviteColumnsNoPK = []string{"role_id", "invited_by", "msg", "expires"}

// ------------------
// -- Project User
// ------------------
const ProjectUserTable = "ProjectUser"
var allProjectUserColumns = []string{"account_id", "project_id", "email", "role_id", "internal", "status"}
var allProjectUserColumnsNoPK = []string{"role_id", "internal", "status"}

// ------------------------
// -- Project User Invite
// ------------------------
const ProjectUserInviteTable = "ProjectUserInvite"
var allProjectUserInviteColumns = []string{"account_id", "project_id", "email", "role_id", "invited_by", "msg", "expires"}
var allProjectUserInviteColumnsNoPK = []string{"role_id", "invited_by", "msg", "expires"}


// TODO: ask to Dani if we need cluster.Consistency = gocql.Quorum
type ScyllaUserProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaUserProvider (address string, port int, keyspace string) * ScyllaUserProvider {
	provider:= ScyllaUserProvider{
	 	ScyllaDB: scylladb.ScyllaDB{
	 		Address: address,
	 		Port: port,
	 		Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

func (sp *ScyllaUserProvider) Disconnect () {

	sp.Lock()
	defer sp.Unlock()

	sp.ScyllaDB.Disconnect()
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) Add(user entities.User) derrors.Error{

	sp.Lock()
	defer sp.Unlock()

	log.Debug().Interface("user", user).Msg("provider add user")
	return sp.UnsafeAdd(UserTable, UserTablePK, user.Email, allUserColumns, user)
}
// Update an existing user in the system
func (sp *ScyllaUserProvider) Update(user entities.User) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeUpdate(UserTable, UserTablePK, user.Email, allUserColumnsNoPK, user)

}
// Exists checks if a user exists on the system.
func (sp *ScyllaUserProvider) Exists(email string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(UserTable, UserTablePK, email)

}
// Get a user.
func (sp *ScyllaUserProvider) Get(email string) (*entities.User, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var user interface{} = &entities.User{}
	err := sp.UnsafeGet(UserTable, UserTablePK, email, allUserColumns, &user)
	if err != nil{
		return nil, err
	}
	return user.(*entities.User), nil

}
// Remove a user.
func (sp *ScyllaUserProvider) Remove(email string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(UserTable, UserTablePK, email)
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) createAccountUserPKMap(accountID string, email string) map[string]interface{}{

	res := map[string]interface{}{
		"account_id": accountID,
		"email": email,
	}
	return res
}
func (sp *ScyllaUserProvider) AddAccountUser(accUser entities.AccountUser) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, accUser.Email)
	if err != nil {
		return err
	}
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(accUser.Email)
	}

	pkColumn := sp.createAccountUserPKMap(accUser.AccountId, accUser.Email)

	return sp.UnsafeCompositeAdd(AccountUserTable, pkColumn, allAccountUserColumns, accUser)
}
func (sp *ScyllaUserProvider) UpdateAccountUser(accUser entities.AccountUser) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createAccountUserPKMap(accUser.AccountId, accUser.Email)

	return sp.UnsafeCompositeUpdate(AccountUserTable, pkColumn, allAccountUserColumnsNoPK, accUser)
}
func (sp *ScyllaUserProvider) RemoveAccountUser(accountID string, email string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createAccountUserPKMap(accountID, email)

	return sp.UnsafeCompositeRemove(AccountUserTable, pkColumn)
}
// ListAccountUser lists all the accounts of a user
func (sp *ScyllaUserProvider) ListAccountUser(email string) ([]entities.AccountUser, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	accounts := make([]entities.AccountUser, 0)

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, email)
	if err != nil {
		return accounts, err
	}
	if ! userExists{
		return accounts, derrors.NewNotFoundError("User").WithParams(email)
	}

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AccountUserTable).Columns(allAccountUserColumns...).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"email": email,
	})

	cqlErr := gocqlx.Select(&accounts, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list AccountUser")
	}

	return accounts, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) AddAccountUserInvite(accUser entities.AccountUserInvite) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, accUser.Email)
	if err != nil {
		return err
	}
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(accUser.Email)
	}

	pkColumn := sp.createAccountUserPKMap(accUser.AccountId, accUser.Email)

	return sp.UnsafeCompositeAdd(AccountUserInviteTable, pkColumn, allAccountUserInviteColumns, accUser)
}
func (sp *ScyllaUserProvider) GetAccountUserInvite(accountID string, email string) (*entities.AccountUserInvite, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK,email)
	if err != nil {
		return nil, err
	}
	if ! userExists{
		return nil, derrors.NewNotFoundError("User").WithParams(email)
	}

	pkColumn := sp.createAccountUserPKMap(accountID, email)

	var invite interface{} = &entities.AccountUserInvite{}

	err = sp.UnsafeCompositeGet(AccountUserInviteTable, pkColumn, allAccountUserInviteColumns, &invite)
	if err != nil {
		return nil, err
	}
	return invite.(*entities.AccountUserInvite), nil
}
func (sp *ScyllaUserProvider) RemoveAccountUserInvite(accountID string, email string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createAccountUserPKMap(accountID, email)

	return sp.UnsafeCompositeRemove(AccountUserInviteTable, pkColumn)
}
func (sp *ScyllaUserProvider) ListAccountUserInvites(email string) ([]entities.AccountUserInvite, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	invites := make([]entities.AccountUserInvite, 0)

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, email)
	if err != nil {
		return invites, err
	}
	if ! userExists{
		return invites, derrors.NewNotFoundError("User").WithParams(email)
	}

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AccountUserInviteTable).Columns(allAccountUserInviteColumns...).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"email": email,
	})

	cqlErr := gocqlx.Select(&invites, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list AccountUserInvite")
	}

	return invites, nil
}

// ---------------------------------------------------------------------------------------------------------------------
func (sp *ScyllaUserProvider) createProjectUserPKMap(accountID string, projectID string, email string) map[string]interface{}{

	res := map[string]interface{}{
		"account_id": accountID,
		"project_id" : projectID,
		"email": email,
	}
	return res
}
func (sp *ScyllaUserProvider) AddProjectUser(projUser entities.ProjectUser) derrors.Error{

	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, projUser.Email)
	if err != nil {
		return err
	}
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(projUser.Email)
	}

	pkColumn := sp.createProjectUserPKMap(projUser.AccountId, projUser.ProjectId, projUser.Email)

	return sp.UnsafeCompositeAdd(ProjectUserTable, pkColumn, allProjectUserColumns, projUser)
}
func (sp *ScyllaUserProvider) UpdateProjectUser(projUser entities.ProjectUser) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createProjectUserPKMap(projUser.AccountId, projUser.ProjectId, projUser.Email)

	return sp.UnsafeCompositeUpdate(ProjectUserTable, pkColumn, allProjectUserColumnsNoPK, projUser)
}
func (sp *ScyllaUserProvider) RemoveProjectUser(accountID string, projectID string, email string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createProjectUserPKMap(accountID, projectID, email)

	return sp.UnsafeCompositeRemove(ProjectUserTable, pkColumn)
}
func (sp *ScyllaUserProvider) ListProjectUser(accountID string, projectID string) ([]entities.ProjectUser, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	projectUsers := make([]entities.ProjectUser, 0)

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ProjectUserTable).Columns(allProjectUserColumns...).Where(qb.Eq("account_id")).
	Where(qb.Eq("project_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"account_id": accountID,
		"project_id": projectID,
	})

	cqlErr := gocqlx.Select(&projectUsers, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list ProjectUser")
	}

	return projectUsers, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) AddProjectUserInvite(invite entities.ProjectUserInvite) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, invite.Email)
	if err != nil {
		return err
	}
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(invite.Email)
	}

	pkColumn := sp.createProjectUserPKMap(invite.AccountId, invite.ProjectId, invite.Email)

	return sp.UnsafeCompositeAdd(ProjectUserInviteTable, pkColumn, allProjectUserInviteColumns, invite)
}
func (sp *ScyllaUserProvider) GetProjectUserInvite(accountID string, projectId string, email string) (*entities.ProjectUserInvite, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK,email)
	if err != nil {
		return nil, err
	}
	if ! userExists{
		return nil, derrors.NewNotFoundError("User").WithParams(email)
	}

	pkColumn := sp.createProjectUserPKMap(accountID, projectId, email)

	var invite interface{} = &entities.ProjectUserInvite{}

	err = sp.UnsafeCompositeGet(ProjectUserInviteTable, pkColumn, allProjectUserInviteColumns, &invite)
	if err != nil {
		return nil, err
	}
	return invite.(*entities.ProjectUserInvite), nil
}
func (sp *ScyllaUserProvider) RemoveProjectUserInvite(accountID string, projectID string, email string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createProjectUserPKMap(accountID, projectID, email)

	return sp.UnsafeCompositeRemove(ProjectUserInviteTable, pkColumn)

}
func (sp *ScyllaUserProvider) ListProjectUserInvites(email string) ([]entities.ProjectUserInvite, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	invites := make([]entities.ProjectUserInvite, 0)

	// ask if the user exists
	userExists, err := sp.UnsafeGenericExist(UserTable, UserTablePK, email)
	if err != nil {
		return invites, err
	}
	if ! userExists{
		return invites, derrors.NewNotFoundError("User").WithParams(email)
	}

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ProjectUserInviteTable).Columns(allProjectUserInviteColumns...).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"email": email,
	})

	cqlErr := gocqlx.Select(&invites, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list ProjectUserInvite")
	}

	return invites, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaUserProvider) Clear() derrors.Error{

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{UserTable, AccountUserTable, AccountUserInviteTable})

}