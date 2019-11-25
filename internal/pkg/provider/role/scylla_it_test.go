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

package role

/*
IT_SCYLLA_HOST=127.0.0.1
IT_NALEJ_KEYSPACE=nalej
RUN_INTEGRATION_TEST=true
*/
/*
before past test:

docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status

Prepare the database...

docker exec -it scylla cqlsh
create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
create table nalej.Roles (organization_id text, role_id text, name text, description text, internal boolean, created int, PRIMARY KEY (role_id));

*/

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/gomega"
	"strconv"

	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("Scylla user provider", func() {

	var numRoles = 20

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
	sp := NewSScyllaRoleProvider(scyllaHost, scyllaPort, nalejKeySpace)

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add user", func() {

		for i := 0; i < numRoles; i++ {
			roleID := fmt.Sprintf("Role_%d", i)
			role := &entities.Role{OrganizationId: "org",
				Name:        "name-" + roleID,
				Created:     int64(i),
				RoleId:      roleID,
				Internal:    true,
				Description: "desc-" + roleID}
			err := sp.Add(*role)
			gomega.Expect(err).To(gomega.Succeed())
		}
	})
})
