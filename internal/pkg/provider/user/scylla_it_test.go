package user

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)
/*
IT_SCYLLA_HOST=127.0.0.1
IT_NALEJ_KEYSPACE=nalej
RUN_INTEGRATION_TEST=true
*/
/*
before past test:

$ docker run --name scylla -p 9042:9042 -d scylladb/scylla -> launch docker image
$ docker exec -it scylla nodetool status -> Check the node is up

Prepare the database...

$ docker exec -it scylla cqlsh
cqlsh> create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
cqlsh> create table nalej.Users (organization_id text, email text, name text, photo_url text, member_since int, PRIMARY KEY (email));

 */


var _ = ginkgo.Describe("Scylla user provider", func(){

	var numUsers = 30

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

	// create a provider and connect it
	sp := NewScyllaUserProvider(scyllaHost, nalejKeySpace)
	err := sp.Connect()
	if err != nil {
		ginkgo.Fail("unable to connect")

	}

	ginkgo.BeforeSuite(func() {
		log.Debug().Msg("clearing table")
		sp.ClearTable()

	})

	// TODO: ask to Dani why the session has been closed
	// defer sp.Disconnect()

	RunTest(sp)

	ginkgo.It("Should be able to add user", func(){

		for i := 0; i < numUsers; i++ {
			email := fmt.Sprintf("email_%d@company.org", i)
			user := &entities.User{OrganizationId: "org",
				Email:       email,
				Name:        "name",
				MemberSince: int64(i)}
			err := sp.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})

})
