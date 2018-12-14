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

const applicationTable = "ApplicationInstances"
const applicationTablePK = "app_instance_id"

const applicationDescriptorTable = "applicationdescriptors"
const applicationDescriptorTablePK = "app_descriptor_id"

const rowNotFound = "not found"

type ScyllaApplicationProvider struct{
	Address string
	Port int
	Keyspace string
	Session *gocql.Session
}

func NewScyllaApplicationProvider (address string, port int, keyspace string) * ScyllaApplicationProvider {
	provider := ScyllaApplicationProvider{address, port, keyspace, nil}
	provider.Connect()
	return &provider
}

func (sp *ScyllaApplicationProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaApplicationProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
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

func (sp *ScyllaApplicationProvider) CheckAndConnect () derrors.Error{

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
// ---------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new application descriptor to the system
func (sp *ScyllaApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.DescriptorExists(descriptor.AppDescriptorId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(descriptor.AppDescriptorId)
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationDescriptorTable).Columns("organization_id","app_descriptor_id", "name",
		"description","configuration_options","environment_variables","labels","rules","groups","services").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(descriptor)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appDescriptor")
	}

	return nil
}

// GetDescriptors retrieves an application descriptor.
func (sp *ScyllaApplicationProvider) GetDescriptor(appDescriptorID string) (* entities.AppDescriptor, derrors.Error) {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
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
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
		}else {
			return nil, derrors.AsError(err, "cannot get appDescriptor")
		}
	}

	return &descriptor, nil
}

// DescriptorExists checks if a given descriptor exists on the system.
func (sp *ScyllaApplicationProvider) DescriptorExists(appDescriptorID string) (bool, derrors.Error) {

	// check connection
	if  err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(applicationDescriptorTable).Columns(applicationDescriptorTablePK).Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else{
			return false, derrors.AsError(err, "cannot determinate if appDescriptor exists")
		}
	}

	return true, nil
}

// DeleteDescriptor removes a given descriptor from the system.
func (sp *ScyllaApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.DescriptorExists(appDescriptorID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}

	// delete app instance
	stmt, _ := qb.Delete(applicationDescriptorTable).Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appDescriptorID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete appDescriptor")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// AddInstance adds a new application instance to the system
func (sp *ScyllaApplicationProvider) AddInstance(instance entities.AppInstance) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(instance.AppInstanceId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(instance.AppDescriptorId)
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationTable).Columns("organization_id","app_descriptor_id","app_instance_id",
		"name","description","configuration_options","environment_variables","labels","rules","groups","services","status").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appInstance")
	}

	return nil

}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaApplicationProvider) InstanceExists(appInstanceID string) (bool, derrors.Error) {

	// check connection
	if  err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(applicationTable).Columns(applicationTablePK).Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationTablePK: appInstanceID,
	})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		}else {
			return false, derrors.AsError(err, "cannot determinate if appInstance exists")
		}
	}

	return true, nil

}

// GetInstance retrieves an application instance.
func (sp *ScyllaApplicationProvider) GetInstance(appInstanceID string) (* entities.AppInstance, derrors.Error) {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	// 2.- Gocqlx
	var app entities.AppInstance
	stmt, names := qb.Select(applicationTable).Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationTablePK: appInstanceID,
	})

	err := q.GetRelease(&app)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("instance").WithParams(appInstanceID)
		}else{
			return nil, derrors.AsError(err,"cannot get appInstance")
		}
	}

	return &app, nil

}

// DeleteInstance removes a given instance from the system.
func (sp *ScyllaApplicationProvider) DeleteInstance(appInstanceID string) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(appInstanceID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("instance").WithParams(appInstanceID)
	}

	// delete app instance
	stmt, _ := qb.Delete(applicationTable).Where(qb.Eq(applicationTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete appInstance")
	}

	return nil
}

// UpdateInstance updates the information of an instance
func (sp *ScyllaApplicationProvider) UpdateInstance(instance entities.AppInstance) derrors.Error {
	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.InstanceExists(instance.AppInstanceId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("instance").WithParams(instance.AppInstanceId)
	}

	// update the application instance
	stmt, names := qb.Update(applicationTable).Set("organization_id","app_descriptor_id",
		"name","description","configuration_options","environment_variables",
		"labels","rules","groups","services","status").Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update appInstance")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// Clear descriptors and instances
func (sp *ScyllaApplicationProvider) Clear() derrors.Error {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// delete app instances
	err := sp.Session.Query("TRUNCATE TABLE applicationinstances").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the applications table")
		return derrors.AsError(err, "cannot truncate applicationinstace table")
	}

	// TRUNCATE TABLE applicationdescriptors
	err = sp.Session.Query("TRUNCATE TABLE applicationdescriptors").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the applications descriptor table")
		return derrors.AsError(err, "cannot truncate applicationdescriptors table")
	}
	return nil


}

