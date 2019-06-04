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
	"time"
)


//const addInstance = "INSERT INTO applications (app_instance_id, app_descriptor_id, configuration_options, description, environment_variables, groups, labels, name, organization_id, rules, services, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
//const selectInstance = "SELECT organization_id, app_descriptor_id, name, description, configuration_options, environment_variables, labels, rules, groups, services, status FROM applications WHERE app_instance_id = ?"

const applicationTable = "ApplicationInstances"
const applicationTablePK = "app_instance_id"

const applicationDescriptorTable = "applicationdescriptors"
const applicationDescriptorTablePK = "app_descriptor_id"

const parametrizedDescriptorTable = "parametrizeddescriptors"
const parametrizedDescriptorTablePK = "app_instance_id"

const instanceParamTable = "instanceParameters" // (app_instance_id, parameters)
const instanceParamTablePK = "app_instance_id"

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


	var count int

	stmt, names := qb.Select(applicationDescriptorTable).CountAll().Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&count)
	if err != nil {
		return false, derrors.AsError(err, "cannot determinate if appDescriptor exists")
	}

	return count == 1, nil
}

func (sp *ScyllaApplicationProvider) unsafeInstanceExists(appInstanceID string) (bool, derrors.Error) {

	var count int

	stmt, names := qb.Select(applicationTable).CountAll().Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationTablePK: appInstanceID,
	})

	err := q.GetRelease(&count)
	if err != nil {
			return false, derrors.AsError(err, "cannot determinate if appInstance exists")
	}

	return count == 1, nil

}

func (sp *ScyllaApplicationProvider) unsafeInstanceParametersExists(appInstanceID string) (bool, derrors.Error) {


	var count int

	stmt, names := qb.Select(instanceParamTable).CountAll().Where(qb.Eq(instanceParamTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		instanceParamTablePK: appInstanceID,
	})

	err := q.GetRelease(&count)
	if err != nil {

			return false, derrors.AsError(err, "cannot determinate if instance parameters exists")

	}

	return count == 1, nil
}

func (sp *ScyllaApplicationProvider) unsafeParametrizedDescriptorExists(appInstanceID string) (bool, derrors.Error) {


	var count int

	stmt, names := qb.Select(parametrizedDescriptorTable).CountAll().Where(qb.Eq(parametrizedDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		parametrizedDescriptorTablePK: appInstanceID,
	})

	err := q.GetRelease(&count)
	if err != nil {

		return false, derrors.AsError(err, "cannot determinate if parametrized descriptor exists")

	}

	return count > 0, nil
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
		"configuration_options","environment_variables","labels","rules","groups", "parameters").ToCql()
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

func (sp *ScyllaApplicationProvider) GetDescriptorParameters(appDescriptorID string) ([]entities.Parameter, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 2.- Gocqlx
	var parameters []entities.Parameter
	stmt, names := qb.Select(applicationDescriptorTable).Columns("parameters").Where(qb.Eq(applicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		applicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&parameters)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
		}else {
			return nil, derrors.AsError(err, "cannot get descriptor parameters")
		}
	}
	if parameters == nil {
		parameters = make([]entities.Parameter, 0)
	}
	return parameters, nil
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
	// TODO: parameters can not be updated, review if that is true
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
		"name","configuration_options","environment_variables","labels","metadata","rules","groups", "status").ToCql()
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
		"name","configuration_options","environment_variables","labels","rules","groups", "status","info").Where(qb.Eq(applicationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(instance)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot update appInstance")
	}

	return nil
}

