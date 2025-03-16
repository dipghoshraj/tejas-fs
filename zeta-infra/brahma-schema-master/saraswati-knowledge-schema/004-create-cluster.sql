CREATE TABLE IF NOT EXISTS clusters( 
    id UUID PRIMARY KEY , 
    name VARCHAR(255) , 
    nodes integer , 
    node_capacity integer , 
    total_capacity integer , 
    used_capacity integer , 
    ingress_node integer , 
    ingress_capacity integer , 
    auto_scaling boolean 
); 
    
CREATE UNIQUE INDEX cluster_name ON clusters (name)