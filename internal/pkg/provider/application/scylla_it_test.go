package application

import (
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
)

/*
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla nodetool status

docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
use nalej;
create type nalej.rule (organization_id text, app_descriptor_id text, rule_id text, name text, source_service_id text, source_port int, access int, auth_services list<text>, device_groups list<text>);
create type nalej.service_group_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, name text, description text, service_instances list<text>, policy int);
create type nalej.credential (username text, password text, email text);
create type nalej.deploy_spec (cpu bigint, memory bigint, replicas int);
create type nalej.storage (size bigint, mount_path text, type int);
create type nalej.endpoint (type int, path text);
create type nalej.port (name text, internal_port int, exposed_port int, endpoint list<FROZEN<endpoint>>);
create type nalej.config_file (organization_id text, app_descriptor_id text, config_file_id text, content blob, mount_path text);

create type nalej.service (organization_id text, app_descriptor_id text, app_instance_id text, service_id text, name text, description text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>, status int);

create table nalej.Applications (organization_id text, app_descriptor_id text, app_instance_id text, name text, description text, configuration_options map<text, text>, environment_variables map<text, text>, labels map<text, text>, rules list<FROZEN<rule>>, groups list<FROZEN<service_group_instance>>, services list<FROZEN<service>>, status int, PRIMARY KEY (app_instance_id));

create table nalej.ApplicationDescriptors (organization_id text, app_descriptor_id text, name text, description text, configuration_options map<text, text>, environment_variables map<text, text>, labels map <text, text>, rules list<FROZEN<rule>>, groups list<FROZEN<service_group_instance>>, services list <FROZEN<service>>, PRIMARY KEY (app_descriptor_id));

insert into nalej.Applications (organization_id, app_descriptor_id, app_instance_id, name, description, configuration_options, environment_variables, labels, rules, groups, services, status )values('organization_id1', 'app_descriptor_id1', 'app_instance_id1', 'Name1', 'description1', {'config1':'config1'},{'env1':'env1'},{'label1':'label1'},[{organization_id: 'organization_id1', app_descriptor_id: 'app_descriptor_id1', rule_id:'RuleID', name:'Name', source_service_id: 'source_service_id1', source_port:10, access: 10, auth_services: ['auth'], device_groups: ['device']}],[{organization_id: 'organization_id1' , app_descriptor_id:'app_descriptor_id1', service_group_id: 'group_id1', app_instance_id: 'app_instance_id1',name: 'name_group', description: 'group description', service_instances:['lista1', 'lista2'], policy:1}],[{organization_id:'organization_id1',app_descriptor_id: 'app_descriptor_id1',app_instance_id:'app_instance1',service_id: 'service_id',name: 'name_service',description: 'description_service',type: 0,image:'./img/img.jpg',credentials: {username: 'Carmen de Lope', password:'******', email:'cdelope@daisho.group'},specs: {cpu:852222, memory:10000, replicas:3},storage:[{size: 123456, mount_path:'../../path', type:0}],exposed_ports:[{name: 'port1', internal_port: 80, exposed_port: 80, endpoint:[{type:0, path: 'path1'},{type:1, path:'path2'}]}],environment_variables:{'HOST_IP':'0.0.0.0', 'HOST_PORT':'9043'},configs:[{organization_id: 'organization_id1', app_descriptor_id:'app_descriptor_id1', config_file_id:'config_file_id1', mount_path:'text'}],labels:{'label1':'label1_value', 'label2':'label2_value'},deploy_after: ['deploy after this'],status: 1}], 1)

*/

var _ = ginkgo.Describe("Scylla application provider", func(){


	if ! utils.RunIntegrationTests() {
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

	// create a provider and connect it
	sp := NewScyllaApplicationProvider(scyllaHost, nalejKeySpace)
	err :=	sp.Connect()

	if err != nil {
		ginkgo.Fail("unable to connect")
	}

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)
})

