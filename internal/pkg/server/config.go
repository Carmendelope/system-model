/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
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
	if conf.UseDBScyllaProviders && conf.UseInMemoryProviders {
		return derrors.NewInvalidArgumentError("only one type of provider must be selected")
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
