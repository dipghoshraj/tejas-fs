CREATE TABLE IF NOT EXISTS `data_objects`(
    id varchar(60) NOT NULL,
    name varchar(100) NOT NULL,
    size integer,
    distributed boolean DEFAULT false,
    ext varchar(20) NOT NULL,
    replica_id varchar(60),
    total_chunks integer DEFAULT 0,
    entry_node_id varchar(60)
)

CREATE Table IF NOT EXISTS `Chunks`(
    id UUID PRIMARY KEY,
    file_id UUID REFERENCES files(file_id),
    sequence_number INT NOT NULL,
    checksum TEXT NOT NULL,
    location TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)