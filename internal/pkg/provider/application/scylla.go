package application

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
	"time"
)

// ------------------------------------

// Application descriptor const
const ApplicationDescriptorTable = "ApplicationDescriptors"
const ApplicationDescriptorTablePK = "app_descriptor_id"

var allApplicationDecriptorColumns = []string{"organization_id", "app_descriptor_id", "name", "configuration_options",
	"environment_variables", "labels", "rules", "groups", "parameters", "inbound_net_interfaces", "outbound_net_interfaces"}

// 'parameter' column is not included in allApplicationDecriptorColumnsNoPK because the value can not be updated
var allApplicationDecriptorColumnsNoPK = []string{"organization_id", "name", "configuration_options", "environment_variables",
	"labels", "rules", "groups", "inbound_net_interfaces", "outbound_net_interfaces"}

// Application Instance const
const ApplicationInstanceTable = "ApplicationInstances"
const ApplicationInstanceTablePK = "app_instance_id"

var allApplicationInstanceColumns = []string{"organization_id", "app_descriptor_id", "app_instance_id",
	"name", "configuration_options", "environment_variables", "labels", "metadata", "rules", "groups", "status",
	"inbound_net_interfaces", "outbound_net_interfaces", "info"}
var allApplicationInstanceColumnsNoPK = []string{"organization_id", "app_descriptor_id",
	"name", "configuration_options", "environment_variables", "labels", "metadata", "rules", "groups", "status",
	"inbound_net_interfaces", "outbound_net_interfaces", "info"}

// Parametrized Descriptor const
const ParametrizedDescriptorTable = "ParametrizedDescriptors"
const ParametrizedDescriptorTablePK = "app_instance_id"

var allParametrizedDescriptorColumns = []string{"organization_id", "app_descriptor_id",
	"app_instance_id", "name", "configuration_options", "environment_variables", "labels", "rules", "groups",
	"inbound_net_interfaces", "outbound_net_interfaces"}

// Parameters const
const InstanceParamTable = "InstanceParameters" // (app_instance_id, parameters)
const InstanceParamTablePK = "app_instance_id"

var allInstanceParamColumns = []string{"app_instance_id", "parameters"}

const AppEndpointsTable = "AppEntrypoints"

var allAppEndPointsColumns = []string{"organization_id", "app_instance_id", "service_group_instance_id",
	"service_instance_id", "port", "endpoint_instance_id", "fqdn", "global_fqdn", "protocol", "type"}

const AppZtNetworkTable = "appztnetworks"

// ------------------------------------

const rowNotFound = "not found"

type ScyllaApplicationProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaApplicationProvider(address string, port int, keyspace string) *ScyllaApplicationProvider {
	provider := ScyllaApplicationProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

func (sp *ScyllaApplicationProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

	sp.ScyllaDB.Disconnect()
}

// ---------------------------------------------- //
// -- Application Descriptor -------------------- //
// ---------------------------------------------- //

// AddDescriptor adds a new application descriptor to the system
func (sp *ScyllaApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {

	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeAdd(ApplicationDescriptorTable, ApplicationDescriptorTablePK, descriptor.AppDescriptorId, allApplicationDecriptorColumns, descriptor)
}

// GetDescriptors retrieves an application descriptor.
func (sp *ScyllaApplicationProvider) GetDescriptor(appDescriptorID string) (*entities.AppDescriptor, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var appDescriptor interface{} = &entities.AppDescriptor{}

	err := sp.UnsafeGet(ApplicationDescriptorTable, ApplicationDescriptorTablePK, appDescriptorID, allApplicationDecriptorColumns, &appDescriptor)
	if err != nil {
		return nil, err
	}
	return appDescriptor.(*entities.AppDescriptor), nil
}

func (sp *ScyllaApplicationProvider) GetDescriptorParameters(appDescriptorID string) ([]entities.Parameter, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	// 2.- Gocqlx
	var parameters []entities.Parameter
	stmt, names := qb.Select(ApplicationDescriptorTable).Columns("parameters").Where(qb.Eq(ApplicationDescriptorTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		ApplicationDescriptorTablePK: appDescriptorID,
	})

	err := q.GetRelease(&parameters)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
		} else {
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

	return sp.UnsafeGenericExist(ApplicationDescriptorTable, ApplicationDescriptorTablePK, appDescriptorID)
}

// UpdateDescriptor updates the information of an application descriptor.
func (sp *ScyllaApplicationProvider) UpdateDescriptor(descriptor entities.AppDescriptor) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// TODO: parameters can not be updated, review if that is true
	return sp.UnsafeUpdate(ApplicationDescriptorTable, ApplicationDescriptorTablePK, descriptor.AppDescriptorId, allApplicationDecriptorColumnsNoPK, descriptor)
}

// DeleteDescriptor removes a given descriptor from the system.
func (sp *ScyllaApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(ApplicationDescriptorTable, ApplicationDescriptorTablePK, appDescriptorID)
}

// -------------------------------------------- //
// -- Application Instance -------------------- //
// -------------------------------------------- //
// AddInstance adds a new application instance to the system
func (sp *ScyllaApplicationProvider) AddInstance(instance entities.AppInstance) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeAdd(ApplicationInstanceTable, ApplicationInstanceTablePK, instance.AppInstanceId, allApplicationInstanceColumns, instance)
}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaApplicationProvider) InstanceExists(appInstanceID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(ApplicationInstanceTable, ApplicationInstanceTablePK, appInstanceID)

}

