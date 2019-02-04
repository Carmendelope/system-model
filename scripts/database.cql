--------------
-- KEYSPACE --
--------------
create KEYSPACE IF NOT EXISTS nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

------------------------
-- USER DEFINED TYPES --
------------------------
create type IF NOT EXISTS nalej.service_group_deployment_specs (num_replicas int, multi_cluster_replica boolean);
create type IF NOT EXISTS nalej.metadata (organization_id text, app_descriptor_id text, app_instance_id text, monitored_instance_id text, type int, instance_id list<text>, desired_replicas int, available_replicas int, unavailable_replicas int, status map<text, int>, info map<text, text>);
create type IF NOT EXISTS nalej.credential (username text, password text, email text, docker_repository text);
create type IF NOT EXISTS nalej.service_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_id text, name text, description text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>, status int, endpoints list<text>, deployed_on_cluster_id text,  run_arguments list<text>, info text);
create type IF NOT EXISTS nalej.security_rule (organization_id text, app_descriptor_id text, rule_id text, name text, source_service_id text, source_port int, access int, auth_services list<text>, device_groups list<text>);
create type IF NOT EXISTS nalej.service_group_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, name text, description text, service_instances list<FROZEN<service_instance>>, policy int, status int, metadata frozen<metadata>, specs FROZEN<service_group_deployment_specs>);
create type IF NOT EXISTS nalej.deploy_spec (cpu bigint, memory bigint, replicas int);
create type IF NOT EXISTS nalej.storage (size bigint, mount_path text, type int);
create type IF NOT EXISTS nalej.endpoint (type int, path text);
create type IF NOT EXISTS nalej.port (name text, internal_port int, exposed_port int, endpoint list<FROZEN<endpoint>>);
create type IF NOT EXISTS nalej.config_file (organization_id text, app_descriptor_id text, config_file_id text, content blob, mount_path text);
create type IF NOT EXISTS nalej.service (organization_id text, app_descriptor_id text, service_id text, name text, description text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>,  run_arguments list<text>);
create type IF NOT EXISTS nalej.service_group (organization_id text, app_descriptor_id text, service_group_id text, name text, description text, services list<text>, policy int, specs FROZEN<service_group_deployment_specs>);
------------
-- TABLES --
------------
create table IF NOT EXISTS nalej.Users (organization_id text, email text, name text, photo_url text, member_since int, PRIMARY KEY (email));
create table IF NOT EXISTS nalej.Roles (organization_id text, role_id text, name text, description text, internal boolean, created int, PRIMARY KEY (role_id));
create table IF NOT EXISTS nalej.organizations (id text, name text, created bigint, PRIMARY KEY (id));
create table IF NOT EXISTS nalej.Organization_Clusters (organization_id text, cluster_id text, PRIMARY KEY (organization_id, cluster_id));
create table IF NOT EXISTS nalej.Organization_Nodes (organization_id text, node_id text, PRIMARY KEY (organization_id, node_id));
create table IF NOT EXISTS nalej.Organization_AppDescriptors (organization_id text, app_descriptor_id text, PRIMARY KEY (organization_id, app_descriptor_id));
create table IF NOT EXISTS nalej.Organization_AppInstances (organization_id text, app_instance_id text, PRIMARY KEY (organization_id, app_instance_id));
create table IF NOT EXISTS nalej.Organization_Users (organization_id text, email text, PRIMARY KEY (organization_id, email));
create table IF NOT EXISTS nalej.Organization_Roles (organization_id text, role_id text, PRIMARY KEY (organization_id, role_id));
create table IF NOT EXISTS nalej.Nodes (organization_id text, cluster_id text, node_id text, ip text, labels map<text, text>, status int, state int, PRIMARY KEY(node_id));
create table IF NOT EXISTS nalej.Clusters (organization_id text, cluster_id text, name text, description text, cluster_type int, hostname text, control_plane_hostname text, multitenant int, status int, labels map<text, text>, cordon boolean, PRIMARY KEY (cluster_id));
create table IF NOT EXISTS nalej.Cluster_Nodes (cluster_id text, node_id text, PRIMARY KEY (cluster_id, node_id));
create table IF NOT EXISTS nalej.ApplicationInstances (organization_id text, app_descriptor_id text, app_instance_id text, name text, configuration_options map<text, text>, environment_variables map<text, text>, labels map<text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group_instance>>, status int, metadata list<FROZEN<metadata>>, PRIMARY KEY (app_instance_id));
create table IF NOT EXISTS nalej.ApplicationDescriptors (organization_id text, app_descriptor_id text, name text, configuration_options map<text, text>, environment_variables map<text, text>, labels map <text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group>>, services list <FROZEN<service>>, PRIMARY KEY (app_descriptor_id));

create table IF NOT EXISTS nalej.Devices (organization_id text, device_group_id text, device_id text, register_since bigint, labels map<text, text>, PRIMARY KEY ( (organization_id, device_group_id), device_id));
create table IF NOT EXISTS nalej.DeviceGroups (organization_id text, device_group_id text, name text, created bigint, labels map<text, text>, primary KEY (organization_id, device_group_id));