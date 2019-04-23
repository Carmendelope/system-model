/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

/*
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status
docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use nalej;
create type IF NOT EXISTS nalej.operating_system_info (name text, version text);
create type IF NOT EXISTS nalej.cpu_info (manufacturer text, model text, architecture text, num_cores int);
create type IF NOT EXISTS nalej.networking_hardware_info (type text, link_capacity int);
create type IF NOT EXISTS nalej.hardware_info (cpus list<FROZEN<cpu_info>>, installed_ram int, net_interfaces list<FROZEN<networking_hardware_info>>);
create type IF NOT EXISTS nalej.storage_hardware_info (type text, total_capacity int);

create table IF NOT EXISTS nalej.Asset (organization_id text, asset_id text, agent_id text, show boolean, created int, labels map<text, text>, os FROZEN<operating_system_info>, hardware FROZEN<hardware_info>, storage FROZEN<storage_hardware_info>, eic_net_ip text, PRIMARY KEY (asset_id));
*/

var _ = ginkgo.Describe("Scylla asset provider", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	provider := NewScyllaAssetProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		provider.Disconnect()
	})

	RunTest(provider)

})