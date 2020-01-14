/*
 * Copyright 2020 Nalej
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

 /*
 docker run --name scylla -p 9042:9042 -d scylladb/scylla
 docker exec -it scylla cqlsh

 create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
 create table nalej.organizationsetting (organization_id text, key text, value text, description text, PRIMARY KEY (organization_id, key));

 IT_SCYLLA_HOST=127.0.0.1
 RUN_INTEGRATION_TEST=true
 IT_NALEJ_KEYSPACE=nalej
 IT_SCYLLA_PORT=9042
 */

package organization_setting

import (
	"github.com/nalej/system-model/internal/pkg/utils"

	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var _ = ginkgo.Describe("Scylla organization setting provider", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")
	}
	scyllaPort, err := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if err != nil {
		ginkgo.Fail("error getting scylla port")
	}
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	sp := NewScyllaOrganizationSettingProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

})
