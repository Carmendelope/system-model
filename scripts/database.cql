--------------
-- KEYSPACE --
--------------
create KEYSPACE IF NOT EXISTS nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};


------------------------
-- USER DEFINED TYPES --
------------------------
create type IF NOT EXISTS nalej.deploy_spec (cpu bigint, memory bigint, replicas int);
create type IF NOT EXISTS nalej.service_group_deployment_specs (replicas int, multi_cluster_replica boolean, deployment_selectors map<text, text>);
create type IF NOT EXISTS nalej.storage (size bigint, mount_path text, type int);
create type IF NOT EXISTS nalej.endpoint (type int, path text);
create type IF NOT EXISTS nalej.endpoint_instance (endpoint_instance_id text, type int, fqdn text, port int);
create type IF NOT EXISTS nalej.port (name text, internal_port int, exposed_port int, endpoint list<FROZEN<endpoint>>);
create type IF NOT EXISTS nalej.metadata (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, monitored_instance_id text, type int, instance_id list<text>, desired_replicas int, available_replicas int, unavailable_replicas int, status map<text, int>, info map<text, text>);
create type IF NOT EXISTS nalej.credential (username text, password text, email text, docker_repository text);
create type IF NOT EXISTS nalej.security_rule (organization_id text, app_descriptor_id text, rule_id text, name text, target_service_group_name text, target_service_name text, target_port int, access int, auth_service_group_name text, auth_services list<text>, device_group_names list<text>, device_group_ids list<text>);
create type IF NOT EXISTS nalej.config_file (organization_id text, app_descriptor_id text, config_file_id text, name text, content blob, mount_path text);
create type IF NOT EXISTS nalej.service_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, service_group_instance_id text, service_id text, service_instance_id text, name text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>, status int, endpoints list<FROZEN<endpoint_instance>>, deployed_on_cluster_id text,  run_arguments list<text>, info text);
create type IF NOT EXISTS nalej.service_group_instance (organization_id text, app_descriptor_id text, app_instance_id text, service_group_id text, service_group_instance_id text, name text, service_instances list<FROZEN<service_instance>>, policy int, status int, metadata frozen<metadata>, specs FROZEN<service_group_deployment_specs>, labels map<text, text>);
create type IF NOT EXISTS nalej.service (organization_id text, app_descriptor_id text, service_group_id text, service_id text, name text, type int, image text, credentials FROZEN <credential>, specs FROZEN<deploy_spec>,storage list<FROZEN<storage>>,exposed_ports list<FROZEN<port>>, environment_variables map<text, text>, configs list<FROZEN<config_file>>, labels map<text, text>,deploy_after list<text>,  run_arguments list<text>);
create type IF NOT EXISTS nalej.service_group (organization_id text, app_descriptor_id text, service_group_id text, name text, services list<FROZEN<service>>, policy int, specs FROZEN<service_group_deployment_specs>, labels map<text, text>);

create type IF NOT EXISTS nalej.operating_system_info (name text, version text, op_class int, architecture text);
create type IF NOT EXISTS nalej.cpu_info (manufacturer text, model text, architecture text, num_cores int);
create type IF NOT EXISTS nalej.networking_hardware_info (type text, link_capacity int);
create type IF NOT EXISTS nalej.hardware_info (cpus list<FROZEN<cpu_info>>, installed_ram bigint, net_interfaces list<FROZEN<networking_hardware_info>>);
create type IF NOT EXISTS nalej.storage_hardware_info (type text, total_capacity int);
create type IF NOT EXISTS nalej.agent_op_summary (operation_id text, timestamp int, status int, info text);
create type IF NOT EXISTS nalej.ec_op_summary (operation_id text, timestamp int, status int, info text);

create type IF NOT EXISTS nalej.inventory_location (geolocation text, geohash text);


create type IF NOT EXISTS nalej.instance_parameter(parameter_name text, value text);
create type IF NOT EXISTS nalej.descriptor_parameter(name text, description text, path text, type int, default_value text, category int, enum_values list<text>, required boolean);

