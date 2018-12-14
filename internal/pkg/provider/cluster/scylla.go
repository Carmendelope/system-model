package cluster

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const clusterTable = "Clusters"
const clusterTablePK = "cluster_id"
const clusterNodeTable = "Cluster_Nodes"

type ScyllaClusterProvider struct {
	Address  string
	Port     int
	Keyspace string
	Session  *gocql.Session
}

const rowNotFound = "not found"

func NewScyllaClusterProvider(address string, port int, keyspace string) *ScyllaClusterProvider {
	provider := ScyllaClusterProvider{address, port, keyspace, nil}
	provider.Connect()
	return &provider
}

// connect to the database
func (sp *ScyllaClusterProvider) Connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaClusterProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}

// disconnect from the database
func (sp *ScyllaClusterProvider) Disconnect() {

	if sp != nil {
		sp.Session.Close()
	}
}

// check that the session is created
func (sp *ScyllaClusterProvider) CheckConnection() derrors.Error {
	if sp.Session == nil {
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaClusterProvider) CheckAndConnect() derrors.Error {

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err = sp.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// Add a new cluster to the system.
func (sp *ScyllaClusterProvider) Add(cluster entities.Cluster) derrors.Error {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// check if the luster exists
	exists, err := sp.Exists(cluster.ClusterId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(cluster.ClusterId)
	}

	// insert the cluster instance
	stmt, names := qb.Insert(clusterTable).Columns("organization_id", "cluster_id", "name", "description",
		"cluster_type", "hostname", "control_plane_hostname", "multitenant", "status", "labels", "cordon").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(cluster)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add cluster")
	}

	return nil
}

// Update an existing cluster in the system
func (sp *ScyllaClusterProvider) Update(cluster entities.Cluster) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// check if the cluster exists
	exists, err := sp.Exists(cluster.ClusterId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(cluster.ClusterId)
	}

	// insert the cluster instance
	stmt, names := qb.Update(clusterTable).Set("organization_id", "name", "description",
		"cluster_type", "hostname", "multitenant", "control_plane_hostname", "status", "labels", "cordon").
		Where(qb.Eq(clusterTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(cluster)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(err,"cannot update cluster")
	}

	return nil
}

// Exists checks if a cluster exists on the system.
func (sp *ScyllaClusterProvider) Exists(clusterID string) (bool, derrors.Error) {

	if err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(clusterTable).Columns(clusterTablePK).Where(qb.Eq(clusterTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		clusterTablePK: clusterID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if cluster exists")
		}
	}

	return true, nil
}

// Get a cluster.
func (sp *ScyllaClusterProvider) Get(clusterID string) (*entities.Cluster, derrors.Error) {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	var cluster entities.Cluster
	stmt, names := qb.Select(clusterTable).Where(qb.Eq(clusterTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		clusterTablePK: clusterID,
	})

	err := q.GetRelease(&cluster)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("cluster").WithParams(clusterID)
		} else {
			return nil, derrors.AsError(err, "cannot get cluster")
		}
	}

	return &cluster, nil
}

// Remove a cluster
func (sp *ScyllaClusterProvider) Remove(clusterID string) derrors.Error {

	// check if the cluster exists
	exists, err := sp.Exists(clusterID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(clusterID)
	}

	// delete cluster instance
	stmt, _ := qb.Delete(clusterTable).Where(qb.Eq(clusterTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, clusterID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot remove cluster")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddNode adds a new node ID to the cluster.
func (sp *ScyllaClusterProvider) AddNode(clusterID string, nodeID string) derrors.Error {

	exists, err := sp.Exists(clusterID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("node").WithParams(clusterID)
	}

	// check if the node exists in the cluster
	exists, err = sp.NodeExists(clusterID, nodeID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("node").WithParams(clusterID, nodeID)
	}

	// insert the node instance
	stmt, names := qb.Insert(clusterNodeTable).Columns("cluster_id", "node_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"cluster_id": clusterID,
		"node_id":    nodeID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add node")
	}

	return nil
}

// NodeExists checks if a node is linked to a cluster.
func (sp *ScyllaClusterProvider) NodeExists(clusterID string, nodeID string) (bool, derrors.Error) {

	if err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(clusterNodeTable).Columns("node_id").Where(qb.Eq("cluster_id")).
		Where(qb.Eq("node_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"cluster_id": clusterID,
		"node_id":    nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if node exists")
		}
	}

	return true, nil
}

// ListNodes returns a list of nodes in a cluster.
func (sp *ScyllaClusterProvider) ListNodes(clusterID string) ([]string, derrors.Error) {

	exists, err := sp.Exists(clusterID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("cluster").WithParams(clusterID)
	}

	stmt, names := qb.Select(clusterNodeTable).Columns("node_id").Where(qb.Eq("cluster_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"cluster_id": clusterID,
	})

	nodes := make([]string, 0)
	cqlErr := gocqlx.Select(&nodes, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list nodes")
	}

	return nodes, nil

}

// DeleteNode removes a node from a cluster.
func (sp *ScyllaClusterProvider) DeleteNode(clusterID string, nodeID string) derrors.Error {

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	// check if the node exists in the cluster
	exists, err := sp.NodeExists(clusterID, nodeID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("node").WithParams(clusterID, nodeID)
	}

	// delete app instance
	stmt, _ := qb.Delete(clusterNodeTable).Where(qb.Eq("cluster_id")).Where(qb.Eq("node_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, clusterID, nodeID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete node")
	}

	return nil
}

func (sp *ScyllaClusterProvider) Clear() derrors.Error {

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// delete clusters table
	err := sp.Session.Query("TRUNCATE TABLE clusters").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the clusters table")
		return derrors.AsError(err, "cannot truncate cluster table")
	}

	err = sp.Session.Query("TRUNCATE TABLE cluster_nodes").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the cluster_nodes table")
		return derrors.AsError(err, "cannot truncate node table")
	}

	return nil
}