// GetInstance retrieves an application instance.
func (sp *ScyllaApplicationProvider) GetInstance(appInstanceID string) (*entities.AppInstance, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var appInstance interface{} = &entities.AppInstance{}

	err := sp.UnsafeGet(ApplicationInstanceTable, ApplicationInstanceTablePK, appInstanceID, allApplicationInstanceColumns, &appInstance)
	if err != nil {
		return nil, err
	}
	return appInstance.(*entities.AppInstance), nil

}

// DeleteInstance removes a given instance from the system.
func (sp *ScyllaApplicationProvider) DeleteInstance(appInstanceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(ApplicationInstanceTable, ApplicationInstanceTablePK, appInstanceID)
}

// UpdateInstance updates the information of an instance
func (sp *ScyllaApplicationProvider) UpdateInstance(instance entities.AppInstance) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeUpdate(ApplicationInstanceTable, ApplicationInstanceTablePK, instance.AppInstanceId, allApplicationInstanceColumnsNoPK, instance)
}

// ------------------------------------------- //
// -- Instance parameters -------------------- //
// ------------------------------------------- //

type InstanceParameterRecord struct {
	// AppInstanceId with the application instance identifier.
	AppInstanceId string                       `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	Parameters    []entities.InstanceParameter `json:"parameters,omitempty" cql:"parameters" `
}

// TODO: Check this method works
// AddInstanceParameters adds deploy parameters of an instance in the system
func (sp *ScyllaApplicationProvider) AddInstanceParameters(appInstanceID string, parameters []entities.InstanceParameter) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	return sp.UnsafeAdd(InstanceParamTable, InstanceParamTablePK, appInstanceID, allInstanceParamColumns, InstanceParameterRecord{appInstanceID, parameters})

}

// GetInstanceParameters retrieves the params of an instance
func (sp *ScyllaApplicationProvider) GetInstanceParameters(appInstanceID string) ([]entities.InstanceParameter, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var parametersRecord interface{} = &InstanceParameterRecord{}

	err := sp.UnsafeGet(InstanceParamTable, InstanceParamTablePK, appInstanceID, allInstanceParamColumns, &parametersRecord)
	if err != nil {
		if err.Type() == derrors.NotFound {
			return []entities.InstanceParameter{}, nil
		} else {
			return nil, err
		}
	}

	return parametersRecord.(*InstanceParameterRecord).Parameters, nil

}

// DeleteInstanceParameters removes the params of an instance
func (sp *ScyllaApplicationProvider) DeleteInstanceParameters(appInstanceID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	err := sp.UnsafeRemove(InstanceParamTable, InstanceParamTablePK, appInstanceID)

	// should not fail when deleting the parameters of an instance (which do not exist)
	if err != nil {
		if err.Type() == derrors.NotFound {
			return nil
		}
		return err
	}
	return nil
}

// --------------------------------- //
// ------ Parametrized Descriptor -- //
// --------------------------------- //
// AddParametrizedDescriptor adds a new parametrized descriptor to the system.
func (sp *ScyllaApplicationProvider) AddParametrizedDescriptor(descriptor entities.ParametrizedDescriptor) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeAdd(ParametrizedDescriptorTable, ParametrizedDescriptorTablePK, descriptor.AppInstanceId, allParametrizedDescriptorColumns, descriptor)

}

// GetParametrizedDescriptor retrieves a parametrized descriptor
func (sp *ScyllaApplicationProvider) GetParametrizedDescriptor(appInstanceID string) (*entities.ParametrizedDescriptor, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	var paramDescriptor interface{} = &entities.ParametrizedDescriptor{}

	err := sp.UnsafeGet(ParametrizedDescriptorTable, ParametrizedDescriptorTablePK, appInstanceID, allParametrizedDescriptorColumns, &paramDescriptor)
	if err != nil {
		return nil, err
	}
	return paramDescriptor.(*entities.ParametrizedDescriptor), nil

}

// ParametrizedDescriptorExists checks if a parametrized descriptor exists on the system.
func (sp *ScyllaApplicationProvider) ParametrizedDescriptorExists(appInstanceID string) (*bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(ParametrizedDescriptorTable, ParametrizedDescriptorTablePK, appInstanceID)

	return &exists, err
}

// DeleteParametrizedDescriptor removes a parametrized Descriptor from the system
func (sp *ScyllaApplicationProvider) DeleteParametrizedDescriptor(appInstanceID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(ParametrizedDescriptorTable, ParametrizedDescriptorTablePK, appInstanceID)
}

// ---------------------------------------------------------------------------------------------------------------------

// Clear descriptors and instances
func (sp *ScyllaApplicationProvider) Clear() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{ApplicationDescriptorTable, ApplicationInstanceTable, ParametrizedDescriptorTable, InstanceParamTable,
		AppEndpointsTable, AppZtNetworkTable})

	err := sp.Session.Query("TRUNCATE TABLE appztnetworkmembers").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the zt network members table")
		return derrors.AsError(err, "cannot truncate AppZtNetworkMembers table")
	}

	return nil
}

// ------------------------------------ //
// -- AppEndpoints -------------------- //
// ------------------------------------ //

func (sp *ScyllaApplicationProvider) createAppEndpointPKMap(organizationID string, appInstanceID string,
	service_group_id string, service_instance_id string, port int32, protocol entities.AppEndpointProtocol) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":           organizationID,
		"app_instance_id":           appInstanceID,
		"service_group_instance_id": service_group_id,
		"service_instance_id":       service_instance_id,
		"port":                      port,
		"protocol":                  protocol,
	}
}

func (sp *ScyllaApplicationProvider) createShortAppEndpointPKMap(organizationID string, appInstanceID string) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID,
	}
}

// AddAppEndPoint adds a new entry point to the system
func (sp *ScyllaApplicationProvider) AddAppEndpoint(appEndPoint entities.AppEndpoint) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// insert the endpoint
	stmt, names := qb.Insert(AppEndpointsTable).Columns(allAppEndPointsColumns...).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(appEndPoint)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add end point")
	}
	return nil
}

// GetAppEndPointByFQDN ()
func (sp *ScyllaApplicationProvider) GetAppEndpointByFQDN(fqdn string) ([]*entities.AppEndpoint, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AppEndpointsTable).Columns(allAppEndPointsColumns...).
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

	return sp.UnsafeCompositeRemove(AppEndpointsTable, sp.createShortAppEndpointPKMap(organizationID, appInstanceID))
}

func (sp *ScyllaApplicationProvider) GetAppEndpointList(organizationID string, appInstanceId string,
	serviceGroupInstanceID string) ([]*entities.AppEndpoint, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	list := make([]*entities.AppEndpoint, 0)

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AppEndpointsTable).
		Columns(allAppEndPointsColumns...).
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

// TODO: no changes apply in these methods because the ZT is going to disappear
// ---------------------------------------------------------------------------------------------------------------------
// AppZtNetwork related methods

func (sp *ScyllaApplicationProvider) AddAppZtNetwork(ztNetwork entities.AppZtNetwork) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// add the zt network
	stmt, names := qb.Insert("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").ToCql()
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
	err := sp.CheckAndConnect()
	if err != nil {
		return nil, err
	}

	stmt, names := qb.Select("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").
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
		} else {
			return nil, derrors.AsError(err, "cannot get appZtNetwork")
		}
	}

	return &ztNetwork, nil
}

// AddZtNetworkProxy add a zt service proxy
func (sp *ScyllaApplicationProvider) AddZtNetworkProxy(proxy entities.ServiceProxy) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// find the service proxy
	stmt, names := qb.Select("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": proxy.OrganizationId,
		"app_instance_id": proxy.AppInstanceId,
	})

	var ztNetwork entities.AppZtNetwork
	cqlErr := gocqlx.Get(&ztNetwork, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return derrors.NewNotFoundError("appZtNetworks").WithParams(proxy.OrganizationId).WithParams(proxy.AppInstanceId)
		} else {
			return derrors.AsError(err, "cannot get appZtNetwork")
		}
	}

	if ztNetwork.AvailableProxies == nil {
		ztNetwork.AvailableProxies = make(map[string]map[string][]entities.ServiceProxy, 0)
	}

	// Add the proxy or overwrite if it is already there
	existingProxies, found := ztNetwork.AvailableProxies[proxy.FQDN]
	if !found {
		aux := map[string][]entities.ServiceProxy{
			proxy.ClusterId: []entities.ServiceProxy{proxy},
		}
		ztNetwork.AvailableProxies[proxy.FQDN] = aux
	} else {
		// search for the entries
		clusterEntries, found := existingProxies[proxy.ClusterId]
		if !found {
			existingProxies[proxy.ClusterId] = []entities.ServiceProxy{proxy}
		} else {
			// add it to the list
			clusterEntries = append(clusterEntries, proxy)
			existingProxies[proxy.ClusterId] = clusterEntries
		}
		ztNetwork.AvailableProxies[proxy.FQDN] = existingProxies
	}

	// update the network proxy entry
	stmt, names = qb.Insert("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").ToCql()
	q = gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(ztNetwork)
	cqlErr = q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appEntryPoint")
	}
	return nil
}

// RemoveZtNetworkProxy remove an existing zt service proxy
func (sp *ScyllaApplicationProvider) RemoveZtNetworkProxy(organizationId string, appInstanceId string, fqdn string, clusterId string, serviceGroupInstanceId string, serviceInstanceId string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// find the service proxy
	stmt, names := qb.Select("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationId,
		"app_instance_id": appInstanceId,
	})

	var ztNetwork entities.AppZtNetwork
	cqlErr := gocqlx.Get(&ztNetwork, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return derrors.NewNotFoundError("appZtNetworks").WithParams(organizationId).WithParams(appInstanceId)
		} else {
			return derrors.AsError(err, "cannot get appZtNetwork")
		}
	}

	// remove it
	existingProxies, found := ztNetwork.AvailableProxies[fqdn]
	if !found {
		return derrors.NewNotFoundError(fmt.Sprintf("impossible to find proxy for fqdn %s", fqdn))
	} else {
		// search for the entries
		clusterEntries, found := existingProxies[clusterId]
		if !found {
			return derrors.NewNotFoundError(fmt.Sprintf("impossible to find proxy for fqdn %s in cluster %s", fqdn, clusterId))
		} else {
			// look for it and remove it
			indexToDelete := -1
			for i, proxy := range clusterEntries {
				if proxy.ServiceInstanceId == serviceInstanceId && proxy.ServiceGroupInstanceId == serviceInstanceId {
					indexToDelete = i
					break
				}
			}
			if indexToDelete == -1 {
				return derrors.NewNotFoundError(fmt.Sprintf("impossible to find proxy for fqdn %s in cluster %s with serviceInstanceId %s",
					fqdn, clusterId, serviceInstanceId))
			}
			if len(clusterEntries) == 1 {
				// remove this cluster entry
				delete(ztNetwork.AvailableProxies[fqdn], clusterId)
			} else {
				ztNetwork.AvailableProxies[fqdn][clusterId] = append(clusterEntries[:indexToDelete], clusterEntries[indexToDelete+1:]...)
			}
		}
		if len(ztNetwork.AvailableProxies[fqdn]) == 0 {
			delete(ztNetwork.AvailableProxies, fqdn)
		}
	}

	// update
	// update the network proxy entry
	stmt, names = qb.Insert("appztnetworks").Columns("organization_id", "app_instance_id", "zt_network_id", "vsa_list", "available_proxies").ToCql()
	q = gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(ztNetwork)
	cqlErr = q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add appEntryPoint")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// AppZtNetworkMembers related methods

// AddZtNetworkMember add a new member for an existing zt network
func (sp *ScyllaApplicationProvider) AddAppZtNetworkMember(member entities.AppZtNetworkMembers) (*entities.AppZtNetworkMembers, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return nil, err
	}

	// if we already have an entry, simply add a new member
	stmt, names := qb.Select("appztnetworkmembers").Columns("organization_id", "app_instance_id",
		"service_group_instance_id", "service_application_instance_id", "zt_network_id", "members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).
		Where(qb.Eq("zt_network_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":                 member.OrganizationId,
		"app_instance_id":                 member.AppInstanceId,
		"service_group_instance_id":       member.ServiceGroupInstanceId,
		"service_application_instance_id": member.ServiceApplicationInstanceId,
		"zt_network_id":                   member.ZtNetworkId,
	})

	var retrievedMembers entities.AppZtNetworkMembers
	cqlErr := gocqlx.Get(&retrievedMembers, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			// insert a new row
			stmt, names := qb.Insert("appztnetworkmembers").Columns("organization_id", "app_instance_id",
				"service_group_instance_id", "service_application_instance_id", "zt_network_id", "members").ToCql()
			// update creation time
			for k, v := range member.Members {
				v.CreatedAt = time.Now().Unix()
				member.Members[k] = v
			}

			q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(member)
			cqlErr := q.Exec()
			if cqlErr != nil {
				return nil, derrors.NewInternalError("appZtNetworkMembers", err).WithParams(member.OrganizationId).
					WithParams(member.AppInstanceId).WithParams(member.ServiceGroupInstanceId).WithParams(member.AppInstanceId).
					WithParams(member.ZtNetworkId)
			}
			return &member, nil
		} else {
			return nil, derrors.AsError(err, "cannot get appZtNetworkMembers")
		}
	}

	// update the map
	for k, v := range member.Members {
		newEntry := v
		newEntry.CreatedAt = time.Now().Unix()
		retrievedMembers.Members[k] = newEntry
	}

	// add the zt network member
	stmt, names = qb.Update("appztnetworkmembers").Set("members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).
		Where(qb.Eq("zt_network_id")).ToCql()
	q = gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(retrievedMembers)
	cqlErr = q.ExecRelease()

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot update app zt network member")
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
	query := sp.Session.Query(stmt, organizationId, appInstanceId, serviceGroupInstanceId, serviceInstanceId, ztNetworkId)
	cqlErr := query.Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete an app zt network member")
	}
	return nil
}

func (sp *ScyllaApplicationProvider) GetAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceApplicationInstanceId string) (*entities.AppZtNetworkMembers, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return nil, err
	}

	// if we already have an entry, simply add a new member
	stmt, names := qb.Select("appztnetworkmembers").Columns("organization_id", "app_instance_id",
		"service_group_instance_id", "service_application_instance_id", "zt_network_id", "members").
		Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).
		Where(qb.Eq("service_group_instance_id")).Where(qb.Eq("service_application_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":                 organizationId,
		"app_instance_id":                 appInstanceId,
		"service_group_instance_id":       serviceGroupInstanceId,
		"service_application_instance_id": serviceApplicationInstanceId,
	})

	var retrievedMembers entities.AppZtNetworkMembers
	cqlErr := gocqlx.Get(&retrievedMembers, q.Query)

	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("get app zt network member")
		} else {
			return nil, derrors.AsError(err, "cannot get app zt network member")
		}
	}
	return &retrievedMembers, nil

}
