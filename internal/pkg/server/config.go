/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
)

// Config structure with the options for the system model.
type Config struct {
	// Address where the API service will listen requests.
	Port int
	// Use in-memory providers
	UseInMemoryProviders bool
}

// Validate the current configuration.
func (conf * Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	return nil
}

// Print the current configuration to the log system.
func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	if conf.UseInMemoryProviders {
		log.Info().Bool("UseInMemoryProviders", conf.UseInMemoryProviders).Msg("Using in-memory providers")
	}
}