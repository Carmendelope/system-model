/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

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

create table if not exists nalej.Connection_Instances
(
    organization_id      text,
    connection_id        text,
    source_instance_id   text,
    source_instance_name text,
    target_instance_id   text,
    target_instance_name text,
    inbound_name         text,
    outbound_name        text,
    outbound_required    boolean,
    status               int,
    ip_range             text,
    zt_network_id        text,
    primary key ((organization_id), source_instance_id, target_instance_id, inbound_name, outbound_name)
);
create index if not exists connectionInstanceTargetIndex ON nalej.Connection_Instances (target_instance_id);
create index IF NOT EXISTS ztNetworId on nalej.connection_instances (zt_network_id);


create table if not exists nalej.Connection_Instance_Links
(
    organization_id    text,
    connection_id      text,
    source_instance_id text,
    source_cluster_id  text,
    target_instance_id text,
    target_cluster_id  text,
    inbound_name       text,
    outbound_name      text,
    status             int,
    primary key ((organization_id), source_instance_id, target_instance_id, inbound_name, outbound_name,
                                    source_cluster_id, target_cluster_id
        )
);

create table if not exists nalej.ztnetworkconnection
(
    organization_id text,
    zt_network_id text,
    app_instance_id text,
    service_id text,
    zt_member text,
    zt_ip text,
    cluster_id text,
    side int,
    primary key ((organization_id, zt_network_id), app_instance_id, service_id, cluster_id)
);

Environment variables
IT_SCYLLA_HOST=127.0.0.1
RUN_INTEGRATION_TEST=true
IT_NALEJ_KEYSPACE=nalej
IT_SCYLLA_PORT=9042

*/

var _ = ginkgo.Describe("Scylla Application Network provider", func() {

	if !utils.RunIntegrationTests() {
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
	provider := NewScyllaApplicationNetworkProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		provider.Disconnect()
	})

	RunTest(provider)

})
