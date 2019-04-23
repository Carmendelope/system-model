/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package scylladb

import (
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

// RowNotFoundMsg corresponds to the error message returned by ScyllaDB if the row is not found.
const RowNotFoundMsg = "not found"

// General purpose structure to be reused to build ScyllaDB providers on top sharing common functionality.
type ScyllaDB struct{
	Address  string
	Port     int
	Keyspace string
	Session  *gocql.Session
}

// Connect to the ScyllaDB cluster.
func (s * ScyllaDB) Connect() derrors.Error{
	// connect to the cluster
	conf := gocql.NewCluster(s.Address)
	conf.Keyspace = s.Keyspace
	conf.Port = s.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}
	s.Session = session
	return nil
}

// Disconnect from the database
func (s * ScyllaDB) Disconnect() {
	if s.Session != nil {
		s.Session.Close()
		s.Session = nil
	}
}

// CheckConnection checks that the session is created
func (s * ScyllaDB) CheckConnection() derrors.Error {
	if s.Session == nil {
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

// CheckAndConnect checks if the connection is set and tries to reconnect otherwise.
func (s * ScyllaDB) CheckAndConnect() derrors.Error {
	err := s.CheckConnection()
	if err != nil {
		log.Info().Msg("session no created, trying to reconnect...")
		// try to reconnect
		err = s.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}