create type IF NOT EXISTS nalej.app_network_member(member_id text, is_proxy boolean, ip text, created_at bigint);
create type IF NOT EXISTS nalej.service_proxy(organization_id text, app_instance_id text, service_group_instance_id text, service_instance_id text, service_group_id text, service_id text, cluster_id text, ip text, fqdn text);


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
create table IF NOT EXISTS nalej.Clusters (organization_id text, cluster_id text, name text, cluster_type int, hostname text, control_plane_hostname text, multitenant int, status int, labels map<text, text>, cordon boolean, PRIMARY KEY (cluster_id));
create table IF NOT EXISTS nalej.Cluster_Nodes (cluster_id text, node_id text, PRIMARY KEY (cluster_id, node_id));
create table IF NOT EXISTS nalej.ApplicationInstances (organization_id text, app_descriptor_id text, app_instance_id text, name text, configuration_options map<text, text>, environment_variables map<text, text>, labels map<text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group_instance>>, status int, metadata list<FROZEN<metadata>>, info text, PRIMARY KEY (app_instance_id));
create table IF NOT EXISTS nalej.ApplicationDescriptors (organization_id text, app_descriptor_id text, name text, configuration_options map<text, text>, environment_variables map<text, text>, labels map <text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group>>, parameters list<FROZEN<descriptor_parameter>>, PRIMARY KEY (app_descriptor_id));
create table IF NOT EXISTS nalej.ParametrizedDescriptors (organization_id text, app_descriptor_id text, app_instance_id text, name text, configuration_options map<text, text>, environment_variables map<text, text>, labels map <text, text>, rules list<FROZEN<security_rule>>, groups list<FROZEN<service_group>>, PRIMARY KEY (app_instance_id));


create table IF NOT EXISTS nalej.AppEntrypoints(organization_id text, app_instance_id text, service_group_instance_id text, service_instance_id text, port int, protocol int, endpoint_instance_id text, type int, fqdn text, global_fqdn text,  PRIMARY KEY ((organization_id, app_instance_id), service_group_instance_id, service_instance_id, port));

create table IF NOT EXISTS nalej.Devices (organization_id text, device_group_id text, device_id text, register_since bigint, labels map<text, text>, os FROZEN<operating_system_info>, hardware FROZEN<hardware_info>, storage list<FROZEN<storage_hardware_info>>, PRIMARY KEY ( (organization_id, device_group_id), device_id));
create table IF NOT EXISTS nalej.DeviceGroups (organization_id text, device_group_id text, name text, created bigint, labels map<text, text>, primary KEY (organization_id, device_group_id));

create table IF NOT EXISTS nalej.AppZtNetworks(organization_id text, app_instance_id text, zt_network_id text, vsa_list map<text,text>, available_proxies map<text,FROZEN<map<text,FROZEN<list<FROZEN<service_proxy>>>>>>,  PRIMARY KEY ((organization_id, app_instance_id), zt_network_id));
create table IF NOT EXISTS nalej.AppZtNetworkMembers(organization_id text, app_instance_id text, service_group_instance_id text, service_application_instance_id text, zt_network_id text, members map<text,FROZEN<app_network_member>>,  PRIMARY KEY ((organization_id, app_instance_id, service_group_instance_id, service_application_instance_id), zt_network_id));


create table IF NOT EXISTS nalej.Asset (organization_id text, edge_controller_id text, asset_id text, agent_id text, show boolean, created int, labels map<text, text>, os FROZEN<operating_system_info>, hardware FROZEN<hardware_info>, storage list<FROZEN<storage_hardware_info>>, eic_net_ip text, last_op_result FROZEN<agent_op_summary>, last_alive_timestamp int, PRIMARY KEY (asset_id));
create table IF NOT EXISTS nalej.Controller (organization_id text, edge_controller_id text, show boolean, created int, name text, labels map<text, text>, last_alive_timestamp int, location FROZEN<inventory_location>, os FROZEN<operating_system_info>, hardware FROZEN<hardware_info>, storage list<FROZEN<storage_hardware_info>>, last_op_result FROZEN<ec_op_summary>, PRIMARY KEY(edge_controller_id));

create table IF NOT EXISTS nalej.InstanceParameters(app_instance_id text, parameters list<FROZEN<instance_parameter>>, PRIMARY KEY (app_instance_id));


-----------
-- INDEX --
-----------
create index IF NOT EXISTS organizationName ON nalej.organizations (name);
create index IF NOT EXISTS deviceGroupName ON nalej.DeviceGroups (name);
create index IF NOT EXISTS entrypointFqdn ON nalej.AppEntrypoints (global_fqdn);
create index IF NOT EXISTS controllerOrg ON nalej.Controller (organization_id);
create index IF NOT EXISTS assetOrg ON nalej.Asset (organization_id);
create index IF NOT EXISTS assetEdgeController ON nalej.Asset (edge_controller_id);