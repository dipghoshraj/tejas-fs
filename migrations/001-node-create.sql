CREATE TABLE IF NOT EXISTS nodes(
    id UUID PRIMARY KEY,
    status text NOT NULL,
    capacity bigint NOT NULL,
    used_space bigint NOT NULL,
    last_heartbeat timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    volume_name varchar(50) NOT NULL,
    port varchar(5)
);