// ------------------------------------------- //
// -- Instance parameters -------------------- //
// ------------------------------------------- //
// AddInstanceParameters adds deploy parameters of an instance in the system
func (sp *ScyllaApplicationProvider)AddInstanceParameters (appInstanceID string, parameters []entities.InstanceParameter) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return  err
	}

	exists, err := sp.unsafeInstanceParametersExists(appInstanceID)
	if err != nil{
		return  err
	}

	if exists{
		return derrors.NewAlreadyExistsError("parameters").WithParams(appInstanceID)
	}

	// insert the application instance
	stmt, names := qb.Insert(instanceParamTable).Columns("app_instance_id", "parameters").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"app_instance_id": appInstanceID,
		"parameters": parameters,
	})
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add instance parameters")
	}

	return nil

	return nil
}
// GetInstanceParameters retrieves the params of an instance
func (sp *ScyllaApplicationProvider) GetInstanceParameters (appInstanceID string) ([]entities.InstanceParameter, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 2.- Gocqlx
	var parameters []entities.InstanceParameter
	stmt, names := qb.Select(instanceParamTable).Columns("parameters").Where(qb.Eq(instanceParamTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		instanceParamTablePK: appInstanceID,
	})

	err := q.GetRelease(&parameters)
	if err != nil {
		if err.Error() == rowNotFound {
			parameters := make ([]entities.InstanceParameter, 0)
			return parameters, nil
		}else{
			return nil, derrors.AsError(err,"cannot get instance parameters")
		}
	}

	return parameters, nil

}


// DeleteInstanceParameters removes the params of an instance
func (sp *ScyllaApplicationProvider) DeleteInstanceParameters (appInstanceID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if  err != nil {
		return err
	}

	stmt, _ := qb.Delete(instanceParamTable).Where(qb.Eq(instanceParamTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete instance parameters")
	}

	return nil
}

// --------------------------------- //
// ------ Parametrized Descriptor -- //
// --------------------------------- //
// AddParametrizedDescriptor adds a new parametrized descriptor to the system.
func (sp *ScyllaApplicationProvider)AddParametrizedDescriptor(descriptor entities.ParametrizedDescriptor) derrors.Error {

	sp.Lock()
	defer sp.Unlock()


	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeParametrizedDescriptorExists(descriptor.AppInstanceId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(descriptor.AppInstanceId)
	}

	// insert the application instance
	stmt, names := qb.Insert(parametrizedDescriptorTable).Columns("organization_id","app_descriptor_id",
		"app_instance_id", "name", "configuration_options","environment_variables","labels","rules","groups").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(descriptor)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add parametrized Descriptor")
	}

	return nil
}

// GetParametrizedDescriptor retrieves a parametrized descriptor
func (sp * ScyllaApplicationProvider) GetParametrizedDescriptor(appInstanceID string) (*entities.ParametrizedDescriptor, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 2.- Gocqlx
	var parametrized entities.ParametrizedDescriptor
	stmt, names := qb.Select(parametrizedDescriptorTable).Where(qb.Eq(parametrizedDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		parametrizedDescriptorTablePK: appInstanceID,
	})

	err := q.GetRelease(&parametrized)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("parametrized descriptor").WithParams(appInstanceID)
		}else{
			return nil, derrors.AsError(err,"cannot get parametrized Descriptor")
		}
	}

	return &parametrized, nil
}

// ParametrizedDescriptorExists checks if a parametrized descriptor exists on the system.
func (sp * ScyllaApplicationProvider) ParametrizedDescriptorExists (appInstanceID string) (*bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if  err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	exists, err := sp.unsafeParametrizedDescriptorExists(appInstanceID)
	if err != nil {
		return nil, err
	}
	return &exists, nil
}

// DeleteParametrizedDescriptor removes a parametrized Descriptor from the system
func (sp * ScyllaApplicationProvider)DeleteParametrizedDescriptor (appInstanceID string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if  err != nil {
		return err
	}

	// check if the application exists
	exists, err := sp.unsafeParametrizedDescriptorExists(appInstanceID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("parametrized descriptor").WithParams(appInstanceID)
	}

	// delete app instance
	stmt, _ := qb.Delete(parametrizedDescriptorTable).Where(qb.Eq(parametrizedDescriptorTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete parametrized descriptor")
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

	err = sp.Session.Query("TRUNCATE TABLE AppEntrypoints").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the applications endpoints table")
		return derrors.AsError(err, "cannot truncate AppEntrypoints table")
	}

	err = sp.Session.Query("TRUNCATE TABLE appztnetworks").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the zt networks table")
		return derrors.AsError(err, "cannot truncate AppZtNetworks table")
	}

	err = sp.Session.Query("TRUNCATE TABLE appztnetworkmembers").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the zt network members table")
		return derrors.AsError(err, "cannot truncate AppZtNetworkMembers table")
	}

	// table instance Parameters
	err = sp.Session.Query("truncate table instanceparameters").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the instance parameters table")
		return derrors.AsError(err, "cannot truncate instanceparameters table")
	}

	// table parametrized descriptor
	err = sp.Session.Query("truncate table parametrizeddescriptors").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the parametrized descriptors table")
		return derrors.AsError(err, "cannot truncate parametrizeddescriptors table")
	}

	return nil
}


