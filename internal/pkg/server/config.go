/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/version"
	"github.com/rs/zerolog/log"
)

// Config structure with the options for the system model.
type Config struct {
	// Debug level is active.
	Debug bool
	// Address where the API service will listen requests.
	Port int
	// Use in-memory providers
	UseInMemoryProviders bool
	// Use scyllaDBProviders
	UseDBScyllaProviders bool
	// Database Address
	ScyllaDBAddress string
	// ScyllaDBPort with the port
	ScyllaDBPort int
	// DataBase KeySpace
	KeySpace string
	// PublicHostDomain
	PublicHostDomain string
}

// Validate the current configuration.
func (conf *Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	if conf.UseDBScyllaProviders {
		if conf.ScyllaDBAddress == "" {
			return derrors.NewInvalidArgumentError("address must be specified to use dbScylla Providers")
		}
		if conf.KeySpace == "" {
			return derrors.NewInvalidArgumentError("keyspace must be specified to use dbScylla Providers")
		}
		if conf.ScyllaDBPort <= 0 {
			return derrors.NewInvalidArgumentError("port must be specified to use dbScylla Providers ")
		}
	}
	if !conf.UseDBScyllaProviders && !conf.UseInMemoryProviders {
		return derrors.NewInvalidArgumentError("a type of provider must be selected")
	}
	return nil
}

// Print the current configuration to the log system.
func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Bool("set", conf.Debug).Msg("Debug")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	if conf.UseInMemoryProviders {
		log.Info().Bool("UseInMemoryProviders", conf.UseInMemoryProviders).Msg("Using in-memory providers")
	}
	if conf.UseDBScyllaProviders {
		log.Info().Bool("UseDBScyllaProviders", conf.UseDBScyllaProviders).Msg("using dbScylla providers")
		log.Info().Str("URL", conf.ScyllaDBAddress).Str("KeySpace", conf.KeySpace).Int("Port", conf.ScyllaDBPort).Msg("ScyllaDB")
	}
	log.Info().Str("PublicHostDomain", conf.PublicHostDomain).Msg("Public Host Domain")
}
