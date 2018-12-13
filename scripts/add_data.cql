-- Cluster
-- organization_id = 'dbb6d661-769a-4e96-9eee-dd04da109e16'
-- cluster_id = 'cluster1-769a-4e96-9eee-dd04da109e00'

-- cluster
INSERT INTO nalej.Clusters (organization_id, cluster_id, name, description, cluster_type, hostname, control_plane_hostname, multitenant, status, cordon) VALUES
('dbb6d661-769a-4e96-9eee-dd04da109e16', 'cluster1-769a-4e96-9eee-dd04da109e00', 'cluster 1 - prueba', 'description cluster 1 - prueba', 1, 'nalej.cluster1.local', 'control.nalej.cluster1.local', 1, 1, true);

-- nodes
insert into nalej.Nodes(organization_id, cluster_id, node_id, ip,  status, state) VALUES ('dbb6d661-769a-4e96-9eee-dd04da109e16', 'cluster1-769a-4e96-9eee-dd04da109e00', 'node0001-769a-4e96-9eee-dd04da109e00', '127.0.0.1', 1,1);
insert into nalej.Nodes(organization_id, cluster_id, node_id, ip,  status, state) VALUES ('dbb6d661-769a-4e96-9eee-dd04da109e16', 'cluster1-769a-4e96-9eee-dd04da109e00', 'node0002-769a-4e96-9eee-dd04da109e00', '127.0.0.2', 1,1);
insert into nalej.Nodes(organization_id, cluster_id, node_id, ip,  status, state) VALUES ('dbb6d661-769a-4e96-9eee-dd04da109e16', 'cluster1-769a-4e96-9eee-dd04da109e00', 'node0003-769a-4e96-9eee-dd04da109e00', '127.0.0.3', 1,1);

-- cluster-nodes
create table IF NOT EXISTS nalej.Cluster_Nodes (cluster_id text, node_id text, PRIMARY KEY (cluster_id, node_id));
insert into nalej.cluster_nodes (cluster_id, node_id) values ('cluster1-769a-4e96-9eee-dd04da109e00','node0001-769a-4e96-9eee-dd04da109e00');
insert into nalej.cluster_nodes (cluster_id, node_id) values ('cluster1-769a-4e96-9eee-dd04da109e00','node0002-769a-4e96-9eee-dd04da109e00');
insert into nalej.cluster_nodes (cluster_id, node_id) values ('cluster1-769a-4e96-9eee-dd04da109e00','node0003-769a-4e96-9eee-dd04da109e00');


-- users
insert into users (email, member_since, name, organization_id) VALUES ('operator@nalej.com', 1544608761, 'Operator', 'dbb6d661-769a-4e96-9eee-dd04da109e16');
