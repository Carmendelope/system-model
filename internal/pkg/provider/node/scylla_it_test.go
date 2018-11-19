package node

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/gomega"

	//"fmt"
	//"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	//"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("Scylla node provider", func(){

	var numNodes = 30

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
	sp := NewScyllaNodeProvider(scyllaHost, nalejKeySpace)
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


		labels := make(map[string]string)
		labels ["lab1"] = "label1"
		for i := 0; i < numNodes; i++ {
			nodeId := fmt.Sprintf("Node%d", i)
			node := &entities.Node{OrganizationId:"organization", ClusterId:"cluster_id", NodeId:nodeId,
			Ip:"0.0.0.0", Labels:labels, Status:0, State:0}
			err := sp.Add(*node)
			gomega.Expect(err).To(gomega.Succeed())
		}

	})


})
