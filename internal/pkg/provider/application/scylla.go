package application

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)


//const addInstance = "INSERT INTO applications (app_instance_id, app_descriptor_id, configuration_options, description, environment_variables, groups, labels, name, organization_id, rules, services, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
//const selectInstance = "SELECT organization_id, app_descriptor_id, name, description, configuration_options, environment_variables, labels, rules, groups, services, status FROM applications WHERE app_instance_id = ?"

const applicationTable = "applications"
const applicationTablePK = "app_instance_id"

const applicationDescriptorTable = "applicationdescriptors"
const applicationDescriptorTablePK = "app_descriptor_id"

type ScyllaApplicationProvider struct{
	Address string
	Keyspace string
	Session *gocql.Session
}

func NewScyllaApplicationProvider (address string, keyspace string) * ScyllaApplicationProvider {
	return &ScyllaApplicationProvider{address, keyspace, nil}
}

func (sp *ScyllaApplicationProvider) Connect() derrors.Error {

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

func (sp *ScyllaApplicationProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}

// check if the session is created
func (sp *ScyllaApplicationProvider) CheckConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new application descriptor to the system
func (sp *ScyllaApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if err != nil {
		log.Info().Msg("unable to add the application descriptor")
		return err
	}

	// check if the application exists
	exists, err := sp.DescriptorExists(descriptor.AppDescriptorId)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewInvalidArgumentError("The descriptor already exists")
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationDescriptorTable).Columns("organization_id","app_descriptor_id", "name",
		"description","configuration_options","environment_variables","labels","rules","groups","services").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(descriptor)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// GetDescriptors retrieves an application descriptor.
func (sp *ScyllaApplicationProvider) GetDescriptor(appDescriptorID string) (* entities.AppDescriptor, derrors.Error) {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		log.Info().Msg("unable to add the descriptor")
		return nil, err
	}

	// 2.- Gocqlx
	var descriptor entities.AppDescriptor
	stmt, names := qb.Select(applicationDescriptorTable).Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&descriptor)
	if err != nil {
		return nil, conversions.ToDerror(err)
	}

	return &descriptor, nil
}

// DescriptorExists checks if a given descriptor exists on the system.
func (sp *ScyllaApplicationProvider) DescriptorExists(appDescriptorID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(applicationDescriptorTable).Columns(applicationDescriptorTablePK).Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&returnedId)
	if err != nil {
		return false, nil
	}

	return returnedId == appDescriptorID, nil
}

// DeleteDescriptor removes a given descriptor from the system.
func (sp *ScyllaApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.DescriptorExists(appDescriptorID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewInvalidArgumentError("Application descriptor does not exit")
	}

	// delete app instance
	stmt, _ := qb.Delete(applicationDescriptorTable).Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appDescriptorID).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to delete the application descriptor")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// AddInstance adds a new application instance to the system
func (sp *ScyllaApplicationProvider) AddInstance(instance entities.AppInstance) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if err != nil {
		log.Info().Msg("unable to add the application instance")
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(instance.AppInstanceId)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if exists {
		return derrors.NewInvalidArgumentError("The application already exists")
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationTable).Columns("organization_id","app_descriptor_id","app_instance_id",
		"name","description","configuration_options","environment_variables","labels","rules","groups","services","status").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil

}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaApplicationProvider) InstanceExists(appInstanceID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(applicationTable).Columns(applicationTablePK).Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationTablePK: appInstanceID,
	})

	err := q.GetRelease(&returnedId)
	if err != nil {
		return false, nil
	}

	return true, nil

}

// GetInstance retrieves an application instance.
func (sp *ScyllaApplicationProvider) GetInstance(appInstanceID string) (* entities.AppInstance, derrors.Error) {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		log.Info().Msg("unable to add the application instance")
		return nil, err
	}

/*
	// 1.-
	var organizationId, appDescriptorId, name, description  string
	var status entities.ApplicationStatus
	configurationOptions := make(map[string]string)
	environmentVariables:= make(map[string]string)
	labels := make(map[string]string)
	groups := make([]entities.ServiceGroupInstance, 0)
	services := make([]entities.ServiceInstance, 0)
	rules := make([]entities.SecurityRule, 0)



	err := sp.Session.Query(selectInstance, appInstanceID).Scan(&organizationId, &appDescriptorId, &name,
		&description, &configurationOptions, &environmentVariables, &labels, &rules, &groups, &services, &status)

	if err != nil {
		return nil, conversions.ToDerror(err)
	}
*/
	// 2.- Gocqlx
	var app entities.AppInstance
	stmt, names := qb.Select(applicationTable).Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationTablePK: appInstanceID,
	})

	err := q.GetRelease(&app)
	if err != nil {
		return nil, conversions.ToDerror(err)
	}

	return &app, nil

	/*return &entities.AppInstance{
		OrganizationId: organizationId,
		AppDescriptorId: appDescriptorId,
		AppInstanceId: appInstanceID,
		Name: name,
		Description:description,
		ConfigurationOptions:configurationOptions,
		EnvironmentVariables:environmentVariables,
		Labels:labels,
		Rules: rules,
		Groups:groups,
		Services: services,
		Status: status}, nil
*/
}

// DeleteInstance removes a given instance from the system.
func (sp *ScyllaApplicationProvider) DeleteInstance(appInstanceID string) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(appInstanceID)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewInvalidArgumentError("Application instance does not exit")
	}

	// delete app instance
	stmt, _ := qb.Delete(applicationTable).Where(qb.Eq(applicationTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appInstanceID).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to delete the application instance")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// UpdateInstance updates the information of an instance
func (sp *ScyllaApplicationProvider) UpdateInstance(instance entities.AppInstance) derrors.Error {
	// check connection
	err := sp.CheckConnection()
	if err != nil {
		log.Info().Msg("unable to update the application instance")
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(instance.AppInstanceId)
	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewInvalidArgumentError("The application does not exist")
	}

	// insert the application instance
	stmt, names := qb.Update(applicationTable).Set("organization_id","app_descriptor_id",
		"name","description","configuration_options","environment_variables",
		"labels","rules","groups","services","status").Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// Clear descriptors and instances
func (sp *ScyllaApplicationProvider) Clear() derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// delete app instances
	err := sp.Session.Query("TRUNCATE TABLE applications").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the applications table")
		return conversions.ToDerror(err)
	}

	// TRUNCATE TABLE applicationdescriptors
	err = sp.Session.Query("TRUNCATE TABLE applicationdescriptors").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the applications descriptor table")
		return conversions.ToDerror(err)
	}
	return nil


}

