package organization

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const organizationTable = "Organizations"
const organizationTablePK = "id"
const organizationClusterTable = "Organization_Clusters"
const organizationNodeTable = "Organization_Nodes"
const organizationDescriptorTable = "organization_appdescriptors"
const organizationInstanceTable = "Organization_appinstances"
const organizationUserTable = "Organization_Users"
const organizationRoleTable = "Organization_Roles"

const rowNotFound = "not found"
type ScyllaOrganizationProvider struct {
	Address string
	Keyspace string
	Session *gocql.Session
}

func NewScyllaOrganizationProvider (address string, keyspace string) * ScyllaOrganizationProvider {
	return &ScyllaOrganizationProvider{ address, keyspace, nil}
}

// connect to the database
func (sp *ScyllaOrganizationProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace

	session, err := conf.CreateSession()
	if err != nil {
		return conversions.ToDerror(err)
	}

	sp.Session = session

	return nil
}

// disconnect from the database
func (sp *ScyllaOrganizationProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}

// check that the session is created
func (sp *ScyllaOrganizationProvider) CheckConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// Add a new organization to the system.
func (sp *ScyllaOrganizationProvider) Add(org entities.Organization) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the organization exists
	exists, err := sp.Exists(org.ID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError(org.ID)
	}

	// insert the organization instance
	stmt, names := qb.Insert(organizationTable).Columns("id", "name", "created").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(org)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Check if an organization exists on the system.
func (sp *ScyllaOrganizationProvider) Exists(organizationID string) (bool, derrors.Error){
	var returnedId string

	stmt, names := qb.Select(organizationTable).Columns(organizationTablePK).Where(qb.Eq(organizationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTablePK: organizationID })

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil
}
// Get an organization.
func (sp *ScyllaOrganizationProvider) Get(organizationID string) (* entities.Organization, derrors.Error){
	// check connection
	if err := sp.CheckConnection(); err != nil {
		return nil, err
	}

	var org entities.Organization
	stmt, names := qb.Select(organizationTable).Where(qb.Eq(organizationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTablePK: organizationID,
	})

	cqlErr := q.GetRelease(&org)
	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound{
			return nil, derrors.NewNotFoundError(organizationID)
		}else {
			return nil, conversions.ToDerror(cqlErr)
		}
	}

	return &org, nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddCluster adds a new cluster ID to the organization.
