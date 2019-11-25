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

package cluster

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
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

create type nalej.cluster_cilium_creds(cilium_id text, cilium_etcd_ca_crt text, cilium_etcd_crt text,cilium_etcd_key text);
create type nalej.cluster_istio_creds(cluster_name text, server_name text, ca_cert text, cluster_token text);
create type nalej.cluster_watch_info(name text, organization_id text, cluster_id text, ip text, network_type int, cilium_certs FROZEN<cluster_cilium_creds>, istio_certs FROZEN<cluster_istio_creds>);

create table IF NOT EXISTS nalej.Clusters (organization_id text, cluster_id text, name text, cluster_type int, hostname text, control_plane_hostname text, multitenant int, status int, labels map<text, text>, cordon boolean, cluster_watch FROZEN <cluster_watch_info>, last_alive_timestamp int, state int, PRIMARY KEY (cluster_id));
create table IF NOT EXISTS nalej.Cluster_Nodes (cluster_id text, node_id text, PRIMARY KEY (cluster_id, node_id));
*/

var _ = ginkgo.Describe("Scylla cluster provider", func() {

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
	sp := NewScyllaClusterProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add clusters", func() {

		numClusters := 100
		for i := 0; i < numClusters; i++ {

			clusterId := fmt.Sprintf("ClusterId_XX%d", i)
			cluster := CreateTestCluster(clusterId)

			err := sp.Add(*cluster)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})

})
