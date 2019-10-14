package node

/*
create table nalej.Nodes (organization_id text, cluster_id text, node_id text, ip text, labels map<text, text>, status int, state int, PRIMARY KEY(node_id));
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

var _ = ginkgo.Describe("Scylla node provider", func() {

	var numNodes = 30

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
	sp := NewScyllaNodeProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add user", func() {

		labels := make(map[string]string)
		labels["lab1"] = "label1"
		for i := 0; i < numNodes; i++ {
			node := &entities.Node{OrganizationId: fmt.Sprintf("Org_%d", i),
				ClusterId: fmt.Sprintf("Cluster%d", i),
				NodeId:    fmt.Sprintf("Node%d", i),
				Ip:        fmt.Sprintf("%d.%d.%d.%d", i, i, i, i),
				Labels:    labels, Status: 0, State: 0}
			err := sp.Add(*node)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})

})