func (sp *ScyllaOrganizationProvider) AddCluster(organizationID string, clusterID string) derrors.Error{
	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the organization exists
	exists, err = sp.ClusterExists(organizationID, clusterID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("cluster").WithParams(organizationID, clusterID)
	}

	// insert the organization instance
	stmt, names := qb.Insert(organizationClusterTable).Columns("organization_id","cluster_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"cluster_id": clusterID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// ClusterExists checks if a cluster is linked to an organization.
func (sp *ScyllaOrganizationProvider) ClusterExists(organizationID string, clusterID string) (bool, derrors.Error){
	var returnedId string

	stmt, names := qb.Select(organizationClusterTable).Columns("cluster_id").Where(qb.Eq("organization_id")).Where(qb.Eq("cluster_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"cluster_id": clusterID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil
}
// ListClusters returns a list of clusters in an organization.
func (sp *ScyllaOrganizationProvider) ListClusters(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationClusterTable).Columns("cluster_id").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	clusters := make ([]string, 0)
	cqlErr := gocqlx.Select(&clusters, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return clusters, nil
}
// DeleteCluster removes a cluster from an organization.
func (sp *ScyllaOrganizationProvider) DeleteCluster(organizationID string, clusterID string) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the cluster exists in the organization
	exists, err := sp.ClusterExists(organizationID, clusterID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("cluster").WithParams(organizationID, clusterID)
	}

	// delete a cluster of an organization
	stmt, _ := qb.Delete(organizationClusterTable).Where(qb.Eq("organization_id")).Where(qb.Eq("cluster_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, clusterID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddNode adds a new node ID to the organization.
func (sp *ScyllaOrganizationProvider) AddNode(organizationID string, nodeID string) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the organization exists
	exists, err = sp.NodeExists(organizationID, nodeID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("node").WithParams(organizationID, nodeID)
	}

	// add a node in the organization instance
	stmt, names := qb.Insert(organizationNodeTable).Columns("organization_id","node_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"node_id": nodeID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// NodeExists checks if a node is linked to an organization.
func (sp *ScyllaOrganizationProvider) NodeExists(organizationID string, nodeID string) (bool, derrors.Error){

	var returnedId string

	stmt, names := qb.Select(organizationNodeTable).Columns("node_id").Where(qb.Eq("organization_id")).Where(qb.Eq("node_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"node_id": nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound{
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}
// ListNodes returns a list of nodes in an organization.
func (sp *ScyllaOrganizationProvider) ListNodes(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationNodeTable).Columns("node_id").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	nodes := make ([]string, 0)
	cqlErr := gocqlx.Select(&nodes, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return nodes, nil
}
// DeleteNode removes a node from an organization.
func (sp *ScyllaOrganizationProvider) DeleteNode(organizationID string, nodeID string) derrors.Error{

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the node exists in the organization
	exists, err := sp.NodeExists(organizationID, nodeID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("node").WithParams(organizationID, nodeID)
	}

	// delete the node of an organization
	stmt, _ := qb.Delete(organizationNodeTable).Where(qb.Eq("organization_id")).Where(qb.Eq("node_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, nodeID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new descriptor ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddDescriptor(organizationID string, appDescriptorID string) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the descriptor exists
	exists, err = sp.DescriptorExists(organizationID, appDescriptorID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("appDescriptor").WithParams(organizationID, appDescriptorID)
	}

	// add an app descriptor in the organization instance
	stmt, names := qb.Insert(organizationDescriptorTable).Columns("organization_id","app_descriptor_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_descriptor_id": appDescriptorID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// DescriptorExists checks if an application descriptor exists on the system.
func (sp *ScyllaOrganizationProvider) DescriptorExists(organizationID string, appDescriptorID string) (bool, derrors.Error){

	var returnedId string

	stmt, names := qb.Select(organizationDescriptorTable).Columns("app_descriptor_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_descriptor_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_descriptor_id": appDescriptorID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound{
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}
// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
func (sp *ScyllaOrganizationProvider) ListDescriptors(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationDescriptorTable).Columns("app_descriptor_id").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	descriptors := make ([]string, 0)
	cqlErr := gocqlx.Select(&descriptors, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return descriptors, nil
}
// DeleteDescriptor removes a descriptor from an organization
func (sp *ScyllaOrganizationProvider) DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error{

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the descriptor exists in the organization
	exists, err := sp.DescriptorExists(organizationID, appDescriptorID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("app descriptor").WithParams(organizationID, appDescriptorID)
	}

	// delete the descriptor of an organization
	stmt, _ := qb.Delete(organizationDescriptorTable).Where(qb.Eq("organization_id")).Where(qb.Eq("app_descriptor_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appDescriptorID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddInstance adds a new application instance ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddInstance(organizationID string, appInstanceID string) derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the instance exists
	exists, err = sp.InstanceExists(organizationID, appInstanceID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("app_instance").WithParams(organizationID, appInstanceID)
	}

	// add an app instance in the organization
	stmt, names := qb.Insert(organizationInstanceTable).Columns("organization_id","app_instance_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaOrganizationProvider) InstanceExists(organizationID string, appInstanceID string) (bool, derrors.Error){

	var returnedId string

	stmt, names := qb.Select(organizationInstanceTable).Columns("app_instance_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound{
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}
// ListInstances returns a the identifiers associate with a given organization.
func (sp *ScyllaOrganizationProvider) ListInstances(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationInstanceTable).Columns("app_instance_id").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	instances := make ([]string, 0)
	cqlErr := gocqlx.Select(&instances, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return instances, nil
}
// DeleteInstance removes an instance from an organization
func (sp *ScyllaOrganizationProvider) DeleteInstance(organizationID string, appInstanceID string) derrors.Error{

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the instance exists in the organization
	exists, err := sp.InstanceExists(organizationID, appInstanceID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("app descriptor").WithParams(organizationID, appInstanceID)
	}

	// delete the instance of an organization
	stmt, _ := qb.Delete(organizationInstanceTable).Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appInstanceID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddUser adds a new user to the organization.
func (sp *ScyllaOrganizationProvider) AddUser(organizationID string, email string) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the user exists
	exists, err = sp.UserExists(organizationID, email)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("user").WithParams(organizationID, email)
	}

	// add an user in the organization
	stmt, names := qb.Insert(organizationUserTable).Columns("organization_id","email").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"email": email})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// UserExists checks if a user is linked to an organization.
func (sp *ScyllaOrganizationProvider) UserExists(organizationID string, email string) (bool, derrors.Error){

	var returnedId string

	stmt, names := qb.Select(organizationUserTable).Columns("email").Where(qb.Eq("organization_id")).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"email": email})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound{
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}
// ListUser returns a list of users in an organization.
func (sp *ScyllaOrganizationProvider) ListUsers(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationUserTable).Columns("email").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	users := make ([]string, 0)
	cqlErr := gocqlx.Select(&users, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return users, nil
}
// DeleteUser removes a user from an organization.
func (sp *ScyllaOrganizationProvider) DeleteUser(organizationID string, email string) derrors.Error{

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the user exists in the organization
	exists, err := sp.UserExists(organizationID, email)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("user").WithParams(organizationID, email)
	}

	// delete the user of an organization
	stmt, _ := qb.Delete(organizationUserTable).Where(qb.Eq("organization_id")).Where(qb.Eq("email")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, email).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddRole adds a new role ID to the organization.
func (sp *ScyllaOrganizationProvider) AddRole(organizationID string, roleID string) derrors.Error{

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	exists, err := sp.Exists(organizationID)
	if err != nil{
		return conversions.ToDerror(err)
	}
	if !exists{
		return derrors.NewAlreadyExistsError("organization").WithParams(organizationID)
	}

	// check if the role exists
	exists, err = sp.RoleExists(organizationID, roleID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewAlreadyExistsError("role").WithParams(organizationID, roleID)
	}

	// add an user in the organization
	stmt, names := qb.Insert(organizationRoleTable).Columns("organization_id","role_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"role_id": roleID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// RoleExists checks if a role is linked to an organization.
func (sp *ScyllaOrganizationProvider) RoleExists(organizationID string, roleID string) (bool, derrors.Error){

	var returnedId string

	stmt, names := qb.Select(organizationRoleTable).Columns("role_id").Where(qb.Eq("organization_id")).Where(qb.Eq("role_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"role_id": roleID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound{
			return false, nil
		}else{
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}
// ListNodes returns a list of roles in an organization.
func (sp *ScyllaOrganizationProvider) ListRoles(organizationID string) ([]string, derrors.Error){

	// 1.-Check if the organization exists
	exists, err := sp.Exists(organizationID)
	if err != nil{
		return nil, conversions.ToDerror(err)
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationRoleTable).Columns("role_id").Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	users := make ([]string, 0)
	cqlErr := gocqlx.Select(&users, q.Query)

	if cqlErr != nil {
		return nil, conversions.ToDerror(cqlErr)
	}

	return users, nil
}
// DeleteRole removes a role from an organization.
func (sp *ScyllaOrganizationProvider) DeleteRole(organizationID string, roleID string) derrors.Error{

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the role exists in the organization
	exists, err := sp.RoleExists(organizationID, roleID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("role").WithParams(organizationID, roleID)
	}

	// delete the role of the organization
	stmt, _ := qb.Delete(organizationRoleTable).Where(qb.Eq("organization_id")).Where(qb.Eq("role_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, roleID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaOrganizationProvider) Clear() derrors.Error{
	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// delete organizations table
	err := sp.Session.Query("TRUNCATE TABLE organizations").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-cluster table
	err = sp.Session.Query("TRUNCATE TABLE organization_clusters").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-nodes table
	err = sp.Session.Query("TRUNCATE TABLE organization_nodes").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-descriptors table
	err = sp.Session.Query("TRUNCATE TABLE organization_appdescriptors").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-instances table
	err = sp.Session.Query("TRUNCATE TABLE organization_appinstances").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-users table
	err = sp.Session.Query("TRUNCATE TABLE organization_users").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}

	// delete organization-roles table
	err = sp.Session.Query("TRUNCATE TABLE organization_roles").Exec()
	if err != nil {
		return conversions.ToDerror(err)
	}
		return nil
}
