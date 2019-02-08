package application

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
	sync.Mutex
}

func NewScyllaApplicationProvider (address string, port int, keyspace string) * ScyllaApplicationProvider {
	provider := ScyllaApplicationProvider{Address:address, Port: port, Keyspace: keyspace, Session: nil}
	provider.connect()
	return &provider
}

func (sp *ScyllaApplicationProvider) connect() derrors.Error {

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

	sp.Lock()
	defer sp.Unlock()

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

// check if the session is created
func (sp *ScyllaApplicationProvider) checkConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaApplicationProvider) checkAndConnect () derrors.Error{

	err := sp.checkConnection()
	if err != nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err = sp.connect()
		if err != nil  {
			return err
		}
	}
	return nil
}

func (sp *ScyllaApplicationProvider) unsafeDescriptorExists(appDescriptorID string) (bool, derrors.Error) {


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

func (sp *ScyllaApplicationProvider) unsafeInstanceExists(appInstanceID string) (bool, derrors.Error) {

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
// ---------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new application descriptor to the system
func (sp *ScyllaApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeDescriptorExists(descriptor.AppDescriptorId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(descriptor.AppDescriptorId)
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationDescriptorTable).Columns("organization_id","app_descriptor_id", "name",
		"configuration_options","environment_variables","labels","rules","groups").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(descriptor)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appDescriptor")
	}

	return nil
}

// GetDescriptors retrieves an application descriptor.
func (sp *ScyllaApplicationProvider) GetDescriptor(appDescriptorID string) (* entities.AppDescriptor, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if  err := sp.checkAndConnect(); err != nil {
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

// UpdateDescriptor updates the information of an application descriptor.
func (sp *ScyllaApplicationProvider) UpdateDescriptor(descriptor entities.AppDescriptor) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}
	// check if the descriptor exists
	exists, err := sp.unsafeDescriptorExists(descriptor.AppDescriptorId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(descriptor.AppDescriptorId)
	}
	// insert the application instance
	stmt, names := qb.Update(applicationDescriptorTable).Set("organization_id", "name",
		"configuration_options","environment_variables","labels","rules","groups").Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(descriptor)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update appDescriptor")
	}

	return nil
}

// DeleteDescriptor removes a given descriptor from the system.
func (sp *ScyllaApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeDescriptorExists(appDescriptorID)
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeInstanceExists(instance.AppInstanceId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(instance.AppDescriptorId)
	}

	// insert the application instance
	stmt, names := qb.Insert(applicationTable).Columns("organization_id","app_descriptor_id","app_instance_id",
		"name","configuration_options","environment_variables","labels","rules","groups", "status").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appInstance")
	}

	return nil

}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaApplicationProvider) InstanceExists(appInstanceID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if  err := sp.checkAndConnect(); err != nil {
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeInstanceExists(appInstanceID)
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeInstanceExists(instance.AppInstanceId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("instance").WithParams(instance.AppInstanceId)
	}

	// update the application instance
	stmt, names := qb.Update(applicationTable).Set("organization_id","app_descriptor_id",
		"name","configuration_options","environment_variables","labels","rules","groups", "status").Where(qb.Eq(applicationTablePK)).ToCql()
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

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
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

