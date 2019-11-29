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
 */

// This is an example of an executable command.

package commands

import (
	"github.com/nalej/system-model/internal/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		config.Debug = debugLevel
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&config.Port, "port", 8800, "Port to launch the System Model")
	runCmd.Flags().BoolVar(&config.UseInMemoryProviders, "useInMemoryProviders", false, "Whether in-memory providers should be used. ONLY for development")
	runCmd.Flags().BoolVar(&config.UseDBScyllaProviders, "useDBScyllaProviders", true, "Whether dbscylla providers should be used")
	runCmd.Flags().StringVar(&config.ScyllaDBAddress, "scyllaDBAddress", "", "address to connect to scylla database")
	runCmd.Flags().IntVar(&config.ScyllaDBPort, "scyllaDBPort", 9042, "port to connect to scylla database")
	runCmd.Flags().StringVar(&config.KeySpace, "scyllaDBKeyspace", "", "keyspace of scylla database")
	runCmd.Flags().StringVar(&config.PublicHostDomain, "publicHost", "nalej.cluster.local", "Public Hostname for the domain")

}
