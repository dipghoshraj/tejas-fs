ALTER TABLE nodes
ADD column cluster_id UUID

ALTER TABLE nodes
ADD CONSTRAINT custer_node_id FOREIGN KEY (cluster_id) REFERENCES clusters (id);
