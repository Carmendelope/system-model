package node


// create table nalej.Nodes (organization_id text, cluster_id text, node_id text, ip text, labels map<text, text>, status int, state int, PRIMARY KEY(node_id))

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

const addNode = "INSERT INTO Nodes (organization_id, cluster_id, node_id, ip, labels, status, state) VALUES (?, ?, ?, ?, ?, ?, ?)"
const updateNode = "UPDATE Nodes SET organization_id = ?, cluster_id = ?, ip = ?, labels = ?, status = ?, state = ? WHERE node_id = ?"
const exitsNode = "SELECT node_id from Nodes where node_id = ?"
const selectNode = "SELECT organization_id, cluster_id, ip, labels, status, state from Nodes where node_id = ?"
const deleteNode = "delete  from Nodes where node_id = ?"

// organizationId string, clusterId string, nodeId string, ip string, labels map[string]string,
//	 status grpc_infrastructure_go.InfraStatus, state entities.NodeState
type ScyllaNodeProvider struct {
	Address string
	Keyspace string
	Session *gocql.Session
}

func NewScyllaNodeProvider (address string, keyspace string) * ScyllaNodeProvider {
	return &ScyllaNodeProvider{ address, keyspace, nil}
}

// connect to the database
func (sp *ScyllaNodeProvider) Connect() derrors.Error {

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

// disconnect from the database
func (sp *ScyllaNodeProvider) Disconnect () {

	if sp != nil {
		sp.Session.Close()
	}
}

// check that the session is created
func (sp *ScyllaNodeProvider) CheckConnection () derrors.Error {
	if sp.Session == nil{
		return derrors.NewGenericError("Session not created")
	}
	return nil
}
// Add a new node to the system.
func (sp *ScyllaNodeProvider) Add (node entities.Node) derrors.Error {

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to add the node")
		return err
	}

	// check if the user exists
	exists, existsErr := sp.Exists(node.NodeId)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(existsErr).DebugReport()).Msg("unable to add the node")
		return conversions.ToDerror(existsErr)
	}

	if  exists {
		log.Info().Str("node", node.NodeId).Msg("unable to add the node, it alredy exists")
		return derrors.NewInvalidArgumentError("Node alredy exists")
	}

	// insert a user
	cqlErr := sp.Session.Query(addNode,node.OrganizationId, node.ClusterId, node.NodeId, node.Ip, node.Labels, node.Status, node.State).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to add the node")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}
// Update an existing node in the system
func (sp *ScyllaNodeProvider) Update(node entities.Node) derrors.Error {

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to update the node")
		return err
	}

	// check if the user exists
	exists, existsErr := sp.Exists(node.NodeId)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(existsErr).DebugReport()).Msg("unable to update the node")
		return conversions.ToDerror(existsErr)
	}
	if  ! exists {
		log.Info().Str("node", node.NodeId).Msg("unable to update the node, node does not exist")
		return derrors.NewInvalidArgumentError("Node does not exist")
	}

	// insert a user
	cqlErr := sp.Session.Query(updateNode,node.OrganizationId, node.ClusterId, node.Ip, node.Labels, node.Status, node.State, node.NodeId).Exec()

	if cqlErr != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to update the node")
		return conversions.ToDerror(cqlErr)
	}
	return nil

}
// Exists checks if a node exists on the system.
func (sp *ScyllaNodeProvider) Exists(nodeID string) (bool, derrors.Error) {

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to check if the node exists")
		return false, err
	}

	// check if exists
	var recoveredNodeId string
	errExists := sp.Session.Query(exitsNode, nodeID).Scan(&recoveredNodeId)

	if err == gocql.ErrNotFound{
		return false, nil
	}

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(errExists).DebugReport()).Msg("failed node exists")
		return false, conversions.ToDerror(errExists)
	}

	return recoveredNodeId != "", nil
	// return true, nil
}
// Get a node.
func (sp *ScyllaNodeProvider) Get(nodeID string) (* entities.Node, derrors.Error) {

	err := sp.CheckConnection()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable get the node")
		return nil, err
	}

	// check if exists
	var organizationId, clusterId, ip string
	var labels map[string]string
	var status, state int

	errSelect := sp.Session.Query(selectNode, nodeID).Scan(&organizationId, &clusterId, &ip, &labels, &status, &state)

	if errSelect != nil {
		log.Info().Str("trace", conversions.ToDerror(errSelect).DebugReport()).Msg("failed getting node")
		return nil, conversions.ToDerror(errSelect)
	}

	return &entities.Node {OrganizationId:organizationId, ClusterId:clusterId, NodeId:nodeID, Ip:ip, Labels: labels,
	Status: entities.InfraStatus(status), State: entities.NodeState(state)}, nil
}
// Remove a node
func (sp *ScyllaNodeProvider) Remove(nodeID string) derrors.Error {

	// check connection
	err := sp.CheckConnection()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable remove the node")
		return  err
	}

	// check if the user exists
	exists, err := sp.Exists(nodeID)

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Str("nodeID", nodeID).Msg("unable to remove the node")
		return conversions.ToDerror(err)
	}
	if ! exists {
		log.Info().Str("node", nodeID).Msg("unable to remove the node, it not exists")
		return derrors.NewInvalidArgumentError("Node does not exit")
	}

	// delete the node
	cqlErr := sp.Session.Query(deleteNode, nodeID).Exec()

	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(cqlErr).DebugReport()).Msg("failed to delete the node")
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

func (sp *ScyllaNodeProvider) ClearTable() derrors.Error {

	err := sp.Session.Query("TRUNCATE TABLE Nodes").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return conversions.ToDerror(err)
	}

	return nil
}