// ---------------------------------------------------------------------------------------------------------------------
// AddAppEndPoint adds a new entry point to the system
func (sp *ScyllaApplicationProvider) AddAppEndpoint (appEntryPoint entities.AppEndpoint) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// insert the appEntryPoint
	stmt, names := qb.Insert("appentrypoints").Columns("organization_id","app_instance_id","service_group_instance_id",
		"service_instance_id","port","endpoint_instance_id","fqdn","global_fqdn","protocol","type").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(appEntryPoint)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appEntryPoint")
	}

	return nil
}

// GetAppEndPointByFQDN ()
func (sp *ScyllaApplicationProvider) GetAppEndpointByFQDN(fqdn string) ([]*entities.AppEndpoint, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select("appentrypoints").Columns("organization_id", "app_instance_id", "service_group_instance_id",
		"service_instance_id", "port", "endpoint_instance_id", "fqdn", "global_fqdn", "protocol", "type").
		Where(qb.Eq("global_fqdn")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"global_fqdn": fqdn,
	})

	entrypoints := make([]*entities.AppEndpoint, 0)
	cqlErr := gocqlx.Select(&entrypoints, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list App entrypoints")
	}

	return entrypoints, nil

}

func (sp *ScyllaApplicationProvider) DeleteAppEndpoints(organizationID string, appInstanceID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil{
		return  err
	}
	// delete app instance
	stmt, _ := qb.Delete("appentrypoints").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete app endpoints")
	}
	return nil
}

func (sp *ScyllaApplicationProvider) GetAppEndpointList(organizationID string , appInstanceId string,
	serviceGroupInstanceID string) ([]*entities.AppEndpoint, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	list := make([]*entities.AppEndpoint, 0)

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select("appentrypoints").
		Columns("organization_id", "app_instance_id", "service_group_instance_id",
			"service_instance_id", "port", "endpoint_instance_id", "fqdn", "global_fqdn", "protocol", "type").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":           organizationID,
		"app_instance_id":           appInstanceId,
		"service_group_instance_id": serviceGroupInstanceID,
	})

	cqlErr := gocqlx.Select(&list, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list app endPoint of a service group")
	}

	return list, nil
}
// ---------------------------------------------------------------------------------------------------------------------
// AppZtNetwork related methods

func (sp *ScyllaApplicationProvider) AddAppZtNetwork(ztNetwork entities.AppZtNetwork) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// add the zt network
	stmt, names := qb.Insert("appztnetworks").Columns("organization_id","app_instance_id","zt_network_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(ztNetwork)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appEntryPoint")
	}

	return nil
}


func (sp *ScyllaApplicationProvider) RemoveAppZtNetwork(organizationID string, appInstanceID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// delete an instance
	stmt, _ := qb.Delete("appztnetworks").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete app zt network")
	}

	return nil
}

func (sp *ScyllaApplicationProvider) GetAppZtNetwork(organizationId string, appInstanceId string) (*entities.AppZtNetwork, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	stmt, names := qb.Select("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"app_instance_id": appInstanceId,
	})

	var ztNetwork entities.AppZtNetwork
	cqlErr := gocqlx.Get(&ztNetwork, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("appZtNetwork").WithParams(organizationId).WithParams(appInstanceId)
		}else {
			return nil, derrors.AsError(err, "cannot get appZtNetwork")
		}
	}

	return &ztNetwork, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// AppZtNetworkMembers related methods

