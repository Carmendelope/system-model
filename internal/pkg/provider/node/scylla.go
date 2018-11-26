package node


import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const nodeTable = "nodes"
const nodeTablePK = "node_id"
const rowNotFound = "not found"

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

// --------------------------------------------------------------------------------------------------------------------

// Add a new node to the system.
func (sp *ScyllaNodeProvider) Add (node entities.Node) derrors.Error {

	//create table nalej.Nodes (organization_id text, cluster_id text, node_id text, ip text, labels map<text, text>, status int, state int, PRIMARY KEY(node_id))
	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(node.NodeId)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if  exists {
		return derrors.NewAlreadyExistsError(node.NodeId)
	}

	// insert a user

	stmt, names := qb.Insert(nodeTable).Columns("organization_id","cluster_id","node_id","ip","labels","status","state").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(node)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// Update an existing node in the system
func (sp *ScyllaNodeProvider) Update(node entities.Node) derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(node.NodeId)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError(node.NodeId)
	}

	// update a user
	stmt, names := qb.Update(nodeTable).Set("organization_id","cluster_id","ip","labels","status","state").Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(node)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

// Exists checks if a node exists on the system.
func (sp *ScyllaNodeProvider) Exists(nodeID string) (bool, derrors.Error) {

	var returnedId string

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(nodeTable).Columns(nodeTablePK).Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		nodeTablePK: nodeID })

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

// Get a node.
func (sp *ScyllaNodeProvider) Get(nodeID string) (* entities.Node, derrors.Error) {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return nil, err
	}

	var node entities.Node
	stmt, names := qb.Select(nodeTable).Where(qb.Eq(nodeTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		nodeTablePK: nodeID,
	})

	err := q.GetRelease(&node)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, conversions.ToDerror(err)
		}else{
			return nil, derrors.NewNotFoundError(nodeID)
		}
	}

	return &node, nil
}

// Remove a node
func (sp *ScyllaNodeProvider) Remove(nodeID string) derrors.Error {

	// check connection
	if err := sp.CheckConnection(); err != nil {
		return err
	}

	// check if the user exists
	exists, err := sp.Exists(nodeID)

	if err != nil {
		return conversions.ToDerror(err)
	}
	if ! exists {
		return derrors.NewNotFoundError("node").WithParams(nodeID)
	}

	// remove a user
	stmt, _ := qb.Delete(nodeTable).Where(qb.Eq(nodeTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, nodeID).Exec()

	if cqlErr != nil {
		return conversions.ToDerror(cqlErr)
	}

	return nil
}

func (sp *ScyllaNodeProvider) Clear() derrors.Error {

	if err := sp.CheckConnection(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE Nodes").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return conversions.ToDerror(err)
	}

	return nil
}