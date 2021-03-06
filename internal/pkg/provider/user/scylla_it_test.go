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
 */

package user

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

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
create table nalej.Users (organization_id text, email text, name text, photo_base64 text, member_since int, PRIMARY KEY (email));

*/

var _ = ginkgo.Describe("Scylla user provider", func() {

	var numUsers = 50

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}

	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))

	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")

	}

	// create a provider and connect it
	sp := NewScyllaUserProvider(scyllaHost, scyllaPort, nalejKeySpace)
	//err :=	sp.Connect()

	//if err != nil {
	//	ginkgo.Fail("unable to connect")
	//}

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add user", func() {

		for i := 0; i < numUsers; i++ {
			email := fmt.Sprintf("email_%d@company.org", i)
			user := &entities.User{OrganizationId: fmt.Sprintf("org_%d", i),
				Email:       email,
				Name:        fmt.Sprintf("name_%d", i),
				MemberSince: int64(i)}
			err := sp.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})

})