// AddZtNetworkMember add a new member for an existing zt network
func (sp *ScyllaApplicationProvider) AddAppZtNetworkMember(member entities.AppZtNetworkMembers) (*entities.AppZtNetworkMembers, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	// if we already have an entry, simply add a new member
	stmt, names := qb.Select("appztnetworkmembers").Columns("organization_id", "app_instance_id",
		"service_group_instance_id", "service_application_instance_id", "zt_network_id","members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).
		Where(qb.Eq("zt_network_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": member.OrganizationId,
		"app_instance_id": member.AppInstanceId,
		"service_group_instance_id": member.ServiceGroupInstanceId,
		"service_application_instance_id": member.ServiceApplicationInstanceId,
		"zt_network_id": member.ZtNetworkId,
	})

	var retrievedMembers entities.AppZtNetworkMembers
	cqlErr := gocqlx.Get(&retrievedMembers, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			// insert a new row
			stmt, names := qb.Insert("appztnetworkmembers").Columns("organization_id", "app_instance_id",
				"service_group_instance_id", "service_application_instance_id", "zt_network_id","members").ToCql()
			// update creation time
			for k, v := range member.Members {
				v.CreatedAt = time.Now().Unix()
				member.Members[k] = v
			}

			log.Debug().Interface("members", member).Msg("members to write")

			q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(member)
			cqlErr := q.Exec()
			if cqlErr != nil {
				return nil, derrors.NewInternalError("appZtNetworkMembers",err).WithParams(member.OrganizationId).
					WithParams(member.AppInstanceId).WithParams(member.ServiceGroupInstanceId).WithParams(member.AppInstanceId).
					WithParams(member.ZtNetworkId)
			}
			return &member, nil
		}else {
			return nil, derrors.AsError(err, "cannot get appZtNetworkMembers")
		}
	}

	log.Debug().Interface("retrieved",retrievedMembers).Msg("retrieved members from DB")

	// update the map
	for k,v := range member.Members {
		newEntry := v
		newEntry.CreatedAt = time.Now().Unix()
		retrievedMembers.Members[k] = newEntry
	}

	log.Debug().Interface("members", retrievedMembers).Msg("members to update")


	// add the zt network member
	stmt, names = qb.Update("appztnetworkmembers").Set("members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).
		Where(qb.Eq("zt_network_id")).ToCql()
	q = gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(retrievedMembers)
	cqlErr = q.ExecRelease()

	if cqlErr != nil {
		return nil,derrors.AsError(cqlErr, "cannot update app zt network member")
	}

	return &member, nil

}

// RemoveZtNetworkMember remove an existing member for a zt network
func (sp *ScyllaApplicationProvider) RemoveAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstanceId string, ztNetworkId string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// delete an instance
	stmt, _ := qb.Delete("appztnetworkmembers").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).
		Where(qb.Eq("zt_network_id")).ToCql()
	query := sp.Session.Query(stmt, organizationId, appInstanceId, serviceGroupInstanceId, serviceInstanceId,ztNetworkId)
	cqlErr := query.Exec()

	log.Debug().Str("query", query.String()).Msg("generated sql statement")



	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete an app zt network member")
	}
	return nil
}

func (sp *ScyllaApplicationProvider) GetAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceApplicationInstanceId string) (*entities.AppZtNetworkMembers, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return nil, err
	}

	// if we already have an entry, simply add a new member
	stmt, names := qb.Select("appztnetworkmembers").Columns("organization_id", "app_instance_id",
		"service_group_instance_id", "service_application_instance_id", "zt_network_id","members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"app_instance_id": appInstanceId,
		"service_group_instance_id": serviceGroupInstanceId,
		"service_application_instance_id": serviceApplicationInstanceId,
	})

	var retrievedMembers entities.AppZtNetworkMembers
	cqlErr := gocqlx.Get(&retrievedMembers, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("get app zt network member")
		}else{
			return nil, derrors.AsError(err,"cannot get app zt network member")
		}
	}
	return &retrievedMembers, nil

}

