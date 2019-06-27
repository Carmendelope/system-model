package application

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

/*
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status
docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use nalej;
create type IF NOT EXISTS nalej.security_rule (organization_id text, app_descriptor_id text, rule_id text, name text, source_service_id text, source_port int, access int, auth_services list<text>, device_groups list<text>);
create type IF NOT EXISTS nalej.service_group_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, name text, description text, service_instances list<text>, policy int);
create type IF NOT EXISTS nalej.credential (username text, password text, email text);
create type IF NOT EXISTS nalej.deploy_spec (cpu bigint, memory bigint, replicas int, multi_cluster_replica_set boolean);
create type IF NOT EXISTS nalej.storage (size bigint, mount_path text, type int);
create type IF NOT EXISTS nalej.endpoint (type int, path text);
create type IF NOT EXISTS nalej.port (name text, internal_port int, exposed_port int, endpoint list<FROZEN<endpoint>>);
create type IF NOT EXISTS nalej.config_file (organization_id text, app_descriptor_id text, config_file_id text, content blob, mount_path text);
create type IF NOT EXISTS nalej.service_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, service_group_instance_id text, service_id text, service_instance_id text, name text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>, status int, endpoints list<FROZEN<endpoint_instance>>, deployed_on_cluster_id text,  run_arguments list<text>, info text, deployment_selectors map<text, text>);
create type IF NOT EXISTS nalej.service (organization_id text, app_descriptor_id text, service_group_id text, service_id text, name text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>,  run_arguments list<text>, deployment_selectors map<text, text>);
create type IF NOT EXISTS nalej.service_group (organization_id text, app_descriptor_id text, service_group_id text, name text, description text, services list<text>, policy int);
create table IF NOT EXISTS nalej.ApplicationInstances (organization_id text, app_descriptor_id text, app_instance_id text, name text, description text, configuration_options map<text, text>, environment_variables map<text, text>, labels map<text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group_instance>>, services list<FROZEN<service_instance>>, status int, PRIMARY KEY (app_instance_id));
create table IF NOT EXISTS nalej.ApplicationDescriptors (organization_id text, app_descriptor_id text, name text, description text, configuration_options map<text, text>, environment_variables map<text, text>, labels map <text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group>>, services list <FROZEN<service>>, PRIMARY KEY (app_descriptor_id));
*/

var _ = ginkgo.Describe("Scylla application provider", func() {

	var numApps = rand.Intn(50) + 1

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
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")

	}

	// create a provider and connect it
	sp := NewScyllaApplicationProvider(scyllaHost, scyllaPort, nalejKeySpace)

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

	ginkgo.It("Should be able to add Applications", func() {

		id := uuid.New().String()

		for i := 0; i < numApps; i++ {
			appID := fmt.Sprintf("00%d", i)

			app := CreateTestApplication(id, appID)

			err := sp.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			_ := sp.DeleteInstance(appID)
		}

	})

	ginkgo.It("Should be able to add Descriptors", func() {

		for i := 0; i < numApps; i++ {

			descriptor := CreateTestApplicationDescriptor(uuid.New().String())

			err := sp.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			_ := sp.DeleteDescriptor(descriptor.AppDescriptorId)
		}

	})
})
